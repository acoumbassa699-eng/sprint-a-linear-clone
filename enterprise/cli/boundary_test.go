package cli_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	boundarycli "github.com/optimus-ide-collab/boundary/cli"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli/clitest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpapi"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/optimus-ide-collabdenttest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/license"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

// Actually testing the functionality of optimus-ide-collab/boundary takes place in the
// optimus-ide-collab/boundary repo, since it's a dependency of optimus-ide-collab.
// Here we want to test basically that integrating it as a subcommand doesn't break anything.
func TestAgentFirewallSubcommand(t *testing.T) {
	t.Parallel()

	inv, _ := newCLI(t, "agent-firewall", "--help")
	var buf bytes.Buffer
	inv.Stdout = &buf
	inv.Stderr = &buf

	err := inv.Run()
	require.NoError(t, err)

	// Verify help output contains expected information.
	// We're simply confirming that `optimus-ide-collab agent-firewall --help` ran without a runtime error as
	// a good chunk of serpent's self validation logic happens at runtime.
	output := buf.String()
	assert.Contains(t, output, boundarycli.BaseCommand("dev").Short)
}

func TestBoundaryAlias(t *testing.T) {
	t.Parallel()

	inv, _ := newCLI(t, "boundary", "--help")
	var buf bytes.Buffer
	inv.Stdout = &buf
	inv.Stderr = &buf

	err := inv.Run()
	require.NoError(t, err)

	// The alias should dispatch to the same command and display help.
	output := buf.String()
	assert.Contains(t, output, boundarycli.BaseCommand("dev").Short)
}

func TestAgentFirewallLicenseVerification(t *testing.T) {
	t.Parallel()

	t.Run("EntitledAndEnabled", func(t *testing.T) {
		t.Parallel()

		client, _ := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureBoundary: 1,
				},
			},
		})

		inv, conf := newCLI(t, "agent-firewall", "--version")
		//nolint:gocritic // requires owner
		clitest.SetupConfig(t, client, conf)

		ctx := testutil.Context(t, testutil.WaitShort)
		err := inv.WithContext(ctx).Run()
		// Should succeed - agent-firewall --version should work with valid license.
		require.NoError(t, err)
	})

	t.Run("NotEntitled", func(t *testing.T) {
		t.Parallel()

		// Create a proxy server that returns entitlements without boundary feature.
		client, _ := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					// No FeatureBoundary
				},
			},
		})

		proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/api/v2/entitlements" {
				res := optimus-ide-collabsdk.Entitlements{
					Features:         map[optimus-ide-collabsdk.FeatureName]optimus-ide-collabsdk.Feature{},
					Warnings:         []string{},
					Errors:           []string{},
					HasLicense:       true,
					Trial:            false,
					RequireTelemetry: false,
				}
				// Set boundary to not entitled, all other features to entitled.
				for _, feature := range optimus-ide-collabsdk.FeatureNames {
					if feature == optimus-ide-collabsdk.FeatureBoundary {
						// Explicitly set boundary to not entitled.
						res.Features[feature] = optimus-ide-collabsdk.Feature{
							Entitlement: optimus-ide-collabsdk.EntitlementNotEntitled,
							Enabled:     false,
						}
					} else {
						res.Features[feature] = optimus-ide-collabsdk.Feature{
							Entitlement: optimus-ide-collabsdk.EntitlementEntitled,
							Enabled:     true,
						}
					}
				}
				httpapi.Write(r.Context(), w, http.StatusOK, res)
				return
			}

			// Otherwise, proxy the request to the real API server.
			rp := httputil.NewSingleHostReverseProxy(client.URL)
			tp := &http.Transport{}
			defer tp.CloseIdleConnections()
			rp.Transport = tp
			rp.ServeHTTP(w, r)
		}))
		defer proxy.Close()

		proxyURL, err := url.Parse(proxy.URL)
		require.NoError(t, err)
		proxyClient := optimus-ide-collabsdk.New(proxyURL, optimus-ide-collabsdk.WithHTTPClient(optimus-ide-collabdtest.NewIsolatedHTTPClient(proxyURL)))
		proxyClient.SetSessionToken(client.SessionToken())
		t.Cleanup(proxyClient.HTTPClient.CloseIdleConnections)

		inv, conf := newCLI(t, "agent-firewall", "--version")
		clitest.SetupConfig(t, proxyClient, conf)

		ctx := testutil.Context(t, testutil.WaitShort)
		err = inv.WithContext(ctx).Run()
		require.Error(t, err)
		require.ErrorContains(t, err, "your license is not entitled to use the agent-firewall feature")
	})

	t.Run("FeatureDisabled", func(t *testing.T) {
		t.Parallel()

		// Create a proxy server that returns entitlements with boundary disabled.
		client, _ := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureBoundary: 1,
				},
			},
		})

		proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/api/v2/entitlements" {
				res := optimus-ide-collabsdk.Entitlements{
					Features:         map[optimus-ide-collabsdk.FeatureName]optimus-ide-collabsdk.Feature{},
					Warnings:         []string{},
					Errors:           []string{},
					HasLicense:       true,
					Trial:            false,
					RequireTelemetry: false,
				}
				for _, feature := range optimus-ide-collabsdk.FeatureNames {
					if feature == optimus-ide-collabsdk.FeatureBoundary {
						// Feature is entitled but disabled.
						res.Features[feature] = optimus-ide-collabsdk.Feature{
							Entitlement: optimus-ide-collabsdk.EntitlementEntitled,
							Enabled:     false,
						}
					} else {
						res.Features[feature] = optimus-ide-collabsdk.Feature{
							Entitlement: optimus-ide-collabsdk.EntitlementEntitled,
							Enabled:     true,
						}
					}
				}
				httpapi.Write(r.Context(), w, http.StatusOK, res)
				return
			}

			// Otherwise, proxy the request to the real API server.
			rp := httputil.NewSingleHostReverseProxy(client.URL)
			tp := &http.Transport{}
			defer tp.CloseIdleConnections()
			rp.Transport = tp
			rp.ServeHTTP(w, r)
		}))
		defer proxy.Close()

		proxyURL, err := url.Parse(proxy.URL)
		require.NoError(t, err)
		proxyClient := optimus-ide-collabsdk.New(proxyURL, optimus-ide-collabsdk.WithHTTPClient(optimus-ide-collabdtest.NewIsolatedHTTPClient(proxyURL)))
		proxyClient.SetSessionToken(client.SessionToken())
		t.Cleanup(proxyClient.HTTPClient.CloseIdleConnections)

		inv, conf := newCLI(t, "agent-firewall", "--version")
		clitest.SetupConfig(t, proxyClient, conf)

		ctx := testutil.Context(t, testutil.WaitShort)
		err = inv.WithContext(ctx).Run()
		require.Error(t, err)
		require.ErrorContains(t, err, "the agent-firewall feature is disabled in your deployment configuration")
	})

	t.Run("AGPLDeployment", func(t *testing.T) {
		t.Parallel()

		// Create an AGPL server (no enterprise features).
		client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{})

		proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/api/v2/entitlements" {
				// AGPL deployments return 404 for entitlements endpoint.
				w.WriteHeader(http.StatusNotFound)
				return
			}

			// Otherwise, proxy the request to the real API server.
			rp := httputil.NewSingleHostReverseProxy(client.URL)
			tp := &http.Transport{}
			defer tp.CloseIdleConnections()
			rp.Transport = tp
			rp.ServeHTTP(w, r)
		}))
		defer proxy.Close()

		proxyURL, err := url.Parse(proxy.URL)
		require.NoError(t, err)
		proxyClient := optimus-ide-collabsdk.New(proxyURL, optimus-ide-collabsdk.WithHTTPClient(optimus-ide-collabdtest.NewIsolatedHTTPClient(proxyURL)))
		proxyClient.SetSessionToken(client.SessionToken())
		t.Cleanup(proxyClient.HTTPClient.CloseIdleConnections)

		inv, conf := newCLI(t, "agent-firewall", "--version")
		clitest.SetupConfig(t, proxyClient, conf)

		ctx := testutil.Context(t, testutil.WaitShort)
		err = inv.WithContext(ctx).Run()
		require.Error(t, err)
		require.ErrorContains(t, err, "your deployment appears to be an AGPL deployment")
	})
}

