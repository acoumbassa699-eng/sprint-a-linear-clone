package optimus-ide-collabd_test

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"cdr.dev/slog/v3/sloggers/slogtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/audit"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbgen"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/provisioner/echo"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/provisionersdk/proto"
)

func TestAuditLogs(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		client := optimus-ide-collabdtest.New(t, nil)
		user := optimus-ide-collabdtest.CreateFirstUser(t, client)

		err := client.CreateTestAuditLog(ctx, optimus-ide-collabsdk.CreateTestAuditLogRequest{
			ResourceID:     user.UserID,
			OrganizationID: user.OrganizationID,
		})
		require.NoError(t, err)

		alogs, err := client.AuditLogs(ctx, optimus-ide-collabsdk.AuditLogsRequest{
			Pagination: optimus-ide-collabsdk.Pagination{
				Limit: 1,
			},
		})
		require.NoError(t, err)

		require.Equal(t, int64(1), alogs.Count)
		require.Len(t, alogs.AuditLogs, 1)
	})

	t.Run("IncludeUser", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		client := optimus-ide-collabdtest.New(t, nil)
		user := optimus-ide-collabdtest.CreateFirstUser(t, client)
		client2, user2 := optimus-ide-collabdtest.CreateAnotherUser(t, client, user.OrganizationID, rbac.RoleOwner())

		err := client2.CreateTestAuditLog(ctx, optimus-ide-collabsdk.CreateTestAuditLogRequest{
			ResourceID:     user2.ID,
			OrganizationID: user.OrganizationID,
		})
		require.NoError(t, err)

		alogs, err := client.AuditLogs(ctx, optimus-ide-collabsdk.AuditLogsRequest{
			Pagination: optimus-ide-collabsdk.Pagination{
				Limit: 1,
			},
		})
		require.NoError(t, err)
		require.Equal(t, int64(1), alogs.Count)
		require.Len(t, alogs.AuditLogs, 1)

		// Make sure the returned user is fully populated.
		foundUser, err := client.User(ctx, user2.ID.String())
		foundUser.OrganizationIDs = []uuid.UUID{} // Not included.
		require.NoError(t, err)
		require.Equal(t, foundUser, *alogs.AuditLogs[0].User)

		// Delete the user and try again.  This is a soft delete so nothing should
		// change.  If users are hard deleted we should get nil, but there is no way
		// to test this at the moment.
		err = client.DeleteUser(ctx, user2.ID)
		require.NoError(t, err)

		alogs, err = client.AuditLogs(ctx, optimus-ide-collabsdk.AuditLogsRequest{
			Pagination: optimus-ide-collabsdk.Pagination{
				Limit: 1,
			},
		})
		require.NoError(t, err)
		require.Equal(t, int64(1), alogs.Count)
		require.Len(t, alogs.AuditLogs, 1)

		foundUser, err = client.User(ctx, user2.ID.String())
		foundUser.OrganizationIDs = []uuid.UUID{} // Not included.
		require.NoError(t, err)
		require.Equal(t, foundUser, *alogs.AuditLogs[0].User)
	})

	t.Run("WorkspaceBuildAuditLink", func(t *testing.T) {
		t.Parallel()

		var (
			ctx      = context.Background()
			client   = optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{IncludeProvisionerDaemon: true})
			user     = optimus-ide-collabdtest.CreateFirstUser(t, client)
			version  = optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, nil)
			template = optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID)
		)

		optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
		workspace := optimus-ide-collabdtest.CreateWorkspace(t, client, template.ID)
		optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, client, workspace.LatestBuild.ID)

		buildResourceInfo := audit.AdditionalFields{
			WorkspaceName: workspace.Name,
			BuildNumber:   strconv.FormatInt(int64(workspace.LatestBuild.BuildNumber), 10),
			BuildReason:   database.BuildReason(string(workspace.LatestBuild.Reason)),
		}

		wriBytes, err := json.Marshal(buildResourceInfo)
		require.NoError(t, err)

		err = client.CreateTestAuditLog(ctx, optimus-ide-collabsdk.CreateTestAuditLogRequest{
			Action:           optimus-ide-collabsdk.AuditActionStop,
			ResourceType:     optimus-ide-collabsdk.ResourceTypeWorkspaceBuild,
			ResourceID:       workspace.LatestBuild.ID,
			AdditionalFields: wriBytes,
			OrganizationID:   user.OrganizationID,
		})
		require.NoError(t, err)

		auditLogs, err := client.AuditLogs(ctx, optimus-ide-collabsdk.AuditLogsRequest{
			Pagination: optimus-ide-collabsdk.Pagination{
				Limit: 1,
			},
		})
		require.NoError(t, err)
		buildNumberString := strconv.FormatInt(int64(workspace.LatestBuild.BuildNumber), 10)
		require.Equal(t, auditLogs.AuditLogs[0].ResourceLink, fmt.Sprintf("/@%s/%s/builds/%s",
			workspace.OwnerName, workspace.Name, buildNumberString))
	})

	t.Run("Organization", func(t *testing.T) {
		t.Parallel()

		logger := slogtest.Make(t, &slogtest.Options{
			IgnoreErrors: true,
		})
		ctx := context.Background()
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
			Logger: &logger,
		})
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
		orgAdmin, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID, rbac.ScopedRoleOrgAdmin(owner.OrganizationID))

		err := client.CreateTestAuditLog(ctx, optimus-ide-collabsdk.CreateTestAuditLogRequest{
			ResourceID:     owner.UserID,
			OrganizationID: owner.OrganizationID,
		})
		require.NoError(t, err)

		// Add an extra audit log in another organization
		err = client.CreateTestAuditLog(ctx, optimus-ide-collabsdk.CreateTestAuditLogRequest{
			ResourceID:     owner.UserID,
			OrganizationID: uuid.New(),
		})
		require.NoError(t, err)

		// Fetching audit logs without an organization selector should only
		// return organization audit logs the org admin is an admin of.
		alogs, err := orgAdmin.AuditLogs(ctx, optimus-ide-collabsdk.AuditLogsRequest{
			Pagination: optimus-ide-collabsdk.Pagination{
				Limit: 5,
			},
		})
		require.NoError(t, err)
		require.Len(t, alogs.AuditLogs, 1)

		// Using the organization selector allows the org admin to fetch audit logs
		alogs, err = orgAdmin.AuditLogs(ctx, optimus-ide-collabsdk.AuditLogsRequest{
			SearchQuery: fmt.Sprintf("organization:%s", owner.OrganizationID.String()),
			Pagination: optimus-ide-collabsdk.Pagination{
				Limit: 5,
			},
		})
		require.NoError(t, err)
		require.Len(t, alogs.AuditLogs, 1)

		// Also try fetching by organization name
		organization, err := orgAdmin.Organization(ctx, owner.OrganizationID)
		require.NoError(t, err)

		alogs, err = orgAdmin.AuditLogs(ctx, optimus-ide-collabsdk.AuditLogsRequest{
			SearchQuery: fmt.Sprintf("organization:%s", organization.Name),
			Pagination: optimus-ide-collabsdk.Pagination{
				Limit: 5,
			},
		})
		require.NoError(t, err)
		require.Len(t, alogs.AuditLogs, 1)
	})

	t.Run("Organization404", func(t *testing.T) {
		t.Parallel()

		logger := slogtest.Make(t, &slogtest.Options{
			IgnoreErrors: true,
		})
		ctx := context.Background()
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
			Logger: &logger,
		})
		owner := optimus-ide-collabdtest.CreateFirstUser(t, client)
		orgAdmin, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID, rbac.ScopedRoleOrgAdmin(owner.OrganizationID))

		_, err := orgAdmin.AuditLogs(ctx, optimus-ide-collabsdk.AuditLogsRequest{
			SearchQuery: fmt.Sprintf("organization:%s", "random-name"),
			Pagination: optimus-ide-collabsdk.Pagination{
				Limit: 5,
			},
		})
		require.Error(t, err)
	})
}

