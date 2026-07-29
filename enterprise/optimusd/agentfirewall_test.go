package optimus-ide-collabd_test

import (
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"

	"cdr.dev/slog/v3/sloggers/slogtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbauthz"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbfake"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbgen"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbtestutil"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbtime"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/util/ptr"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/optimus-ide-collabdenttest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/license"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

func TestAgentFirewallSessionByID(t *testing.T) {
	t.Parallel()

	// seedBoundarySession inserts a boundary session linked to a workspace agent.
	// Uses the raw DB store to avoid dbauthz permission and ordering constraints during setup.
	seedBoundarySession := func(t *testing.T, rawDB database.Store, ownerID, orgID uuid.UUID) (database.BoundarySession, database.WorkspaceTable) {
		t.Helper()

		resp := dbfake.WorkspaceBuild(t, rawDB, database.WorkspaceTable{
			OwnerID:        ownerID,
			OrganizationID: orgID,
		}).WithAgent().Do()

		require.NotEmpty(t, resp.Agents, "expected at least one agent")

		session := dbgen.BoundarySession(t, rawDB, database.BoundarySession{
			WorkspaceAgentID:    resp.Agents[0].ID,
			OwnerID:             uuid.NullUUID{UUID: ownerID, Valid: true},
			ConfinedProcessName: "claude-code",
		})
		return session, resp.Workspace
	}

	t.Run("Owner", func(t *testing.T) {
		t.Parallel()

		db, pubsub := dbtestutil.NewDB(t)
		ownerClient, _, owner := optimus-ide-collabdenttest.NewWithDatabase(t, &optimus-ide-collabdenttest.Options{
			Options: &optimus-ide-collabdtest.Options{
				Database: db,
				Pubsub:   pubsub,
			},
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureBoundary: 1,
				},
			},
		})
		ctx := testutil.Context(t, testutil.WaitLong)

		session, ws := seedBoundarySession(t, db, owner.UserID, owner.OrganizationID)

		//nolint:gocritic // Testing owner role.
		got, err := ownerClient.AgentFirewallSessionByID(ctx, session.ID)
		require.NoError(t, err)
		require.Equal(t, session.ID, got.ID)
		require.Equal(t, ws.OwnerID, got.OwnerID)
		require.Equal(t, session.ConfinedProcessName, got.ConfinedProcess)
		require.Equal(t, ws.ID, got.WorkspaceID)
	})

	t.Run("Auditor", func(t *testing.T) {
		t.Parallel()

		db, pubsub := dbtestutil.NewDB(t)
		ownerClient, _, owner := optimus-ide-collabdenttest.NewWithDatabase(t, &optimus-ide-collabdenttest.Options{
			Options: &optimus-ide-collabdtest.Options{
				Database: db,
				Pubsub:   pubsub,
			},
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureBoundary: 1,
				},
			},
		})
		ctx := testutil.Context(t, testutil.WaitLong)

		session, _ := seedBoundarySession(t, db, owner.UserID, owner.OrganizationID)

		auditorClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, ownerClient, owner.OrganizationID, rbac.RoleAuditor())

		got, err := auditorClient.AgentFirewallSessionByID(ctx, session.ID)
		require.NoError(t, err)
		require.Equal(t, session.ID, got.ID)
		require.Equal(t, "claude-code", got.ConfinedProcess)
	})

	t.Run("MemberDenied", func(t *testing.T) {
		t.Parallel()

		db, pubsub := dbtestutil.NewDB(t)
		ownerClient, _, owner := optimus-ide-collabdenttest.NewWithDatabase(t, &optimus-ide-collabdenttest.Options{
			Options: &optimus-ide-collabdtest.Options{
				Database: db,
				Pubsub:   pubsub,
			},
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureBoundary: 1,
				},
			},
		})
		ctx := testutil.Context(t, testutil.WaitLong)

		session, _ := seedBoundarySession(t, db, owner.UserID, owner.OrganizationID)

		memberClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, ownerClient, owner.OrganizationID)

		_, err := memberClient.AgentFirewallSessionByID(ctx, session.ID)
		var sdkErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &sdkErr)
		require.Equal(t, http.StatusNotFound, sdkErr.StatusCode())
	})

	t.Run("NotFound", func(t *testing.T) {
		t.Parallel()

		ownerClient, _ := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureBoundary: 1,
				},
			},
		})
		ctx := testutil.Context(t, testutil.WaitLong)

		//nolint:gocritic // Testing owner role.
		_, err := ownerClient.AgentFirewallSessionByID(ctx, uuid.New())
		var sdkErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &sdkErr)
		require.Equal(t, http.StatusNotFound, sdkErr.StatusCode())
	})
}

