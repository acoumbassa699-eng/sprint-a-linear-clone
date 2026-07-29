package agentapi

import (
	"context"
	"io"
	"net"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"golang.org/x/xerrors"
	"storj.io/drpc/drpcmux"
	"storj.io/drpc/drpcserver"
	"tailscale.com/tailcfg"

	"cdr.dev/slog/v3"
	agentproto "github.com/optimus-ide-collab/optimus-ide-collab/v2/agent/proto"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/agentapi/metadatabatcher"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/agentapi/resourcesmonitor"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/appearance"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/boundaryusage"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/connectionlog"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/pubsub"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/externalauth"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/notifications"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/portsharing"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/prometheusmetrics"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/tracing"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/workspacestats"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/wspubsub"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk/agentsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk/drpcsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/tailnet"
	tailnetproto "github.com/optimus-ide-collab/optimus-ide-collab/v2/tailnet/proto"
	"github.com/optimus-ide-collab/quartz"
)

const workspaceCacheRefreshInterval = 5 * time.Minute

// API implements the DRPC agent API interface from agent/proto. This struct is
// instantiated once per agent connection and kept alive for the duration of the
// session.
type API struct {
	opts Options
	*ManifestAPI
	*AnnouncementBannerAPI
	*StatsAPI
	*LifecycleAPI
	*AppsAPI
	*MetadataAPI
	*ResourcesMonitoringAPI
	*LogsAPI
	*ScriptsAPI
	*ConnLogAPI
	*SubAgentAPI
	*BoundaryLogsAPI
	*ContextAPI
	*tailnet.DRPCService

	cachedWorkspaceFields *CachedWorkspaceFields

	mu sync.Mutex
}

var _ agentproto.DRPCAgentServer = &API{}

type Options struct {
	AgentID           uuid.UUID
	OwnerID           uuid.UUID
	WorkspaceID       uuid.UUID
	OrganizationID    uuid.UUID
	TemplateVersionID uuid.UUID

	AuthenticatedCtx      context.Context
	Log                   slog.Logger
	Clock                 quartz.Clock
	Database              database.Store
	NotificationsEnqueuer notifications.Enqueuer
	Pubsub                pubsub.Pubsub
	// ContextDirtyMarker is the chatd-backed hydrate/dirty fan-out invoked
	// from PushContextState. Nil when chatd is disabled.
	ContextDirtyMarker                ContextDirtyMarker
	ConnectionLogger                  *atomic.Pointer[connectionlog.ConnectionLogger]
	DerpMapFn                         func() *tailcfg.DERPMap
	TailnetCoordinator                *atomic.Pointer[tailnet.Coordinator]
	StatsReporter                     *workspacestats.Reporter
	MetadataBatcher                   *metadatabatcher.Batcher
	AppearanceFetcher                 *atomic.Pointer[appearance.Fetcher]
	PublishWorkspaceUpdateFn          func(ctx context.Context, userID uuid.UUID, event wspubsub.WorkspaceEvent)
	PublishWorkspaceAgentLogsUpdateFn func(ctx context.Context, workspaceAgentID uuid.UUID, msg agentsdk.LogsNotifyMessage)
	NetworkTelemetryHandler           func(batch []*tailnetproto.TelemetryEvent)
	BoundaryUsageTracker              *boundaryusage.Tracker
	LifecycleMetrics                  *LifecycleMetrics
	PortSharer                        *atomic.Pointer[portsharing.PortSharer]

	AccessURL                 *url.URL
	AppHostname               string
	AgentStatsRefreshInterval time.Duration
	DisableDirectConnections  bool
	DerpForceWebSockets       bool
	DerpMapUpdateFrequency    time.Duration
	ExternalAuthConfigs       []*externalauth.Config
	Experiments               optimus-ide-collabsdk.Experiments

	UpdateAgentMetricsFn func(ctx context.Context, labels prometheusmetrics.AgentMetricLabels, metrics []*agentproto.Stats_Metric)
}

