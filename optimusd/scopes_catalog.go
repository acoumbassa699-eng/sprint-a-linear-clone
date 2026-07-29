package optimus-ide-collabd

import (
	"net/http"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpapi"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

// listExternalScopes returns the curated list of API key scopes (resource:action)
// requestable via the API.
//
// @Summary List API key scopes
// @ID list-api-key-scopes
// @Tags Authorization
// @Produce json
// @Success 200 {object} optimus-ide-collabsdk.ExternalAPIKeyScopes
// @Router /api/v2/auth/scopes [get]
func (*API) listExternalScopes(rw http.ResponseWriter, r *http.Request) {
	scopes := rbac.ExternalScopeNames()
	external := make([]optimus-ide-collabsdk.APIKeyScope, 0, len(scopes))
	for _, scope := range scopes {
		external = append(external, optimus-ide-collabsdk.APIKeyScope(scope))
	}

	httpapi.Write(r.Context(), rw, http.StatusOK, optimus-ide-collabsdk.ExternalAPIKeyScopes{
		External: external,
	})
}
