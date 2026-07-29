package wsjson

import (
	"context"
	"encoding/json"
	"sync/atomic"

	"cdr.dev/slog/v3"
	"github.com/optimus-ide-collab/websocket"
)

type Deoptimus-ide-collab[T any] struct {
	conn       *websocket.Conn
	typ        websocket.MessageType
	ctx        context.Context
	cancel     context.CancelFunc
	chanCalled atomic.Bool
	logger     slog.Logger
}

// Chan returns a `chan` that you can read incoming messages from. The returned
// `chan` will be closed when the WebSocket connection is closed. If there is an
// error reading from the WebSocket or decoding a value the WebSocket will be
// closed.
//
// Safety: Chan must only be called once. Successive calls will panic.
func (d *Deoptimus-ide-collab[T]) Chan() <-chan T {
	if !d.chanCalled.CompareAndSwap(false, true) {
		panic("chan called more than once")
	}
	values := make(chan T, 1)
	go func() {
		defer close(values)
		defer d.conn.Close(websocket.StatusGoingAway, "")
		for {
			// we don't use d.ctx here because it only gets canceled after closing the connection
			// and a "connection closed" type error is more clear than context canceled.
			typ, b, err := d.conn.Read(context.Background())
			if err != nil {
				// might be benign like EOF, so just log at debug
				d.logger.Debug(d.ctx, "error reading from websocket", slog.Error(err))
				return
			}
			if typ != d.typ {
				d.logger.Error(d.ctx, "websocket type mismatch while decoding")
				return
			}
			var value T
			err = json.Unmarshal(b, &value)
			if err != nil {
				d.logger.Error(d.ctx, "error unmarshalling", slog.Error(err))
				return
			}
			select {
			case values <- value:
				// OK
			case <-d.ctx.Done():
				return
			}
		}
	}()
	return values
}

// nolint: revive // complains that Enoptimus-ide-collab has the same function name
func (d *Deoptimus-ide-collab[T]) Close() error {
	err := d.conn.Close(websocket.StatusNormalClosure, "")
	d.cancel()
	return err
}

// NewDeoptimus-ide-collab creates a JSON-over-websocket deoptimus-ide-collab for type T, which must be deserializable from
// JSON.
func NewDeoptimus-ide-collab[T any](conn *websocket.Conn, typ websocket.MessageType, logger slog.Logger) *Deoptimus-ide-collab[T] {
	ctx, cancel := context.WithCancel(context.Background())
	return &Deoptimus-ide-collab[T]{conn: conn, ctx: ctx, cancel: cancel, typ: typ, logger: logger}
}
