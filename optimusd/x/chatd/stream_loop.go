package chatd

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/xerrors"

	"cdr.dev/slog/v3"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/db2sdk"
	optimus-ide-collabdpubsub "github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/pubsub"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/x/chatd/chatstate"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

type streamLoop struct {
	chatID uuid.UUID
	db     database.Store
	logger slog.Logger
	state  streamLocalState
}

type streamLocalState struct {
	snapshotVersion int64
	historyVersion  int64
	queueVersion    int64
	retryVersion    int64

	knownMessages map[int64]int64

	status database.ChatStatus

	errorHistoryVersion          int64
	actionRequiredHistoryVersion int64

	workerID          uuid.NullUUID
	generationAttempt int64
	lastPartSeq       int64

	afterMessageID         int64
	initialMessageSyncDone bool
}

type streamSyncHint struct {
	snapshotVersion   int64
	historyVersion    int64
	queueVersion      int64
	retryVersion      int64
	status            database.ChatStatus
	workerID          uuid.NullUUID
	generationAttempt int64
}

type streamDBSnapshot struct {
	chat database.Chat

	historyChanged  bool
	changedMessages []database.ChatMessage
	historyReset    bool
	fullHistory     []database.ChatMessage

	queueChanged bool
	queue        []database.ChatQueuedMessage

	actionRequired *optimus-ide-collabsdk.ChatStreamActionRequired
}

func newStreamLoop(chat database.Chat, db database.Store, logger slog.Logger, afterMessageID int64) *streamLoop {
	return &streamLoop{
		chatID: chat.ID,
		db:     db,
		logger: logger,
		state: streamLocalState{
			knownMessages:  make(map[int64]int64),
			afterMessageID: afterMessageID,
		},
	}
}

func streamSyncHintFromUpdate(update optimus-ide-collabdpubsub.ChatStateUpdateMessage) streamSyncHint {
	hint := streamSyncHint{
		snapshotVersion:   update.SnapshotVersion,
		historyVersion:    update.HistoryVersion,
		queueVersion:      update.QueueVersion,
		retryVersion:      update.RetryStateVersion,
		status:            database.ChatStatus(update.Status),
		generationAttempt: update.GenerationAttempt,
	}
	if update.WorkerID != nil {
		hint.workerID = uuid.NullUUID{UUID: *update.WorkerID, Valid: true}
	}
	return hint
}

func (l *streamLoop) sync(ctx context.Context, hint streamSyncHint) ([]optimus-ide-collabsdk.ChatStreamEvent, streamRelayTarget, bool, error) {
	if !l.shouldFetch(hint) {
		return nil, l.currentRelayTarget(), false, nil
	}
	return l.syncDB(ctx)
}

func (l *streamLoop) syncDB(ctx context.Context) ([]optimus-ide-collabsdk.ChatStreamEvent, streamRelayTarget, bool, error) {
	snapshot, err := l.loadDBSnapshot(ctx)
	if err != nil {
		return nil, l.currentRelayTarget(), false, err
	}
	if snapshot.chat.SnapshotVersion <= l.state.snapshotVersion {
		return nil, l.currentRelayTarget(), false, nil
	}
	return l.applyDBSnapshot(snapshot), l.currentRelayTarget(), true, nil
}

func (l *streamLoop) shouldFetch(hint streamSyncHint) bool {
	if hint.snapshotVersion <= l.state.snapshotVersion {
		return false
	}
	if hint.historyVersion > l.state.historyVersion {
		return true
	}
	if hint.queueVersion > l.state.queueVersion {
		return true
	}
	if hint.retryVersion > l.state.retryVersion {
		return true
	}
	if hint.status != l.state.status {
		return true
	}
	if !sameNullUUID(hint.workerID, l.state.workerID) {
		return true
	}
	if hint.generationAttempt != l.state.generationAttempt {
		return true
	}
	return false
}

