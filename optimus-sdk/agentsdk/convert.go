package agentsdk

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/xerrors"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/agent/proto"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/tailnet"
)

// ManifestFromProto converts the proto manifest to the SDK Manifest.
// Secrets are intentionally NOT included on the returned Manifest:
// keeping them off of the SDK type makes it impossible for any code
// path that only holds a *Manifest to leak secret values via
// logging, JSON encoding, fmt verbs, or debug endpoints.
func ManifestFromProto(manifest *proto.Manifest) (Manifest, error) {
	parentID := uuid.Nil
	if pid := manifest.GetParentId(); pid != nil {
		var err error
		parentID, err = uuid.FromBytes(pid)
		if err != nil {
			return Manifest{}, xerrors.Errorf("error converting workspace agent parent ID: %w", err)
		}
	}
	apps, err := AppsFromProto(manifest.Apps)
	if err != nil {
		return Manifest{}, xerrors.Errorf("error converting workspace agent apps: %w", err)
	}
	scripts, err := AgentScriptsFromProto(manifest.Scripts)
	if err != nil {
		return Manifest{}, xerrors.Errorf("error converting workspace agent scripts: %w", err)
	}
	agentID, err := uuid.FromBytes(manifest.AgentId)
	if err != nil {
		return Manifest{}, xerrors.Errorf("error converting workspace agent ID: %w", err)
	}
	workspaceID, err := uuid.FromBytes(manifest.WorkspaceId)
	if err != nil {
		return Manifest{}, xerrors.Errorf("error converting workspace ID: %w", err)
	}
	devcontainers, err := DevcontainersFromProto(manifest.Devcontainers)
	if err != nil {
		return Manifest{}, xerrors.Errorf("error converting workspace agent devcontainers: %w", err)
	}
	return Manifest{
		ParentID:                 parentID,
		AgentID:                  agentID,
		AgentName:                manifest.AgentName,
		OwnerName:                manifest.OwnerUsername,
		WorkspaceID:              workspaceID,
		WorkspaceName:            manifest.WorkspaceName,
		Apps:                     apps,
		Scripts:                  scripts,
		DERPMap:                  tailnet.DERPMapFromProto(manifest.DerpMap),
		DERPForceWebSockets:      manifest.DerpForceWebsockets,
		GitAuthConfigs:           int(manifest.GitAuthConfigs),
		EnvironmentVariables:     manifest.EnvironmentVariables,
		Directory:                manifest.Directory,
		VSCodePortProxyURI:       manifest.VsCodePortProxyUri,
		MOTDFile:                 manifest.MotdPath,
		DisableDirectConnections: manifest.DisableDirectConnections,
		Metadata:                 MetadataDescriptionsFromProto(manifest.Metadata),
		Devcontainers:            devcontainers,
	}, nil
}

// ProtoFromManifest converts the SDK Manifest to the proto manifest.
// It does not populate the proto's Secrets field because the SDK
// Manifest intentionally does not carry secrets (see ManifestFromProto).
func ProtoFromManifest(manifest Manifest) (*proto.Manifest, error) {
	apps, err := ProtoFromApps(manifest.Apps)
	if err != nil {
		return nil, xerrors.Errorf("convert workspace apps: %w", err)
	}
	return &proto.Manifest{
		ParentId:      manifest.ParentID[:],
		AgentId:       manifest.AgentID[:],
		AgentName:     manifest.AgentName,
		OwnerUsername: manifest.OwnerName,
		WorkspaceId:   manifest.WorkspaceID[:],
		WorkspaceName: manifest.WorkspaceName,
		// #nosec G115 - Safe conversion for GitAuthConfigs which is expected to be small and positive
		GitAuthConfigs:           uint32(manifest.GitAuthConfigs),
		EnvironmentVariables:     manifest.EnvironmentVariables,
		Directory:                manifest.Directory,
		VsCodePortProxyUri:       manifest.VSCodePortProxyURI,
		MotdPath:                 manifest.MOTDFile,
		DisableDirectConnections: manifest.DisableDirectConnections,
		DerpForceWebsockets:      manifest.DERPForceWebSockets,
		DerpMap:                  tailnet.DERPMapToProto(manifest.DERPMap),
		Scripts:                  ProtoFromScripts(manifest.Scripts),
		Apps:                     apps,
		Metadata:                 ProtoFromMetadataDescriptions(manifest.Metadata),
		Devcontainers:            ProtoFromDevcontainers(manifest.Devcontainers),
	}, nil
}

