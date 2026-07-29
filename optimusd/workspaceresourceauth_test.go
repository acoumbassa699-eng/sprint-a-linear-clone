package optimus-ide-collabd_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbauthz"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbgen"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk/agentsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/provisioner/echo"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/provisionersdk/proto"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

func TestPostWorkspaceAuthAzureInstanceIdentity(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		instanceID := newTestInstanceID(t)
		certificates, metadataClient := optimus-ide-collabdtest.NewAzureInstanceIdentity(t, instanceID)
		client, _ := setupInstanceIDWorkspace(t, &optimus-ide-collabdtest.Options{
			AzureCertificates: certificates,
		}, workspaceAgentsForInstanceID(instanceID, "dev"))

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		agentClient := agentsdk.New(client.URL, agentsdk.WithAzureInstanceIdentity())
		agentClient.SDK.HTTPClient = metadataClient

		err := agentClient.RefreshToken(ctx)
		require.NoError(t, err)
	})

	t.Run("Ambiguous/AzureWithSelector", func(t *testing.T) {
		t.Parallel()

		instanceID := newTestInstanceID(t)
		certificates, metadataClient := optimus-ide-collabdtest.NewAzureInstanceIdentity(t, instanceID)
		client, store := setupInstanceIDWorkspace(t, &optimus-ide-collabdtest.Options{
			AzureCertificates: certificates,
		}, workspaceAgentsForInstanceID(instanceID, "alpha", "beta"))

		expectedAgent := requireWorkspaceAgentByInstanceIDAndName(t, store, instanceID, "alpha")
		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		agentClient := agentsdk.New(client.URL, agentsdk.WithAzureInstanceIdentity(
			agentsdk.WithInstanceIdentityAgentName("alpha"),
		))
		agentClient.SDK.HTTPClient = metadataClient

		err := agentClient.RefreshToken(ctx)
		require.NoError(t, err)
		require.Equal(t, expectedAgent.AuthToken.String(), agentClient.SDK.SessionToken())
	})
}

