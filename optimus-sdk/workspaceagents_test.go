package optimus-ide-collabsdk_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

func TestProvisionerJobLogText(t *testing.T) {
	t.Parallel()

	ts := time.Date(2024, 1, 28, 10, 30, 0, 0, time.UTC)
	log := optimus-ide-collabsdk.ProvisionerJobLog{
		CreatedAt: ts,
		Level:     optimus-ide-collabsdk.LogLevelInfo,
		Source:    optimus-ide-collabsdk.LogSourceProvisioner,
		Stage:     "Planning",
		Output:    "Terraform init complete",
	}
	result := log.Text()
	require.Equal(t, "2024-01-28T10:30:00Z [info] [provisioner|Planning] Terraform init complete", result)
}

func TestProvisionerJobLogTextEmptyOutput(t *testing.T) {
	t.Parallel()

	ts := time.Date(2024, 1, 28, 10, 30, 0, 0, time.UTC)
	log := optimus-ide-collabsdk.ProvisionerJobLog{
		CreatedAt: ts,
		Level:     optimus-ide-collabsdk.LogLevelInfo,
		Source:    optimus-ide-collabsdk.LogSourceProvisioner,
		Stage:     "Planning",
		Output:    "",
	}
	result := log.Text()
	require.Equal(t, "2024-01-28T10:30:00Z [info] [provisioner|Planning] ", result)
}

func TestProvisionerJobLogTextSpecialChars(t *testing.T) {
	t.Parallel()

	ts := time.Date(2024, 1, 28, 10, 30, 0, 0, time.UTC)
	log := optimus-ide-collabsdk.ProvisionerJobLog{
		CreatedAt: ts,
		Level:     optimus-ide-collabsdk.LogLevelInfo,
		Source:    optimus-ide-collabsdk.LogSourceProvisioner,
		Stage:     "Applying",
		Output:    "\033[32mSuccess!\033[0m Unicode: 你好世界",
	}
	result := log.Text()
	require.Equal(t, "2024-01-28T10:30:00Z [info] [provisioner|Applying] \033[32mSuccess!\033[0m Unicode: 你好世界", result)
}

func TestWorkspaceAgentLogText(t *testing.T) {
	t.Parallel()

	ts := time.Date(2024, 1, 28, 10, 30, 0, 0, time.UTC)
	log := optimus-ide-collabsdk.WorkspaceAgentLog{
		CreatedAt: ts,
		Level:     optimus-ide-collabsdk.LogLevelInfo,
		Output:    "Agent started successfully",
		SourceID:  uuid.New(),
	}
	result := log.Text("main", "startup_script")
	require.Equal(t, "2024-01-28T10:30:00Z [info] [agent.main|startup_script] Agent started successfully", result)
}

func TestWorkspaceAgentLogTextEmptySourceAndAgent(t *testing.T) {
	t.Parallel()

	ts := time.Date(2024, 1, 28, 10, 30, 0, 0, time.UTC)
	log := optimus-ide-collabsdk.WorkspaceAgentLog{
		CreatedAt: ts,
		Level:     optimus-ide-collabsdk.LogLevelWarn,
		Output:    "Warning message",
		SourceID:  uuid.New(),
	}
	result := log.Text("", "")
	require.Equal(t, "2024-01-28T10:30:00Z [warn] [agent] Warning message", result)
}

func TestWorkspaceAgentLogTextMultiline(t *testing.T) {
	t.Parallel()

	ts := time.Date(2024, 1, 28, 10, 30, 0, 0, time.UTC)
	log := optimus-ide-collabsdk.WorkspaceAgentLog{
		CreatedAt: ts,
		Level:     optimus-ide-collabsdk.LogLevelInfo,
		Output:    "Line 1\nLine 2\nLine 3",
		SourceID:  uuid.New(),
	}
	result := log.Text("main", "startup_script")
	require.Equal(t, "2024-01-28T10:30:00Z [info] [agent.main|startup_script] Line 1\nLine 2\nLine 3", result)
}

func TestWorkspaceAgentLogTextSpecialChars(t *testing.T) {
	t.Parallel()

	ts := time.Date(2024, 1, 28, 10, 30, 0, 0, time.UTC)
	log := optimus-ide-collabsdk.WorkspaceAgentLog{
		CreatedAt: ts,
		Level:     optimus-ide-collabsdk.LogLevelDebug,
		Output:    "\033[31mError!\033[0m 🚀 Unicode: 日本語",
		SourceID:  uuid.New(),
	}
	result := log.Text("main", "startup_script")
	require.Equal(t, "2024-01-28T10:30:00Z [debug] [agent.main|startup_script] \033[31mError!\033[0m 🚀 Unicode: 日本語", result)
}