func (l *streamLoop) loadDBSnapshot(ctx context.Context) (streamDBSnapshot, error) {
	var snapshot streamDBSnapshot
	machine := chatstate.NewChatMachine(l.db, nil, l.chatID)
	err := machine.ReadLock(ctx, func(tx database.Store) error {
		chat, err := tx.GetChatByID(ctx, l.chatID)
		if err != nil {
			return xerrors.Errorf("get chat for stream: %w", err)
		}
		snapshot.chat = chat

		if chat.HistoryVersion > l.state.historyVersion {
			snapshot.historyChanged = true
			snapshot.changedMessages, err = tx.GetChatMessagesByRevisionForStream(ctx, database.GetChatMessagesByRevisionForStreamParams{
				ChatID:        l.chatID,
				AfterRevision: l.state.historyVersion,
			})
			if err != nil {
				return xerrors.Errorf("get changed chat messages: %w", err)
			}
			for _, msg := range snapshot.changedMessages {
				if msg.Deleted {
					snapshot.historyReset = true
					break
				}
			}
			if snapshot.historyReset {
				snapshot.fullHistory, err = tx.GetChatMessagesByChatID(ctx, database.GetChatMessagesByChatIDParams{
					ChatID:  l.chatID,
					AfterID: 0,
				})
				if err != nil {
					return xerrors.Errorf("get full chat history: %w", err)
				}
			}
		}

		if chat.QueueVersion > l.state.queueVersion {
			snapshot.queueChanged = true
			snapshot.queue, err = tx.GetChatQueuedMessages(ctx, l.chatID)
			if err != nil {
				return xerrors.Errorf("get chat queue: %w", err)
			}
		}

		if chat.Status == database.ChatStatusRequiresAction {
			history := snapshot.fullHistory
			if len(history) == 0 {
				history, err = tx.GetChatMessagesByChatID(ctx, database.GetChatMessagesByChatIDParams{
					ChatID:  l.chatID,
					AfterID: 0,
				})
				if err != nil {
					return xerrors.Errorf("get requires_action history: %w", err)
				}
			}
			actionRequired, err := l.actionRequiredFromHistory(chat, history)
			if err != nil {
				return err
			}
			snapshot.actionRequired = actionRequired
		}
		return nil
	})
	if err != nil {
		return streamDBSnapshot{}, err
	}
	return snapshot, nil
}

func (*streamLoop) actionRequiredFromHistory(chat database.Chat, messages []database.ChatMessage) (*optimus-ide-collabsdk.ChatStreamActionRequired, error) {
	dynamicToolNames, err := parseDynamicToolNames(chat.DynamicTools)
	if err != nil {
		return nil, xerrors.Errorf("parse dynamic tools for stream: %w", err)
	}
	_, pending, err := unresolvedToolCallsFromHistory(messages, dynamicToolNames)
	if err != nil {
		return nil, xerrors.Errorf("derive pending dynamic tool calls: %w", err)
	}
	toolCalls := make([]optimus-ide-collabsdk.ChatStreamToolCall, 0, len(pending))
	for _, call := range pending {
		toolCalls = append(toolCalls, optimus-ide-collabsdk.ChatStreamToolCall{
			ToolCallID: call.ToolCallID,
			ToolName:   call.ToolName,
			Args:       call.Args,
		})
	}
	return &optimus-ide-collabsdk.ChatStreamActionRequired{ToolCalls: toolCalls}, nil
}

