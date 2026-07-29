package optimus-ide-collabd_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbauthz"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbfake"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/provisionersdk/proto"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

func TestPostWorkspaceAgentPortShare(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
	defer cancel()
	ownerClient, db := optimus-ide-collabdtest.NewWithDatabase(t, nil)
	owner := optimus-ide-collabdtest.CreateFirstUser(t, ownerClient)
	client, user := optimus-ide-collabdtest.CreateAnotherUser(t, ownerClient, owner.OrganizationID)

	tmpDir := t.TempDir()
	r := dbfake.WorkspaceBuild(t, db, database.WorkspaceTable{
		OrganizationID: owner.OrganizationID,
		OwnerID:        user.ID,
	}).WithAgent(func(agents []*proto.Agent) []*proto.Agent {
		agents[0].Directory = tmpDir
		return agents
	}).Do()
	agents, err := db.GetWorkspaceAgentsInLatestBuildByWorkspaceID(dbauthz.As(ctx, optimus-ide-collabdtest.AuthzUserSubjectWithDB(ctx, t, db, user)), r.Workspace.ID)
	require.NoError(t, err)

	// owner level should fail
	_, err = client.UpsertWorkspaceAgentPortShare(ctx, r.Workspace.ID, optimus-ide-collabsdk.UpsertWorkspaceAgentPortShareRequest{
		AgentName:  agents[0].Name,
		Port:       8080,
		ShareLevel: optimus-ide-collabsdk.WorkspaceAgentPortShareLevel("owner"),
		Protocol:   optimus-ide-collabsdk.WorkspaceAgentPortShareProtocolHTTP,
	})
	require.Error(t, err)

	// invalid level should fail
	_, err = client.UpsertWorkspaceAgentPortShare(ctx, r.Workspace.ID, optimus-ide-collabsdk.UpsertWorkspaceAgentPortShareRequest{
		AgentName:  agents[0].Name,
		Port:       8080,
		ShareLevel: optimus-ide-collabsdk.WorkspaceAgentPortShareLevel("invalid"),
		Protocol:   optimus-ide-collabsdk.WorkspaceAgentPortShareProtocolHTTP,
	})
	require.Error(t, err)

	// invalid protocol should fail
	_, err = client.UpsertWorkspaceAgentPortShare(ctx, r.Workspace.ID, optimus-ide-collabsdk.UpsertWorkspaceAgentPortShareRequest{
		AgentName:  agents[0].Name,
		Port:       8080,
		ShareLevel: optimus-ide-collabsdk.WorkspaceAgentPortShareLevelPublic,
		Protocol:   optimus-ide-collabsdk.WorkspaceAgentPortShareProtocol("invalid"),
	})
	require.Error(t, err)

	// invalid port should fail
	_, err = client.UpsertWorkspaceAgentPortShare(ctx, r.Workspace.ID, optimus-ide-collabsdk.UpsertWorkspaceAgentPortShareRequest{
		AgentName:  agents[0].Name,
		Port:       0,
		ShareLevel: optimus-ide-collabsdk.WorkspaceAgentPortShareLevelPublic,
		Protocol:   optimus-ide-collabsdk.WorkspaceAgentPortShareProtocolHTTP,
	})
	require.Error(t, err)
	_, err = client.UpsertWorkspaceAgentPortShare(ctx, r.Workspace.ID, optimus-ide-collabsdk.UpsertWorkspaceAgentPortShareRequest{
		AgentName:  agents[0].Name,
		Port:       90000000,
		ShareLevel: optimus-ide-collabsdk.WorkspaceAgentPortShareLevelPublic,
	})
	require.Error(t, err)

	// OK, ignoring template max port share level because we are AGPL
	ps, err := client.UpsertWorkspaceAgentPortShare(ctx, r.Workspace.ID, optimus-ide-collabsdk.UpsertWorkspaceAgentPortShareRequest{
		AgentName:  agents[0].Name,
		Port:       8080,
		ShareLevel: optimus-ide-collabsdk.WorkspaceAgentPortShareLevelPublic,
		Protocol:   optimus-ide-collabsdk.WorkspaceAgentPortShareProtocolHTTPS,
	})
	require.NoError(t, err)
	require.EqualValues(t, optimus-ide-collabsdk.WorkspaceAgentPortShareLevelPublic, ps.ShareLevel)
	require.EqualValues(t, optimus-ide-collabsdk.WorkspaceAgentPortShareProtocolHTTPS, ps.Protocol)

	// list
	list, err := client.GetWorkspaceAgentPortShares(ctx, r.Workspace.ID)
	require.NoError(t, err)
	require.Len(t, list.Shares, 1)
	require.EqualValues(t, agents[0].Name, list.Shares[0].AgentName)
	require.EqualValues(t, 8080, list.Shares[0].Port)
	require.EqualValues(t, optimus-ide-collabsdk.WorkspaceAgentPortShareLevelPublic, list.Shares[0].ShareLevel)
	require.EqualValues(t, optimus-ide-collabsdk.WorkspaceAgentPortShareProtocolHTTPS, list.Shares[0].Protocol)

	// update share level and protocol
	ps, err = client.UpsertWorkspaceAgentPortShare(ctx, r.Workspace.ID, optimus-ide-collabsdk.UpsertWorkspaceAgentPortShareRequest{
		AgentName:  agents[0].Name,
		Port:       8080,
		ShareLevel: optimus-ide-collabsdk.WorkspaceAgentPortShareLevelAuthenticated,
		Protocol:   optimus-ide-collabsdk.WorkspaceAgentPortShareProtocolHTTP,
	})
	require.NoError(t, err)
	require.EqualValues(t, optimus-ide-collabsdk.WorkspaceAgentPortShareLevelAuthenticated, ps.ShareLevel)
	require.EqualValues(t, optimus-ide-collabsdk.WorkspaceAgentPortShareProtocolHTTP, ps.Protocol)

	// list
	list, err = client.GetWorkspaceAgentPortShares(ctx, r.Workspace.ID)
	require.NoError(t, err)
	require.Len(t, list.Shares, 1)
	require.EqualValues(t, agents[0].Name, list.Shares[0].AgentName)
	require.EqualValues(t, 8080, list.Shares[0].Port)
	require.EqualValues(t, optimus-ide-collabsdk.WorkspaceAgentPortShareLevelAuthenticated, list.Shares[0].ShareLevel)
	require.EqualValues(t, optimus-ide-collabsdk.WorkspaceAgentPortShareProtocolHTTP, list.Shares[0].Protocol)

	// list 2 ordered by port
	ps, err = client.UpsertWorkspaceAgentPortShare(ctx, r.Workspace.ID, optimus-ide-collabsdk.UpsertWorkspaceAgentPortShareRequest{
		AgentName:  agents[0].Name,
		Port:       8081,
		ShareLevel: optimus-ide-collabsdk.WorkspaceAgentPortShareLevelPublic,
		Protocol:   optimus-ide-collabsdk.WorkspaceAgentPortShareProtocolHTTPS,
	})
	require.NoError(t, err)
	list, err = client.GetWorkspaceAgentPortShares(ctx, r.Workspace.ID)
	require.NoError(t, err)
	require.Len(t, list.Shares, 2)
	require.EqualValues(t, 8080, list.Shares[0].Port)
	require.EqualValues(t, 8081, list.Shares[1].Port)
}

