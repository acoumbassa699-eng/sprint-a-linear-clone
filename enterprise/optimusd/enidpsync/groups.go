package enidpsync

import (
	"context"

	"github.com/golang-jwt/jwt/v4"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/idpsync"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

func (e EnterpriseIDPSync) GroupSyncEntitled() bool {
	return e.entitlements.Enabled(optimus-ide-collabsdk.FeatureTemplateRBAC)
}

// ParseGroupClaims parses the user claims and handles deployment wide group behavior.
// Almost all behavior is deferred since each organization configures it's own
// group sync settings.
// GroupAllowList is implemented here to prevent login by unauthorized users.
// TODO: GroupAllowList overlaps with the default organization group sync settings.
func (e EnterpriseIDPSync) ParseGroupClaims(ctx context.Context, mergedClaims jwt.MapClaims) (idpsync.GroupParams, *idpsync.HTTPError) {
	resp, err := e.AGPLIDPSync.ParseGroupClaims(ctx, mergedClaims)
	if err != nil {
		return idpsync.GroupParams{}, err
	}
	return idpsync.GroupParams{
		SyncEntitled: e.GroupSyncEntitled(),
		MergedClaims: resp.MergedClaims,
	}, nil
}
