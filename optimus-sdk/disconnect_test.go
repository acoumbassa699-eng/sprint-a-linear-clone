package optimus-ide-collabsdk_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

func TestDisconnectReason_Valid(t *testing.T) {
	t.Parallel()

	cases := []struct {
		reason optimus-ide-collabsdk.DisconnectReason
		valid  bool
	}{
		{optimus-ide-collabsdk.DisconnectReasonUnknown, true},
		{optimus-ide-collabsdk.DisconnectReasonGraceful, true},
		{optimus-ide-collabsdk.DisconnectReasonClientClosed, true},
		{optimus-ide-collabsdk.DisconnectReasonServerShutdown, true},
		{optimus-ide-collabsdk.DisconnectReasonNetworkError, true},
		{optimus-ide-collabsdk.DisconnectReasonProtocolError, true},
		{optimus-ide-collabsdk.DisconnectReasonWorkspaceStopped, true},
		{optimus-ide-collabsdk.DisconnectReasonControlPlaneLost, true},
		{optimus-ide-collabsdk.DisconnectReason("not_a_real_reason"), false},
	}

	for _, c := range cases {
		require.Equal(t, c.valid, c.reason.Valid(), "reason=%q", c.reason)
	}
}

func TestDisconnectReason_Expected(t *testing.T) {
	t.Parallel()

	expected := map[optimus-ide-collabsdk.DisconnectReason]bool{
		optimus-ide-collabsdk.DisconnectReasonGraceful:         true,
		optimus-ide-collabsdk.DisconnectReasonClientClosed:     true,
		optimus-ide-collabsdk.DisconnectReasonServerShutdown:   true,
		optimus-ide-collabsdk.DisconnectReasonWorkspaceStopped: true,

		optimus-ide-collabsdk.DisconnectReasonUnknown:          false,
		optimus-ide-collabsdk.DisconnectReasonNetworkError:     false,
		optimus-ide-collabsdk.DisconnectReasonProtocolError:    false,
		optimus-ide-collabsdk.DisconnectReasonControlPlaneLost: false,
	}

	for reason, want := range expected {
		require.Equal(t, want, reason.Expected(), "reason=%q", reason)
	}

	// Unknown values default to not-expected so that uncategorized
	// emit sites surface in the "investigate" bucket.
	require.False(t, optimus-ide-collabsdk.DisconnectReason("not_a_real_reason").Expected())
}

func TestDisconnectInitiator_Valid(t *testing.T) {
	t.Parallel()

	cases := []struct {
		initiator optimus-ide-collabsdk.DisconnectInitiator
		valid     bool
	}{
		{optimus-ide-collabsdk.DisconnectInitiatorUnknown, true},
		{optimus-ide-collabsdk.DisconnectInitiatorClient, true},
		{optimus-ide-collabsdk.DisconnectInitiatorAgent, true},
		{optimus-ide-collabsdk.DisconnectInitiatorServer, true},
		{optimus-ide-collabsdk.DisconnectInitiatorNetwork, true},
		{optimus-ide-collabsdk.DisconnectInitiator("nobody"), false},
	}

	for _, c := range cases {
		require.Equal(t, c.valid, c.initiator.Valid(), "initiator=%q", c.initiator)
	}
}

func TestConnectionMethod_Valid(t *testing.T) {
	t.Parallel()

	cases := []struct {
		method optimus-ide-collabsdk.ConnectionMethod
		valid  bool
	}{
		{optimus-ide-collabsdk.ConnectionMethodUnknown, true},
		{optimus-ide-collabsdk.ConnectionMethodDirect, true},
		{optimus-ide-collabsdk.ConnectionMethodDERP, true},
		{optimus-ide-collabsdk.ConnectionMethod("magic"), false},
	}

	for _, c := range cases {
		require.Equal(t, c.valid, c.method.Valid(), "method=%q", c.method)
	}
}
