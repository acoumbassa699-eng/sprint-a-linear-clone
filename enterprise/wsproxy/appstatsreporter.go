package wsproxy

import (
	"context"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/workspaceapps"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/wsproxy/wsproxysdk"
)

var _ workspaceapps.StatsReporter = (*appStatsReporter)(nil)

type appStatsReporter struct {
	Client *wsproxysdk.Client
}

func (r *appStatsReporter) ReportAppStats(ctx context.Context, stats []workspaceapps.StatsReport) error {
	err := r.Client.ReportAppStats(ctx, wsproxysdk.ReportAppStatsRequest{
		Stats: stats,
	})
	return err
}