func New(opts Options, workspace database.Workspace, agent database.WorkspaceAgent) *API {
	if opts.Clock == nil {
		opts.Clock = quartz.NewReal()
	}

	api := &API{
		opts: opts,
		mu:   sync.Mutex{},
	}

	api.ManifestAPI = &ManifestAPI{
		AccessURL:                opts.AccessURL,
		AppHostname:              opts.AppHostname,
		ExternalAuthConfigs:      opts.ExternalAuthConfigs,
		DisableDirectConnections: opts.DisableDirectConnections,
		DerpForceWebSockets:      opts.DerpForceWebSockets,
		AgentFn:                  api.agent,
		Database:                 opts.Database,
		DerpMapFn:                opts.DerpMapFn,
		WorkspaceID:              opts.WorkspaceID,
	}

	// Don't cache details for prebuilds, though the cached fields will eventually be updated
	// by the refresh routine once the prebuild workspace is claimed.
	api.cachedWorkspaceFields = &CachedWorkspaceFields{}
	if !workspace.IsPrebuild() {
		api.cachedWorkspaceFields.UpdateValues(workspace)
	}

	api.AnnouncementBannerAPI = &AnnouncementBannerAPI{
		appearanceFetcher: opts.AppearanceFetcher,
	}

	api.ResourcesMonitoringAPI = &ResourcesMonitoringAPI{
		AgentID:               opts.AgentID,
		WorkspaceID:           opts.WorkspaceID,
		Clock:                 opts.Clock,
		Database:              opts.Database,
		NotificationsEnqueuer: opts.NotificationsEnqueuer,
		Debounce:              30 * time.Minute,

		Config: resourcesmonitor.Config{
			NumDatapoints:      20,
			CollectionInterval: 10 * time.Second,

			Alert: resourcesmonitor.AlertConfig{
				MinimumNOKsPercent:     20,
				ConsecutiveNOKsPercent: 50,
			},
		},
	}

	api.StatsAPI = &StatsAPI{
		AgentID:                   agent.ID,
		AgentName:                 agent.Name,
		Workspace:                 api.cachedWorkspaceFields,
		Database:                  opts.Database,
		Log:                       opts.Log,
		StatsReporter:             opts.StatsReporter,
		AgentStatsRefreshInterval: opts.AgentStatsRefreshInterval,
		Experiments:               opts.Experiments,
	}

	api.LifecycleAPI = &LifecycleAPI{
		AgentFn:                  api.agent,
		WorkspaceID:              opts.WorkspaceID,
		Database:                 opts.Database,
		Log:                      opts.Log,
		PublishWorkspaceUpdateFn: api.publishWorkspaceUpdate,
		Metrics:                  opts.LifecycleMetrics,
	}

	api.AppsAPI = &AppsAPI{
		AgentID:                  agent.ID,
		AgentFn:                  api.agent,
		Database:                 opts.Database,
		Log:                      opts.Log,
		Workspace:                api.cachedWorkspaceFields,
		PublishWorkspaceUpdateFn: api.publishWorkspaceUpdate,
		Clock:                    opts.Clock,
		NotificationsEnqueuer:    opts.NotificationsEnqueuer,
	}

	api.MetadataAPI = &MetadataAPI{
		AgentID:   agent.ID,
		Workspace: api.cachedWorkspaceFields,
		Database:  opts.Database,
		Log:       opts.Log,
		Batcher:   opts.MetadataBatcher,
	}

	api.LogsAPI = &LogsAPI{
		AgentFn:                           api.agent,
		Database:                          opts.Database,
		Log:                               opts.Log,
		PublishWorkspaceUpdateFn:          api.publishWorkspaceUpdate,
		PublishWorkspaceAgentLogsUpdateFn: opts.PublishWorkspaceAgentLogsUpdateFn,
	}

	api.ScriptsAPI = &ScriptsAPI{
		Database: opts.Database,
	}

	api.ConnLogAPI = &ConnLogAPI{
		AgentID:          agent.ID,
		AgentName:        agent.Name,
		ConnectionLogger: opts.ConnectionLogger,
		Database:         opts.Database,
		Workspace:        api.cachedWorkspaceFields,
		Log:              opts.Log,
	}

	api.DRPCService = &tailnet.DRPCService{
		CoordPtr:                opts.TailnetCoordinator,
		Logger:                  opts.Log,
		DerpMapUpdateFrequency:  opts.DerpMapUpdateFrequency,
		DerpMapFn:               opts.DerpMapFn,
		NetworkTelemetryHandler: opts.NetworkTelemetryHandler,
	}

	api.SubAgentAPI = &SubAgentAPI{
		OwnerID:        opts.OwnerID,
		OrganizationID: opts.OrganizationID,
		AgentFn:        api.agent,
		Log:            opts.Log,
		Clock:          opts.Clock,
		Database:       opts.Database,
		PortSharer:     opts.PortSharer,
	}

	api.BoundaryLogsAPI = &BoundaryLogsAPI{
		Log:                  opts.Log,
		Database:             opts.Database,
		AgentID:              opts.AgentID,
		WorkspaceID:          opts.WorkspaceID,
		OwnerID:              opts.OwnerID,
		TemplateID:           workspace.TemplateID,
		TemplateVersionID:    opts.TemplateVersionID,
		BoundaryUsageTracker: opts.BoundaryUsageTracker,
	}

	api.ContextAPI = &ContextAPI{
		AgentID:     agent.ID,
		Workspace:   api.cachedWorkspaceFields,
		Log:         opts.Log,
		Clock:       opts.Clock,
		Database:    opts.Database,
		DirtyMarker: opts.ContextDirtyMarker,
	}

	// Start background cache refresh loop to handle workspace changes
	// like prebuild claims where owner_id and other fields may be modified in the DB.
	go api.startCacheRefreshLoop(opts.AuthenticatedCtx)

	return api
}

