package optimus-ide-collabd

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/xerrors"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbauthz"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpapi"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/workspaceapps/appurl"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

// PrimaryRegion exposes the user facing values of a workspace proxy to
// be used by a user.
func (api *API) PrimaryRegion(ctx context.Context) (optimus-ide-collabsdk.Region, error) {
	deploymentIDStr, err := api.Database.GetDeploymentID(ctx)
	if xerrors.Is(err, sql.ErrNoRows) {
		// This shouldn't happen but it's pretty easy to avoid this causing
		// issues by falling back to a nil UUID.
		deploymentIDStr = uuid.Nil.String()
	} else if err != nil {
		return optimus-ide-collabsdk.Region{}, xerrors.Errorf("get deployment ID: %w", err)
	}
	deploymentID, err := uuid.Parse(deploymentIDStr)
	if err != nil {
		// This also shouldn't happen but we fallback to nil UUID.
		deploymentID = uuid.Nil
	}

	proxy, err := api.Database.GetDefaultProxyConfig(ctx)
	if err != nil {
		return optimus-ide-collabsdk.Region{}, xerrors.Errorf("get default proxy config: %w", err)
	}

	return optimus-ide-collabsdk.Region{
		ID:               deploymentID,
		Name:             "primary",
		DisplayName:      proxy.DisplayName,
		IconURL:          proxy.IconURL,
		Healthy:          true,
		PathAppURL:       api.AccessURL.String(),
		WildcardHostname: appurl.SubdomainAppHost(api.AppHostname, api.AccessURL),
	}, nil
}

// PrimaryWorkspaceProxy returns the primary workspace proxy for the site.
func (api *API) PrimaryWorkspaceProxy(ctx context.Context) (database.WorkspaceProxy, error) {
	region, err := api.PrimaryRegion(ctx)
	if err != nil {
		return database.WorkspaceProxy{}, err
	}

	// The default proxy is an edge case because these values are computed
	// rather then being stored in the database.
	return database.WorkspaceProxy{
		ID:               region.ID,
		Name:             region.Name,
		DisplayName:      region.DisplayName,
		Icon:             region.IconURL,
		Url:              region.PathAppURL,
		WildcardHostname: region.WildcardHostname,
		Deleted:          false,
	}, nil
}

// @Summary Get site-wide regions for workspace connections
// @ID get-site-wide-regions-for-workspace-connections
// @Security Optimus-IDE-CollabSessionToken
// @Produce json
// @Tags WorkspaceProxies
// @Success 200 {object} optimus-ide-collabsdk.RegionsResponse[optimus-ide-collabsdk.Region]
// @Router /api/v2/regions [get]
func (api *API) regions(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	//nolint:gocritic // this route intentionally requests resources that users
	// cannot usually access in order to give them a full list of available
	// regions.
	ctx = dbauthz.AsSystemRestricted(ctx)

	region, err := api.PrimaryRegion(ctx)
	if err != nil {
		httpapi.InternalServerError(rw, err)
		return
	}

	httpapi.Write(ctx, rw, http.StatusOK, optimus-ide-collabsdk.RegionsResponse[optimus-ide-collabsdk.Region]{
		Regions: []optimus-ide-collabsdk.Region{region},
	})
}