func TestGetWorkspaceAgentPortShares(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
	defer cancel()

	ownerClient, db := optimus-ide-collabdtest.NewWithDatabase(t, nil)
	owner := optimus-ide-collabdtest.CreateFirstUser(t, ownerClient)
	client, user := optimus-ide-collabdtest.CreateAnotherUser(t, ownerClient, owner.OrganizationID)

	tmpDir := t.TempDir()
	r := dbfake.WorkspaceBuild(t, db, database.WorkspaceTable{
		OrganizationID: owner.OrganizationID,
		OwnerID:        user.ID,
	}).WithAgent(func(agents []*proto.Agent) []*proto.Agent {
		agents[0].Directory = tmpDir
		return agents
	}).Do()
	agents, err := db.GetWorkspaceAgentsInLatestBuildByWorkspaceID(dbauthz.As(ctx, optimus-ide-collabdtest.AuthzUserSubjectWithDB(ctx, t, db, user)), r.Workspace.ID)
	require.NoError(t, err)

	_, err = client.UpsertWorkspaceAgentPortShare(ctx, r.Workspace.ID, optimus-ide-collabsdk.UpsertWorkspaceAgentPortShareRequest{
		AgentName:  agents[0].Name,
		Port:       8080,
		ShareLevel: optimus-ide-collabsdk.WorkspaceAgentPortShareLevelPublic,
		Protocol:   optimus-ide-collabsdk.WorkspaceAgentPortShareProtocolHTTP,
	})
	require.NoError(t, err)

	ps, err := client.GetWorkspaceAgentPortShares(ctx, r.Workspace.ID)
	require.NoError(t, err)
	require.Len(t, ps.Shares, 1)
	require.EqualValues(t, agents[0].Name, ps.Shares[0].AgentName)
	require.EqualValues(t, 8080, ps.Shares[0].Port)
	require.EqualValues(t, optimus-ide-collabsdk.WorkspaceAgentPortShareLevelPublic, ps.Shares[0].ShareLevel)
}