// TestInsertBoundaryLogs_AgentAuth verifies that a workspace agent context
// can insert boundary logs through the dbauthz layer. Create is user-scoped
// in the member role; the agent's owner ID must match the resource owner.
func TestInsertBoundaryLogs_AgentAuth(t *testing.T) {
	t.Parallel()

	rawDB, _ := dbtestutil.NewDB(t)
	authorizer := rbac.NewStrictAuthorizer(prometheus.NewRegistry())
	authzDB := dbauthz.New(rawDB, authorizer, slogtest.Make(t, nil), &atomic.Pointer[dbauthz.AccessControlStore]{})

	ctx := testutil.Context(t, testutil.WaitLong)

	// Seed a workspace with an agent.
	user := dbgen.User(t, rawDB, database.User{})
	org := dbgen.Organization(t, rawDB, database.Organization{})
	_ = dbgen.OrganizationMember(t, rawDB, database.OrganizationMember{
		UserID:         user.ID,
		OrganizationID: org.ID,
	})
	tmpl := dbgen.Template(t, rawDB, database.Template{
		OrganizationID: org.ID,
		CreatedBy:      user.ID,
	})
	tmplVer := dbgen.TemplateVersion(t, rawDB, database.TemplateVersion{
		TemplateID:     uuid.NullUUID{Valid: true, UUID: tmpl.ID},
		OrganizationID: org.ID,
		CreatedBy:      user.ID,
	})
	ws := dbgen.Workspace(t, rawDB, database.WorkspaceTable{
		OrganizationID: org.ID,
		TemplateID:     tmpl.ID,
		OwnerID:        user.ID,
	})
	job := dbgen.ProvisionerJob(t, rawDB, nil, database.ProvisionerJob{
		Type: database.ProvisionerJobTypeWorkspaceBuild,
	})
	build := dbgen.WorkspaceBuild(t, rawDB, database.WorkspaceBuild{
		JobID:             job.ID,
		WorkspaceID:       ws.ID,
		TemplateVersionID: tmplVer.ID,
	})
	resource := dbgen.WorkspaceResource(t, rawDB, database.WorkspaceResource{
		JobID: build.JobID,
	})
	agent := dbgen.WorkspaceAgent(t, rawDB, database.WorkspaceAgent{
		ResourceID: resource.ID,
	})

	// Insert a boundary session using the raw DB (no auth check).
	now := time.Now().UTC()
	sessionID := uuid.New()
	_, err := rawDB.InsertBoundarySession(ctx, database.InsertBoundarySessionParams{
		ID:                  sessionID,
		WorkspaceAgentID:    agent.ID,
		OwnerID:             uuid.NullUUID{UUID: user.ID, Valid: true},
		ConfinedProcessName: "claude-code",
		StartedAt:           now,
		UpdatedAt:           now,
	})
	require.NoError(t, err)

	// Build a workspace agent RBAC subject.
	memberRole, err := rbac.RoleByName(rbac.RoleMember())
	require.NoError(t, err)
	agentSubject := rbac.Subject{
		ID:    user.ID.String(),
		Roles: rbac.Roles{memberRole},
		Scope: rbac.WorkspaceAgentScope(rbac.WorkspaceAgentScopeParams{
			WorkspaceID: ws.ID,
			OwnerID:     user.ID,
			TemplateID:  tmpl.ID,
			VersionID:   tmplVer.ID,
		}),
	}.WithCachedASTValue()
	agentCtx := dbauthz.As(ctx, agentSubject)

	// Insert boundary logs through dbauthz with the correct owner.
	// User-scoped create succeeds because the agent subject ID matches.
	logID := uuid.New()
	_, err = authzDB.InsertBoundaryLogs(agentCtx, database.InsertBoundaryLogsParams{
		SessionID:      sessionID,
		OwnerID:        user.ID,
		ID:             []uuid.UUID{logID},
		SequenceNumber: []int32{1},
		CapturedAt:     []time.Time{now},
		CreatedAt:      []time.Time{now},
		Proto:          []string{"tcp"},
		Method:         []string{"connect"},
		Detail:         []string{"example.com:443"},
		MatchedRule:    []string{"allow-all"},
	})
	require.NoError(t, err, "agent should be able to insert boundary logs for own owner")

	// Verify the logs were actually persisted.
	got, err := rawDB.GetBoundaryLogByID(ctx, logID)
	require.NoError(t, err)
	require.Equal(t, sessionID, got.SessionID)

	// Inserting with a different owner ID must fail (user-scoped create).
	otherUser := dbgen.User(t, rawDB, database.User{})
	_, err = authzDB.InsertBoundaryLogs(agentCtx, database.InsertBoundaryLogsParams{
		SessionID:      sessionID,
		OwnerID:        otherUser.ID,
		ID:             []uuid.UUID{uuid.New()},
		SequenceNumber: []int32{2},
		CapturedAt:     []time.Time{now},
		CreatedAt:      []time.Time{now},
		Proto:          []string{"tcp"},
		Method:         []string{"connect"},
		Detail:         []string{"evil.com:443"},
		MatchedRule:    []string{"allow-all"},
	})
	require.Error(t, err, "agent must not insert boundary logs for a different owner")
}

