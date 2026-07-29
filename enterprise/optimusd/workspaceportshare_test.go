package optimus-ide-collabd_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/util/ptr"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/optimus-ide-collabdenttest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/license"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

func TestWorkspacePortSharePublic(t *testing.T) {
	t.Parallel()

	ownerClient, owner := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
		Options: &optimus-ide-collabdtest.Options{IncludeProvisionerDaemon: true},
		LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
			Features: license.Features{optimus-ide-collabsdk.FeatureControlSharedPorts: 1},
		},
	})
	client, user := optimus-ide-collabdtest.CreateAnotherUser(t, ownerClient, owner.OrganizationID, rbac.RoleTemplateAdmin())
	r := setupWorkspaceAgent(t, client, optimus-ide-collabsdk.CreateFirstUserResponse{
		UserID:         user.ID,
		OrganizationID: owner.OrganizationID,
	}, 0)
	ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitShort)
	defer cancel()

	templ, err := client.Template(ctx, r.workspace.TemplateID)
	require.NoError(t, err)
	require.Equal(t, templ.MaxPortShareLevel, optimus-ide-collabsdk.WorkspaceAgentPortShareLevelOwner)

	// Try to update port share with template max port share level owner.
	_, err = client.UpsertWorkspaceAgentPortShare(ctx, r.workspace.ID, optimus-ide-collabsdk.UpsertWorkspaceAgentPortShareRequest{
		AgentName:  r.sdkAgent.Name,
		Port:       8080,
		ShareLevel: optimus-ide-collabsdk.WorkspaceAgentPortShareLevelPublic,
		Protocol:   optimus-ide-collabsdk.WorkspaceAgentPortShareProtocolHTTP,
	})
	require.Error(t, err, "Port sharing level not allowed")

	// Update the template max port share level to public
	client.UpdateTemplateMeta(ctx, r.workspace.TemplateID, optimus-ide-collabsdk.UpdateTemplateMeta{
		MaxPortShareLevel: ptr.Ref(optimus-ide-collabsdk.WorkspaceAgentPortShareLevelPublic),
	})

	// OK
	ps, err := client.UpsertWorkspaceAgentPortShare(ctx, r.workspace.ID, optimus-ide-collabsdk.UpsertWorkspaceAgentPortShareRequest{
		AgentName:  r.sdkAgent.Name,
		Port:       8080,
		ShareLevel: optimus-ide-collabsdk.WorkspaceAgentPortShareLevelPublic,
		Protocol:   optimus-ide-collabsdk.WorkspaceAgentPortShareProtocolHTTP,
	})
	require.NoError(t, err)
	require.EqualValues(t, optimus-ide-collabsdk.WorkspaceAgentPortShareLevelPublic, ps.ShareLevel)
}

func TestWorkspacePortShareOrganization(t *testing.T) {
	t.Parallel()

	ownerClient, owner := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
		Options: &optimus-ide-collabdtest.Options{IncludeProvisionerDaemon: true},
		LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
			Features: license.Features{optimus-ide-collabsdk.FeatureControlSharedPorts: 1},
		},
	})
	client, user := optimus-ide-collabdtest.CreateAnotherUser(t, ownerClient, owner.OrganizationID, rbac.RoleTemplateAdmin())
	r := setupWorkspaceAgent(t, client, optimus-ide-collabsdk.CreateFirstUserResponse{
		UserID:         user.ID,
		OrganizationID: owner.OrganizationID,
	}, 0)
	ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitShort)
	defer cancel()

	templ, err := client.Template(ctx, r.workspace.TemplateID)
	require.NoError(t, err)
	require.Equal(t, templ.MaxPortShareLevel, optimus-ide-collabsdk.WorkspaceAgentPortShareLevelOwner)

	// Try to update port share with template max port share level owner
	_, err = client.UpsertWorkspaceAgentPortShare(ctx, r.workspace.ID, optimus-ide-collabsdk.UpsertWorkspaceAgentPortShareRequest{
		AgentName:  r.sdkAgent.Name,
		Port:       8080,
		ShareLevel: optimus-ide-collabsdk.WorkspaceAgentPortShareLevelOrganization,
		Protocol:   optimus-ide-collabsdk.WorkspaceAgentPortShareProtocolHTTP,
	})
	require.Error(t, err, "Port sharing level not allowed")

	// Update the template max port share level to organization
	client.UpdateTemplateMeta(ctx, r.workspace.TemplateID, optimus-ide-collabsdk.UpdateTemplateMeta{
		MaxPortShareLevel: ptr.Ref(optimus-ide-collabsdk.WorkspaceAgentPortShareLevelOrganization),
	})

	// Try to share a port publicly with template max port share level organization
	_, err = client.UpsertWorkspaceAgentPortShare(ctx, r.workspace.ID, optimus-ide-collabsdk.UpsertWorkspaceAgentPortShareRequest{
		AgentName:  r.sdkAgent.Name,
		Port:       8080,
		ShareLevel: optimus-ide-collabsdk.WorkspaceAgentPortShareLevelPublic,
		Protocol:   optimus-ide-collabsdk.WorkspaceAgentPortShareProtocolHTTP,
	})
	require.Error(t, err, "Port sharing level not allowed")

	// OK
	ps, err := client.UpsertWorkspaceAgentPortShare(ctx, r.workspace.ID, optimus-ide-collabsdk.UpsertWorkspaceAgentPortShareRequest{
		AgentName:  r.sdkAgent.Name,
		Port:       8080,
		ShareLevel: optimus-ide-collabsdk.WorkspaceAgentPortShareLevelOrganization,
		Protocol:   optimus-ide-collabsdk.WorkspaceAgentPortShareProtocolHTTP,
	})
	require.NoError(t, err)
	require.EqualValues(t, optimus-ide-collabsdk.WorkspaceAgentPortShareLevelOrganization, ps.ShareLevel)
}