func (l *streamLoop) applyDBSnapshot(snapshot streamDBSnapshot) []optimus-ide-collabsdk.ChatStreamEvent {
	chat := snapshot.chat
	events := make([]optimus-ide-collabsdk.ChatStreamEvent, 0)
	historyChanged := chat.HistoryVersion > l.state.historyVersion
	generationChanged := chat.GenerationAttempt != l.state.generationAttempt

	if historyChanged {
		events = append(events, l.messageEvents(snapshot)...)
	}
	if !l.state.initialMessageSyncDone {
		l.state.initialMessageSyncDone = true
	}

	if chat.QueueVersion > l.state.queueVersion {
		events = append(events, optimus-ide-collabsdk.ChatStreamEvent{
			Type:           optimus-ide-collabsdk.ChatStreamEventTypeQueueUpdate,
			ChatID:         l.chatID,
			QueuedMessages: db2sdk.ChatQueuedMessages(snapshot.queue),
		})
	}

	if chat.Status != l.state.status {
		events = append(events, optimus-ide-collabsdk.ChatStreamEvent{
			Type:   optimus-ide-collabsdk.ChatStreamEventTypeStatus,
			ChatID: l.chatID,
			Status: &optimus-ide-collabsdk.ChatStreamStatus{Status: optimus-ide-collabsdk.ChatStatus(chat.Status)},
		})
	}

	if chat.Status == database.ChatStatusError && chat.HistoryVersion > l.state.errorHistoryVersion {
		events = append(events, optimus-ide-collabsdk.ChatStreamEvent{
			Type:   optimus-ide-collabsdk.ChatStreamEventTypeError,
			ChatID: l.chatID,
			Error:  l.chatError(chat),
		})
		l.state.errorHistoryVersion = chat.HistoryVersion
	}

	if chat.Status == database.ChatStatusRequiresAction && chat.HistoryVersion > l.state.actionRequiredHistoryVersion {
		actionRequired := snapshot.actionRequired
		if actionRequired == nil {
			actionRequired = &optimus-ide-collabsdk.ChatStreamActionRequired{}
		}
		events = append(events, optimus-ide-collabsdk.ChatStreamEvent{
			Type:           optimus-ide-collabsdk.ChatStreamEventTypeActionRequired,
			ChatID:         l.chatID,
			ActionRequired: actionRequired,
		})
		l.state.actionRequiredHistoryVersion = chat.HistoryVersion
	}

	if chat.RetryStateVersion > l.state.retryVersion {
		if retry := l.retryEvent(chat); retry != nil {
			events = append(events, *retry)
		}
	}

	if historyChanged || (generationChanged && chat.GenerationAttempt != 0) {
		l.state.lastPartSeq = 0
		events = append(events, optimus-ide-collabsdk.ChatStreamEvent{
			Type:   optimus-ide-collabsdk.ChatStreamEventTypePreviewReset,
			ChatID: l.chatID,
		})
	}

	l.state.snapshotVersion = chat.SnapshotVersion
	l.state.historyVersion = chat.HistoryVersion
	l.state.queueVersion = chat.QueueVersion
	l.state.retryVersion = chat.RetryStateVersion
	l.state.status = chat.Status
	l.state.workerID = chat.WorkerID
	l.state.generationAttempt = chat.GenerationAttempt
	return events
}

func (l *streamLoop) messageEvents(snapshot streamDBSnapshot) []optimus-ide-collabsdk.ChatStreamEvent {
	if snapshot.historyReset {
		events := []optimus-ide-collabsdk.ChatStreamEvent{{
			Type:   optimus-ide-collabsdk.ChatStreamEventTypeHistoryReset,
			ChatID: l.chatID,
		}}
		clear(l.state.knownMessages)
		for _, msg := range snapshot.fullHistory {
			l.state.knownMessages[msg.ID] = msg.Revision
			sdkMsg := db2sdk.ChatMessage(msg)
			events = append(events, optimus-ide-collabsdk.ChatStreamEvent{
				Type:    optimus-ide-collabsdk.ChatStreamEventTypeMessage,
				ChatID:  l.chatID,
				Message: &sdkMsg,
			})
		}
		return events
	}

	events := make([]optimus-ide-collabsdk.ChatStreamEvent, 0, len(snapshot.changedMessages))
	for _, msg := range snapshot.changedMessages {
		knownRevision := l.state.knownMessages[msg.ID]
		if knownRevision >= msg.Revision {
			continue
		}
		l.state.knownMessages[msg.ID] = msg.Revision
		if !l.state.initialMessageSyncDone && msg.ID <= l.state.afterMessageID {
			continue
		}
		sdkMsg := db2sdk.ChatMessage(msg)
		events = append(events, optimus-ide-collabsdk.ChatStreamEvent{
			Type:    optimus-ide-collabsdk.ChatStreamEventTypeMessage,
			ChatID:  l.chatID,
			Message: &sdkMsg,
		})
	}
	return events
}

