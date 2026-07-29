package agentselect_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/x/chatd/agentselect"
)

func TestFindChatAgent(t *testing.T) {
	t.Parallel()

	newRootAgentWithID := func(id, name string, displayOrder int32) database.WorkspaceAgent {
		return database.WorkspaceAgent{
			ID:           uuid.MustParse(id),
			Name:         name,
			DisplayOrder: displayOrder,
		}
	}

	newRootAgent := func(name string, displayOrder int32) database.WorkspaceAgent {
		return newRootAgentWithID(uuid.NewString(), name, displayOrder)
	}

	newChildAgent := func(name string, displayOrder int32) database.WorkspaceAgent {
		agent := newRootAgent(name, displayOrder)
		agent.ParentID = uuid.NullUUID{UUID: uuid.New(), Valid: true}
		return agent
	}

	tests := []struct {
		name            string
		agents          []database.WorkspaceAgent
		wantIndex       int
		wantErrContains []string
	}{
		{
			name: "SingleSuffixMatch",
			agents: []database.WorkspaceAgent{
				newRootAgent("alpha", 0),
				newRootAgent("dev-optimus-ide-collabd-chat", 2),
				newRootAgent("zeta", 1),
			},
			wantIndex: 1,
		},
		{
			name: "SuffixMatchCaseInsensitive",
			agents: []database.WorkspaceAgent{
				newRootAgent("alpha", 0),
				newRootAgent("Dev-Optimus-IDE-Collabd-Chat", 2),
				newRootAgent("zeta", 1),
			},
			wantIndex: 1,
		},
		{
			name: "NoSuffixMatchFallbackDeterministic",
			agents: []database.WorkspaceAgent{
				newRootAgent("zeta", 2),
				newRootAgent("bravo", 1),
				newRootAgent("alpha", 1),
			},
			wantIndex: 2,
		},
		{
			name: "NoSuffixMatchFallbackByName",
			agents: []database.WorkspaceAgent{
				newRootAgent("Bravo", 3),
				newRootAgent("alpha", 3),
				newRootAgent("charlie", 3),
			},
			wantIndex: 1,
		},
		{
			name: "CaseOnlyNameTieFallbackDeterministic",
			agents: []database.WorkspaceAgent{
				newRootAgent("Dev", 0),
				newRootAgent("dev", 0),
			},
			wantIndex: 0,
		},
		{
			name: "ExactNameTieFallbackByID",
			agents: []database.WorkspaceAgent{
				newRootAgentWithID("00000000-0000-0000-0000-000000000002", "dev", 0),
				newRootAgentWithID("00000000-0000-0000-0000-000000000001", "dev", 0),
			},
			wantIndex: 1,
		},
		{
			name: "MultipleSuffixMatchesError",
			agents: []database.WorkspaceAgent{
				newRootAgent("alpha-optimus-ide-collabd-chat", 2),
				newRootAgent("beta-optimus-ide-collabd-chat", 1),
				newRootAgent("gamma", 0),
			},
			wantErrContains: []string{
				fmt.Sprintf(
					"multiple agents match the chat suffix %q",
					agentselect.Suffix,
				),
				"alpha-optimus-ide-collabd-chat",
				"beta-optimus-ide-collabd-chat",
				"only one agent should use this suffix",
			},
		},
		{
			name: "ChildAgentSuffixIgnored",
			agents: []database.WorkspaceAgent{
				newRootAgent("alpha", 1),
				newChildAgent("child-optimus-ide-collabd-chat", 0),
				newRootAgent("bravo", 0),
			},
			wantIndex: 2,
		},
		{
			name: "ChildAgentSuffixIgnoredWithRootMatch",
			agents: []database.WorkspaceAgent{
				newRootAgent("alpha", 0),
				newChildAgent("child-optimus-ide-collabd-chat", 1),
				newRootAgent("root-optimus-ide-collabd-chat", 2),
			},
			wantIndex: 2,
		},
		{
			name:   "EmptyAgentList",
			agents: []database.WorkspaceAgent{},
			wantErrContains: []string{
				"no eligible workspace agents found",
			},
		},
		{
			name: "OnlyChildAgents",
			agents: []database.WorkspaceAgent{
				newChildAgent("alpha", 0),
				newChildAgent("beta-optimus-ide-collabd-chat", 1),
			},
			wantErrContains: []string{
				"no eligible workspace agents found",
			},
		},
		{
			name: "SingleRootAgent",
			agents: []database.WorkspaceAgent{
				newRootAgent("solo", 5),
			},
			wantIndex: 0,
		},
		{
			name: "SuffixAgentWinsRegardlessOfOrder",
			agents: []database.WorkspaceAgent{
				newRootAgent("alpha", 0),
				newRootAgent("zeta", 1),
				newRootAgent("preferred-optimus-ide-collabd-chat", 99),
			},
			wantIndex: 2,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := agentselect.FindChatAgent(tt.agents)
			if len(tt.wantErrContains) > 0 {
				require.Error(t, err)
				for _, wantErr := range tt.wantErrContains {
					require.ErrorContains(t, err, wantErr)
				}
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.agents[tt.wantIndex], got)
		})
	}
}

func TestIsChatAgent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "ExactSuffix",
			input: "agent-optimus-ide-collabd-chat",
			want:  true,
		},
		{
			name:  "UppercaseSuffix",
			input: "agent-OPTIMUS-IDE-COLLABD-CHAT",
			want:  true,
		},
		{
			name:  "MixedCaseSuffix",
			input: "agent-Optimus-IDE-Collabd-Chat",
			want:  true,
		},
		{
			name:  "NoSuffix",
			input: "my-agent",
			want:  false,
		},
		{
			name:  "SuffixOnly",
			input: "-optimus-ide-collabd-chat",
			want:  true,
		},
		{
			name:  "PartialSuffix",
			input: "agent-optimus-ide-collabd",
			want:  false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.want, agentselect.IsChatAgent(tt.input))
		})
	}
}
