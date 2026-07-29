package optimus-ide-collabd

import (
	"net/http"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpapi"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

// @Summary API root handler
// @ID api-root-handler
// @Produce json
// @Tags General
// @Success 200 {object} optimus-ide-collabsdk.Response
// @Router /api/v2/ [get]
func apiRoot(w http.ResponseWriter, r *http.Request) {
	httpapi.Write(r.Context(), w, http.StatusOK, optimus-ide-collabsdk.Response{
		//nolint:gocritic
		Message: "👋",
	})
}