func TestWorkspaceAgentDevcontainerEquals(t *testing.T) {
	t.Parallel()

	agentID := uuid.New()

	base := optimus-ide-collabsdk.WorkspaceAgentDevcontainer{
		ID:              uuid.New(),
		Name:            "test-dc",
		WorkspaceFolder: "/workspace",
		Status:          optimus-ide-collabsdk.WorkspaceAgentDevcontainerStatusRunning,
		Dirty:           false,
		Container:       &optimus-ide-collabsdk.WorkspaceAgentContainer{ID: "container-123"},
		Agent:           &optimus-ide-collabsdk.WorkspaceAgentDevcontainerAgent{ID: agentID, Name: "agent-1"},
		Error:           "",
	}

	tests := []struct {
		name      string
		modify    func(*optimus-ide-collabsdk.WorkspaceAgentDevcontainer)
		wantEqual bool
	}{
		{
			name:      "identical",
			modify:    func(d *optimus-ide-collabsdk.WorkspaceAgentDevcontainer) {},
			wantEqual: true,
		},
		{
			name:      "different ID",
			modify:    func(d *optimus-ide-collabsdk.WorkspaceAgentDevcontainer) { d.ID = uuid.New() },
			wantEqual: false,
		},
		{
			name:      "different Name",
			modify:    func(d *optimus-ide-collabsdk.WorkspaceAgentDevcontainer) { d.Name = "other-dc" },
			wantEqual: false,
		},
		{
			name:      "different WorkspaceFolder",
			modify:    func(d *optimus-ide-collabsdk.WorkspaceAgentDevcontainer) { d.WorkspaceFolder = "/other" },
			wantEqual: false,
		},
		{
			name: "different SubagentID (one valid, one nil)",
			modify: func(d *optimus-ide-collabsdk.WorkspaceAgentDevcontainer) {
				d.SubagentID = uuid.NullUUID{Valid: true, UUID: uuid.New()}
			},
			wantEqual: false,
		},
		{
			name: "different SubagentID UUIDs",
			modify: func(d *optimus-ide-collabsdk.WorkspaceAgentDevcontainer) {
				d.SubagentID = uuid.NullUUID{Valid: true, UUID: uuid.New()}
			},
			wantEqual: false,
		},
		{
			name: "different Status",
			modify: func(d *optimus-ide-collabsdk.WorkspaceAgentDevcontainer) {
				d.Status = optimus-ide-collabsdk.WorkspaceAgentDevcontainerStatusStopped
			},
			wantEqual: false,
		},
		{
			name:      "different Dirty",
			modify:    func(d *optimus-ide-collabsdk.WorkspaceAgentDevcontainer) { d.Dirty = true },
			wantEqual: false,
		},
		{
			name:      "different Container (one nil)",
			modify:    func(d *optimus-ide-collabsdk.WorkspaceAgentDevcontainer) { d.Container = nil },
			wantEqual: false,
		},
		{
			name: "different Container IDs",
			modify: func(d *optimus-ide-collabsdk.WorkspaceAgentDevcontainer) {
				d.Container = &optimus-ide-collabsdk.WorkspaceAgentContainer{ID: "different-container"}
			},
			wantEqual: false,
		},
		{
			name:      "different Agent (one nil)",
			modify:    func(d *optimus-ide-collabsdk.WorkspaceAgentDevcontainer) { d.Agent = nil },
			wantEqual: false,
		},
		{
			name: "different Agent values",
			modify: func(d *optimus-ide-collabsdk.WorkspaceAgentDevcontainer) {
				d.Agent = &optimus-ide-collabsdk.WorkspaceAgentDevcontainerAgent{ID: agentID, Name: "agent-2"}
			},
			wantEqual: false,
		},
		{
			name:      "different Error",
			modify:    func(d *optimus-ide-collabsdk.WorkspaceAgentDevcontainer) { d.Error = "some error" },
			wantEqual: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			modified := base
			tt.modify(&modified)
			require.Equal(t, tt.wantEqual, base.Equals(modified))
		})
	}
}

func TestWorkspaceAgentDevcontainerIsTerraformDefined(t *testing.T) {
	t.Parallel()

	t.Run("SubagentID Valid", func(t *testing.T) {
		t.Parallel()

		dc := optimus-ide-collabsdk.WorkspaceAgentDevcontainer{
			ID:              uuid.New(),
			Name:            "test-dc",
			WorkspaceFolder: "/workspace",
			SubagentID:      uuid.NullUUID{Valid: true, UUID: uuid.New()},
		}

		require.True(t, dc.IsTerraformDefined())
	})

	t.Run("SubagentID Null", func(t *testing.T) {
		t.Parallel()

		dc := optimus-ide-collabsdk.WorkspaceAgentDevcontainer{
			ID:              uuid.New(),
			Name:            "test-dc",
			WorkspaceFolder: "/workspace",
			SubagentID:      uuid.NullUUID{Valid: false},
		}

		require.False(t, dc.IsTerraformDefined())
	})
}
