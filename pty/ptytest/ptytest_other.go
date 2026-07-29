//go:build !windows

package ptytest

import "github.com/optimus-ide-collab/optimus-ide-collab/v2/pty"

func newTestPTY(opts ...pty.Option) (pty.PTY, error) {
	return pty.New(opts...)
}