func MetadataDescriptionsFromProto(descriptions []*proto.WorkspaceAgentMetadata_Description) []optimus-ide-collabsdk.WorkspaceAgentMetadataDescription {
	ret := make([]optimus-ide-collabsdk.WorkspaceAgentMetadataDescription, len(descriptions))
	for i, description := range descriptions {
		ret[i] = MetadataDescriptionFromProto(description)
	}
	return ret
}

func ProtoFromMetadataDescriptions(descriptions []optimus-ide-collabsdk.WorkspaceAgentMetadataDescription) []*proto.WorkspaceAgentMetadata_Description {
	ret := make([]*proto.WorkspaceAgentMetadata_Description, len(descriptions))
	for i, d := range descriptions {
		ret[i] = ProtoFromMetadataDescription(d)
	}
	return ret
}

func MetadataDescriptionFromProto(description *proto.WorkspaceAgentMetadata_Description) optimus-ide-collabsdk.WorkspaceAgentMetadataDescription {
	return optimus-ide-collabsdk.WorkspaceAgentMetadataDescription{
		DisplayName: description.DisplayName,
		Key:         description.Key,
		Script:      description.Script,
		Interval:    int64(description.Interval.AsDuration()),
		Timeout:     int64(description.Timeout.AsDuration()),
	}
}

func ProtoFromMetadataDescription(d optimus-ide-collabsdk.WorkspaceAgentMetadataDescription) *proto.WorkspaceAgentMetadata_Description {
	return &proto.WorkspaceAgentMetadata_Description{
		DisplayName: d.DisplayName,
		Key:         d.Key,
		Script:      d.Script,
		Interval:    durationpb.New(time.Duration(d.Interval)),
		Timeout:     durationpb.New(time.Duration(d.Timeout)),
	}
}

func ProtoFromMetadataResult(r optimus-ide-collabsdk.WorkspaceAgentMetadataResult) *proto.WorkspaceAgentMetadata_Result {
	return &proto.WorkspaceAgentMetadata_Result{
		CollectedAt: timestamppb.New(r.CollectedAt),
		Age:         r.Age,
		Value:       r.Value,
		Error:       r.Error,
	}
}

func MetadataResultFromProto(r *proto.WorkspaceAgentMetadata_Result) optimus-ide-collabsdk.WorkspaceAgentMetadataResult {
	return optimus-ide-collabsdk.WorkspaceAgentMetadataResult{
		CollectedAt: r.GetCollectedAt().AsTime(),
		Age:         r.GetAge(),
		Value:       r.GetValue(),
		Error:       r.GetError(),
	}
}

func MetadataFromProto(m *proto.Metadata) Metadata {
	return Metadata{
		Key:                          m.GetKey(),
		WorkspaceAgentMetadataResult: MetadataResultFromProto(m.GetResult()),
	}
}

func AgentScriptsFromProto(protoScripts []*proto.WorkspaceAgentScript) ([]optimus-ide-collabsdk.WorkspaceAgentScript, error) {
	ret := make([]optimus-ide-collabsdk.WorkspaceAgentScript, len(protoScripts))
	for i, protoScript := range protoScripts {
		app, err := AgentScriptFromProto(protoScript)
		if err != nil {
			return nil, xerrors.Errorf("parse script %v: %w", i, err)
		}
		ret[i] = app
	}
	return ret, nil
}

func ProtoFromScripts(scripts []optimus-ide-collabsdk.WorkspaceAgentScript) []*proto.WorkspaceAgentScript {
	ret := make([]*proto.WorkspaceAgentScript, len(scripts))
	for i, script := range scripts {
		ret[i] = ProtoFromScript(script)
	}
	return ret
}

func AgentScriptFromProto(protoScript *proto.WorkspaceAgentScript) (optimus-ide-collabsdk.WorkspaceAgentScript, error) {
	id, err := uuid.FromBytes(protoScript.Id)
	if err != nil {
		return optimus-ide-collabsdk.WorkspaceAgentScript{}, xerrors.Errorf("parse id: %w", err)
	}

	logSourceID, err := uuid.FromBytes(protoScript.LogSourceId)
	if err != nil {
		return optimus-ide-collabsdk.WorkspaceAgentScript{}, xerrors.Errorf("parse log source id: %w", err)
	}

	return optimus-ide-collabsdk.WorkspaceAgentScript{
		ID:               id,
		LogSourceID:      logSourceID,
		LogPath:          protoScript.LogPath,
		Script:           protoScript.Script,
		Cron:             protoScript.Cron,
		RunOnStart:       protoScript.RunOnStart,
		RunOnStop:        protoScript.RunOnStop,
		StartBlocksLogin: protoScript.StartBlocksLogin,
		Timeout:          protoScript.Timeout.AsDuration(),
		DisplayName:      protoScript.DisplayName,
	}, nil
}

