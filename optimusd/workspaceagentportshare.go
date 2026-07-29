package optimus-ide-collabd

import (
	"database/sql"
	"errors"
	"net/http"
	"slices"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpapi"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpmw"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

// @Summary Upsert workspace agent port share
// @ID upsert-workspace-agent-port-share
// @Security Optimus-IDE-CollabSessionToken
// @Accept json
// @Produce json
// @Tags PortSharing
// @Param workspace path string true "Workspace ID" format(uuid)
// @Param request body optimus-ide-collabsdk.UpsertWorkspaceAgentPortShareRequest true "Upsert port sharing level request"
// @Success 200 {object} optimus-ide-collabsdk.WorkspaceAgentPortShare
// @Router /api/v2/workspaces/{workspace}/port-share [post]
func (api *API) postWorkspaceAgentPortShare(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	workspace := httpmw.WorkspaceParam(r)
	portSharer := *api.PortSharer.Load()
	var req optimus-ide-collabsdk.UpsertWorkspaceAgentPortShareRequest
	if !httpapi.Read(ctx, rw, r, &req) {
		return
	}

	if !req.ShareLevel.ValidPortShareLevel() {
		httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message: "Port sharing level not allowed.",
			Validations: []optimus-ide-collabsdk.ValidationError{
				{
					Field:  "share_level",
					Detail: "Port sharing level not allowed.",
				},
			},
		})
		return
	}

	if req.Port < 9 || req.Port > 65535 {
		httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message: "Port must be between 9 and 65535.",
			Validations: []optimus-ide-collabsdk.ValidationError{
				{
					Field:  "port",
					Detail: "Port must be between 9 and 65535.",
				},
			},
		})
		return
	}
	if !req.Protocol.ValidPortProtocol() {
		httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message: "Port protocol not allowed.",
		})
		return
	}

	template, err := api.Database.GetTemplateByID(ctx, workspace.TemplateID)
	if err != nil {
		httpapi.InternalServerError(rw, err)
		return
	}

	err = portSharer.AuthorizedLevel(template, req.ShareLevel)
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message: err.Error(),
		})
		return
	}

	agents, err := api.Database.GetWorkspaceAgentsInLatestBuildByWorkspaceID(ctx, workspace.ID)
	if err != nil {
		httpapi.InternalServerError(rw, err)
		return
	}

	found := false
	for _, agent := range agents {
		if agent.Name == req.AgentName {
			found = true
			break
		}
	}
	if !found {
		httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message: "Agent not found.",
		})
		return
	}

	psl, err := api.Database.UpsertWorkspaceAgentPortShare(ctx, database.UpsertWorkspaceAgentPortShareParams{
		WorkspaceID: workspace.ID,
		AgentName:   req.AgentName,
		Port:        req.Port,
		ShareLevel:  database.AppSharingLevel(req.ShareLevel),
		Protocol:    database.PortShareProtocol(req.Protocol),
	})
	if err != nil {
		httpapi.InternalServerError(rw, err)
		return
	}

	httpapi.Write(ctx, rw, http.StatusOK, convertPortShare(psl))
}

// @Summary Get workspace agent port shares
// @ID get-workspace-agent-port-shares
// @Security Optimus-IDE-CollabSessionToken
// @Produce json
// @Tags PortSharing
// @Param workspace path string true "Workspace ID" format(uuid)
// @Success 200 {object} optimus-ide-collabsdk.WorkspaceAgentPortShares
// @Router /api/v2/workspaces/{workspace}/port-share [get]
func (api *API) workspaceAgentPortShares(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	workspace := httpmw.WorkspaceParam(r)

	shares, err := api.Database.ListWorkspaceAgentPortShares(ctx, workspace.ID)
	if err != nil {
		httpapi.InternalServerError(rw, err)
		return
	}

	httpapi.Write(ctx, rw, http.StatusOK, optimus-ide-collabsdk.WorkspaceAgentPortShares{
		Shares: convertPortShares(shares),
	})
}

// @Summary Delete workspace agent port share
// @ID delete-workspace-agent-port-share
// @Security Optimus-IDE-CollabSessionToken
// @Accept json
// @Tags PortSharing
// @Param workspace path string true "Workspace ID" format(uuid)
// @Param request body optimus-ide-collabsdk.DeleteWorkspaceAgentPortShareRequest true "Delete port sharing level request"
// @Success 200
// @Router /api/v2/workspaces/{workspace}/port-share [delete]
func (api *API) deleteWorkspaceAgentPortShare(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	workspace := httpmw.WorkspaceParam(r)
	var req optimus-ide-collabsdk.DeleteWorkspaceAgentPortShareRequest
	if !httpapi.Read(ctx, rw, r, &req) {
		return
	}

	_, err := api.Database.GetWorkspaceAgentPortShare(ctx, database.GetWorkspaceAgentPortShareParams{
		WorkspaceID: workspace.ID,
		AgentName:   req.AgentName,
		Port:        req.Port,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httpapi.Write(ctx, rw, http.StatusNotFound, optimus-ide-collabsdk.Response{
				Message: "Port share not found.",
			})
			return
		}

		httpapi.InternalServerError(rw, err)
		return
	}

	err = api.Database.DeleteWorkspaceAgentPortShare(ctx, database.DeleteWorkspaceAgentPortShareParams{
		WorkspaceID: workspace.ID,
		AgentName:   req.AgentName,
		Port:        req.Port,
	})
	if err != nil {
		httpapi.InternalServerError(rw, err)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func convertPortShares(shares []database.WorkspaceAgentPortShare) []optimus-ide-collabsdk.WorkspaceAgentPortShare {
	converted := []optimus-ide-collabsdk.WorkspaceAgentPortShare{}
	for _, share := range shares {
		converted = append(converted, convertPortShare(share))
	}
	slices.SortFunc(converted, func(i, j optimus-ide-collabsdk.WorkspaceAgentPortShare) int {
		return (int)(i.Port - j.Port)
	})
	return converted
}

func convertPortShare(share database.WorkspaceAgentPortShare) optimus-ide-collabsdk.WorkspaceAgentPortShare {
	return optimus-ide-collabsdk.WorkspaceAgentPortShare{
		WorkspaceID: share.WorkspaceID,
		AgentName:   share.AgentName,
		Port:        share.Port,
		ShareLevel:  optimus-ide-collabsdk.WorkspaceAgentPortShareLevel(share.ShareLevel),
		Protocol:    optimus-ide-collabsdk.WorkspaceAgentPortShareProtocol(share.Protocol),
	}
}
