package enterprise_test

import (
	"net"
	"testing"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbtestutil"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpmw"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/workspaceapps/apptest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/optimus-ide-collabdenttest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/license"
	"github.com/optimus-ide-collab/serpent"
)

func TestWorkspaceApps(t *testing.T) {
	t.Parallel()

	apptest.Run(t, true, func(t *testing.T, opts *apptest.DeploymentOptions) *apptest.Deployment {
		deploymentValues := optimus-ide-collabdtest.DeploymentValues(t)
		deploymentValues.DisablePathApps = serpent.Bool(opts.DisablePathApps)
		deploymentValues.Dangerous.AllowPathAppSharing = serpent.Bool(opts.DangerousAllowPathAppSharing)
		deploymentValues.Dangerous.AllowPathAppSiteOwnerAccess = serpent.Bool(opts.DangerousAllowPathAppSiteOwnerAccess)
		deploymentValues.Experiments = []string{
			"*",
		}

		if opts.DisableSubdomainApps {
			opts.AppHost = ""
		}

		flushStatsCollectorCh := make(chan chan<- struct{}, 1)
		opts.StatsCollectorOptions.Flush = flushStatsCollectorCh
		flushStats := func() {
			flushStatsCollectorDone := make(chan struct{}, 1)
			flushStatsCollectorCh <- flushStatsCollectorDone
			<-flushStatsCollectorDone
		}

		db, pubsub := dbtestutil.NewDB(t)

		client, _, _, user := optimus-ide-collabdenttest.NewWithAPI(t, &optimus-ide-collabdenttest.Options{
			Options: &optimus-ide-collabdtest.Options{
				DeploymentValues:         deploymentValues,
				AppHostname:              opts.AppHost,
				IncludeProvisionerDaemon: true,
				RealIPConfig: &httpmw.RealIPConfig{
					TrustedOrigins: []*net.IPNet{{
						IP:   net.ParseIP("127.0.0.1"),
						Mask: net.CIDRMask(8, 32),
					}},
					TrustedHeaders: []string{
						"CF-Connecting-IP",
					},
				},
				WorkspaceAppsStatsCollectorOptions: opts.StatsCollectorOptions,
				Database:                           db,
				Pubsub:                             pubsub,
			},
			LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
				Features: license.Features{
					optimus-ide-collabsdk.FeatureMultipleOrganizations: 1,
				},
			},
		})

		return &apptest.Deployment{
			Options:        opts,
			SDKClient:      client,
			FirstUser:      user,
			PathAppBaseURL: client.URL,
			FlushStats:     flushStats,
		}
	})
}
