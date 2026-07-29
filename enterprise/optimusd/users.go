package optimus-ide-collabd

import (
	"net/http"
	"time"

	"golang.org/x/xerrors"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/audit"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpapi"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpmw"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/schedule"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

const TimeFormatHHMM = "15:04"

func (api *API) autostopRequirementEnabledMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		feature, ok := api.Entitlements.Feature(optimus-ide-collabsdk.FeatureAdvancedTemplateScheduling)
		if !ok || !feature.Entitlement.Entitled() {
			httpapi.Write(r.Context(), rw, http.StatusForbidden, optimus-ide-collabsdk.Response{
				Message: "Advanced template scheduling (and user quiet hours schedule) is an Enterprise feature. Contact sales!",
			})
			return
		}
		if !feature.Enabled {
			httpapi.Write(r.Context(), rw, http.StatusForbidden, optimus-ide-collabsdk.Response{
				Message: "Advanced template scheduling (and user quiet hours schedule) is not enabled.",
			})
			return
		}

		next.ServeHTTP(rw, r)
	})
}

// @Summary Get user quiet hours schedule
// @ID get-user-quiet-hours-schedule
// @Security Optimus-IDE-CollabSessionToken
// @Produce json
// @Tags Enterprise
// @Param user path string true "User ID" format(uuid)
// @Success 200 {array} optimus-ide-collabsdk.UserQuietHoursScheduleResponse
// @Router /api/v2/users/{user}/quiet-hours [get]
func (api *API) userQuietHoursSchedule(rw http.ResponseWriter, r *http.Request) {
	var (
		ctx  = r.Context()
		user = httpmw.UserParam(r)
	)

	opts, err := (*api.UserQuietHoursScheduleStore.Load()).Get(ctx, api.Database, user.ID)
	if err != nil {
		httpapi.InternalServerError(rw, err)
		return
	}
	if opts.Schedule == nil {
		httpapi.ResourceNotFound(rw)
		return
	}

	httpapi.Write(ctx, rw, http.StatusOK, optimus-ide-collabsdk.UserQuietHoursScheduleResponse{
		RawSchedule: opts.Schedule.String(),
		UserSet:     opts.UserSet,
		UserCanSet:  opts.UserCanSet,
		Time:        opts.Schedule.TimeParsed().Format(TimeFormatHHMM),
		Timezone:    opts.Schedule.Location().String(),
		Next:        opts.Schedule.Next(time.Now().In(opts.Schedule.Location())),
	})
}

// @Summary Update user quiet hours schedule
// @ID update-user-quiet-hours-schedule
// @Security Optimus-IDE-CollabSessionToken
// @Accept json
// @Produce json
// @Tags Enterprise
// @Param user path string true "User ID" format(uuid)
// @Param request body optimus-ide-collabsdk.UpdateUserQuietHoursScheduleRequest true "Update schedule request"
// @Success 200 {array} optimus-ide-collabsdk.UserQuietHoursScheduleResponse
// @Router /api/v2/users/{user}/quiet-hours [put]
func (api *API) putUserQuietHoursSchedule(rw http.ResponseWriter, r *http.Request) {
	var (
		ctx               = r.Context()
		user              = httpmw.UserParam(r)
		params            optimus-ide-collabsdk.UpdateUserQuietHoursScheduleRequest
		aReq, commitAudit = audit.InitRequest[database.User](rw, &audit.RequestParams{
			Audit:   api.Auditor,
			Log:     api.Logger,
			Request: r,
			Action:  database.AuditActionWrite,
		})
	)
	defer commitAudit()
	aReq.Old = user

	if !httpapi.Read(ctx, rw, r, &params) {
		return
	}

	opts, err := (*api.UserQuietHoursScheduleStore.Load()).Set(ctx, api.Database, user.ID, params.Schedule)
	if xerrors.Is(err, schedule.ErrUserCannotSetQuietHoursSchedule) {
		httpapi.Write(ctx, rw, http.StatusForbidden, optimus-ide-collabsdk.Response{
			Message: "Users cannot set custom quiet hours schedule due to deployment configuration.",
		})
		return
	} else if err != nil {
		// TODO(@dean): some of these errors are related to bad syntax, so it
		// would be nice to 400 instead
		httpapi.InternalServerError(rw, err)
		return
	}

	httpapi.Write(ctx, rw, http.StatusOK, optimus-ide-collabsdk.UserQuietHoursScheduleResponse{
		RawSchedule: opts.Schedule.String(),
		UserSet:     opts.UserSet,
		UserCanSet:  opts.UserCanSet,
		Time:        opts.Schedule.TimeParsed().Format(TimeFormatHHMM),
		Timezone:    opts.Schedule.Location().String(),
		Next:        opts.Schedule.Next(time.Now().In(opts.Schedule.Location())),
	})
}
