package optimus-ide-collabd_test

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	agentproto "github.com/optimus-ide-collab/optimus-ide-collab/v2/agent/proto"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbauthz"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbfake"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbtestutil"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbtime"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk/agentsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

// TestReportBoundaryLogsAgentRBAC guards against regressions where
// a pre-insert read (e.g. GetBoundarySessionByID) would be silently denied for
// agents and prevent session creation.
func TestReportBoundaryLogsAgentRBAC(t *testing.T) {
	t.Parallel()

	store, ps := dbtestutil.NewDB(t)
	client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{Database: store, Pubsub: ps})
	user := optimus-ide-collabdtest.CreateFirstUser(t, client)
	r := dbfake.WorkspaceBuild(t, store, database.WorkspaceTable{
		OrganizationID: user.OrganizationID,
		OwnerID:        user.UserID,
	}).WithAgent().Do()

	ctx := testutil.Context(t, testutil.WaitLong)

	// Connect as a real workspace agent.
	ac := agentsdk.New(client.URL, agentsdk.WithFixedToken(r.AgentToken))
	conn, err := ac.ConnectRPC(ctx)
	require.NoError(t, err)
	defer conn.Close()

	agentClient := agentproto.NewDRPCAgentClient(conn)
	sessionID := uuid.New()

	_, err = agentClient.ReportBoundaryLogs(ctx, &agentproto.ReportBoundaryLogsRequest{
		SessionId:           sessionID.String(),
		ConfinedProcessName: "claude-code",
		Logs: []*agentproto.BoundaryLog{
			{
				Allowed:        true,
				Time:           timestamppb.New(dbtime.Now()),
				SequenceNumber: 0,
				Resource: &agentproto.BoundaryLog_HttpRequest_{
					HttpRequest: &agentproto.BoundaryLog_HttpRequest{
						Method:      "GET",
						Url:         "https://example.com",
						MatchedRule: "domain=example.com",
					},
				},
			},
		},
	})
	require.NoError(t, err)

	// Verify persistence via the raw store: because ReportBoundaryLogs swallows
	// DB errors and returns success regardless, only a direct read proves the
	// session and log were actually persisted under agent RBAC.
	sess, err := store.GetBoundarySessionByID(ctx, sessionID)
	require.NoError(t, err, "session must be persisted")
	require.Equal(t, r.Agents[0].ID, sess.WorkspaceAgentID)

	logs, err := store.ListBoundaryLogsBySessionID(ctx, database.ListBoundaryLogsBySessionIDParams{
		SessionID: sessionID,
	})
	require.NoError(t, err)
	require.Len(t, logs, 1, "log must be persisted")

	// Assert that the agent subject cannot read boundary sessions.
	memberRole, err := rbac.RoleByName(rbac.RoleMember())
	require.NoError(t, err)
	agentSubject := rbac.Subject{
		ID:    r.Workspace.OwnerID.String(),
		Roles: rbac.Roles{memberRole},
		Scope: rbac.WorkspaceAgentScope(rbac.WorkspaceAgentScopeParams{
			WorkspaceID: r.Workspace.ID,
			OwnerID:     r.Workspace.OwnerID,
			TemplateID:  r.Workspace.TemplateID,
			VersionID:   r.Build.TemplateVersionID,
		}),
	}.WithCachedASTValue()

	auth := rbac.NewStrictCachingAuthorizer(prometheus.NewRegistry())
	acsPtr := &atomic.Pointer[dbauthz.AccessControlStore]{}
	var acs dbauthz.AccessControlStore = dbauthz.AGPLTemplateAccessControlStore{}
	acsPtr.Store(&acs)
	authzStore := dbauthz.New(store, auth, testutil.Logger(t), acsPtr)

	agentCtx := dbauthz.As(context.Background(), agentSubject)
	_, err = authzStore.GetBoundarySessionByID(agentCtx, sessionID)
	require.True(t, dbauthz.IsNotAuthorizedError(err),
		"agents must not be able to read boundary sessions, got: %v", err)
}
