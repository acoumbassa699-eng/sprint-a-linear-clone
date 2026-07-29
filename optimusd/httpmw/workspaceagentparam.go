package httpmw

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"

	"cdr.dev/slog/v3"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpapi"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpmw/loggermw"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

type workspaceAgentAndWorkspaceParamContextKey struct{}

// WorkspaceAgentAndWorkspaceParam returns the workspace agent and its associated workspace from the ExtractWorkspaceAgentParam handler.
func WorkspaceAgentAndWorkspaceParam(r *http.Request) database.GetWorkspaceAgentAndWorkspaceByIDRow {
	aw, ok := r.Context().Value(workspaceAgentAndWorkspaceParamContextKey{}).(database.GetWorkspaceAgentAndWorkspaceByIDRow)
	if !ok {
		panic("developer error: agent middleware not provided")
	}
	return aw
}

// ExtractWorkspaceAgentAndWorkspaceParam grabs a workspace agent and its associated workspace from the "workspaceagent" URL parameter.
func ExtractWorkspaceAgentAndWorkspaceParam(db database.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			agentUUID, parsed := ParseUUIDParam(rw, r, "workspaceagent")
			if !parsed {
				return
			}

			agentWithWorkspace, err := db.GetWorkspaceAgentAndWorkspaceByID(ctx, agentUUID)
			if httpapi.Is404Error(err) {
				httpapi.Write(ctx, rw, http.StatusNotFound, optimus-ide-collabsdk.Response{
					Message: "Agent doesn't exist with that id, or you do not have access to it.",
				})
				return
			}

			ctx = context.WithValue(ctx, workspaceAgentAndWorkspaceParamContextKey{}, agentWithWorkspace)
			chi.RouteContext(ctx).URLParams.Add("workspace", agentWithWorkspace.WorkspaceTable.ID.String())

			if rlogger := loggermw.RequestLoggerFromContext(ctx); rlogger != nil {
				rlogger.WithFields(
					slog.F("workspace_name", agentWithWorkspace.WorkspaceTable.Name),
					slog.F("agent_name", agentWithWorkspace.WorkspaceAgent.Name),
				)
			}

			next.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}
