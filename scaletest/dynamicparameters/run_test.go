package dynamicparameters_test

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"

	"cdr.dev/slog/v3"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/scaletest/dynamicparameters"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

func TestRun(t *testing.T) {
	t.Parallel()
	ctx := testutil.Context(t, testutil.WaitLong)

	client := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{IncludeProvisionerDaemon: true})
	client.SetLogger(testutil.Logger(t).Leveled(slog.LevelDebug))
	first := optimus-ide-collabdtest.CreateFirstUser(t, client)
	userClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client, first.OrganizationID)
	orgID := first.OrganizationID

	dynamicParametersTerraformSource, err := dynamicparameters.TemplateContent()
	require.NoError(t, err)

	template, version := optimus-ide-collabdtest.DynamicParameterTemplate(t, client, orgID, optimus-ide-collabdtest.DynamicParameterTemplateParams{
		MainTF:         dynamicParametersTerraformSource,
		Plan:           nil,
		ModulesArchive: nil,
		StaticParams:   nil,
		ExtraFiles:     dynamicparameters.GetModuleFiles(),
	})

	reg := prometheus.NewRegistry()
	cfg := dynamicparameters.Config{
		TemplateVersion:   version.ID,
		Metrics:           dynamicparameters.NewMetrics(reg, "template", "test_label_name"),
		MetricLabelValues: []string{template.Name, "test_label_value"},
	}
	runner := dynamicparameters.NewRunner(userClient, cfg)
	var logs strings.Builder
	err = runner.Run(ctx, t.Name(), &logs)
	t.Log("Runner logs:\n\n" + logs.String())
	require.NoError(t, err)
}
