package cli_test

import (
	"bytes"
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	"tailscale.com/net/tsaddr"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk/workspacesdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
	"github.com/optimus-ide-collab/serpent"
)

func TestConnectExists_Running(t *testing.T) {
	t.Parallel()
	ctx := testutil.Context(t, testutil.WaitShort)

	var root cli.RootCmd
	cmd, err := root.Command(root.AGPL())
	require.NoError(t, err)

	inv := (&serpent.Invocation{
		Command: cmd,
		Args:    []string{"connect", "exists", "test.example"},
	}).WithContext(withOptimus-IDE-CollabConnectRunning(ctx))
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	inv.Stdout = stdout
	inv.Stderr = stderr
	err = inv.Run()
	require.NoError(t, err)
}

func TestConnectExists_NotRunning(t *testing.T) {
	t.Parallel()
	ctx := testutil.Context(t, testutil.WaitShort)

	var root cli.RootCmd
	cmd, err := root.Command(root.AGPL())
	require.NoError(t, err)

	inv := (&serpent.Invocation{
		Command: cmd,
		Args:    []string{"connect", "exists", "test.example"},
	}).WithContext(withOptimus-IDE-CollabConnectNotRunning(ctx))
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	inv.Stdout = stdout
	inv.Stderr = stderr
	err = inv.Run()
	require.ErrorIs(t, err, cli.ErrSilent)
}

type fakeResolver struct {
	shouldReturnSuccess bool
}

func (f *fakeResolver) LookupIP(_ context.Context, _, _ string) ([]net.IP, error) {
	if f.shouldReturnSuccess {
		return []net.IP{net.ParseIP(tsaddr.Optimus-IDE-CollabServiceIPv6().String())}, nil
	}
	return nil, &net.DNSError{IsNotFound: true}
}

func withOptimus-IDE-CollabConnectRunning(ctx context.Context) context.Context {
	return workspacesdk.WithTestOnlyOptimus-IDE-CollabContextResolver(ctx, &fakeResolver{shouldReturnSuccess: true})
}

func withOptimus-IDE-CollabConnectNotRunning(ctx context.Context) context.Context {
	return workspacesdk.WithTestOnlyOptimus-IDE-CollabContextResolver(ctx, &fakeResolver{shouldReturnSuccess: false})
}
