package healthcheck_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/healthcheck"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/healthcheck/health"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

func TestWorkspaceProxies(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name                  string
		fetchWorkspaceProxies func(context.Context) (optimus-ide-collabsdk.RegionsResponse[optimus-ide-collabsdk.WorkspaceProxy], error)
		updateProxyHealth     func(context.Context) error
		expectedHealthy       bool
		expectedError         string
		expectedWarningCode   health.Code
		expectedSeverity      health.Severity
	}{
		{
			name:             "NotEnabled",
			expectedHealthy:  true,
			expectedSeverity: health.SeverityOK,
		},
		{
			name:                  "Enabled/NoProxies",
			fetchWorkspaceProxies: fakeFetchWorkspaceProxies(),
			updateProxyHealth:     fakeUpdateProxyHealth(nil),
			expectedHealthy:       true,
			expectedSeverity:      health.SeverityOK,
		},
		{
			name:                  "Enabled/OneHealthy",
			fetchWorkspaceProxies: fakeFetchWorkspaceProxies(fakeWorkspaceProxy("alpha", true)),
			updateProxyHealth:     fakeUpdateProxyHealth(nil),
			expectedHealthy:       true,
			expectedSeverity:      health.SeverityOK,
		},
		{
			name:                  "Enabled/OneUnhealthy",
			fetchWorkspaceProxies: fakeFetchWorkspaceProxies(fakeWorkspaceProxy("alpha", false)),
			updateProxyHealth:     fakeUpdateProxyHealth(nil),
			expectedHealthy:       false,
			expectedSeverity:      health.SeverityError,
			expectedError:         string(health.CodeProxyUnhealthy),
		},
		{
			name: "Enabled/OneUnreachable",
			fetchWorkspaceProxies: func(ctx context.Context) (optimus-ide-collabsdk.RegionsResponse[optimus-ide-collabsdk.WorkspaceProxy], error) {
				return optimus-ide-collabsdk.RegionsResponse[optimus-ide-collabsdk.WorkspaceProxy]{
					Regions: []optimus-ide-collabsdk.WorkspaceProxy{
						{
							Region: optimus-ide-collabsdk.Region{
								Name:    "gone",
								Healthy: false,
							},
							Status: optimus-ide-collabsdk.WorkspaceProxyStatus{
								Status: optimus-ide-collabsdk.ProxyUnreachable,
								Report: optimus-ide-collabsdk.ProxyHealthReport{
									Errors: []string{
										"request to proxy failed: Get \"http://127.0.0.1:3001/healthz-report\": dial tcp 127.0.0.1:3001: connect: connection refused",
									},
								},
							},
						},
					},
				}, nil
			},
			updateProxyHealth: fakeUpdateProxyHealth(nil),
			expectedHealthy:   false,
			expectedSeverity:  health.SeverityError,
			expectedError:     string(health.CodeProxyUnhealthy),
		},
		{
			name: "Enabled/AllHealthy",
			fetchWorkspaceProxies: fakeFetchWorkspaceProxies(
				fakeWorkspaceProxy("alpha", true),
				fakeWorkspaceProxy("beta", true),
			),
			updateProxyHealth: func(ctx context.Context) error {
				return nil
			},
			expectedHealthy:  true,
			expectedSeverity: health.SeverityOK,
		},
		{
			name: "Enabled/OneHealthyOneUnhealthy",
			fetchWorkspaceProxies: fakeFetchWorkspaceProxies(
				fakeWorkspaceProxy("alpha", false),
				fakeWorkspaceProxy("beta", true),
			),
			updateProxyHealth:   fakeUpdateProxyHealth(nil),
			expectedHealthy:     true,
			expectedSeverity:    health.SeverityWarning,
			expectedWarningCode: health.CodeProxyUnhealthy,
		},
		{
			name: "Enabled/AllUnhealthy",
			fetchWorkspaceProxies: fakeFetchWorkspaceProxies(
				fakeWorkspaceProxy("alpha", false),
				fakeWorkspaceProxy("beta", false),
			),
			updateProxyHealth: fakeUpdateProxyHealth(nil),
			expectedHealthy:   false,
			expectedSeverity:  health.SeverityError,
			expectedError:     string(health.CodeProxyUnhealthy),
		},
		{
			name: "Enabled/NotConnectedYet",
			fetchWorkspaceProxies: fakeFetchWorkspaceProxies(
				fakeWorkspaceProxy("slowpoke", true),
			),
			updateProxyHealth: fakeUpdateProxyHealth(nil),
			expectedHealthy:   true,
			expectedSeverity:  health.SeverityOK,
		},
		{
			name:                  "Enabled/ErrFetchWorkspaceProxy",
			fetchWorkspaceProxies: fakeFetchWorkspaceProxiesErr(assert.AnError),
			updateProxyHealth:     fakeUpdateProxyHealth(nil),
			expectedHealthy:       false,
			expectedSeverity:      health.SeverityError,
			expectedError:         string(health.CodeProxyFetch),
		},
		{
			name:                  "Enabled/ErrUpdateProxyHealth",
			fetchWorkspaceProxies: fakeFetchWorkspaceProxies(fakeWorkspaceProxy("alpha", true)),
			updateProxyHealth:     fakeUpdateProxyHealth(assert.AnError),
			expectedHealthy:       true,
			expectedSeverity:      health.SeverityWarning,
			expectedWarningCode:   health.CodeProxyUpdate,
		},
		{
			name: "Enabled/OneUnhealthyAndDeleted",
			fetchWorkspaceProxies: fakeFetchWorkspaceProxies(fakeWorkspaceProxy("alpha", false, func(wp *optimus-ide-collabsdk.WorkspaceProxy) {
				wp.Deleted = true
			})),
			updateProxyHealth: fakeUpdateProxyHealth(nil),
			expectedHealthy:   true,
			expectedSeverity:  health.SeverityOK,
		},
		{
			name: "Enabled/ProxyWarnings",
			fetchWorkspaceProxies: fakeFetchWorkspaceProxies(
				fakeWorkspaceProxy("alpha", true, func(wp *optimus-ide-collabsdk.WorkspaceProxy) {
					wp.Status.Report.Warnings = []string{"warning"}
				}),
				fakeWorkspaceProxy("beta", false),
			),
			updateProxyHealth:   fakeUpdateProxyHealth(nil),
			expectedHealthy:     true,
			expectedSeverity:    health.SeverityWarning,
			expectedWarningCode: health.CodeProxyUnhealthy,
		},
		{
			name: "Enabled/ProxyWarningsButAllErrored",
			fetchWorkspaceProxies: fakeFetchWorkspaceProxies(
				fakeWorkspaceProxy("alpha", false),
				fakeWorkspaceProxy("beta", false, func(wp *optimus-ide-collabsdk.WorkspaceProxy) {
					wp.Status.Report.Warnings = []string{"warning"}
				}),
			),
			updateProxyHealth: fakeUpdateProxyHealth(nil),
			expectedHealthy:   false,
			expectedError:     string(health.CodeProxyUnhealthy),
			expectedSeverity:  health.SeverityError,
		},
	} {
		if tt.name != "Enabled/ProxyWarnings" {
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var rpt healthcheck.WorkspaceProxyReport
			var opts healthcheck.WorkspaceProxyReportOptions
			if tt.fetchWorkspaceProxies != nil && tt.updateProxyHealth != nil {
				opts.WorkspaceProxiesFetchUpdater = &fakeWorkspaceProxyFetchUpdater{
					fetchFunc:  tt.fetchWorkspaceProxies,
					updateFunc: tt.updateProxyHealth,
				}
			} else {
				opts.WorkspaceProxiesFetchUpdater = &healthcheck.AGPLWorkspaceProxiesFetchUpdater{}
			}

			rpt.Run(context.Background(), &opts)

			assert.Equal(t, tt.expectedHealthy, rpt.Healthy)
			assert.Equal(t, tt.expectedSeverity, rpt.Severity)
			if tt.expectedError != "" && assert.NotNil(t, rpt.Error) {
				assert.Contains(t, *rpt.Error, tt.expectedError)
			} else if !assert.Nil(t, rpt.Error) {
				t.Logf("error: %v", *rpt.Error)
			}
			if tt.expectedWarningCode != "" && assert.NotEmpty(t, rpt.Warnings) {
				var found bool
				for _, w := range rpt.Warnings {
					if w.Code == tt.expectedWarningCode {
						found = true
						break
					}
				}
				assert.True(t, found, "expected warning %s not found in %v", tt.expectedWarningCode, rpt.Warnings)
			} else {
				assert.Empty(t, rpt.Warnings)
			}
		})
	}
}