func TestPostWorkspaceAuthAWSInstanceIdentity(t *testing.T) {
	t.Parallel()

	t.Run("Ambiguous/SingleAgent", func(t *testing.T) {
		t.Parallel()

		instanceID := newTestInstanceID(t)
		certificates, metadataClient := optimus-ide-collabdtest.NewAWSInstanceIdentity(t, instanceID)
		client, _ := setupInstanceIDWorkspace(t, &optimus-ide-collabdtest.Options{
			AWSCertificates: certificates,
		}, workspaceAgentsForInstanceID(instanceID, "dev"))

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		agentClient := agentsdk.New(client.URL, agentsdk.WithAWSInstanceIdentity())
		agentClient.SDK.HTTPClient = metadataClient

		err := agentClient.RefreshToken(ctx)
		require.NoError(t, err)
	})

	t.Run("RecycledInstanceID", func(t *testing.T) {
		t.Parallel()

		instanceID := newTestInstanceID(t)
		certificates, metadataClient := optimus-ide-collabdtest.NewAWSInstanceIdentity(t, instanceID)
		setup := setupInstanceIDWorkspaceWithResources(t, &optimus-ide-collabdtest.Options{
			AWSCertificates: certificates,
		}, workspaceAgentsForInstanceID(instanceID, "dev"))

		successorVersion := optimus-ide-collabdtest.CreateTemplateVersion(t, setup.client, setup.user.OrganizationID, &echo.Responses{
			Parse: echo.ParseComplete,
			ProvisionGraph: []*proto.Response{{
				Type: &proto.Response_Graph{
					Graph: &proto.GraphComplete{
						Resources: []*proto.Resource{{
							Name:   "resource",
							Type:   "instance",
							Agents: workspaceAgentsForInstanceID(newTestInstanceID(t), "dev"),
						}},
					},
				},
			}},
		}, func(req *optimus-ide-collabsdk.CreateTemplateVersionRequest) {
			req.TemplateID = setup.template.ID
		})
		optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, setup.client, successorVersion.ID)
		build := optimus-ide-collabdtest.CreateWorkspaceBuild(t, setup.client, setup.workspace, database.WorkspaceTransitionStart, func(req *optimus-ide-collabsdk.CreateWorkspaceBuildRequest) {
			req.TemplateVersionID = successorVersion.ID
		})
		optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, setup.client, build.ID)

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		agentClient := agentsdk.New(setup.client.URL, agentsdk.WithAWSInstanceIdentity())
		agentClient.SDK.HTTPClient = metadataClient

		err := agentClient.RefreshToken(ctx)
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		// The prior build's agent is soft-deleted when the successor
		// build completes (SoftDeletePriorWorkspaceAgents), so the
		// auth query finds no candidates at all and returns 404.
		require.Equal(t, http.StatusNotFound, apiErr.StatusCode())
	})

	t.Run("Ambiguous/MultipleAgentsNoSelector", func(t *testing.T) {
		t.Parallel()

		instanceID := newTestInstanceID(t)
		certificates, metadataClient := optimus-ide-collabdtest.NewAWSInstanceIdentity(t, instanceID)
		client, _ := setupInstanceIDWorkspace(t, &optimus-ide-collabdtest.Options{
			AWSCertificates: certificates,
		}, workspaceAgentsForInstanceID(instanceID, "alpha", "beta"))

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		agentClient := agentsdk.New(client.URL, agentsdk.WithAWSInstanceIdentity())
		agentClient.SDK.HTTPClient = metadataClient

		err := agentClient.RefreshToken(ctx)
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusConflict, apiErr.StatusCode())
		require.Contains(t, apiErr.Message, "OPTIMUS-IDE-COLLAB_AGENT_NAME")
		require.Contains(t, apiErr.Message, "alpha, beta")
	})

	t.Run("Ambiguous/EmptyAgentNameTreatedAsUnset", func(t *testing.T) {
		t.Parallel()

		instanceID := newTestInstanceID(t)
		certificates, metadataClient := optimus-ide-collabdtest.NewAWSInstanceIdentity(t, instanceID)
		client, _ := setupInstanceIDWorkspace(t, &optimus-ide-collabdtest.Options{
			AWSCertificates: certificates,
		}, workspaceAgentsForInstanceID(instanceID, "alpha", "beta"))

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		res := postAWSInstanceIdentity(ctx, t, client, metadataClient, "")
		defer res.Body.Close()

		require.Equal(t, http.StatusConflict, res.StatusCode)
		err := optimus-ide-collabsdk.ReadBodyAsError(res)
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusConflict, apiErr.StatusCode())
		require.Contains(t, apiErr.Message, "OPTIMUS-IDE-COLLAB_AGENT_NAME")
		require.Contains(t, apiErr.Message, "alpha, beta")
	})

	t.Run("Ambiguous/WhitespaceAgentNameTreatedAsUnset", func(t *testing.T) {
		t.Parallel()

		instanceID := newTestInstanceID(t)
		certificates, metadataClient := optimus-ide-collabdtest.NewAWSInstanceIdentity(t, instanceID)
		client, _ := setupInstanceIDWorkspace(t, &optimus-ide-collabdtest.Options{
			AWSCertificates: certificates,
		}, workspaceAgentsForInstanceID(instanceID, "alpha", "beta"))

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		res := postAWSInstanceIdentity(ctx, t, client, metadataClient, "   ")
		defer res.Body.Close()

		require.Equal(t, http.StatusConflict, res.StatusCode)
		err := optimus-ide-collabsdk.ReadBodyAsError(res)
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusConflict, apiErr.StatusCode())
		require.Contains(t, apiErr.Message, "OPTIMUS-IDE-COLLAB_AGENT_NAME")
		require.Contains(t, apiErr.Message, "alpha, beta")
	})

	t.Run("Ambiguous/MultipleAgentsWithSelector", func(t *testing.T) {
		t.Parallel()

		instanceID := newTestInstanceID(t)
		certificates, metadataClient := optimus-ide-collabdtest.NewAWSInstanceIdentity(t, instanceID)
		client, store := setupInstanceIDWorkspace(t, &optimus-ide-collabdtest.Options{
			AWSCertificates: certificates,
		}, workspaceAgentsForInstanceID(instanceID, "alpha", "beta"))

		expectedAgent := requireWorkspaceAgentByInstanceIDAndName(t, store, instanceID, "alpha")
		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		agentClient := agentsdk.New(client.URL, agentsdk.WithAWSInstanceIdentity(
			agentsdk.WithInstanceIdentityAgentName("alpha"),
		))
		agentClient.SDK.HTTPClient = metadataClient

		err := agentClient.RefreshToken(ctx)
		require.NoError(t, err)
		require.Equal(t, expectedAgent.AuthToken.String(), agentClient.SDK.SessionToken())
	})

	t.Run("Ambiguous/MultipleAgentsUnknownSelector", func(t *testing.T) {
		t.Parallel()

		instanceID := newTestInstanceID(t)
		certificates, metadataClient := optimus-ide-collabdtest.NewAWSInstanceIdentity(t, instanceID)
		client, _ := setupInstanceIDWorkspace(t, &optimus-ide-collabdtest.Options{
			AWSCertificates: certificates,
		}, workspaceAgentsForInstanceID(instanceID, "alpha", "beta"))

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		agentClient := agentsdk.New(client.URL, agentsdk.WithAWSInstanceIdentity(
			agentsdk.WithInstanceIdentityAgentName("nonexistent"),
		))
		agentClient.SDK.HTTPClient = metadataClient

		err := agentClient.RefreshToken(ctx)
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusNotFound, apiErr.StatusCode())
	})

	t.Run("Ambiguous/SubAgentExcluded", func(t *testing.T) {
		t.Parallel()

		instanceID := newTestInstanceID(t)
		certificates, metadataClient := optimus-ide-collabdtest.NewAWSInstanceIdentity(t, instanceID)
		client, store := setupInstanceIDWorkspace(t, &optimus-ide-collabdtest.Options{
			AWSCertificates: certificates,
		}, workspaceAgentsForInstanceID(instanceID, "dev"))

		rootAgent := requireWorkspaceAgentByInstanceIDAndName(t, store, instanceID, "dev")
		_ = dbgen.WorkspaceSubAgent(t, store, rootAgent, database.WorkspaceAgent{
			Name: "sub",
			AuthInstanceID: sql.NullString{
				String: instanceID,
				Valid:  true,
			},
		})

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		agentClient := agentsdk.New(client.URL, agentsdk.WithAWSInstanceIdentity())
		agentClient.SDK.HTTPClient = metadataClient

		err := agentClient.RefreshToken(ctx)
		require.NoError(t, err)
		require.Equal(t, rootAgent.AuthToken.String(), agentClient.SDK.SessionToken())
	})
}

