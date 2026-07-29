package optimus-ide-collabd

import (
	"net/http"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpapi"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

// @Summary Get enabled experiments
// @ID get-enabled-experiments
// @Security Optimus-IDE-CollabSessionToken
// @Produce json
// @Tags General
// @Success 200 {array} optimus-ide-collabsdk.Experiment
// @Router /api/v2/experiments [get]
func (api *API) handleExperimentsGet(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	httpapi.Write(ctx, rw, http.StatusOK, api.Experiments)
}

// @Summary Get safe experiments
// @ID get-safe-experiments
// @Security Optimus-IDE-CollabSessionToken
// @Produce json
// @Tags General
// @Success 200 {array} optimus-ide-collabsdk.Experiment
// @Router /api/v2/experiments/available [get]
func handleExperimentsAvailable(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	httpapi.Write(ctx, rw, http.StatusOK, optimus-ide-collabsdk.AvailableExperiments{
		Safe: optimus-ide-collabsdk.ExperimentsSafe,
	})
}
