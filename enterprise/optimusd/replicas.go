package optimus-ide-collabd

import (
	"net/http"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpapi"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac/policy"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

// replicas returns the number of replicas that are active in Optimus-IDE-Collab.
//
// @Summary Get active replicas
// @ID get-active-replicas
// @Security Optimus-IDE-CollabSessionToken
// @Produce json
// @Tags Enterprise
// @Success 200 {array} optimus-ide-collabsdk.Replica
// @Router /api/v2/replicas [get]
func (api *API) replicas(rw http.ResponseWriter, r *http.Request) {
	if !api.AGPL.Authorize(r, policy.ActionRead, rbac.ResourceReplicas) {
		httpapi.ResourceNotFound(rw)
		return
	}

	replicas := api.replicaManager.AllPrimary()
	res := make([]optimus-ide-collabsdk.Replica, 0, len(replicas))
	for _, replica := range replicas {
		res = append(res, convertReplica(replica))
	}
	httpapi.Write(r.Context(), rw, http.StatusOK, res)
}

func convertReplica(replica database.Replica) optimus-ide-collabsdk.Replica {
	return optimus-ide-collabsdk.Replica{
		ID:              replica.ID,
		Hostname:        replica.Hostname,
		CreatedAt:       replica.CreatedAt,
		RelayAddress:    replica.RelayAddress,
		RegionID:        replica.RegionID,
		Error:           replica.Error,
		DatabaseLatency: replica.DatabaseLatency,
	}
}