func TestPostWorkspaceAuthGoogleInstanceIdentity(t *testing.T) {
	t.Parallel()

	t.Run("Expired", func(t *testing.T) {
		t.Parallel()

		instanceID := newTestInstanceID(t)
		validator, metadata := optimus-ide-collabdtest.NewGoogleInstanceIdentity(t, instanceID, true)
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
			GoogleTokenValidator: validator,
		})

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		agentClient := agentsdk.New(client.URL, agentsdk.WithGoogleInstanceIdentity("", metadata))
		err := agentClient.RefreshToken(ctx)
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusUnauthorized, apiErr.StatusCode())
	})

	t.Run("InstanceNotFound", func(t *testing.T) {
		t.Parallel()

		instanceID := newTestInstanceID(t)
		validator, metadata := optimus-ide-collabdtest.NewGoogleInstanceIdentity(t, instanceID, false)
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
			GoogleTokenValidator: validator,
		})

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		agentClient := agentsdk.New(client.URL, agentsdk.WithGoogleInstanceIdentity("", metadata))
		err := agentClient.RefreshToken(ctx)
		var apiErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, http.StatusNotFound, apiErr.StatusCode())
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		instanceID := newTestInstanceID(t)
		validator, metadata := optimus-ide-collabdtest.NewGoogleInstanceIdentity(t, instanceID, false)
		client, _ := setupInstanceIDWorkspace(t, &optimus-ide-collabdtest.Options{
			GoogleTokenValidator: validator,
		}, workspaceAgentsForInstanceID(instanceID, "dev"))

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		agentClient := agentsdk.New(client.URL, agentsdk.WithGoogleInstanceIdentity("", metadata))
		err := agentClient.RefreshToken(ctx)
		require.NoError(t, err)
	})

	t.Run("Ambiguous/GoogleWithSelector", func(t *testing.T) {
		t.Parallel()

		instanceID := newTestInstanceID(t)
		validator, metadata := optimus-ide-collabdtest.NewGoogleInstanceIdentity(t, instanceID, false)
		client, store := setupInstanceIDWorkspace(t, &optimus-ide-collabdtest.Options{
			GoogleTokenValidator: validator,
		}, workspaceAgentsForInstanceID(instanceID, "alpha", "beta"))

		expectedAgent := requireWorkspaceAgentByInstanceIDAndName(t, store, instanceID, "alpha")
		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		agentClient := agentsdk.New(client.URL, agentsdk.WithGoogleInstanceIdentity(
			"",
			metadata,
			agentsdk.WithInstanceIdentityAgentName("alpha"),
		))
		err := agentClient.RefreshToken(ctx)
		require.NoError(t, err)
		require.Equal(t, expectedAgent.AuthToken.String(), agentClient.SDK.SessionToken())
	})
}