func (l *streamLoop) chatError(chat database.Chat) *optimus-ide-collabsdk.ChatError {
	if !chat.LastError.Valid || len(chat.LastError.RawMessage) == 0 {
		return &optimus-ide-collabsdk.ChatError{
			Message: "The chat request failed unexpectedly.",
			Kind:    optimus-ide-collabsdk.ChatErrorKindGeneric,
		}
	}
	var payload optimus-ide-collabsdk.ChatError
	if err := json.Unmarshal(chat.LastError.RawMessage, &payload); err != nil {
		l.logger.Warn(context.Background(), "failed to parse chat stream last_error",
			slog.F("chat_id", l.chatID),
			slog.Error(err),
		)
		return &optimus-ide-collabsdk.ChatError{
			Message: "The chat request failed unexpectedly.",
			Kind:    optimus-ide-collabsdk.ChatErrorKindGeneric,
		}
	}
	if payload.Message == "" {
		payload.Message = "The chat request failed unexpectedly."
	}
	if payload.Kind == "" {
		payload.Kind = optimus-ide-collabsdk.ChatErrorKindGeneric
	}
	return &payload
}

func (l *streamLoop) retryEvent(chat database.Chat) *optimus-ide-collabsdk.ChatStreamEvent {
	if !chat.RetryState.Valid || len(chat.RetryState.RawMessage) == 0 {
		return nil
	}
	var retry optimus-ide-collabsdk.ChatStreamRetry
	if err := json.Unmarshal(chat.RetryState.RawMessage, &retry); err != nil {
		l.logger.Warn(context.Background(), "failed to parse chat stream retry_state",
			slog.F("chat_id", l.chatID),
			slog.Error(err),
		)
		return nil
	}
	return &optimus-ide-collabsdk.ChatStreamEvent{
		Type:   optimus-ide-collabsdk.ChatStreamEventTypeRetry,
		ChatID: l.chatID,
		Retry:  &retry,
	}
}

func (l *streamLoop) part(part streamPart) (event optimus-ide-collabsdk.ChatStreamEvent, accepted bool, err error) {
	if part.HistoryVersion != l.state.historyVersion || part.GenerationAttempt != l.state.generationAttempt {
		return optimus-ide-collabsdk.ChatStreamEvent{}, false, nil
	}
	if part.Seq <= l.state.lastPartSeq {
		return optimus-ide-collabsdk.ChatStreamEvent{}, false, nil
	}
	if part.Seq != l.state.lastPartSeq+1 {
		err := xerrors.Errorf(
			"chat stream message part sequence gap: got %d after %d",
			part.Seq,
			l.state.lastPartSeq,
		)
		l.logger.Error(context.Background(), "chat stream message part sequence gap",
			slog.F("chat_id", l.chatID),
			slog.F("history_version", part.HistoryVersion),
			slog.F("generation_attempt", part.GenerationAttempt),
			slog.F("last_seq", l.state.lastPartSeq),
			slog.F("seq", part.Seq),
			slog.Error(err),
		)
		return optimus-ide-collabsdk.ChatStreamEvent{}, false, err
	}
	l.state.lastPartSeq = part.Seq
	return optimus-ide-collabsdk.ChatStreamEvent{
		Type:   optimus-ide-collabsdk.ChatStreamEventTypeMessagePart,
		ChatID: l.chatID,
		MessagePart: &optimus-ide-collabsdk.ChatStreamMessagePart{
			Role:              part.Role,
			Part:              part.Part,
			HistoryVersion:    part.HistoryVersion,
			GenerationAttempt: part.GenerationAttempt,
			Seq:               part.Seq,
		},
	}, true, nil
}

func (l *streamLoop) currentRelayTarget() streamRelayTarget {
	return streamRelayTarget{
		workerID:          l.state.workerID,
		historyVersion:    l.state.historyVersion,
		generationAttempt: l.state.generationAttempt,
	}
}

func sameNullUUID(a, b uuid.NullUUID) bool {
	if a.Valid != b.Valid {
		return false
	}
	if !a.Valid {
		return true
	}
	return a.UUID == b.UUID
}

func cloneHeader(header http.Header) http.Header {
	if header == nil {
		return nil
	}
	return header.Clone()
}
