package optimus-ide-collabd

import (
	"net/http"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/db2sdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpapi"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpmw"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/util/slice"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

// @Summary Get organizations
// @ID get-organizations
// @Security Optimus-IDE-CollabSessionToken
// @Produce json
// @Tags Organizations
// @Success 200 {object} []optimus-ide-collabsdk.Organization
// @Router /api/v2/organizations [get]
func (api *API) organizations(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	organizations, err := api.Database.GetOrganizations(ctx, database.GetOrganizationsParams{})
	if httpapi.Is404Error(err) {
		httpapi.ResourceNotFound(rw)
		return
	}
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Internal error fetching organizations.",
			Detail:  err.Error(),
		})
		return
	}

	httpapi.Write(ctx, rw, http.StatusOK, slice.List(organizations, db2sdk.Organization))
}

// @Summary Get organization by ID
// @ID get-organization-by-id
// @Security Optimus-IDE-CollabSessionToken
// @Produce json
// @Tags Organizations
// @Param organization path string true "Organization ID" format(uuid)
// @Success 200 {object} optimus-ide-collabsdk.Organization
// @Router /api/v2/organizations/{organization} [get]
func (*API) organization(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	organization := httpmw.OrganizationParam(r)

	httpapi.Write(ctx, rw, http.StatusOK, db2sdk.Organization(organization))
}
