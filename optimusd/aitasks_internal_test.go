package optimus-ide-collabd

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/util/ptr"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

func TestDeriveTaskCurrentState_Unit(t *testing.T) {
	t.Parallel()

	now := time.Now()
	tests := []struct {
		name               string
		task               database.Task
		agentLifecycle     *optimus-ide-collabsdk.WorkspaceAgentLifecycle
		appHealth          *optimus-ide-collabsdk.WorkspaceAppHealth
		latestAppStatus    *optimus-ide-collabsdk.WorkspaceAppStatus
		latestBuild        optimus-ide-collabsdk.WorkspaceBuild
		expectCurrentState bool
		expectedTimestamp  time.Time
		expectedState      optimus-ide-collabsdk.TaskState
		expectedMessage    string
	}{
		{
			name: "NoAppStatus",
			task: database.Task{
				ID:     uuid.New(),
				Status: database.TaskStatusActive,
			},
			agentLifecycle:  nil,
			appHealth:       nil,
			latestAppStatus: nil,
			latestBuild: optimus-ide-collabsdk.WorkspaceBuild{
				Transition: optimus-ide-collabsdk.WorkspaceTransitionStart,
				CreatedAt:  now,
			},
			expectCurrentState: false,
		},
		{
			name: "BuildStartTransition_AppStatus_NewerThanBuild",
			task: database.Task{
				ID:     uuid.New(),
				Status: database.TaskStatusActive,
			},
			agentLifecycle: nil,
			appHealth:      nil,
			latestAppStatus: &optimus-ide-collabsdk.WorkspaceAppStatus{
				State:     optimus-ide-collabsdk.WorkspaceAppStatusStateWorking,
				Message:   "Task is working",
				CreatedAt: now.Add(1 * time.Minute),
			},
			latestBuild: optimus-ide-collabsdk.WorkspaceBuild{
				Transition: optimus-ide-collabsdk.WorkspaceTransitionStart,
				CreatedAt:  now,
			},
			expectCurrentState: true,
			expectedTimestamp:  now.Add(1 * time.Minute),
			expectedState:      optimus-ide-collabsdk.TaskState(optimus-ide-collabsdk.WorkspaceAppStatusStateWorking),
			expectedMessage:    "Task is working",
		},
		{
			name: "BuildStartTransition_StaleAppStatus_OlderThanBuild",
			task: database.Task{
				ID:     uuid.New(),
				Status: database.TaskStatusActive,
			},
			agentLifecycle: nil,
			appHealth:      nil,
			latestAppStatus: &optimus-ide-collabsdk.WorkspaceAppStatus{
				State:     optimus-ide-collabsdk.WorkspaceAppStatusStateComplete,
				Message:   "Previous task completed",
				CreatedAt: now.Add(-1 * time.Minute),
			},
			latestBuild: optimus-ide-collabsdk.WorkspaceBuild{
				Transition: optimus-ide-collabsdk.WorkspaceTransitionStart,
				CreatedAt:  now,
			},
			expectCurrentState: false,
		},
		{
			name: "BuildStopTransition",
			task: database.Task{
				ID:     uuid.New(),
				Status: database.TaskStatusActive,
			},
			agentLifecycle: nil,
			appHealth:      nil,
			latestAppStatus: &optimus-ide-collabsdk.WorkspaceAppStatus{
				State:     optimus-ide-collabsdk.WorkspaceAppStatusStateComplete,
				Message:   "Task completed before stop",
				CreatedAt: now.Add(-1 * time.Minute),
			},
			latestBuild: optimus-ide-collabsdk.WorkspaceBuild{
				Transition: optimus-ide-collabsdk.WorkspaceTransitionStop,
				CreatedAt:  now,
			},
			expectCurrentState: true,
			expectedTimestamp:  now.Add(-1 * time.Minute),
			expectedState:      optimus-ide-collabsdk.TaskState(optimus-ide-collabsdk.WorkspaceAppStatusStateComplete),
			expectedMessage:    "Task completed before stop",
		},
		{
			name: "TaskInitializing_WorkspacePending",
			task: database.Task{
				ID:     uuid.New(),
				Status: database.TaskStatusInitializing,
			},
			agentLifecycle:  nil,
			appHealth:       nil,
			latestAppStatus: nil,
			latestBuild: optimus-ide-collabsdk.WorkspaceBuild{
				Status:    optimus-ide-collabsdk.WorkspaceStatusPending,
				CreatedAt: now,
			},
			expectCurrentState: true,
			expectedTimestamp:  now,
			expectedState:      optimus-ide-collabsdk.TaskStateWorking,
			expectedMessage:    "Workspace is pending",
		},
		{
			name: "TaskInitializing_WorkspaceStarting",
			task: database.Task{
				ID:     uuid.New(),
				Status: database.TaskStatusInitializing,
			},
			agentLifecycle:  nil,
			appHealth:       nil,
			latestAppStatus: nil,
			latestBuild: optimus-ide-collabsdk.WorkspaceBuild{
				Status:    optimus-ide-collabsdk.WorkspaceStatusStarting,
				CreatedAt: now,
			},
			expectCurrentState: true,
			expectedTimestamp:  now,
			expectedState:      optimus-ide-collabsdk.TaskStateWorking,
			expectedMessage:    "Workspace is starting",
		},
		{
			name: "TaskInitializing_AgentConnecting",
			task: database.Task{
				ID:     uuid.New(),
				Status: database.TaskStatusInitializing,
			},
			agentLifecycle:  ptr.Ref(optimus-ide-collabsdk.WorkspaceAgentLifecycleCreated),
			appHealth:       nil,
			latestAppStatus: nil,
			latestBuild: optimus-ide-collabsdk.WorkspaceBuild{
				Status:    optimus-ide-collabsdk.WorkspaceStatusRunning,
				CreatedAt: now,
			},
			expectCurrentState: true,
			expectedTimestamp:  now,
			expectedState:      optimus-ide-collabsdk.TaskStateWorking,
			expectedMessage:    "Agent is connecting",
		},
		{
			name: "TaskInitializing_AgentStarting",
			task: database.Task{
				ID:     uuid.New(),
				Status: database.TaskStatusInitializing,
			},
			agentLifecycle:  ptr.Ref(optimus-ide-collabsdk.WorkspaceAgentLifecycleStarting),
			appHealth:       nil,
			latestAppStatus: nil,
			latestBuild: optimus-ide-collabsdk.WorkspaceBuild{
				Status:    optimus-ide-collabsdk.WorkspaceStatusRunning,
				CreatedAt: now,
			},
			expectCurrentState: true,
			expectedTimestamp:  now,
			expectedState:      optimus-ide-collabsdk.TaskStateWorking,
			expectedMessage:    "Agent is starting",
		},
		{
			name: "TaskInitializing_AppInitializing",
			task: database.Task{
				ID:     uuid.New(),
				Status: database.TaskStatusInitializing,
			},
			agentLifecycle:  ptr.Ref(optimus-ide-collabsdk.WorkspaceAgentLifecycleReady),
			appHealth:       ptr.Ref(optimus-ide-collabsdk.WorkspaceAppHealthInitializing),
			latestAppStatus: nil,
			latestBuild: optimus-ide-collabsdk.WorkspaceBuild{
				Status:    optimus-ide-collabsdk.WorkspaceStatusRunning,
				CreatedAt: now,
			},
			expectCurrentState: true,
			expectedTimestamp:  now,
			expectedState:      optimus-ide-collabsdk.TaskStateWorking,
			expectedMessage:    "App is initializing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ws := optimus-ide-collabsdk.Workspace{
				LatestBuild:     tt.latestBuild,
				LatestAppStatus: tt.latestAppStatus,
			}

			currentState := deriveTaskCurrentState(tt.task, ws, tt.agentLifecycle, tt.appHealth)

			if tt.expectCurrentState {
				require.NotNil(t, currentState)
				assert.Equal(t, tt.expectedTimestamp.UTC(), currentState.Timestamp.UTC())
				assert.Equal(t, tt.expectedState, currentState.State)
				assert.Equal(t, tt.expectedMessage, currentState.Message)
			} else {
				assert.Nil(t, currentState)
			}
		})
	}
}