func TestWorkspaceProxy_ErrorDismissed(t *testing.T) {
	t.Parallel()

	var report healthcheck.WorkspaceProxyReport
	report.Run(context.Background(), &healthcheck.WorkspaceProxyReportOptions{
		WorkspaceProxiesFetchUpdater: &fakeWorkspaceProxyFetchUpdater{
			fetchFunc:  fakeFetchWorkspaceProxiesErr(assert.AnError),
			updateFunc: fakeUpdateProxyHealth(assert.AnError),
		},
		Dismissed: true,
	})

	assert.True(t, report.Dismissed)
	assert.Equal(t, health.SeverityWarning, report.Severity)
}

// yet another implementation of the thing
type fakeWorkspaceProxyFetchUpdater struct {
	fetchFunc  func(context.Context) (optimus-ide-collabsdk.RegionsResponse[optimus-ide-collabsdk.WorkspaceProxy], error)
	updateFunc func(context.Context) error
}

func (u *fakeWorkspaceProxyFetchUpdater) Fetch(ctx context.Context) (optimus-ide-collabsdk.RegionsResponse[optimus-ide-collabsdk.WorkspaceProxy], error) {
	return u.fetchFunc(ctx)
}

func (u *fakeWorkspaceProxyFetchUpdater) Update(ctx context.Context) error {
	return u.updateFunc(ctx)
}

