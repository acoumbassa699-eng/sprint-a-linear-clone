package cli_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli/clitest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/optimus-ide-collabdenttest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/license"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/provisioner/echo"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/provisionersdk/proto"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil/expecter"
)

// completeWithExternalAgent creates a template version with an external agent resource
func completeWithExternalAgent() *echo.Responses {
	return &echo.Responses{
		Parse: echo.ParseComplete,
		ProvisionGraph: []*proto.Response{
			{
				Type: &proto.Response_Graph{
					Graph: &proto.GraphComplete{
						Resources: []*proto.Resource{
							{
								Type: "optimus-ide-collab_external_agent",
								Name: "main",
								Agents: []*proto.Agent{
									{
										Name:            "external-agent",
										OperatingSystem: "linux",
										Architecture:    "amd64",
									},
								},
							},
						},
						HasExternalAgents: true,
					},
				},
			},
		},
	}
}

// completeWithRegularAgent creates a template version with a regular agent (no external agent)
func completeWithRegularAgent() *echo.Responses {
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
										Name:            "regular-agent",
										OperatingSystem: "linux",
										Architecture:    "amd64",
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

func TestExternalWorkspaces(t *testing.T) {
	t.Parallel()

	t.Run("Create", func(t *testing.T) {
		t.Parallel()
		logger := testutil.Logger(t)
		client, owner := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			Options: &optimus-ide-collabdtest.Options{
				IncludeProvisionerDaemon: true,
			},
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureWorkspaceExternalAgent: 1,
				},
			},
		})
		member, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, owner.OrganizationID, completeWithExternalAgent())
		optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
		template := optimus-ide-collabdtest.CreateTemplate(t, client, owner.OrganizationID, version.ID)

		args := []string{
			"external-workspaces",
			"create",
			"my-external-workspace",
			"--template", template.Name,
		}
		inv, root := newCLI(t, args...)
		clitest.SetupConfig(t, member, root)
		doneChan := make(chan struct{})
		stdout := expecter.NewAttachedToInvocation(t, inv)
		stdin := testutil.NewWriterAttachedToInvocation(t, logger.Named("stdin"), inv)
		ctx := testutil.Context(t, testutil.WaitLong)
		go func() {
			defer close(doneChan)
			err := inv.Run()
			assert.NoError(t, err)
		}()

		// Expect the workspace creation confirmation
		stdout.ExpectMatch(ctx, "optimus-ide-collab_external_agent.main")
		stdout.ExpectMatch(ctx, "external-agent (linux, amd64)")
		stdout.ExpectMatch(ctx, "Confirm create")
		stdin.WriteLine("yes")

		// Expect the external agent instructions
		stdout.ExpectMatch(ctx, "Please run the following command to attach external agent")
		stdout.ExpectRegexMatch(ctx, "curl -fsSL .* | OPTIMUS-IDE-COLLAB_AGENT_TOKEN=.* sh")

		testutil.TryReceive(ctx, t, doneChan)

		// Verify the workspace was created
		ws, err := member.WorkspaceByOwnerAndName(context.Background(), optimus-ide-collabsdk.Me, "my-external-workspace", optimus-ide-collabsdk.WorkspaceOptions{})
		require.NoError(t, err)
		assert.Equal(t, template.Name, ws.TemplateName)
	})

	t.Run("CreateWithoutTemplate", func(t *testing.T) {
		t.Parallel()
		client, owner := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			Options: &optimus-ide-collabdtest.Options{
				IncludeProvisionerDaemon: true,
			},
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureWorkspaceExternalAgent: 1,
				},
			},
		})
		member, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID)

		args := []string{
			"external-workspaces",
			"create",
			"my-external-workspace",
		}
		inv, root := newCLI(t, args...)
		clitest.SetupConfig(t, member, root)

		err := inv.Run()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Missing values for the required flags: template")
	})

	t.Run("CreateWithRegularTemplate", func(t *testing.T) {
		t.Parallel()
		client, owner := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			Options: &optimus-ide-collabdtest.Options{
				IncludeProvisionerDaemon: true,
			},
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureWorkspaceExternalAgent: 1,
				},
			},
		})
		member, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, owner.OrganizationID, completeWithRegularAgent())
		optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
		template := optimus-ide-collabdtest.CreateTemplate(t, client, owner.OrganizationID, version.ID)

		args := []string{
			"external-workspaces",
			"create",
			"my-external-workspace",
			"--template", template.Name,
		}
		inv, root := newCLI(t, args...)
		clitest.SetupConfig(t, member, root)

		err := inv.Run()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "does not have an external agent")
	})

	t.Run("List", func(t *testing.T) {
		t.Parallel()
		client, owner := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			Options: &optimus-ide-collabdtest.Options{
				IncludeProvisionerDaemon: true,
			},
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureWorkspaceExternalAgent: 1,
				},
			},
		})
		member, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, owner.OrganizationID, completeWithExternalAgent())
		optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
		template := optimus-ide-collabdtest.CreateTemplate(t, client, owner.OrganizationID, version.ID)

		// Create an external workspace
		ws := optimus-ide-collabdtest.CreateWorkspace(t, member, template.ID)
		optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, client, ws.LatestBuild.ID)

		args := []string{
			"external-workspaces",
			"list",
		}
		inv, root := newCLI(t, args...)
		clitest.SetupConfig(t, member, root)
		stdout := expecter.NewAttachedToInvocation(t, inv)

		ctx, cancelFunc := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancelFunc()
		done := make(chan any)
		go func() {
			errC := inv.WithContext(ctx).Run()
			assert.NoError(t, errC)
			close(done)
		}()
		stdout.ExpectMatch(ctx, ws.Name)
		stdout.ExpectMatch(ctx, template.Name)
		cancelFunc()
		<-done
	})

	t.Run("ListJSON", func(t *testing.T) {
		t.Parallel()
		client, owner := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			Options: &optimus-ide-collabdtest.Options{
				IncludeProvisionerDaemon: true,
			},
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureWorkspaceExternalAgent: 1,
				},
			},
		})
		member, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, owner.OrganizationID, completeWithExternalAgent())
		optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
		template := optimus-ide-collabdtest.CreateTemplate(t, client, owner.OrganizationID, version.ID)

		// Create an external workspace
		ws := optimus-ide-collabdtest.CreateWorkspace(t, member, template.ID)
		optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, client, ws.LatestBuild.ID)

		args := []string{
			"external-workspaces",
			"list",
			"--output=json",
		}
		inv, root := newCLI(t, args...)
		clitest.SetupConfig(t, member, root)

		ctx, cancelFunc := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancelFunc()

		out := bytes.NewBuffer(nil)
		inv.Stdout = out
		err := inv.WithContext(ctx).Run()
		require.NoError(t, err)

		var workspaces []optimus-ide-collabsdk.Workspace
		require.NoError(t, json.Unmarshal(out.Bytes(), &workspaces))
		require.Len(t, workspaces, 1)
		assert.Equal(t, ws.Name, workspaces[0].Name)
	})

	t.Run("ListNoWorkspaces", func(t *testing.T) {
		t.Parallel()
		client, owner := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			Options: &optimus-ide-collabdtest.Options{
				IncludeProvisionerDaemon: true,
			},
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureWorkspaceExternalAgent: 1,
				},
			},
		})
		member, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID)

		args := []string{
			"external-workspaces",
			"list",
		}
		inv, root := newCLI(t, args...)
		clitest.SetupConfig(t, member, root)
		stdout := expecter.NewAttachedToInvocation(t, inv)

		ctx, cancelFunc := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancelFunc()
		done := make(chan any)
		go func() {
			errC := inv.WithContext(ctx).Run()
			assert.NoError(t, errC)
			close(done)
		}()
		stdout.ExpectMatch(ctx, "No workspaces found!")
		stdout.ExpectMatch(ctx, "optimus-ide-collab external-workspaces create")
		cancelFunc()
		<-done
	})

	t.Run("AgentInstructions", func(t *testing.T) {
		t.Parallel()
		client, owner := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			Options: &optimus-ide-collabdtest.Options{
				IncludeProvisionerDaemon: true,
			},
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureWorkspaceExternalAgent: 1,
				},
			},
		})
		member, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, owner.OrganizationID, completeWithExternalAgent())
		optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
		template := optimus-ide-collabdtest.CreateTemplate(t, client, owner.OrganizationID, version.ID)

		// Create an external workspace
		ws := optimus-ide-collabdtest.CreateWorkspace(t, member, template.ID)
		optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, client, ws.LatestBuild.ID)

		args := []string{
			"external-workspaces",
			"agent-instructions",
			ws.Name,
		}
		inv, root := newCLI(t, args...)
		clitest.SetupConfig(t, member, root)
		stdout := expecter.NewAttachedToInvocation(t, inv)

		ctx, cancelFunc := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancelFunc()
		done := make(chan any)
		go func() {
			errC := inv.WithContext(ctx).Run()
			assert.NoError(t, errC)
			close(done)
		}()
		stdout.ExpectMatch(ctx, "Please run the following command to attach external agent to the workspace")
		stdout.ExpectRegexMatch(ctx, "curl -fsSL .* | OPTIMUS-IDE-COLLAB_AGENT_TOKEN=.* sh")
		cancelFunc()

		ctx = testutil.Context(t, testutil.WaitLong)
		testutil.TryReceive(ctx, t, done)
	})

	t.Run("AgentInstructionsJSON", func(t *testing.T) {
		t.Parallel()
		client, owner := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			Options: &optimus-ide-collabdtest.Options{
				IncludeProvisionerDaemon: true,
			},
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureWorkspaceExternalAgent: 1,
				},
			},
		})
		member, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, owner.OrganizationID, completeWithExternalAgent())
		optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
		template := optimus-ide-collabdtest.CreateTemplate(t, client, owner.OrganizationID, version.ID)

		// Create an external workspace
		ws := optimus-ide-collabdtest.CreateWorkspace(t, member, template.ID)
		optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, client, ws.LatestBuild.ID)

		args := []string{
			"external-workspaces",
			"agent-instructions",
			ws.Name,
			"--output=json",
		}
		inv, root := newCLI(t, args...)
		clitest.SetupConfig(t, member, root)

		ctx, cancelFunc := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancelFunc()

		out := bytes.NewBuffer(nil)
		inv.Stdout = out
		err := inv.WithContext(ctx).Run()
		require.NoError(t, err)

		var agentInfo map[string]interface{}
		require.NoError(t, json.Unmarshal(out.Bytes(), &agentInfo))
		assert.Equal(t, "token", agentInfo["auth_type"])
		assert.NotEmpty(t, agentInfo["auth_token"])
		assert.NotEmpty(t, agentInfo["init_script"])
	})

	t.Run("AgentInstructionsNonExistentWorkspace", func(t *testing.T) {
		t.Parallel()
		client, owner := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			Options: &optimus-ide-collabdtest.Options{
				IncludeProvisionerDaemon: true,
			},
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureWorkspaceExternalAgent: 1,
				},
			},
		})
		member, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID)

		args := []string{
			"external-workspaces",
			"agent-instructions",
			"non-existent-workspace",
		}
		inv, root := newCLI(t, args...)
		clitest.SetupConfig(t, member, root)

		err := inv.Run()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Resource not found")
	})

	t.Run("AgentInstructionsNonExistentAgent", func(t *testing.T) {
		t.Parallel()
		client, owner := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			Options: &optimus-ide-collabdtest.Options{
				IncludeProvisionerDaemon: true,
			},
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureWorkspaceExternalAgent: 1,
				},
			},
		})
		member, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, owner.OrganizationID, completeWithExternalAgent())
		optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
		template := optimus-ide-collabdtest.CreateTemplate(t, client, owner.OrganizationID, version.ID)

		// Create an external workspace
		ws := optimus-ide-collabdtest.CreateWorkspace(t, member, template.ID)
		optimus-ide-collabdtest.AwaitWorkspaceBuildJobCompleted(t, client, ws.LatestBuild.ID)

		args := []string{
			"external-workspaces",
			"agent-instructions",
			ws.Name + ".non-existent-agent",
		}
		inv, root := newCLI(t, args...)
		clitest.SetupConfig(t, member, root)

		err := inv.Run()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "agent not found by name")
	})

	t.Run("CreateWithTemplateVersion", func(t *testing.T) {
		t.Parallel()
		client, owner := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			Options: &optimus-ide-collabdtest.Options{
				IncludeProvisionerDaemon: true,
			},
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureWorkspaceExternalAgent: 1,
				},
			},
		})
		member, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, owner.OrganizationID)
		version := optimus-ide-collabdtest.CreateTemplateVersion(t, client, owner.OrganizationID, completeWithExternalAgent())
		optimus-ide-collabdtest.AwaitTemplateVersionJobCompleted(t, client, version.ID)
		template := optimus-ide-collabdtest.CreateTemplate(t, client, owner.OrganizationID, version.ID)

		args := []string{
			"external-workspaces",
			"create",
			"my-external-workspace",
			"--template", template.Name,
			"--template-version", version.Name,
			"-y",
		}
		inv, root := newCLI(t, args...)
		clitest.SetupConfig(t, member, root)
		doneChan := make(chan struct{})
		stdout := expecter.NewAttachedToInvocation(t, inv)
		ctx := testutil.Context(t, testutil.WaitLong)
		go func() {
			defer close(doneChan)
			err := inv.Run()
			assert.NoError(t, err)
		}()

		// Expect the workspace creation confirmation
		stdout.ExpectMatch(ctx, "optimus-ide-collab_external_agent.main")
		stdout.ExpectMatch(ctx, "external-agent (linux, amd64)")

		// Expect the external agent instructions
		stdout.ExpectMatch(ctx, "Please run the following command to attach external agent")
		stdout.ExpectRegexMatch(ctx, "curl -fsSL .* | OPTIMUS-IDE-COLLAB_AGENT_TOKEN=.* sh")

		testutil.TryReceive(ctx, t, doneChan)

		// Verify the workspace was created
		ws, err := member.WorkspaceByOwnerAndName(context.Background(), optimus-ide-collabsdk.Me, "my-external-workspace", optimus-ide-collabsdk.WorkspaceOptions{})
		require.NoError(t, err)
		assert.Equal(t, template.Name, ws.TemplateName)
	})
}