func TestAuditLogsFilter(t *testing.T) {
	t.Parallel()

	t.Run("Filter", func(t *testing.T) {
		t.Parallel()

		var (
			ctx      = context.Background()
			client   = optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{IncludeProvisionerDaemon: true})
			user     = optimus-ide-collabdtest.CreateFirstUser(t, client)
			version  = optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, completeWithAgentAndApp())
			template = optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID)
		)

		optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
		workspace := optimus-ide-collabdtest.CreateWorkspace(t, client, template.ID)
		workspace.LatestBuild = optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, client, workspace.LatestBuild.ID)

		// Create two logs with "Create"
		err := client.CreateTestAuditLog(ctx, optimus-ide-collabsdk.CreateTestAuditLogRequest{
			OrganizationID: user.OrganizationID,
			Action:         optimus-ide-collabsdk.AuditActionCreate,
			ResourceType:   optimus-ide-collabsdk.ResourceTypeTemplate,
			ResourceID:     template.ID,
			Time:           time.Date(2022, 8, 15, 14, 30, 45, 100, time.UTC), // 2022-8-15 14:30:45
		})
		require.NoError(t, err)
		err = client.CreateTestAuditLog(ctx, optimus-ide-collabsdk.CreateTestAuditLogRequest{
			OrganizationID: user.OrganizationID,
			Action:         optimus-ide-collabsdk.AuditActionCreate,
			ResourceType:   optimus-ide-collabsdk.ResourceTypeUser,
			ResourceID:     user.UserID,
			Time:           time.Date(2022, 8, 16, 14, 30, 45, 100, time.UTC), // 2022-8-16 14:30:45
		})
		require.NoError(t, err)

		// Create one log with "Delete"
		err = client.CreateTestAuditLog(ctx, optimus-ide-collabsdk.CreateTestAuditLogRequest{
			OrganizationID: user.OrganizationID,
			Action:         optimus-ide-collabsdk.AuditActionDelete,
			ResourceType:   optimus-ide-collabsdk.ResourceTypeUser,
			ResourceID:     user.UserID,
			Time:           time.Date(2022, 8, 15, 14, 30, 45, 100, time.UTC), // 2022-8-15 14:30:45
		})
		require.NoError(t, err)

		// Create one log with "Start"
		err = client.CreateTestAuditLog(ctx, optimus-ide-collabsdk.CreateTestAuditLogRequest{
			OrganizationID: user.OrganizationID,
			Action:         optimus-ide-collabsdk.AuditActionStart,
			ResourceType:   optimus-ide-collabsdk.ResourceTypeWorkspaceBuild,
			ResourceID:     workspace.LatestBuild.ID,
			Time:           time.Date(2022, 8, 15, 14, 30, 45, 100, time.UTC), // 2022-8-15 14:30:45
		})
		require.NoError(t, err)

		// Create one log with "Stop"
		err = client.CreateTestAuditLog(ctx, optimus-ide-collabsdk.CreateTestAuditLogRequest{
			OrganizationID: user.OrganizationID,
			Action:         optimus-ide-collabsdk.AuditActionStop,
			ResourceType:   optimus-ide-collabsdk.ResourceTypeWorkspaceBuild,
			ResourceID:     workspace.LatestBuild.ID,
			Time:           time.Date(2022, 8, 15, 14, 30, 45, 100, time.UTC), // 2022-8-15 14:30:45
		})
		require.NoError(t, err)

		// Create one log with "Connect" and "Disconect".
		connectRequestID := uuid.New()
		err = client.CreateTestAuditLog(ctx, optimus-ide-collabsdk.CreateTestAuditLogRequest{
			OrganizationID: user.OrganizationID,
			Action:         optimus-ide-collabsdk.AuditActionConnect,
			RequestID:      connectRequestID,
			ResourceType:   optimus-ide-collabsdk.ResourceTypeWorkspaceAgent,
			ResourceID:     workspace.LatestBuild.Resources[0].Agents[0].ID,
			Time:           time.Date(2022, 8, 15, 14, 30, 45, 100, time.UTC), // 2022-8-15 14:30:45
		})
		require.NoError(t, err)

		err = client.CreateTestAuditLog(ctx, optimus-ide-collabsdk.CreateTestAuditLogRequest{
			OrganizationID: user.OrganizationID,
			Action:         optimus-ide-collabsdk.AuditActionDisconnect,
			RequestID:      connectRequestID,
			ResourceType:   optimus-ide-collabsdk.ResourceTypeWorkspaceAgent,
			ResourceID:     workspace.LatestBuild.Resources[0].Agents[0].ID,
			Time:           time.Date(2022, 8, 15, 14, 35, 0o0, 100, time.UTC), // 2022-8-15 14:35:00
		})
		require.NoError(t, err)

		// Create one log with "Open" and "Close".
		openRequestID := uuid.New()
		err = client.CreateTestAuditLog(ctx, optimus-ide-collabsdk.CreateTestAuditLogRequest{
			OrganizationID: user.OrganizationID,
			Action:         optimus-ide-collabsdk.AuditActionOpen,
			RequestID:      openRequestID,
			ResourceType:   optimus-ide-collabsdk.ResourceTypeWorkspaceApp,
			ResourceID:     workspace.LatestBuild.Resources[0].Agents[0].Apps[0].ID,
			Time:           time.Date(2022, 8, 15, 14, 30, 45, 100, time.UTC), // 2022-8-15 14:30:45
		})
		require.NoError(t, err)
		err = client.CreateTestAuditLog(ctx, optimus-ide-collabsdk.CreateTestAuditLogRequest{
			OrganizationID: user.OrganizationID,
			Action:         optimus-ide-collabsdk.AuditActionClose,
			RequestID:      openRequestID,
			ResourceType:   optimus-ide-collabsdk.ResourceTypeWorkspaceApp,
			ResourceID:     workspace.LatestBuild.Resources[0].Agents[0].Apps[0].ID,
			Time:           time.Date(2022, 8, 15, 14, 35, 0o0, 100, time.UTC), // 2022-8-15 14:35:00
		})
		require.NoError(t, err)

		// Test cases
		testCases := []struct {
			Name           string
			SearchQuery    string
			ExpectedResult int
			ExpectedError  bool
		}{
			{
				Name:           "FilterByCreateAction",
				SearchQuery:    "action:create",
				ExpectedResult: 2,
			},
			{
				Name:           "FilterByDeleteAction",
				SearchQuery:    "action:delete",
				ExpectedResult: 1,
			},
			{
				Name:           "FilterByUserResourceType",
				SearchQuery:    "resource_type:user",
				ExpectedResult: 2,
			},
			{
				Name:           "FilterByTemplateResourceType",
				SearchQuery:    "resource_type:template",
				ExpectedResult: 1,
			},
			{
				Name:           "FilterByEmail",
				SearchQuery:    "email:" + optimus-ide-collabdtest.FirstUserParams.Email,
				ExpectedResult: 9,
			},
			{
				Name:           "FilterByUsername",
				SearchQuery:    "username:" + optimus-ide-collabdtest.FirstUserParams.Username,
				ExpectedResult: 9,
			},
			{
				Name:           "FilterByResourceID",
				SearchQuery:    "resource_id:" + user.UserID.String(),
				ExpectedResult: 2,
			},
			{
				Name:          "FilterInvalidSingleValue",
				SearchQuery:   "invalid",
				ExpectedError: true,
			},
			{
				Name:          "FilterWithInvalidResourceType",
				SearchQuery:   "resource_type:invalid",
				ExpectedError: true,
			},
			{
				Name:          "FilterWithInvalidAction",
				SearchQuery:   "action:invalid",
				ExpectedError: true,
			},
			{
				Name:           "FilterOnCreateSingleDay",
				SearchQuery:    "action:create date_from:2022-08-15 date_to:2022-08-15",
				ExpectedResult: 1,
			},
			{
				Name:           "FilterOnCreateDateFrom",
				SearchQuery:    "action:create date_from:2022-08-15",
				ExpectedResult: 2,
			},
			{
				Name:           "FilterOnCreateDateTo",
				SearchQuery:    "action:create date_to:2022-08-15",
				ExpectedResult: 1,
			},
			{
				Name:           "FilterOnWorkspaceBuildStart",
				SearchQuery:    "resource_type:workspace_build action:start",
				ExpectedResult: 1,
			},
			{
				Name:           "FilterOnWorkspaceBuildStop",
				SearchQuery:    "resource_type:workspace_build action:stop",
				ExpectedResult: 1,
			},
			{
				Name:           "FilterOnWorkspaceBuildStartByInitiator",
				SearchQuery:    "resource_type:workspace_build action:start build_reason:initiator",
				ExpectedResult: 1,
			},
			{
				Name:           "FilterOnWorkspaceAgentConnect",
				SearchQuery:    "resource_type:workspace_agent action:connect",
				ExpectedResult: 1,
			},
			{
				Name:           "FilterOnWorkspaceAgentDisconnect",
				SearchQuery:    "resource_type:workspace_agent action:disconnect",
				ExpectedResult: 1,
			},
			{
				Name:           "FilterOnWorkspaceAgentConnectionRequestID",
				SearchQuery:    "resource_type:workspace_agent request_id:" + connectRequestID.String(),
				ExpectedResult: 2,
			},
			{
				Name:           "FilterOnWorkspaceAppOpen",
				SearchQuery:    "resource_type:workspace_app action:open",
				ExpectedResult: 1,
			},
			{
				Name:           "FilterOnWorkspaceAppClose",
				SearchQuery:    "resource_type:workspace_app action:close",
				ExpectedResult: 1,
			},
			{
				Name:           "FilterOnWorkspaceAppOpenRequestID",
				SearchQuery:    "resource_type:workspace_app request_id:" + openRequestID.String(),
				ExpectedResult: 2,
			},
		}

		for _, testCase := range testCases {
			// Test filtering
			t.Run(testCase.Name, func(t *testing.T) {
				t.Parallel()
				auditLogs, err := client.AuditLogs(ctx, optimus-ide-collabsdk.AuditLogsRequest{
					SearchQuery: testCase.SearchQuery,
				})
				if testCase.ExpectedError {
					require.Error(t, err, "expected error")
				} else {
					require.NoError(t, err, "fetch audit logs")
					require.Len(t, auditLogs.AuditLogs, testCase.ExpectedResult, "expected audit logs returned")
					require.Equal(t, testCase.ExpectedResult, int(auditLogs.Count), "expected audit log count returned")
				}
			})
		}
	})
}