func ProtoFromScript(s optimus-ide-collabsdk.WorkspaceAgentScript) *proto.WorkspaceAgentScript {
	return &proto.WorkspaceAgentScript{
		Id:               s.ID[:],
		LogSourceId:      s.LogSourceID[:],
		LogPath:          s.LogPath,
		Script:           s.Script,
		Cron:             s.Cron,
		RunOnStart:       s.RunOnStart,
		RunOnStop:        s.RunOnStop,
		StartBlocksLogin: s.StartBlocksLogin,
		Timeout:          durationpb.New(s.Timeout),
		DisplayName:      s.DisplayName,
	}
}

func AppsFromProto(protoApps []*proto.WorkspaceApp) ([]optimus-ide-collabsdk.WorkspaceApp, error) {
	ret := make([]optimus-ide-collabsdk.WorkspaceApp, len(protoApps))
	for i, protoApp := range protoApps {
		app, err := AppFromProto(protoApp)
		if err != nil {
			return nil, xerrors.Errorf("parse app %v (%q): %w", i, protoApp.Slug, err)
		}
		ret[i] = app
	}
	return ret, nil
}

func ProtoFromApps(apps []optimus-ide-collabsdk.WorkspaceApp) ([]*proto.WorkspaceApp, error) {
	ret := make([]*proto.WorkspaceApp, len(apps))
	var err error
	for i, a := range apps {
		ret[i], err = ProtoFromApp(a)
		if err != nil {
			return nil, err
		}
	}
	return ret, nil
}

func AppFromProto(protoApp *proto.WorkspaceApp) (optimus-ide-collabsdk.WorkspaceApp, error) {
	id, err := uuid.FromBytes(protoApp.Id)
	if err != nil {
		return optimus-ide-collabsdk.WorkspaceApp{}, xerrors.Errorf("parse id: %w", err)
	}

	sharingLevel := optimus-ide-collabsdk.WorkspaceAppSharingLevel(strings.ToLower(protoApp.SharingLevel.String()))
	if _, ok := optimus-ide-collabsdk.MapWorkspaceAppSharingLevels[sharingLevel]; !ok {
		return optimus-ide-collabsdk.WorkspaceApp{}, xerrors.Errorf("unknown app sharing level: %v (%q)", protoApp.SharingLevel, protoApp.SharingLevel.String())
	}

	health := optimus-ide-collabsdk.WorkspaceAppHealth(strings.ToLower(protoApp.Health.String()))
	if _, ok := optimus-ide-collabsdk.MapWorkspaceAppHealths[health]; !ok {
		return optimus-ide-collabsdk.WorkspaceApp{}, xerrors.Errorf("unknown app health: %v (%q)", protoApp.Health, protoApp.Health.String())
	}

	return optimus-ide-collabsdk.WorkspaceApp{
		ID:            id,
		URL:           protoApp.Url,
		External:      protoApp.External,
		Slug:          protoApp.Slug,
		DisplayName:   protoApp.DisplayName,
		Command:       protoApp.Command,
		Icon:          protoApp.Icon,
		Subdomain:     protoApp.Subdomain,
		SubdomainName: protoApp.SubdomainName,
		SharingLevel:  sharingLevel,
		Healthcheck: optimus-ide-collabsdk.Healthcheck{
			URL:       protoApp.Healthcheck.Url,
			Interval:  int32(protoApp.Healthcheck.Interval.AsDuration().Seconds()),
			Threshold: protoApp.Healthcheck.Threshold,
		},
		Health: health,
		Hidden: protoApp.Hidden,
	}, nil
}

func ProtoFromApp(a optimus-ide-collabsdk.WorkspaceApp) (*proto.WorkspaceApp, error) {
	sharingLevel, ok := proto.WorkspaceApp_SharingLevel_value[strings.ToUpper(string(a.SharingLevel))]
	if !ok {
		return nil, xerrors.Errorf("unknown sharing level %s", a.SharingLevel)
	}
	health, ok := proto.WorkspaceApp_Health_value[strings.ToUpper(string(a.Health))]
	if !ok {
		return nil, xerrors.Errorf("unknown health %s", a.Health)
	}
	return &proto.WorkspaceApp{
		Id:            a.ID[:],
		Url:           a.URL,
		External:      a.External,
		Slug:          a.Slug,
		DisplayName:   a.DisplayName,
		Command:       a.Command,
		Icon:          a.Icon,
		Subdomain:     a.Subdomain,
		SubdomainName: a.SubdomainName,
		SharingLevel:  proto.WorkspaceApp_SharingLevel(sharingLevel),
		Healthcheck: &proto.WorkspaceApp_Healthcheck{
			Url:       a.Healthcheck.URL,
			Interval:  durationpb.New(time.Duration(a.Healthcheck.Interval) * time.Second),
			Threshold: a.Healthcheck.Threshold,
		},
		Health: proto.WorkspaceApp_Health(health),
		Hidden: a.Hidden,
	}, nil
}