func TestAgentFirewallSessionLogs(t *testing.T) {
	t.Parallel()

	type logOpt struct {
		SeqNum int32
		Proto  string
		Method string
		Detail string
		// Rule is the matched rule. Non-empty means the request was allowed.
		Rule string
	}

	// Creates a boundary session and returns a helper to insert logs.
	setupSession := func(t *testing.T, db database.Store, ownerID, orgID uuid.UUID) (database.BoundarySession, func(opts ...logOpt)) {
		t.Helper()

		resp := dbfake.WorkspaceBuild(t, db, database.WorkspaceTable{
			OwnerID:        ownerID,
			OrganizationID: orgID,
		}).WithAgent().Do()
		require.NotEmpty(t, resp.Agents, "expected at least one agent")

		session := dbgen.BoundarySession(t, db, database.BoundarySession{
			WorkspaceAgentID:    resp.Agents[0].ID,
			OwnerID:             uuid.NullUUID{UUID: ownerID, Valid: true},
			ConfinedProcessName: "claude-code",
		})

		insertLogs := func(opts ...logOpt) {
			t.Helper()
			//nolint:gocritic // Test seeding requires system context.
			sysCtx := dbauthz.AsSystemRestricted(t.Context())

			ids := make([]uuid.UUID, len(opts))
			seqNums := make([]int32, len(opts))
			capturedAts := make([]time.Time, len(opts))
			createdAts := make([]time.Time, len(opts))
			protos := make([]string, len(opts))
			methods := make([]string, len(opts))
			details := make([]string, len(opts))
			matchedRules := make([]string, len(opts))

			now := dbtime.Now()
			for i, o := range opts {
				ids[i] = uuid.New()
				seqNums[i] = o.SeqNum
				capturedAts[i] = now
				createdAts[i] = now
				protos[i] = o.Proto
				methods[i] = o.Method
				details[i] = o.Detail
				matchedRules[i] = o.Rule
			}

			_, err := db.InsertBoundaryLogs(sysCtx, database.InsertBoundaryLogsParams{
				ID:             ids,
				SessionID:      session.ID,
				OwnerID:        ownerID,
				SequenceNumber: seqNums,
				CapturedAt:     capturedAts,
				CreatedAt:      createdAts,
				Proto:          protos,
				Method:         methods,
				Detail:         details,
				MatchedRule:    matchedRules,
			})
			require.NoError(t, err, "insert boundary logs")
		}
		return session, insertLogs
	}

	// Creates an enterprise client with FeatureBoundary enabled.
	newEntClient := func(t *testing.T) (*optimus-ide-collabsdk.Client, database.Store, optimus-ide-collabsdk.CreateFirstUserResponse) {
		t.Helper()
		db, pubsub := dbtestutil.NewDB(t)
		client, _, firstUser := optimus-ide-collabdenttest.NewWithDatabase(t, &optimus-ide-collabdenttest.Options{
			Options: &optimus-ide-collabdtest.Options{
				Database: db,
				Pubsub:   pubsub,
			},
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureBoundary: 1,
				},
			},
		})
		return client, db, firstUser
	}

	t.Run("OK", func(t *testing.T) {
		t.Parallel()
		client, db, owner := newEntClient(t)

		session, insertLogs := setupSession(t, db, owner.UserID, owner.OrganizationID)
		insertLogs(
			logOpt{SeqNum: 0, Proto: "http", Method: "GET", Detail: "https://github.com/optimus-ide-collab/optimus-ide-collab", Rule: "domain=github.com"},
			logOpt{SeqNum: 1, Proto: "http", Method: "POST", Detail: "https://evil.com/exfil"},
			logOpt{SeqNum: 2, Proto: "http", Method: "POST", Detail: "https://api.anthropic.com/v1/messages", Rule: "domain=api.anthropic.com"},
		)

		ctx := testutil.Context(t, testutil.WaitLong)

		//nolint:gocritic // Testing owner role.
		resp, err := client.AgentFirewallSessionLogs(ctx, session.ID, optimus-ide-collabsdk.AgentFirewallSessionLogsParams{})
		require.NoError(t, err)
		require.Len(t, resp.Results, 3)

		require.Equal(t, int32(0), resp.Results[0].SequenceNumber)
		require.Equal(t, int32(1), resp.Results[1].SequenceNumber)
		require.Equal(t, int32(2), resp.Results[2].SequenceNumber)

		// Allowed request: MatchedRule is non-NULL.
		require.True(t, resp.Results[0].Allowed)
		require.Equal(t, "GET", resp.Results[0].Method)
		require.Equal(t, "https://github.com/optimus-ide-collab/optimus-ide-collab", resp.Results[0].Detail)
		require.NotNil(t, resp.Results[0].MatchedRule)
		require.Equal(t, "domain=github.com", *resp.Results[0].MatchedRule)

		// Denied request: no matched rule.
		require.False(t, resp.Results[1].Allowed)
		require.Equal(t, "POST", resp.Results[1].Method)
		require.Equal(t, "https://evil.com/exfil", resp.Results[1].Detail)
		require.Nil(t, resp.Results[1].MatchedRule)

		// Second allowed request.
		require.True(t, resp.Results[2].Allowed)
		require.Equal(t, "POST", resp.Results[2].Method)
		require.Equal(t, "https://api.anthropic.com/v1/messages", resp.Results[2].Detail)
		require.NotNil(t, resp.Results[2].MatchedRule)
		require.Equal(t, "domain=api.anthropic.com", *resp.Results[2].MatchedRule)
	})

	// Table-driven tests for sequence number filtering and limit.
	filterTests := []struct {
		name     string
		params   optimus-ide-collabsdk.AgentFirewallSessionLogsParams
		wantSeqs []int32
	}{
		{
			name:     "SeqAfterIncludesBound",
			params:   optimus-ide-collabsdk.AgentFirewallSessionLogsParams{SeqAfter: ptr.Ref(int64(0))},
			wantSeqs: []int32{0, 1, 2},
		},
		{
			name:     "SeqBeforeExcludesBound",
			params:   optimus-ide-collabsdk.AgentFirewallSessionLogsParams{SeqBefore: ptr.Ref(int64(2))},
			wantSeqs: []int32{0, 1},
		},
		{
			name:     "BetweenBoundsInclusiveExclusive",
			params:   optimus-ide-collabsdk.AgentFirewallSessionLogsParams{SeqAfter: ptr.Ref(int64(0)), SeqBefore: ptr.Ref(int64(2))},
			wantSeqs: []int32{0, 1},
		},
		{
			name:     "LimitCapsResults",
			params:   optimus-ide-collabsdk.AgentFirewallSessionLogsParams{Limit: ptr.Ref(int32(2))},
			wantSeqs: []int32{0, 1},
		},
	}
	for _, tc := range filterTests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client, db, owner := newEntClient(t)

			session, insertLogs := setupSession(t, db, owner.UserID, owner.OrganizationID)
			// Insert in reverse order to prove the endpoint sorts by
			// sequence_number regardless of DB insertion order.
			insertLogs(
				logOpt{SeqNum: 2, Proto: "http", Method: "GET", Detail: "https://c.com", Rule: "domain=c.com"},
				logOpt{SeqNum: 1, Proto: "http", Method: "GET", Detail: "https://b.com", Rule: "domain=b.com"},
				logOpt{SeqNum: 0, Proto: "http", Method: "GET", Detail: "https://a.com", Rule: "domain=a.com"},
			)

			ctx := testutil.Context(t, testutil.WaitLong)
			//nolint:gocritic // Testing owner role.
			resp, err := client.AgentFirewallSessionLogs(ctx, session.ID, tc.params)
			require.NoError(t, err)
			require.Len(t, resp.Results, len(tc.wantSeqs))
			for i, wantSeq := range tc.wantSeqs {
				require.Equal(t, wantSeq, resp.Results[i].SequenceNumber)
			}
		})
	}

	t.Run("BetweenTwoInterceptions", func(t *testing.T) {
		t.Parallel()
		client, db, owner := newEntClient(t)

		session, insertLogs := setupSession(t, db, owner.UserID, owner.OrganizationID)
		insertLogs(
			logOpt{SeqNum: 5, Proto: "http", Method: "POST", Detail: "https://api.anthropic.com/v1/messages", Rule: "domain=api.anthropic.com"},
			logOpt{SeqNum: 6, Proto: "http", Method: "GET", Detail: "https://github.com/optimus-ide-collab/optimus-ide-collab/pulls", Rule: "domain=github.com"},
			logOpt{SeqNum: 7, Proto: "http", Method: "GET", Detail: "https://evil.com/exfil"},
			logOpt{SeqNum: 11, Proto: "http", Method: "POST", Detail: "https://api.anthropic.com/v1/messages", Rule: "domain=api.anthropic.com"},
			logOpt{SeqNum: 12, Proto: "http", Method: "POST", Detail: "https://api.anthropic.com/v1/messages", Rule: "domain=api.anthropic.com"},
		)

		ctx := testutil.Context(t, testutil.WaitLong)
		//nolint:gocritic // Testing owner role.
		resp, err := client.AgentFirewallSessionLogs(ctx, session.ID, optimus-ide-collabsdk.AgentFirewallSessionLogsParams{
			SeqAfter:  ptr.Ref(int64(5)),
			SeqBefore: ptr.Ref(int64(12)),
		})
		require.NoError(t, err)
		require.Len(t, resp.Results, 4, "should return events at seq 5, 6, 7, 11")
		require.Equal(t, int32(5), resp.Results[0].SequenceNumber)
		require.Equal(t, int32(6), resp.Results[1].SequenceNumber)
		require.Equal(t, int32(7), resp.Results[2].SequenceNumber)
		require.Equal(t, int32(11), resp.Results[3].SequenceNumber)
	})

	t.Run("EmptySession", func(t *testing.T) {
		t.Parallel()
		client, db, owner := newEntClient(t)

		session, _ := setupSession(t, db, owner.UserID, owner.OrganizationID)

		ctx := testutil.Context(t, testutil.WaitLong)
		//nolint:gocritic // Testing owner role.
		resp, err := client.AgentFirewallSessionLogs(ctx, session.ID, optimus-ide-collabsdk.AgentFirewallSessionLogsParams{})
		require.NoError(t, err)
		require.Empty(t, resp.Results)
	})

	t.Run("NonexistentSession", func(t *testing.T) {
		t.Parallel()
		client, _, _ := newEntClient(t)

		ctx := testutil.Context(t, testutil.WaitLong)
		//nolint:gocritic // Testing owner role.
		resp, err := client.AgentFirewallSessionLogs(ctx, uuid.New(), optimus-ide-collabsdk.AgentFirewallSessionLogsParams{})
		require.NoError(t, err)
		require.Empty(t, resp.Results)
	})

	t.Run("MemberDenied", func(t *testing.T) {
		t.Parallel()

		db, pubsub := dbtestutil.NewDB(t)
		ownerClient, _, owner := optimus-ide-collabdenttest.NewWithDatabase(t, &optimus-ide-collabdenttest.Options{
			Options: &optimus-ide-collabdtest.Options{
				Database: db,
				Pubsub:   pubsub,
			},
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureBoundary: 1,
				},
			},
		})

		memberClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, ownerClient, owner.OrganizationID)

		ctx := testutil.Context(t, testutil.WaitLong)
		_, err := memberClient.AgentFirewallSessionLogs(ctx, uuid.New(), optimus-ide-collabsdk.AgentFirewallSessionLogsParams{})
		var sdkErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &sdkErr)
		require.Equal(t, http.StatusNotFound, sdkErr.StatusCode())
	})

	t.Run("AuditorAllowed", func(t *testing.T) {
		t.Parallel()

		db, pubsub := dbtestutil.NewDB(t)
		ownerClient, _, owner := optimus-ide-collabdenttest.NewWithDatabase(t, &optimus-ide-collabdenttest.Options{
			Options: &optimus-ide-collabdtest.Options{
				Database: db,
				Pubsub:   pubsub,
			},
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureBoundary: 1,
				},
			},
		})

		auditorClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, ownerClient, owner.OrganizationID, rbac.RoleAuditor())

		session, insertLogs := setupSession(t, db, owner.UserID, owner.OrganizationID)
		insertLogs(
			logOpt{SeqNum: 0, Proto: "http", Method: "GET", Detail: "https://a.com", Rule: "domain=a.com"},
		)

		ctx := testutil.Context(t, testutil.WaitLong)
		resp, err := auditorClient.AgentFirewallSessionLogs(ctx, session.ID, optimus-ide-collabsdk.AgentFirewallSessionLogsParams{})
		require.NoError(t, err)
		require.Len(t, resp.Results, 1)
	})
}
