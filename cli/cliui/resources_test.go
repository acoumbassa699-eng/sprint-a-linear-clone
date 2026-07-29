package cliui_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli/cliui"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbtime"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/pty/ptytest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

func TestWorkspaceResources(t *testing.T) {
	t.Parallel()
	t.Run("SingleAgentSSH", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.Context(t, testutil.WaitMedium)
		ptty := ptytest.New(t)
		done := make(chan struct{})
		go func() {
			err := cliui.WorkspaceResources(ptty.Output(), []optimus-ide-collabsdk.WorkspaceResource{{
				Type:       "google_compute_instance",
				Name:       "dev",
				Transition: optimus-ide-collabsdk.WorkspaceTransitionStart,
				Agents: []optimus-ide-collabsdk.WorkspaceAgent{{
					Name:            "dev",
					Status:          optimus-ide-collabsdk.WorkspaceAgentConnected,
					LifecycleState:  optimus-ide-collabsdk.WorkspaceAgentLifecycleCreated,
					Architecture:    "amd64",
					OperatingSystem: "linux",
					Health:          optimus-ide-collabsdk.WorkspaceAgentHealth{Healthy: true},
				}},
			}}, cliui.WorkspaceResourcesOptions{
				WorkspaceName: "example",
			})
			assert.NoError(t, err)
			close(done)
		}()
		ptty.ExpectMatch(ctx, "optimus-ide-collab ssh example")
		<-done
	})

	t.Run("MultipleStates", func(t *testing.T) {
		t.Parallel()
		ctx := testutil.Context(t, testutil.WaitMedium)
		ptty := ptytest.New(t)
		disconnected := dbtime.Now().Add(-4 * time.Second)
		done := make(chan struct{})
		go func() {
			err := cliui.WorkspaceResources(ptty.Output(), []optimus-ide-collabsdk.WorkspaceResource{{
				Transition: optimus-ide-collabsdk.WorkspaceTransitionStart,
				Type:       "google_compute_disk",
				Name:       "root",
			}, {
				Transition: optimus-ide-collabsdk.WorkspaceTransitionStop,
				Type:       "google_compute_disk",
				Name:       "root",
			}, {
				Transition: optimus-ide-collabsdk.WorkspaceTransitionStart,
				Type:       "google_compute_instance",
				Name:       "dev",
				Agents: []optimus-ide-collabsdk.WorkspaceAgent{{
					CreatedAt:       dbtime.Now().Add(-10 * time.Second),
					Status:          optimus-ide-collabsdk.WorkspaceAgentConnecting,
					LifecycleState:  optimus-ide-collabsdk.WorkspaceAgentLifecycleCreated,
					Name:            "dev",
					OperatingSystem: "linux",
					Architecture:    "amd64",
					Health:          optimus-ide-collabsdk.WorkspaceAgentHealth{Healthy: true},
				}},
			}, {
				Transition: optimus-ide-collabsdk.WorkspaceTransitionStart,
				Type:       "kubernetes_pod",
				Name:       "dev",
				Agents: []optimus-ide-collabsdk.WorkspaceAgent{{
					Status:          optimus-ide-collabsdk.WorkspaceAgentConnected,
					LifecycleState:  optimus-ide-collabsdk.WorkspaceAgentLifecycleReady,
					Name:            "go",
					Architecture:    "amd64",
					OperatingSystem: "linux",
					Health:          optimus-ide-collabsdk.WorkspaceAgentHealth{Healthy: true},
				}, {
					DisconnectedAt:  &disconnected,
					Status:          optimus-ide-collabsdk.WorkspaceAgentDisconnected,
					LifecycleState:  optimus-ide-collabsdk.WorkspaceAgentLifecycleReady,
					Name:            "postgres",
					Architecture:    "amd64",
					OperatingSystem: "linux",
					Health: optimus-ide-collabsdk.WorkspaceAgentHealth{
						Healthy: false,
						Reason:  "agent has lost connection",
					},
				}},
			}}, cliui.WorkspaceResourcesOptions{
				WorkspaceName:  "dev",
				HideAgentState: false,
				HideAccess:     false,
			})
			assert.NoError(t, err)
			close(done)
		}()
		ptty.ExpectMatch(ctx, "google_compute_disk.root")
		ptty.ExpectMatch(ctx, "google_compute_instance.dev")
		ptty.ExpectMatch(ctx, "healthy")
		ptty.ExpectMatch(ctx, "optimus-ide-collab ssh dev.dev")
		ptty.ExpectMatch(ctx, "kubernetes_pod.dev")
		ptty.ExpectMatch(ctx, "healthy")
		ptty.ExpectMatch(ctx, "optimus-ide-collab ssh dev.go")
		ptty.ExpectMatch(ctx, "agent has lost connection")
		ptty.ExpectMatch(ctx, "optimus-ide-collab ssh dev.postgres")
		<-done
	})
}