func ServiceBannerFromProto(sbp *proto.ServiceBanner) optimus-ide-collabsdk.BannerConfig {
	return optimus-ide-collabsdk.BannerConfig{
		Enabled:         sbp.GetEnabled(),
		Message:         sbp.GetMessage(),
		BackgroundColor: sbp.GetBackgroundColor(),
	}
}

func ProtoFromServiceBanner(sb optimus-ide-collabsdk.BannerConfig) *proto.ServiceBanner {
	return &proto.ServiceBanner{
		Enabled:         sb.Enabled,
		Message:         sb.Message,
		BackgroundColor: sb.BackgroundColor,
	}
}

func BannerConfigFromProto(sbp *proto.BannerConfig) optimus-ide-collabsdk.BannerConfig {
	return optimus-ide-collabsdk.BannerConfig{
		Enabled:         sbp.GetEnabled(),
		Message:         sbp.GetMessage(),
		BackgroundColor: sbp.GetBackgroundColor(),
	}
}

func ProtoFromBannerConfig(sb optimus-ide-collabsdk.BannerConfig) *proto.BannerConfig {
	return &proto.BannerConfig{
		Enabled:         sb.Enabled,
		Message:         sb.Message,
		BackgroundColor: sb.BackgroundColor,
	}
}

func ProtoFromSubsystems(ss []optimus-ide-collabsdk.AgentSubsystem) ([]proto.Startup_Subsystem, error) {
	ret := make([]proto.Startup_Subsystem, len(ss))
	for i, s := range ss {
		pi, ok := proto.Startup_Subsystem_value[strings.ToUpper(string(s))]
		if !ok {
			return nil, xerrors.Errorf("unknown subsystem: %s", s)
		}
		ret[i] = proto.Startup_Subsystem(pi)
	}
	return ret, nil
}

func ProtoFromAppHealthsRequest(req PostAppHealthsRequest) (*proto.BatchUpdateAppHealthRequest, error) {
	pReq := &proto.BatchUpdateAppHealthRequest{}
	for id, h := range req.Healths {
		hp, ok := proto.AppHealth_value[strings.ToUpper(string(h))]
		if !ok {
			return nil, xerrors.Errorf("unknown app health: %s", h)
		}

		// Copy the ID, otherwise all updates will have the same ID (the last
		// one in the list).
		var idCopy uuid.UUID
		copy(idCopy[:], id[:])
		pReq.Updates = append(pReq.Updates, &proto.BatchUpdateAppHealthRequest_HealthUpdate{
			Id:     idCopy[:],
			Health: proto.AppHealth(hp),
		})
	}
	return pReq, nil
}

func ProtoFromLog(log Log) (*proto.Log, error) {
	lvl, ok := proto.Log_Level_value[strings.ToUpper(string(log.Level))]
	if !ok {
		return nil, xerrors.Errorf("unknown log level: %s", log.Level)
	}
	return &proto.Log{
		CreatedAt: timestamppb.New(log.CreatedAt),
		Output:    SanitizeLogOutput(log.Output),
		Level:     proto.Log_Level(lvl),
	}, nil
}

func ProtoFromLifecycle(req PostLifecycleRequest) (*proto.Lifecycle, error) {
	s, ok := proto.Lifecycle_State_value[strings.ToUpper(string(req.State))]
	if !ok {
		return nil, xerrors.Errorf("unknown lifecycle state: %s", req.State)
	}
	return &proto.Lifecycle{
		State:     proto.Lifecycle_State(s),
		ChangedAt: timestamppb.New(req.ChangedAt),
	}, nil
}

func LifecycleStateFromProto(s proto.Lifecycle_State) (optimus-ide-collabsdk.WorkspaceAgentLifecycle, error) {
	caps, ok := proto.Lifecycle_State_name[int32(s)]
	if !ok {
		return "", xerrors.Errorf("unknown lifecycle state: %d", s)
	}
	return optimus-ide-collabsdk.WorkspaceAgentLifecycle(strings.ToLower(caps)), nil
}

