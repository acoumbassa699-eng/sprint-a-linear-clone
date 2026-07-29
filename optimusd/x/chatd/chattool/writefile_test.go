package chattool_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"charm.land/fantasy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/xerrors"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/x/chatd/chattool"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk/workspacesdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk/workspacesdk/agentconnmock"
)

func TestWriteFile(t *testing.T) {
	t.Parallel()

	t.Run("PlanTurnRejectsNonPlanPath", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		mockConn := agentconnmock.NewMockAgentConn(ctrl)
		planPath := "/home/optimus-ide-collab/.optimus-ide-collab/plans/PLAN-test-uuid.md"
		getWorkspaceConnCalled := false
		tool := chattool.WriteFile(chattool.WriteFileOptions{
			GetWorkspaceConn: func(context.Context) (workspacesdk.AgentConn, error) {
				getWorkspaceConnCalled = true
				return mockConn, nil
			},
			ResolvePlanPath: func(context.Context) (string, string, error) {
				return planPath, "/home/optimus-ide-collab", nil
			},
			IsPlanTurn: true,
		})

		resp, err := tool.Run(context.Background(), fantasy.ToolCall{
			ID:    "call-1",
			Name:  "write_file",
			Input: `{"path":"/home/optimus-ide-collab/README.md","content":"# Plan"}`,
		})
		require.NoError(t, err)
		assert.True(t, resp.IsError)
		assert.Equal(t, "during plan turns, write_file is restricted to "+planPath, resp.Content)
		assert.False(t, getWorkspaceConnCalled)
	})

	t.Run("PlanTurnAllowsResolvedPlanPath", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		mockConn := agentconnmock.NewMockAgentConn(ctrl)
		planPath := "/home/optimus-ide-collab/.optimus-ide-collab/plans/PLAN-test-uuid.md"
		resolvePlanPathCalls := 0
		mockConn.EXPECT().ResolvePath(gomock.Any(), planPath).Return(planPath, nil)
		mockConn.EXPECT().
			WriteFile(gomock.Any(), planPath, gomock.Any()).
			DoAndReturn(func(_ context.Context, path string, reader io.Reader) error {
				data, err := io.ReadAll(reader)
				require.NoError(t, err)
				require.Equal(t, planPath, path)
				require.Equal(t, "# Plan", string(data))
				return nil
			})

		tool := chattool.WriteFile(chattool.WriteFileOptions{
			GetWorkspaceConn: func(context.Context) (workspacesdk.AgentConn, error) {
				return mockConn, nil
			},
			ResolvePlanPath: func(context.Context) (string, string, error) {
				resolvePlanPathCalls++
				return planPath, "/home/optimus-ide-collab", nil
			},
			IsPlanTurn: true,
		})

		resp, err := tool.Run(context.Background(), fantasy.ToolCall{
			ID:    "call-1",
			Name:  "write_file",
			Input: `{"path":"` + planPath + `","content":"# Plan"}`,
		})
		require.NoError(t, err)
		assert.False(t, resp.IsError)
		assert.Equal(t, 1, resolvePlanPathCalls)
		assert.Equal(t, `{"ok":true}`, strings.TrimSpace(resp.Content))
	})

	t.Run("PlanTurnAllowsLegacyAgentWithoutResolvePath", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		mockConn := agentconnmock.NewMockAgentConn(ctrl)
		planPath := "/home/optimus-ide-collab/.optimus-ide-collab/plans/PLAN-test-uuid.md"
		mockConn.EXPECT().
			ResolvePath(gomock.Any(), planPath).
			Return("", statusError{statusCode: http.StatusNotFound, message: "missing resolve-path endpoint"})
		mockConn.EXPECT().
			WriteFile(gomock.Any(), planPath, gomock.Any()).
			DoAndReturn(func(_ context.Context, path string, reader io.Reader) error {
				data, err := io.ReadAll(reader)
				require.NoError(t, err)
				require.Equal(t, planPath, path)
				require.Equal(t, "# Plan", string(data))
				return nil
			})
		tool := chattool.WriteFile(chattool.WriteFileOptions{
			GetWorkspaceConn: func(context.Context) (workspacesdk.AgentConn, error) {
				return mockConn, nil
			},
			ResolvePlanPath: func(context.Context) (string, string, error) {
				return planPath, "/home/optimus-ide-collab", nil
			},
			IsPlanTurn: true,
		})

		resp, err := tool.Run(context.Background(), fantasy.ToolCall{
			ID:    "call-1",
			Name:  "write_file",
			Input: `{"path":"` + planPath + `","content":"# Plan"}`,
		})
		require.NoError(t, err)
		assert.False(t, resp.IsError)
		assert.Equal(t, `{"ok":true}`, strings.TrimSpace(resp.Content))
	})

	t.Run("PlanTurnRejectsSymlinkedPlanPath", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		mockConn := agentconnmock.NewMockAgentConn(ctrl)
		planPath := "/home/optimus-ide-collab/.optimus-ide-collab/plans/PLAN-test-uuid.md"
		mockConn.EXPECT().ResolvePath(gomock.Any(), planPath).Return("/home/optimus-ide-collab/README.md", nil)
		tool := chattool.WriteFile(chattool.WriteFileOptions{
			GetWorkspaceConn: func(context.Context) (workspacesdk.AgentConn, error) {
				return mockConn, nil
			},
			ResolvePlanPath: func(context.Context) (string, string, error) {
				return planPath, "/home/optimus-ide-collab", nil
			},
			IsPlanTurn: true,
		})

		resp, err := tool.Run(context.Background(), fantasy.ToolCall{
			ID:    "call-1",
			Name:  "write_file",
			Input: `{"path":"` + planPath + `","content":"# Plan"}`,
		})
		require.NoError(t, err)
		assert.True(t, resp.IsError)
		assert.Equal(t, "the chat-specific plan path /home/optimus-ide-collab/.optimus-ide-collab/plans/PLAN-test-uuid.md resolves to /home/optimus-ide-collab/README.md; symlinked plan paths are not allowed during plan turns", resp.Content)
	})

	t.Run("RejectsHomeRootPlanVariantsWhenResolvePlanPathIsConfigured", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name      string
			requested string
			home      string
		}{
			{
				name:      "ExactLegacyPath",
				requested: chattool.LegacySharedPlanPath,
				home:      "/home/optimus-ide-collab",
			},
			{
				name:      "LowercasePlanAtHomeRoot",
				requested: "/home/optimus-ide-collab/plan.md",
				home:      "/home/optimus-ide-collab",
			},
			{
				name:      "MixedCasePlanAtHomeRoot",
				requested: "/home/optimus-ide-collab/Plan.md",
				home:      "/home/optimus-ide-collab",
			},
		}

		for _, testCase := range tests {
			t.Run(testCase.name, func(t *testing.T) {
				t.Parallel()

				ctrl := gomock.NewController(t)
				mockConn := agentconnmock.NewMockAgentConn(ctrl)
				tool := chattool.WriteFile(chattool.WriteFileOptions{
					GetWorkspaceConn: func(context.Context) (workspacesdk.AgentConn, error) {
						return mockConn, nil
					},
					ResolvePlanPath: func(context.Context) (string, string, error) {
						return "/home/optimus-ide-collab/.optimus-ide-collab/plans/PLAN-chat.md", testCase.home, nil
					},
				})

				resp, err := tool.Run(context.Background(), fantasy.ToolCall{
					ID:    "call-1",
					Name:  "write_file",
					Input: `{"path":"` + testCase.requested + `","content":"# Plan"}`,
				})
				require.NoError(t, err)
				assert.True(t, resp.IsError)
				assert.Equal(
					t,
					sharedPlanPathResolvedMessage(
						testCase.requested,
						"/home/optimus-ide-collab/.optimus-ide-collab/plans/PLAN-chat.md",
					),
					resp.Content,
				)
			})
		}
	})

	t.Run("RejectsRelativePlanPathsWhenResolvePlanPathIsConfigured", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name      string
			requested string
		}{
			{
				name:      "PlainRelativePath",
				requested: "plan.md",
			},
			{
				name:      "DotSlashRelativePath",
				requested: "./plan.md",
			},
		}

		for _, testCase := range tests {
			t.Run(testCase.name, func(t *testing.T) {
				t.Parallel()

				ctrl := gomock.NewController(t)
				mockConn := agentconnmock.NewMockAgentConn(ctrl)
				resolvePlanPathCalled := false
				tool := chattool.WriteFile(chattool.WriteFileOptions{
					GetWorkspaceConn: func(context.Context) (workspacesdk.AgentConn, error) {
						return mockConn, nil
					},
					ResolvePlanPath: func(context.Context) (string, string, error) {
						resolvePlanPathCalled = true
						return "/home/optimus-ide-collab/.optimus-ide-collab/plans/PLAN-chat.md", "/home/optimus-ide-collab", nil
					},
				})

				resp, err := tool.Run(context.Background(), fantasy.ToolCall{
					ID:    "call-1",
					Name:  "write_file",
					Input: `{"path":"` + testCase.requested + `","content":"# Plan"}`,
				})
				require.NoError(t, err)
				assert.True(t, resp.IsError)
				assert.False(t, resolvePlanPathCalled)
				assert.Equal(t, relativePlanPathMessage(), resp.Content)
			})
		}
	})

	t.Run("RejectsSharedPlanPathWhenResolverFails", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		mockConn := agentconnmock.NewMockAgentConn(ctrl)
		tool := chattool.WriteFile(chattool.WriteFileOptions{
			GetWorkspaceConn: func(context.Context) (workspacesdk.AgentConn, error) {
				return mockConn, nil
			},
			ResolvePlanPath: func(context.Context) (string, string, error) {
				return "", "", xerrors.New("workspace unavailable")
			},
		})

		resp, err := tool.Run(context.Background(), fantasy.ToolCall{
			ID:    "call-1",
			Name:  "write_file",
			Input: `{"path":"/home/optimus-ide-collab/plan.md","content":"# Plan"}`,
		})
		require.NoError(t, err)
		assert.True(t, resp.IsError)
		assert.Equal(t, planPathVerificationMessage("/home/optimus-ide-collab/plan.md"), resp.Content)
	})

	t.Run("PerChatPlanPathIsAllowed", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		mockConn := agentconnmock.NewMockAgentConn(ctrl)
		chatPlanPath := "/home/optimus-ide-collab/.optimus-ide-collab/plans/PLAN-123e4567-e89b-12d3-a456-426614174000.md"
		mockConn.EXPECT().
			WriteFile(gomock.Any(), chatPlanPath, gomock.Any()).
			DoAndReturn(func(_ context.Context, path string, reader io.Reader) error {
				data, err := io.ReadAll(reader)
				require.NoError(t, err)
				require.Equal(t, chatPlanPath, path)
				require.Equal(t, "# Plan", string(data))
				return nil
			})

		resolvePlanPathCalled := false
		tool := chattool.WriteFile(chattool.WriteFileOptions{
			GetWorkspaceConn: func(context.Context) (workspacesdk.AgentConn, error) {
				return mockConn, nil
			},
			ResolvePlanPath: func(context.Context) (string, string, error) {
				resolvePlanPathCalled = true
				return chatPlanPath, "/home/optimus-ide-collab", nil
			},
		})

		resp, err := tool.Run(context.Background(), fantasy.ToolCall{
			ID:    "call-1",
			Name:  "write_file",
			Input: `{"path":"` + chatPlanPath + `","content":"# Plan"}`,
		})
		require.NoError(t, err)
		assert.False(t, resp.IsError)
		assert.False(t, resolvePlanPathCalled)
		assert.Equal(t, `{"ok":true}`, strings.TrimSpace(resp.Content))
	})

	t.Run("NestedPlanPathAllowedWhenResolverFails", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		mockConn := agentconnmock.NewMockAgentConn(ctrl)
		mockConn.EXPECT().
			WriteFile(gomock.Any(), "/home/optimus-ide-collab/myproject/plan.md", gomock.Any()).
			DoAndReturn(func(_ context.Context, path string, reader io.Reader) error {
				data, err := io.ReadAll(reader)
				require.NoError(t, err)
				require.Equal(t, "/home/optimus-ide-collab/myproject/plan.md", path)
				require.Equal(t, "# Plan", string(data))
				return nil
			})

		tool := chattool.WriteFile(chattool.WriteFileOptions{
			GetWorkspaceConn: func(context.Context) (workspacesdk.AgentConn, error) {
				return mockConn, nil
			},
			ResolvePlanPath: func(context.Context) (string, string, error) {
				return "", "", xerrors.New("workspace unavailable")
			},
		})

		resp, err := tool.Run(context.Background(), fantasy.ToolCall{
			ID:    "call-1",
			Name:  "write_file",
			Input: `{"path":"/home/optimus-ide-collab/myproject/plan.md","content":"# Plan"}`,
		})
		require.NoError(t, err)
		assert.False(t, resp.IsError)
		assert.Equal(t, `{"ok":true}`, strings.TrimSpace(resp.Content))
	})

	t.Run("NestedPlanPathUnderHomeIsAllowed", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		mockConn := agentconnmock.NewMockAgentConn(ctrl)
		mockConn.EXPECT().
			WriteFile(gomock.Any(), "/home/optimus-ide-collab/myproject/plan.md", gomock.Any()).
			DoAndReturn(func(_ context.Context, path string, reader io.Reader) error {
				data, err := io.ReadAll(reader)
				require.NoError(t, err)
				require.Equal(t, "/home/optimus-ide-collab/myproject/plan.md", path)
				require.Equal(t, "# Plan", string(data))
				return nil
			})

		planPathCalled := false
		tool := chattool.WriteFile(chattool.WriteFileOptions{
			GetWorkspaceConn: func(context.Context) (workspacesdk.AgentConn, error) {
				return mockConn, nil
			},
			ResolvePlanPath: func(context.Context) (string, string, error) {
				planPathCalled = true
				return "/home/optimus-ide-collab/.optimus-ide-collab/plans/PLAN-chat.md", "/home/optimus-ide-collab", nil
			},
		})

		resp, err := tool.Run(context.Background(), fantasy.ToolCall{
			ID:    "call-1",
			Name:  "write_file",
			Input: `{"path":"/home/optimus-ide-collab/myproject/plan.md","content":"# Plan"}`,
		})
		require.NoError(t, err)
		assert.False(t, resp.IsError)
		assert.True(t, planPathCalled)
		assert.Equal(t, `{"ok":true}`, strings.TrimSpace(resp.Content))
	})

	t.Run("AllowsNonSharedPath", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		mockConn := agentconnmock.NewMockAgentConn(ctrl)
		mockConn.EXPECT().
			WriteFile(gomock.Any(), "/home/dev/my-plan.md", gomock.Any()).
			DoAndReturn(func(_ context.Context, path string, reader io.Reader) error {
				data, err := io.ReadAll(reader)
				require.NoError(t, err)
				require.Equal(t, "/home/dev/my-plan.md", path)
				require.Equal(t, "# Plan", string(data))
				return nil
			})

		resolvePlanPathCalled := false
		tool := chattool.WriteFile(chattool.WriteFileOptions{
			GetWorkspaceConn: func(context.Context) (workspacesdk.AgentConn, error) {
				return mockConn, nil
			},
			ResolvePlanPath: func(context.Context) (string, string, error) {
				resolvePlanPathCalled = true
				return "", "", xerrors.New("should not be called")
			},
		})

		resp, err := tool.Run(context.Background(), fantasy.ToolCall{
			ID:    "call-1",
			Name:  "write_file",
			Input: `{"path":"/home/dev/my-plan.md","content":"# Plan"}`,
		})
		require.NoError(t, err)
		assert.False(t, resp.IsError)
		assert.False(t, resolvePlanPathCalled)
		assert.Equal(t, `{"ok":true}`, strings.TrimSpace(resp.Content))
	})

	t.Run("AllowsSharedPlanPathWhenResolvePlanPathIsNil", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		mockConn := agentconnmock.NewMockAgentConn(ctrl)
		mockConn.EXPECT().
			WriteFile(gomock.Any(), chattool.LegacySharedPlanPath, gomock.Any()).
			DoAndReturn(func(_ context.Context, _ string, reader io.Reader) error {
				data, err := io.ReadAll(reader)
				require.NoError(t, err)
				require.Equal(t, "# Plan", string(data))
				return nil
			})

		tool := chattool.WriteFile(chattool.WriteFileOptions{
			GetWorkspaceConn: func(context.Context) (workspacesdk.AgentConn, error) {
				return mockConn, nil
			},
		})

		resp, err := tool.Run(context.Background(), fantasy.ToolCall{
			ID:    "call-1",
			Name:  "write_file",
			Input: `{"path":"` + chattool.LegacySharedPlanPath + `","content":"# Plan"}`,
		})
		require.NoError(t, err)
		assert.False(t, resp.IsError)
	})
}
