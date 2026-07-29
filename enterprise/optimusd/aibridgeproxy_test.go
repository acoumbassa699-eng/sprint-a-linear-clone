package optimus-ide-collabd_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/optimus-ide-collabdenttest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/license"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
	"github.com/optimus-ide-collab/serpent"
)

func TestAIBridgeProxyCertificateRetrieval(t *testing.T) {
	t.Parallel()

	t.Run("DisabledReturns404", func(t *testing.T) {
		t.Parallel()

		dv := optimus-ide-collabdtest.DeploymentValues(t)
		dv.AI.BridgeConfig.Enabled = serpent.Bool(true)
		// Proxy is disabled by default, so we don't need to set it explicitly.
		client, _ := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			Options: &optimus-ide-collabdtest.Options{
				DeploymentValues: dv,
			},
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureAIBridge: 1,
				},
			},
		})

		ctx := testutil.Context(t, testutil.WaitLong)

		// Make a request to the proxy CA cert endpoint.
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, client.URL.String()+"/api/v2/ai-gateway/proxy/ca-cert.pem", nil)
		require.NoError(t, err)
		req.Header.Set(optimus-ide-collabsdk.SessionTokenHeader, client.SessionToken())

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("RequiresLicenseFeature", func(t *testing.T) {
		t.Parallel()

		dv := optimus-ide-collabdtest.DeploymentValues(t)
		dv.AI.BridgeConfig.Enabled = serpent.Bool(true)
		client, _ := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			Options: &optimus-ide-collabdtest.Options{
				DeploymentValues: dv,
			},
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				// No aibridge feature.
				Features: license.Features{},
			},
		})

		ctx := testutil.Context(t, testutil.WaitLong)

		// Make a request to the proxy CA cert endpoint.
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, client.URL.String()+"/api/v2/ai-gateway/proxy/ca-cert.pem", nil)
		require.NoError(t, err)
		req.Header.Set(optimus-ide-collabsdk.SessionTokenHeader, client.SessionToken())

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("RequiresAuthentication", func(t *testing.T) {
		t.Parallel()

		dv := optimus-ide-collabdtest.DeploymentValues(t)
		dv.AI.BridgeConfig.Enabled = serpent.Bool(true)
		client, _ := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			Options: &optimus-ide-collabdtest.Options{
				DeploymentValues: dv,
			},
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureAIBridge: 1,
				},
			},
		})

		ctx := testutil.Context(t, testutil.WaitLong)

		// Make a request to the proxy CA cert endpoint without authentication.
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, client.URL.String()+"/api/v2/ai-gateway/proxy/ca-cert.pem", nil)
		require.NoError(t, err)

		// No session token header set.
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}
