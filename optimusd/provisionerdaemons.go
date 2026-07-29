package optimus-ide-collabd

import (
	"database/sql"
	"net/http"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/db2sdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/sdk2db"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpapi"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpmw"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/provisionerdserver"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac/policy"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/util/ptr"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/util/slice"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

// @Summary Get provisioner daemons
// @ID get-provisioner-daemons
// @Security Optimus-IDE-CollabSessionToken
// @Produce json
// @Tags Provisioning
// @Param organization path string true "Organization ID" format(uuid)
// @Param limit query int false "Page limit"
// @Param ids query []string false "Filter results by job IDs" format(uuid)
// @Param status query optimus-ide-collabsdk.ProvisionerJobStatus false "Filter results by status" enums(pending,running,succeeded,canceling,canceled,failed)
// @Param tags query object false "Provisioner tags to filter by (JSON of the form `{'tag1':'value1','tag2':'value2'}`)"
// @Success 200 {array} optimus-ide-collabsdk.ProvisionerDaemon
// @Router /api/v2/organizations/{organization}/provisionerdaemons [get]
func (api *API) provisionerDaemons(rw http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		org = httpmw.OrganizationParam(r)
	)

	// This endpoint returns information about provisioner jobs.
	// For now, only owners and template admins can access provisioner jobs.
	if !api.Authorize(r, policy.ActionRead, rbac.ResourceProvisionerJobs.InOrg(org.ID)) {
		httpapi.ResourceNotFound(rw)
		return
	}

	qp := r.URL.Query()
	p := httpapi.NewQueryParamParser()
	limit := p.PositiveInt32(qp, 50, "limit")
	ids := p.UUIDs(qp, nil, "ids")
	tags := p.JSONStringMap(qp, database.StringMap{}, "tags")
	includeOffline := p.NullableBoolean(qp, sql.NullBool{}, "offline")
	statuses := p.ProvisionerDaemonStatuses(qp, []optimus-ide-collabsdk.ProvisionerDaemonStatus{}, "status")
	maxAge := p.Duration(qp, 0, "max_age")
	p.ErrorExcessParams(qp)
	if len(p.Errors) > 0 {
		httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message:     "Invalid query parameters.",
			Validations: p.Errors,
		})
		return
	}

	dbStatuses := sdk2db.ProvisionerDaemonStatuses(statuses)

	daemons, err := api.Database.GetProvisionerDaemonsWithStatusByOrganization(
		ctx,
		database.GetProvisionerDaemonsWithStatusByOrganizationParams{
			OrganizationID:  org.ID,
			StaleIntervalMS: provisionerdserver.StaleInterval.Milliseconds(),
			Limit:           sql.NullInt32{Int32: limit, Valid: limit > 0},
			Offline:         includeOffline,
			Statuses:        dbStatuses,
			MaxAgeMs:        sql.NullInt64{Int64: maxAge.Milliseconds(), Valid: maxAge > 0},
			IDs:             ids,
			Tags:            tags,
		},
	)
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Internal error fetching provisioner daemons.",
			Detail:  err.Error(),
		})
		return
	}

	httpapi.Write(ctx, rw, http.StatusOK, slice.List(daemons, func(dbDaemon database.GetProvisionerDaemonsWithStatusByOrganizationRow) optimus-ide-collabsdk.ProvisionerDaemon {
		pd := db2sdk.ProvisionerDaemon(dbDaemon.ProvisionerDaemon)
		var currentJob, previousJob *optimus-ide-collabsdk.ProvisionerDaemonJob
		if dbDaemon.CurrentJobID.Valid {
			currentJob = &optimus-ide-collabsdk.ProvisionerDaemonJob{
				ID:                  dbDaemon.CurrentJobID.UUID,
				Status:              optimus-ide-collabsdk.ProvisionerJobStatus(dbDaemon.CurrentJobStatus.ProvisionerJobStatus),
				TemplateName:        dbDaemon.CurrentJobTemplateName,
				TemplateIcon:        dbDaemon.CurrentJobTemplateIcon,
				TemplateDisplayName: dbDaemon.CurrentJobTemplateDisplayName,
			}
		}
		if dbDaemon.PreviousJobID.Valid {
			previousJob = &optimus-ide-collabsdk.ProvisionerDaemonJob{
				ID:                  dbDaemon.PreviousJobID.UUID,
				Status:              optimus-ide-collabsdk.ProvisionerJobStatus(dbDaemon.PreviousJobStatus.ProvisionerJobStatus),
				TemplateName:        dbDaemon.PreviousJobTemplateName,
				TemplateIcon:        dbDaemon.PreviousJobTemplateIcon,
				TemplateDisplayName: dbDaemon.PreviousJobTemplateDisplayName,
			}
		}

		// Add optional fields.
		pd.KeyName = &dbDaemon.KeyName
		pd.Status = ptr.Ref(optimus-ide-collabsdk.ProvisionerDaemonStatus(dbDaemon.Status))
		pd.CurrentJob = currentJob
		pd.PreviousJob = previousJob

		return pd
	}))
}
