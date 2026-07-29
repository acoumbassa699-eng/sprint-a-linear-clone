package enidpsync

import (
	"context"

	"github.com/golang-jwt/jwt/v4"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/idpsync"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

func (e EnterpriseIDPSync) OrganizationSyncEntitled() bool {
	return e.entitlements.Enabled(optimus-ide-collabsdk.FeatureMultipleOrganizations)
}

func (e EnterpriseIDPSync) OrganizationSyncEnabled(ctx context.Context, db database.Store) bool {
	if !e.OrganizationSyncEntitled() {
		return false
	}

	// If this logic is ever updated, make sure to update the corresponding
	// checkIDPOrgSync in optimus-ide-collabd/telemetry/telemetry.go.
	settings, err := e.OrganizationSyncSettings(ctx, db)
	if err == nil && settings.Field != "" {
		return true
	}
	return false
}

func (e EnterpriseIDPSync) ParseOrganizationClaims(ctx context.Context, mergedClaims jwt.MapClaims) (idpsync.OrganizationParams, *idpsync.HTTPError) {
	if !e.OrganizationSyncEntitled() {
		// Default to agpl if multi-org is not enabled
		return e.AGPLIDPSync.ParseOrganizationClaims(ctx, mergedClaims)
	}

	return idpsync.OrganizationParams{
		// Return true if entitled
		SyncEntitled: true,
		MergedClaims: mergedClaims,
	}, nil
}