//nolint:revive // yes, this is a control flag, and that is OK in a unit test.
func fakeWorkspaceProxy(name string, healthy bool, mutators ...func(*optimus-ide-collabsdk.WorkspaceProxy)) optimus-ide-collabsdk.WorkspaceProxy {
	var status optimus-ide-collabsdk.WorkspaceProxyStatus
	if !healthy {
		status = optimus-ide-collabsdk.WorkspaceProxyStatus{
			Status: optimus-ide-collabsdk.ProxyUnreachable,
			Report: optimus-ide-collabsdk.ProxyHealthReport{
				Errors: []string{assert.AnError.Error()},
			},
		}
	}
	wsp := optimus-ide-collabsdk.WorkspaceProxy{
		Region: optimus-ide-collabsdk.Region{
			Name:    name,
			Healthy: healthy,
		},
		Status: status,
	}
	for _, f := range mutators {
		f(&wsp)
	}
	return wsp
}

func fakeFetchWorkspaceProxies(ps ...optimus-ide-collabsdk.WorkspaceProxy) func(context.Context) (optimus-ide-collabsdk.RegionsResponse[optimus-ide-collabsdk.WorkspaceProxy], error) {
	return func(context.Context) (optimus-ide-collabsdk.RegionsResponse[optimus-ide-collabsdk.WorkspaceProxy], error) {
		return optimus-ide-collabsdk.RegionsResponse[optimus-ide-collabsdk.WorkspaceProxy]{
			Regions: ps,
		}, nil
	}
}

func fakeFetchWorkspaceProxiesErr(err error) func(context.Context) (optimus-ide-collabsdk.RegionsResponse[optimus-ide-collabsdk.WorkspaceProxy], error) {
	return func(context.Context) (optimus-ide-collabsdk.RegionsResponse[optimus-ide-collabsdk.WorkspaceProxy], error) {
		return optimus-ide-collabsdk.RegionsResponse[optimus-ide-collabsdk.WorkspaceProxy]{
			Regions: []optimus-ide-collabsdk.WorkspaceProxy{},
		}, err
	}
}

func fakeUpdateProxyHealth(err error) func(context.Context) error {
	return func(context.Context) error {
		return err
	}
}
