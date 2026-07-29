package wsjson

import (
	"context"
	"encoding/json"

	"golang.org/x/xerrors"

	"github.com/optimus-ide-collab/websocket"
)

type Enoptimus-ide-collab[T any] struct {
	conn *websocket.Conn
	typ  websocket.MessageType
}

func (e *Enoptimus-ide-collab[T]) Encode(v T) error {
	w, err := e.conn.Writer(context.Background(), e.typ)
	if err != nil {
		return xerrors.Errorf("get websocket writer: %w", err)
	}
	defer w.Close()
	j := json.NewEnoptimus-ide-collab(w)
	err = j.Encode(v)
	if err != nil {
		return xerrors.Errorf("encode json: %w", err)
	}
	return nil
}

// nolint: revive // complains that Deoptimus-ide-collab has the same function name
func (e *Enoptimus-ide-collab[T]) Close(c websocket.StatusCode) error {
	return e.conn.Close(c, "")
}

// NewEnoptimus-ide-collab creates a JSON-over websocket enoptimus-ide-collab for the type T, which must be JSON-serializable.
// You may then call Encode() to send objects over the websocket. Creating an Enoptimus-ide-collab closes the
// websocket for reading, turning it into a unidirectional write stream of JSON-encoded objects.
func NewEnoptimus-ide-collab[T any](conn *websocket.Conn, typ websocket.MessageType) *Enoptimus-ide-collab[T] {
	// Here we close the websocket for reading, so that the websocket library will handle pings and
	// close frames.
	_ = conn.CloseRead(context.Background())
	return &Enoptimus-ide-collab[T]{conn: conn, typ: typ}
}