func completeWithAgentAndApp() *echo.Responses {
	return &echo.Responses{
		Parse: echo.ParseComplete,
		ProvisionGraph: []*proto.Response{
			{
				Type: &proto.Response_Graph{
					Graph: &proto.GraphComplete{
						Resources: []*proto.Resource{
							{
								Type: "compute",
								Name: "main",
								Agents: []*proto.Agent{
									{
										Name:            "smith",
										OperatingSystem: "linux",
										Architecture:    "i386",
										Apps: []*proto.App{
											{
												Slug:        "app",
												DisplayName: "App",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// TestDeprecatedConnEvents tests the deprecated connection and disconnection
// events in the audit logs. These events are no longer created, but need to be
// returned by the API.
func TestDeprecatedConnEvents(t *testing.T) {
	t.Parallel()
	var (
		ctx            = context.Background()
		client, _, api = optimus-ide-collabdtest.NewWithAPI(t, &optimus-ide-collabdtest.Options{IncludeProvisionerDaemon: true})
		user           = optimus-ide-collabdtest.CreateFirstUser(t, client)
		version        = optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, completeWithAgentAndApp())
		template       = optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID)
	)

	optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
	workspace := optimus-ide-collabdtest.CreateWorkspace(t, client, template.ID)
	workspace.LatestBuild = optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, client, workspace.LatestBuild.ID)

	type additionalFields struct {
		audit.AdditionalFields
		ConnectionType string `json:"connection_type"`
	}

	sshFields := additionalFields{
		AdditionalFields: audit.AdditionalFields{
			WorkspaceName:  workspace.Name,
			BuildNumber:    "999",
			BuildReason:    "initiator",
			WorkspaceOwner: workspace.OwnerName,
			WorkspaceID:    workspace.ID,
		},
		ConnectionType: "SSH",
	}

	sshFieldsBytes, err := json.Marshal(sshFields)
	require.NoError(t, err)

	appFields := audit.AdditionalFields{
		WorkspaceName: workspace.Name,
		// Deliberately empty
		BuildNumber:    "",
		BuildReason:    "",
		WorkspaceOwner: workspace.OwnerName,
		WorkspaceID:    workspace.ID,
	}

	appFieldsBytes, err := json.Marshal(appFields)
	require.NoError(t, err)

	dbgen.AuditLog(t, api.Database, database.AuditLog{
		OrganizationID:   user.OrganizationID,
		Action:           database.AuditActionConnect,
		ResourceType:     database.ResourceTypeWorkspaceAgent,
		ResourceID:       workspace.LatestBuild.Resources[0].Agents[0].ID,
		ResourceTarget:   workspace.LatestBuild.Resources[0].Agents[0].Name,
		Time:             time.Date(2022, 8, 15, 14, 30, 45, 100, time.UTC), // 2022-8-15 14:30:45
		AdditionalFields: sshFieldsBytes,
	})

	dbgen.AuditLog(t, api.Database, database.AuditLog{
		OrganizationID:   user.OrganizationID,
		Action:           database.AuditActionDisconnect,
		ResourceType:     database.ResourceTypeWorkspaceAgent,
		ResourceID:       workspace.LatestBuild.Resources[0].Agents[0].ID,
		ResourceTarget:   workspace.LatestBuild.Resources[0].Agents[0].Name,
		Time:             time.Date(2022, 8, 15, 14, 35, 0o0, 100, time.UTC), // 2022-8-15 14:35:00
		AdditionalFields: sshFieldsBytes,
	})

	dbgen.AuditLog(t, api.Database, database.AuditLog{
		OrganizationID:   user.OrganizationID,
		UserID:           user.UserID,
		Action:           database.AuditActionOpen,
		ResourceType:     database.ResourceTypeWorkspaceApp,
		ResourceID:       workspace.LatestBuild.Resources[0].Agents[0].Apps[0].ID,
		ResourceTarget:   workspace.LatestBuild.Resources[0].Agents[0].Apps[0].Slug,
		Time:             time.Date(2022, 8, 15, 14, 30, 45, 100, time.UTC), // 2022-8-15 14:30:45
		AdditionalFields: appFieldsBytes,
	})

	connLog, err := client.AuditLogs(ctx, optimus-ide-collabsdk.AuditLogsRequest{
		SearchQuery: "action:connect",
	})
	require.NoError(t, err)
	require.Len(t, connLog.AuditLogs, 1)
	var sshOutFields additionalFields
	err = json.Unmarshal(connLog.AuditLogs[0].AdditionalFields, &sshOutFields)
	require.NoError(t, err)
	require.Equal(t, sshFields, sshOutFields)

	dcLog, err := client.AuditLogs(ctx, optimus-ide-collabsdk.AuditLogsRequest{
		SearchQuery: "action:disconnect",
	})
	require.NoError(t, err)
	require.Len(t, dcLog.AuditLogs, 1)
	err = json.Unmarshal(dcLog.AuditLogs[0].AdditionalFields, &sshOutFields)
	require.NoError(t, err)
	require.Equal(t, sshFields, sshOutFields)

	openLog, err := client.AuditLogs(ctx, optimus-ide-collabsdk.AuditLogsRequest{
		SearchQuery: "action:open",
	})
	require.NoError(t, err)
	require.Len(t, openLog.AuditLogs, 1)
	var appOutFields audit.AdditionalFields
	err = json.Unmarshal(openLog.AuditLogs[0].AdditionalFields, &appOutFields)
	require.NoError(t, err)
	require.Equal(t, appFields, appOutFields)
}