func ProtoFromLifecycleState(s optimus-ide-collabsdk.WorkspaceAgentLifecycle) (proto.Lifecycle_State, error) {
	caps, ok := proto.Lifecycle_State_value[strings.ToUpper(string(s))]
	if !ok {
		return 0, xerrors.Errorf("unknown lifecycle state: %s", s)
	}
	return proto.Lifecycle_State(caps), nil
}

func DevcontainersFromProto(pdcs []*proto.WorkspaceAgentDevcontainer) ([]optimus-ide-collabsdk.WorkspaceAgentDevcontainer, error) {
	ret := make([]optimus-ide-collabsdk.WorkspaceAgentDevcontainer, len(pdcs))
	for i, pdc := range pdcs {
		dc, err := DevcontainerFromProto(pdc)
		if err != nil {
			return nil, xerrors.Errorf("parse devcontainer %v: %w", i, err)
		}
		ret[i] = dc
	}
	return ret, nil
}

func DevcontainerFromProto(pdc *proto.WorkspaceAgentDevcontainer) (optimus-ide-collabsdk.WorkspaceAgentDevcontainer, error) {
	id, err := uuid.FromBytes(pdc.Id)
	if err != nil {
		return optimus-ide-collabsdk.WorkspaceAgentDevcontainer{}, xerrors.Errorf("parse id: %w", err)
	}
	var subagentID uuid.NullUUID
	if pdc.SubagentId != nil {
		subagentID.Valid = true
		subagentID.UUID, err = uuid.FromBytes(pdc.SubagentId)
		if err != nil {
			return optimus-ide-collabsdk.WorkspaceAgentDevcontainer{}, xerrors.Errorf("parse subagent id: %w", err)
		}
	}
	return optimus-ide-collabsdk.WorkspaceAgentDevcontainer{
		ID:              id,
		Name:            pdc.Name,
		WorkspaceFolder: pdc.WorkspaceFolder,
		ConfigPath:      pdc.ConfigPath,
		SubagentID:      subagentID,
	}, nil
}

func ProtoFromDevcontainers(dcs []optimus-ide-collabsdk.WorkspaceAgentDevcontainer) []*proto.WorkspaceAgentDevcontainer {
	ret := make([]*proto.WorkspaceAgentDevcontainer, len(dcs))
	for i, dc := range dcs {
		ret[i] = ProtoFromDevcontainer(dc)
	}
	return ret
}

func ProtoFromDevcontainer(dc optimus-ide-collabsdk.WorkspaceAgentDevcontainer) *proto.WorkspaceAgentDevcontainer {
	var subagentID []byte
	if dc.SubagentID.Valid {
		subagentID = dc.SubagentID.UUID[:]
	}

	return &proto.WorkspaceAgentDevcontainer{
		Id:              dc.ID[:],
		Name:            dc.Name,
		WorkspaceFolder: dc.WorkspaceFolder,
		ConfigPath:      dc.ConfigPath,
		SubagentId:      subagentID,
	}
}

func ProtoFromPatchAppStatus(pas PatchAppStatus) (*proto.UpdateAppStatusRequest, error) {
	state, ok := proto.UpdateAppStatusRequest_AppStatusState_value[strings.ToUpper(string(pas.State))]
	if !ok {
		return nil, xerrors.Errorf("Invalid state: %s", pas.State)
	}
	return &proto.UpdateAppStatusRequest{
		Slug:    pas.AppSlug,
		State:   proto.UpdateAppStatusRequest_AppStatusState(state),
		Message: pas.Message,
		Uri:     pas.URI,
	}, nil
}

func SecretsFromProto(protoSecrets []*proto.WorkspaceSecret) []WorkspaceSecret {
	ret := make([]WorkspaceSecret, len(protoSecrets))
	for i, s := range protoSecrets {
		ret[i] = WorkspaceSecret{
			EnvName:  s.EnvName,
			FilePath: s.FilePath,
			Value:    s.Value,
		}
	}
	return ret
}

func ProtoFromSecrets(secrets []WorkspaceSecret) []*proto.WorkspaceSecret {
	ret := make([]*proto.WorkspaceSecret, len(secrets))
	for i, s := range secrets {
		ret[i] = &proto.WorkspaceSecret{
			EnvName:  s.EnvName,
			FilePath: s.FilePath,
			Value:    s.Value,
		}
	}
	return ret
}
