package optimus-ide-collabd

import (
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/db2sdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpapi"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpmw"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/provisionerdserver"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/provisionerkey"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

// @Summary Create provisioner key
// @ID create-provisioner-key
// @Security Optimus-IDE-CollabSessionToken
// @Produce json
// @Tags Enterprise
// @Param organization path string true "Organization ID"
// @Success 201 {object} optimus-ide-collabsdk.CreateProvisionerKeyResponse
// @Router /api/v2/organizations/{organization}/provisionerkeys [post]
func (api *API) postProvisionerKey(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	organization := httpmw.OrganizationParam(r)

	var req optimus-ide-collabsdk.CreateProvisionerKeyRequest
	if !httpapi.Read(ctx, rw, r, &req) {
		return
	}

	if req.Name == "" {
		httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message: "Name is required",
			Validations: []optimus-ide-collabsdk.ValidationError{
				{
					Field:  "name",
					Detail: "Name is required",
				},
			},
		})
		return
	}

	if len(req.Name) > 64 {
		httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message: "Name must be at most 64 characters",
			Validations: []optimus-ide-collabsdk.ValidationError{
				{
					Field:  "name",
					Detail: "Name must be at most 64 characters",
				},
			},
		})
		return
	}

	if slices.ContainsFunc(optimus-ide-collabsdk.ReservedProvisionerKeyNames(), func(s string) bool {
		return strings.EqualFold(req.Name, s)
	}) {
		httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message: fmt.Sprintf("Name cannot be reserved name '%s'", req.Name),
			Validations: []optimus-ide-collabsdk.ValidationError{
				{
					Field:  "name",
					Detail: fmt.Sprintf("Name cannot be reserved name '%s'", req.Name),
				},
			},
		})
		return
	}

	params, token, err := provisionerkey.New(organization.ID, req.Name, req.Tags)
	if err != nil {
		httpapi.InternalServerError(rw, err)
		return
	}

	_, err = api.Database.InsertProvisionerKey(ctx, params)
	if database.IsUniqueViolation(err, database.UniqueProvisionerKeysOrganizationIDNameIndex) {
		httpapi.Write(ctx, rw, http.StatusConflict, optimus-ide-collabsdk.Response{
			Message: fmt.Sprintf("Provisioner key with name '%s' already exists in organization", req.Name),
		})
		return
	}
	if err != nil {
		httpapi.InternalServerError(rw, err)
		return
	}

	httpapi.Write(ctx, rw, http.StatusCreated, optimus-ide-collabsdk.CreateProvisionerKeyResponse{
		Key: token,
	})
}

// @Summary List provisioner key
// @ID list-provisioner-key
// @Security Optimus-IDE-CollabSessionToken
// @Produce json
// @Tags Enterprise
// @Param organization path string true "Organization ID"
// @Success 200 {object} []optimus-ide-collabsdk.ProvisionerKey
// @Router /api/v2/organizations/{organization}/provisionerkeys [get]
func (api *API) provisionerKeys(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	organization := httpmw.OrganizationParam(r)

	pks, err := api.Database.ListProvisionerKeysByOrganizationExcludeReserved(ctx, organization.ID)
	if err != nil {
		httpapi.InternalServerError(rw, err)
		return
	}

	httpapi.Write(ctx, rw, http.StatusOK, convertProvisionerKeys(pks))
}

