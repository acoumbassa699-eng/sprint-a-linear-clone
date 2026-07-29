package optimus-ide-collabdenttest_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/optimus-ide-collabdenttest"
)

func TestEnterpriseEndpointsDocumented(t *testing.T) {
	t.Parallel()

	swaggerComments, err := optimus-ide-collabdtest.ParseSwaggerComments(
		"..", "../../../optimus-ide-collabd", "../../../optimus-ide-collabd/workspaceconnwatcher")
	require.NoError(t, err, "can't parse swagger comments")
	require.NotEmpty(t, swaggerComments, "swagger comments must be present")

	//nolint: dogsled
	_, _, api, _ := optimus-ide-collabdenttest.NewWithAPI(t, nil)
	optimus-ide-collabdtest.VerifySwaggerDefinitions(t, api.AGPL.APIHandler, swaggerComments, optimus-ide-collabdtest.WithSwaggerRoutePrefix("/api/v2"))
}