func (a *API) Server(ctx context.Context) (*drpcserver.Server, error) {
	mux := drpcmux.New()
	err := agentproto.DRPCRegisterAgent(mux, a)
	if err != nil {
		return nil, xerrors.Errorf("register agent API protocol in DRPC mux: %w", err)
	}

	err = tailnetproto.DRPCRegisterTailnet(mux, a)
	if err != nil {
		return nil, xerrors.Errorf("register tailnet API protocol in DRPC mux: %w", err)
	}

	return drpcserver.NewWithOptions(&tracing.DRPCHandler{Handler: mux},
		drpcserver.Options{
			Manager: drpcsdk.DefaultDRPCOptions(nil),
			Log: func(err error) {
				if xerrors.Is(err, io.EOF) {
					return
				}
				a.opts.Log.Debug(ctx, "drpc server error", slog.Error(err))
			},
		},
	), nil
}

func (a *API) Serve(ctx context.Context, l net.Listener) error {
	server, err := a.Server(ctx)
	if err != nil {
		return xerrors.Errorf("create agent API server: %w", err)
	}

	if err := a.ResourcesMonitoringAPI.InitMonitors(ctx); err != nil {
		return xerrors.Errorf("initialize resource monitoring: %w", err)
	}

	return server.Serve(ctx, l)
}

func (a *API) agent(ctx context.Context) (database.WorkspaceAgent, error) {
	agent, err := a.opts.Database.GetWorkspaceAgentByID(ctx, a.opts.AgentID)
	if err != nil {
		return database.WorkspaceAgent{}, xerrors.Errorf("get workspace agent by id %q: %w", a.opts.AgentID, err)
	}
	return agent, nil
}

// refreshCachedWorkspace periodically updates the cached workspace fields.
// This ensures that changes like prebuild claims (which modify owner_id, name, etc.)
// are eventually reflected in the cache without requiring agent reconnection.
func (a *API) refreshCachedWorkspace(ctx context.Context) {
	ws, err := a.opts.Database.GetWorkspaceByID(ctx, a.opts.WorkspaceID)
	if err != nil {
		// Do not clear the cache on transient DB errors. Stale data is
		// preferable to no data, which forces callers to fall back to
		// expensive queries like GetWorkspaceByAgentID.
		a.opts.Log.Warn(ctx, "failed to refresh cached workspace fields", slog.Error(err))
		return
	}

	if ws.IsPrebuild() {
		return
	}

	// If we still have the same values, skip the update and logging calls.
	if a.cachedWorkspaceFields.identity.Equal(database.WorkspaceIdentityFromWorkspace(ws)) {
		return
	}
	// Update fields that can change during workspace lifecycle (e.g., AutostartSchedule)
	a.cachedWorkspaceFields.UpdateValues(ws)

	a.opts.Log.Debug(ctx, "refreshed cached workspace fields",
		slog.F("workspace_id", ws.ID),
		slog.F("owner_id", ws.OwnerID),
		slog.F("name", ws.Name))
}

// startCacheRefreshLoop runs a background goroutine that periodically refreshes
// the cached workspace fields. This is primarily needed to handle prebuild claims
// where the owner_id and other fields change while the agent connection persists.
func (a *API) startCacheRefreshLoop(ctx context.Context) {
	// Refresh every 5 minutes. This provides a reasonable balance between:
	// - Keeping cache fresh for prebuild claims and other workspace updates
	// - Minimizing unnecessary database queries
	ticker := a.opts.Clock.TickerFunc(ctx, workspaceCacheRefreshInterval, func() error {
		a.refreshCachedWorkspace(ctx)
		return nil
	}, "cache_refresh")

	// We need to wait on the ticker exiting.
	_ = ticker.Wait()

	a.opts.Log.Debug(ctx, "cache refresh loop exited, invalidating the workspace cache on agent API",
		slog.F("workspace_id", a.cachedWorkspaceFields.identity.ID),
		slog.F("owner_id", a.cachedWorkspaceFields.identity.OwnerUsername),
		slog.F("name", a.cachedWorkspaceFields.identity.Name))
	a.cachedWorkspaceFields.Clear()
}

func (a *API) publishWorkspaceUpdate(ctx context.Context, agentID uuid.UUID, kind wspubsub.WorkspaceEventKind) error {
	a.opts.PublishWorkspaceUpdateFn(ctx, a.opts.OwnerID, wspubsub.WorkspaceEvent{
		Kind:        kind,
		WorkspaceID: a.opts.WorkspaceID,
		AgentID:     &agentID,
	})
	return nil
}