type instanceIDWorkspaceSetup struct {
	client    *optimus-ide-collabsdk.Client
	store     database.Store
	user      optimus-ide-collabsdk.CreateFirstUserResponse
	template  optimus-ide-collabsdk.Template
	workspace optimus-ide-collabsdk.Workspace
}

func setupInstanceIDWorkspace(t *testing.T, opts *optimus-ide-collabdtest.Options, agents []*proto.Agent) (*optimus-ide-collabsdk.Client, database.Store) {
	t.Helper()

	setup := setupInstanceIDWorkspaceWithResources(t, opts, agents)
	return setup.client, setup.store
}

func setupInstanceIDWorkspaceWithResources(
	t *testing.T,
	opts *optimus-ide-collabdtest.Options,
	agents []*proto.Agent,
) instanceIDWorkspaceSetup {
	t.Helper()

	actualOpts := &optimus-ide-collabdtest.Options{}
	if opts != nil {
		*actualOpts = *opts
	}
	actualOpts.IncludeProvisionerDaemon = true

	client, store := optimus-ide-collabdtest.NewWithDatabase(t, actualOpts)
	user := optimus-ide-collabdtest.CreateFirstUser(t, client)
	version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, user.OrganizationID, &echo.Responses{
		Parse: echo.ParseComplete,
		ProvisionGraph: []*proto.Response{{
			Type: &proto.Response_Graph{
				Graph: &proto.GraphComplete{
					Resources: []*proto.Resource{{
						Name:   "resource",
						Type:   "instance",
						Agents: agents,
					}},
				},
			},
		}},
	})
	template := optimus-ide-collabdtest.CreateTemplate(t, client, user.OrganizationID, version.ID)
	optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
	workspace := optimus-ide-collabdtest.CreateWorkspace(t, client, template.ID)
	optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, client, workspace.LatestBuild.ID)

	return instanceIDWorkspaceSetup{
		client:    client,
		store:     store,
		user:      user,
		template:  template,
		workspace: workspace,
	}
}

func workspaceAgentsForInstanceID(instanceID string, names ...string) []*proto.Agent {
	agents := make([]*proto.Agent, 0, len(names))
	for _, name := range names {
		agents = append(agents, &proto.Agent{
			Name: name,
			Auth: &proto.Agent_InstanceId{InstanceId: instanceID},
		})
	}
	return agents
}

func requireWorkspaceAgentByInstanceIDAndName(t testing.TB, store database.Store, instanceID string, name string) database.WorkspaceAgent {
	t.Helper()

	ctx := dbauthz.AsSystemRestricted(testutil.Context(t, testutil.WaitLong))
	agents, err := store.GetWorkspaceAgentsByInstanceID(ctx, instanceID)
	require.NoError(t, err)
	for _, agent := range agents {
		if agent.Name == name {
			return agent
		}
	}
	require.FailNow(t, "workspace agent not found", "instance ID %q, name %q", instanceID, name)
	return database.WorkspaceAgent{}
}

const awsInstanceIdentityMetadataURL = "http://169.254.169.254/latest/dynamic/instance-identity"

func postAWSInstanceIdentity(
	ctx context.Context,
	t testing.TB,
	client *optimus-ide-collabsdk.Client,
	metadataClient *http.Client,
	agentName string,
) *http.Response {
	t.Helper()

	signature := readAWSInstanceMetadata(ctx, t, metadataClient, "signature")
	document := readAWSInstanceMetadata(ctx, t, metadataClient, "document")
	reqBody, err := json.Marshal(map[string]string{
		"signature":  signature,
		"document":   document,
		"agent_name": agentName,
	})
	require.NoError(t, err)

	res, err := client.RequestWithoutSessionToken(
		ctx,
		http.MethodPost,
		"/api/v2/workspaceagents/aws-instance-identity",
		reqBody,
	)
	require.NoError(t, err)
	return res
}

func readAWSInstanceMetadata(
	ctx context.Context,
	t testing.TB,
	metadataClient *http.Client,
	path string,
) string {
	t.Helper()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		awsInstanceIdentityMetadataURL+"/"+path,
		nil,
	)
	require.NoError(t, err)
	res, err := metadataClient.Do(req)
	require.NoError(t, err)
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	return string(body)
}

func newTestInstanceID(t testing.TB) string {
	t.Helper()
	return fmt.Sprintf("instance-%d", time.Now().UnixNano())
}