func TestDeleteWorkspaceAgentPortShare(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
	defer cancel()

	ownerClient, db := optimus-ide-collabdtest.NewWithDatabase(t, nil)
	owner := optimus-ide-collabdtest.CreateFirstUser(t, ownerClient)
	client, user := optimus-ide-collabdtest.CreateAnotherUser(t, ownerClient, owner.OrganizationID)

	tmpDir := t.TempDir()
	r := dbfake.WorkspaceBuild(t, db, database.WorkspaceTable{
		OrganizationID: owner.OrganizationID,
		OwnerID:        user.ID,
	}).WithAgent(func(agents []*proto.Agent) []*proto.Agent {
		agents[0].Directory = tmpDir
		return agents
	}).Do()
	agents, err := db.GetWorkspaceAgentsInLatestBuildByWorkspaceID(dbauthz.As(ctx, optimus-ide-collabdtest.AuthzUserSubjectWithDB(ctx, t, db, user)), r.Workspace.ID)
	require.NoError(t, err)

	// create
	ps, err := client.UpsertWorkspaceAgentPortShare(ctx, r.Workspace.ID, optimus-ide-collabsdk.UpsertWorkspaceAgentPortShareRequest{
		AgentName:  agents[0].Name,
		Port:       8080,
		ShareLevel: optimus-ide-collabsdk.WorkspaceAgentPortShareLevelPublic,
		Protocol:   optimus-ide-collabsdk.WorkspaceAgentPortShareProtocolHTTP,
	})
	require.NoError(t, err)
	require.EqualValues(t, optimus-ide-collabsdk.WorkspaceAgentPortShareLevelPublic, ps.ShareLevel)

	// delete
	err = client.DeleteWorkspaceAgentPortShare(ctx, r.Workspace.ID, optimus-ide-collabsdk.DeleteWorkspaceAgentPortShareRequest{
		AgentName: agents[0].Name,
		Port:      8080,
	})
	require.NoError(t, err)

	// delete missing
	err = client.DeleteWorkspaceAgentPortShare(ctx, r.Workspace.ID, optimus-ide-collabsdk.DeleteWorkspaceAgentPortShareRequest{
		AgentName: agents[0].Name,
		Port:      8080,
	})
	require.Error(t, err)

	_, err = db.GetWorkspaceAgentPortShare(dbauthz.As(ctx, optimus-ide-collabdtest.AuthzUserSubjectWithDB(ctx, t, db, user)), database.GetWorkspaceAgentPortShareParams{
		WorkspaceID: r.Workspace.ID,
		AgentName:   agents[0].Name,
		Port:        8080,
	})
	require.Error(t, err)
}