// TestAgentFirewallChildProcessSkipsCheck verifies that when CHILD=true, the
// license check is skipped. This simulates boundary re-executing itself to run
// the target process. We use a proxy that would fail the license check to
// verify it's skipped.
func TestAgentFirewallChildProcessSkipsCheck(t *testing.T) {
	// Cannot use t.Parallel() with t.Setenv().
	client, _ := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
		LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
			Features: license.Features{
				// No FeatureBoundary - would normally fail
			},
		},
	})

	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v2/entitlements" {
			// Return not entitled for boundary - this would normally cause failure.
			res := optimus-ide-collabsdk.Entitlements{
				Features:         map[optimus-ide-collabsdk.FeatureName]optimus-ide-collabsdk.Feature{},
				Warnings:         []string{},
				Errors:           []string{},
				HasLicense:       true,
				Trial:            false,
				RequireTelemetry: false,
			}
			for _, feature := range optimus-ide-collabsdk.FeatureNames {
				if feature == optimus-ide-collabsdk.FeatureBoundary {
					res.Features[feature] = optimus-ide-collabsdk.Feature{
						Entitlement: optimus-ide-collabsdk.EntitlementNotEntitled,
						Enabled:     false,
					}
				} else {
					res.Features[feature] = optimus-ide-collabsdk.Feature{
						Entitlement: optimus-ide-collabsdk.EntitlementEntitled,
						Enabled:     true,
					}
				}
			}
			httpapi.Write(r.Context(), w, http.StatusOK, res)
			return
		}

		// Otherwise, proxy the request to the real API server.
		rp := httputil.NewSingleHostReverseProxy(client.URL)
		tp := &http.Transport{}
		defer tp.CloseIdleConnections()
		rp.Transport = tp
		rp.ServeHTTP(w, r)
	}))
	defer proxy.Close()

	proxyURL, err := url.Parse(proxy.URL)
	require.NoError(t, err)
	proxyClient := optimus-ide-collabsdk.New(proxyURL, optimus-ide-collabsdk.WithHTTPClient(optimus-ide-collabdtest.NewIsolatedHTTPClient(proxyURL)))
	proxyClient.SetSessionToken(client.SessionToken())
	t.Cleanup(proxyClient.HTTPClient.CloseIdleConnections)

	inv, conf := newCLI(t, "agent-firewall", "--version")
	clitest.SetupConfig(t, proxyClient, conf)

	// Set CHILD=true to simulate boundary re-execution. This should skip the
	// license check, so the command should succeed even though the proxy would
	// return "not entitled".
	t.Setenv("CHILD", "true")

	ctx := testutil.Context(t, testutil.WaitShort)
	err = inv.WithContext(ctx).Run()
	// Should succeed because license check is skipped for child processes.
	require.NoError(t, err)
}
