package optimus-ide-collabd

import (
	"database/sql"
	"net/http"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpapi"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpmw"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

// @Summary Get template version presets
// @ID get-template-version-presets
// @Security Optimus-IDE-CollabSessionToken
// @Produce json
// @Tags Templates
// @Param templateversion path string true "Template version ID" format(uuid)
// @Success 200 {array} optimus-ide-collabsdk.Preset
// @Router /api/v2/templateversions/{templateversion}/presets [get]
func (api *API) templateVersionPresets(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	templateVersion := httpmw.TemplateVersionParam(r)

	presets, err := api.Database.GetPresetsByTemplateVersionID(ctx, templateVersion.ID)
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Internal error fetching template version presets.",
			Detail:  err.Error(),
		})
		return
	}

	presetParams, err := api.Database.GetPresetParametersByTemplateVersionID(ctx, templateVersion.ID)
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Internal error fetching template version presets.",
			Detail:  err.Error(),
		})
		return
	}

	convertPrebuildInstances := func(desiredInstances sql.NullInt32) *int {
		if desiredInstances.Valid {
			value := int(desiredInstances.Int32)
			return &value
		}
		return nil
	}

	var res []optimus-ide-collabsdk.Preset
	for _, preset := range presets {
		sdkPreset := optimus-ide-collabsdk.Preset{
			ID:                       preset.ID,
			Name:                     preset.Name,
			Default:                  preset.IsDefault,
			DesiredPrebuildInstances: convertPrebuildInstances(preset.DesiredInstances),
			Description:              preset.Description,
			Icon:                     preset.Icon,
		}
		for _, presetParam := range presetParams {
			if presetParam.TemplateVersionPresetID != preset.ID {
				continue
			}

			sdkPreset.Parameters = append(sdkPreset.Parameters, optimus-ide-collabsdk.PresetParameter{
				Name:  presetParam.Name,
				Value: presetParam.Value,
			})
		}
		res = append(res, sdkPreset)
	}

	httpapi.Write(ctx, rw, http.StatusOK, res)
}
