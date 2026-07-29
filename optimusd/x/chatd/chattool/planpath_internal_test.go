package chattool

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsAbsolutePath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		path string
		want bool
	}{
		{"/home/optimus-ide-collab/PLAN.md", true},
		{"/workspace/project/plan.md", true},
		{"plan.md", false},
		{"./plan.md", false},
		{"../plan.md", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, isAbsolutePath(tt.path))
		})
	}
}

func TestLooksLikePlanFileName(t *testing.T) {
	t.Parallel()

	require.True(t, looksLikePlanFileName("plan.md"))
	require.True(t, looksLikePlanFileName("./Plan.md"))
	require.True(t, looksLikePlanFileName("/home/optimus-ide-collab/PLAN.md"))
	require.False(t, looksLikePlanFileName("/home/optimus-ide-collab/README.md"))
}

func TestLooksLikeLegacySharedPlanPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		requested string
		want      bool
	}{
		{
			name:      "ExactMatch",
			requested: "/home/optimus-ide-collab/PLAN.md",
			want:      true,
		},
		{
			name:      "CaseInsensitive",
			requested: "/home/optimus-ide-collab/plan.md",
			want:      true,
		},
		{
			name:      "MixedCase",
			requested: "/home/optimus-ide-collab/Plan.md",
			want:      true,
		},
		{
			name:      "NestedPath",
			requested: "/home/optimus-ide-collab/myproject/plan.md",
			want:      false,
		},
		{
			name:      "DifferentHome",
			requested: "/Users/dev/PLAN.md",
			want:      false,
		},
		{
			name:      "PerChatPath",
			requested: "/home/optimus-ide-collab/.optimus-ide-collab/plans/PLAN-123e4567-e89b-12d3-a456-426614174000.md",
			want:      false,
		},
		{
			name:      "EmptyString",
			requested: "",
			want:      false,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, testCase.want, looksLikeLegacySharedPlanPath(testCase.requested))
		})
	}
}

func TestRejectSharedPlanPath(t *testing.T) {
	t.Parallel()

	resp, rejected := rejectSharedPlanPath(
		LegacySharedPlanPath,
		"/Users/dev",
		"/Users/dev/.optimus-ide-collab/plans/PLAN-chat.md",
		nil,
	)

	require.True(t, rejected)
	require.True(t, resp.IsError)
	require.Equal(
		t,
		sharedPlanPathMessage(
			LegacySharedPlanPath,
			"/Users/dev/.optimus-ide-collab/plans/PLAN-chat.md",
		),
		resp.Content,
	)
}

func TestSharedPlanPathMessage(t *testing.T) {
	t.Parallel()

	require.Equal(
		t,
		"the plan path /home/optimus-ide-collab/plan.md is no longer supported at the home root; use the chat-specific plan path: /home/optimus-ide-collab/.optimus-ide-collab/plans/PLAN-chat.md",
		sharedPlanPathMessage(
			"/home/optimus-ide-collab/plan.md",
			"/home/optimus-ide-collab/.optimus-ide-collab/plans/PLAN-chat.md",
		),
	)
	require.Equal(
		t,
		"the plan path /home/optimus-ide-collab/plan.md could not be verified because the workspace is currently unavailable to resolve the chat-specific plan path, try again shortly",
		planPathVerificationMessage("/home/optimus-ide-collab/plan.md"),
	)
}
