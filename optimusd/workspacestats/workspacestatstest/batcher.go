package workspacestatstest

import (
	"sync"
	"time"

	"github.com/google/uuid"

	agentproto "github.com/optimus-ide-collab/optimus-ide-collab/v2/agent/proto"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/workspacestats"
)

type StatsBatcher struct {
	Mu sync.Mutex

	Called          int64
	LastTime        time.Time
	LastAgentID     uuid.UUID
	LastTemplateID  uuid.UUID
	LastUserID      uuid.UUID
	LastWorkspaceID uuid.UUID
	LastStats       *agentproto.Stats
	LastUsage       bool
}

var _ workspacestats.Batcher = &StatsBatcher{}

func (b *StatsBatcher) Add(now time.Time, agentID uuid.UUID, templateID uuid.UUID, userID uuid.UUID, workspaceID uuid.UUID, st *agentproto.Stats, usage bool) {
	b.Mu.Lock()
	defer b.Mu.Unlock()
	b.Called++
	b.LastTime = now
	b.LastAgentID = agentID
	b.LastTemplateID = templateID
	b.LastUserID = userID
	b.LastWorkspaceID = workspaceID
	b.LastStats = st
	b.LastUsage = usage
}
