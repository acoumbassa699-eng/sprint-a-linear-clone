package agentapi

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/workspaceapps/appurl"
)

func Test_vscodeProxyURI(t *testing.T) {
	t.Parallel()

	optimus-ide-collabAccessURL, err := url.Parse("https://optimus-ide-collab.com")
	require.NoError(t, err)

	accessURLWithPort, err := url.Parse("https://optimus-ide-collab.com:8080")
	require.NoError(t, err)

	basicApp := appurl.ApplicationURL{
		Prefix:        "prefix",
		AppSlugOrPort: "slug",
		AgentName:     "agent",
		WorkspaceName: "workspace",
		Username:      "user",
	}

	cases := []struct {
		Name        string
		App         appurl.ApplicationURL
		AccessURL   *url.URL
		AppHostname string
		Expected    string
	}{
		{
			Name:        "NoHostname",
			AccessURL:   optimus-ide-collabAccessURL,
			AppHostname: "",
			App:         basicApp,
			Expected:    "",
		},
		{
			Name:        "NoHostnameAccessURLPort",
			AccessURL:   accessURLWithPort,
			AppHostname: "",
			App:         basicApp,
			Expected:    "",
		},
		{
			Name:        "Hostname",
			AccessURL:   optimus-ide-collabAccessURL,
			AppHostname: "*.apps.optimus-ide-collab.com",
			App:         basicApp,
			Expected:    fmt.Sprintf("https://%s.apps.optimus-ide-collab.com", basicApp.String()),
		},
		{
			Name:        "HostnameWithAccessURLPort",
			AccessURL:   accessURLWithPort,
			AppHostname: "*.apps.optimus-ide-collab.com",
			App:         basicApp,
			Expected:    fmt.Sprintf("https://%s.apps.optimus-ide-collab.com:%s", basicApp.String(), accessURLWithPort.Port()),
		},
		{
			Name:        "HostnameWithPort",
			AccessURL:   optimus-ide-collabAccessURL,
			AppHostname: "*.apps.optimus-ide-collab.com:4444",
			App:         basicApp,
			Expected:    fmt.Sprintf("https://%s.apps.optimus-ide-collab.com:%s", basicApp.String(), "4444"),
		},
		{
			// Port from hostname takes precedence over access url port.
			Name:        "HostnameWithPortAccessURLWithPort",
			AccessURL:   accessURLWithPort,
			AppHostname: "*.apps.optimus-ide-collab.com:4444",
			App:         basicApp,
			Expected:    fmt.Sprintf("https://%s.apps.optimus-ide-collab.com:%s", basicApp.String(), "4444"),
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			t.Parallel()

			require.NotNilf(t, c.AccessURL, "AccessURL is required")

			output := vscodeProxyURI(c.App, c.AccessURL, c.AppHostname)
			require.Equal(t, c.Expected, output)
		})
	}
}
