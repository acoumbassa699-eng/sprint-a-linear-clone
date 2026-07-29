package cliutil_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli/cliutil"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

func TestWarnMatchedProvisioners(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name   string
		mp     *optimus-ide-collabsdk.MatchedProvisioners
		job    optimus-ide-collabsdk.ProvisionerJob
		expect string
	}{
		{
			name: "no_match",
			mp: &optimus-ide-collabsdk.MatchedProvisioners{
				Count:     0,
				Available: 0,
			},
			job: optimus-ide-collabsdk.ProvisionerJob{
				Status: optimus-ide-collabsdk.ProvisionerJobPending,
			},
			expect: `there are no provisioners that accept the required tags`,
		},
		{
			name: "no_available",
			mp: &optimus-ide-collabsdk.MatchedProvisioners{
				Count:     1,
				Available: 0,
			},
			job: optimus-ide-collabsdk.ProvisionerJob{
				Status: optimus-ide-collabsdk.ProvisionerJobPending,
			},
			expect: `Provisioners that accept the required tags have not responded for longer than expected`,
		},
		{
			name: "match",
			mp: &optimus-ide-collabsdk.MatchedProvisioners{
				Count:     1,
				Available: 1,
			},
			job: optimus-ide-collabsdk.ProvisionerJob{
				Status: optimus-ide-collabsdk.ProvisionerJobPending,
			},
		},
		{
			name: "not_pending",
			mp:   &optimus-ide-collabsdk.MatchedProvisioners{},
			job: optimus-ide-collabsdk.ProvisionerJob{
				Status: optimus-ide-collabsdk.ProvisionerJobRunning,
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var w strings.Builder
			cliutil.WarnMatchedProvisioners(&w, tt.mp, tt.job)
			if tt.expect != "" {
				require.Contains(t, w.String(), tt.expect)
			} else {
				require.Empty(t, w.String())
			}
		})
	}
}