// @Summary List provisioner key daemons
// @ID list-provisioner-key-daemons
// @Security Optimus-IDE-CollabSessionToken
// @Produce json
// @Tags Enterprise
// @Param organization path string true "Organization ID"
// @Success 200 {object} []optimus-ide-collabsdk.ProvisionerKeyDaemons
// @Router /api/v2/organizations/{organization}/provisionerkeys/daemons [get]
func (api *API) provisionerKeyDaemons(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	organization := httpmw.OrganizationParam(r)

	pks, err := api.Database.ListProvisionerKeysByOrganization(ctx, organization.ID)
	if err != nil {
		httpapi.InternalServerError(rw, err)
		return
	}
	sdkKeys := convertProvisionerKeys(pks)

	// For the default organization, we insert three rows for the special
	// provisioner key types (built-in, user-auth, and psk). We _don't_ insert
	// those into the database for any other org, but we still need to include the
	// user-auth key in this list, so we just insert it manually.
	if !slices.ContainsFunc(sdkKeys, func(key optimus-ide-collabsdk.ProvisionerKey) bool {
		return key.ID == optimus-ide-collabsdk.ProvisionerKeyUUIDUserAuth
	}) {
		sdkKeys = append(sdkKeys, optimus-ide-collabsdk.ProvisionerKey{
			ID:   optimus-ide-collabsdk.ProvisionerKeyUUIDUserAuth,
			Name: optimus-ide-collabsdk.ProvisionerKeyNameUserAuth,
			Tags: map[string]string{},
		})
	}

	daemons, err := api.Database.GetProvisionerDaemonsByOrganization(ctx, database.GetProvisionerDaemonsByOrganizationParams{OrganizationID: organization.ID})
	if err != nil {
		httpapi.InternalServerError(rw, err)
		return
	}
	// provisionerdserver.DefaultHeartbeatInterval*3 matches the healthcheck report staleInterval.
	recentDaemons := db2sdk.RecentProvisionerDaemons(time.Now(), provisionerdserver.DefaultHeartbeatInterval*3, daemons)

	pkDaemons := []optimus-ide-collabsdk.ProvisionerKeyDaemons{}
	for _, key := range sdkKeys {
		// The key.OrganizationID for the `user-auth` key is hardcoded to
		// the default org in the database and we are overwriting it here
		// to be the correct org we used to query the list.
		// This will be changed when we update the `user-auth` keys to be
		// directly tied to a user ID.
		if key.ID.String() == optimus-ide-collabsdk.ProvisionerKeyIDUserAuth {
			key.OrganizationID = organization.ID
		}
		daemons := []optimus-ide-collabsdk.ProvisionerDaemon{}
		for _, daemon := range recentDaemons {
			if daemon.KeyID == key.ID {
				daemons = append(daemons, daemon)
			}
		}
		pkDaemons = append(pkDaemons, optimus-ide-collabsdk.ProvisionerKeyDaemons{
			Key:     key,
			Daemons: daemons,
		})
	}

	httpapi.Write(ctx, rw, http.StatusOK, pkDaemons)
}

// @Summary Delete provisioner key
// @ID delete-provisioner-key
// @Security Optimus-IDE-CollabSessionToken
// @Tags Enterprise
// @Param organization path string true "Organization ID"
// @Param provisionerkey path string true "Provisioner key name"
// @Success 204
// @Router /api/v2/organizations/{organization}/provisionerkeys/{provisionerkey} [delete]
func (api *API) deleteProvisionerKey(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	provisionerKey := httpmw.ProvisionerKeyParam(r)

	if provisionerKey.ID.String() == optimus-ide-collabsdk.ProvisionerKeyIDBuiltIn ||
		provisionerKey.ID.String() == optimus-ide-collabsdk.ProvisionerKeyIDUserAuth ||
		provisionerKey.ID.String() == optimus-ide-collabsdk.ProvisionerKeyIDPSK {
		httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message: fmt.Sprintf("Cannot delete reserved '%s' provisioner key", provisionerKey.Name),
		})
		return
	}

	err := api.Database.DeleteProvisionerKey(ctx, provisionerKey.ID)
	if err != nil {
		httpapi.InternalServerError(rw, err)
		return
	}

	httpapi.Write(ctx, rw, http.StatusNoContent, nil)
}

// @Summary Fetch provisioner key details
// @ID fetch-provisioner-key-details
// @Security Optimus-IDE-CollabProvisionerKey
// @Produce json
// @Tags Enterprise
// @Param provisionerkey path string true "Provisioner Key"
// @Success 200 {object} optimus-ide-collabsdk.ProvisionerKey
// @Router /api/v2/provisionerkeys/{provisionerkey} [get]
func (*API) fetchProvisionerKey(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	pk, ok := httpmw.ProvisionerKeyAuthOptional(r)
	// extra check but this one should never happen as it is covered by the auth middleware
	if !ok {
		httpapi.Write(ctx, rw, http.StatusForbidden, optimus-ide-collabsdk.Response{
			Message: fmt.Sprintf("unable to auth: please provide the %s header", optimus-ide-collabsdk.ProvisionerDaemonKey),
		})
		return
	}

	httpapi.Write(ctx, rw, http.StatusOK, convertProvisionerKey(pk))
}

func convertProvisionerKey(dbKey database.ProvisionerKey) optimus-ide-collabsdk.ProvisionerKey {
	return optimus-ide-collabsdk.ProvisionerKey{
		ID:             dbKey.ID,
		CreatedAt:      dbKey.CreatedAt,
		OrganizationID: dbKey.OrganizationID,
		Name:           dbKey.Name,
		Tags:           optimus-ide-collabsdk.ProvisionerKeyTags(dbKey.Tags),
		// HashedSecret - never include the access token in the API response
	}
}

func convertProvisionerKeys(dbKeys []database.ProvisionerKey) []optimus-ide-collabsdk.ProvisionerKey {
	keys := make([]optimus-ide-collabsdk.ProvisionerKey, 0, len(dbKeys))
	for _, dbKey := range dbKeys {
		keys = append(keys, convertProvisionerKey(dbKey))
	}

	slices.SortFunc(keys, func(key1, key2 optimus-ide-collabsdk.ProvisionerKey) int {
		return key1.CreatedAt.Compare(key2.CreatedAt)
	})

	return keys
}
