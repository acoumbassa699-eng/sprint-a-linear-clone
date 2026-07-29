package chattool

import (
	"context"
	"errors"
	"fmt"

	"charm.land/fantasy"
	"github.com/google/uuid"

	"cdr.dev/slog/v3"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbauthz"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

const workspaceQuotaErrorTitle = "Workspace quota reached"

type buildFailureAction string

const (
	buildFailureActionCreate buildFailureAction = "create"
	buildFailureActionStart  buildFailureAction = "start"
)

type workspaceBuildError struct {
	message string
	code    optimus-ide-collabsdk.JobErrorCode
}

func (e *workspaceBuildError) Error() string {
	return e.message
}

func buildErrorCode(err error) optimus-ide-collabsdk.JobErrorCode {
	var buildErr *workspaceBuildError
	if errors.As(err, &buildErr) {
		return buildErr.code
	}
	return ""
}

// quotaErrorResult is the structured response returned when a workspace
// build fails because the user's workspace quota is exhausted.
type quotaErrorResult struct {
	ErrorCode optimus-ide-collabsdk.JobErrorCode `json:"error_code"`
	// Error is the raw build failure string used for debugging and
	// frontend error detection.
	Error string `json:"error"`
	// Title is a short user-facing summary.
	Title string `json:"title"`
	// Message explains the failure and inlines the recovery guidance
	// the model should relay to the user.
	Message string             `json:"message"`
	BuildID string             `json:"build_id,omitempty"`
	Quota   *quotaErrorDetails `json:"quota,omitempty"`
}

type quotaErrorDetails struct {
	CreditsConsumed int64 `json:"credits_consumed"`
	Budget          int64 `json:"budget"`
}

func newQuotaError(
	msg string,
	buildID uuid.UUID,
	action buildFailureAction,
	quota *quotaErrorDetails,
) quotaErrorResult {
	verb := "create"
	if action == buildFailureActionStart {
		verb = "start"
	}
	message := fmt.Sprintf(
		"Optimus-IDE-Collab could not %s this workspace because your workspace quota is "+
			"full. Delete a workspace you no longer need to free quota, or "+
			"ask an administrator to raise your group quota allowance.",
		verb,
	)

	r := quotaErrorResult{
		ErrorCode: optimus-ide-collabsdk.InsufficientQuota,
		Error:     msg,
		Title:     workspaceQuotaErrorTitle,
		Message:   message,
		Quota:     quota,
	}
	if buildID != uuid.Nil {
		r.BuildID = buildID.String()
	}
	return r
}

func workspaceQuotaDetails(
	ctx context.Context,
	logger slog.Logger,
	db database.Store,
	ownerID uuid.UUID,
	organizationID uuid.UUID,
) *quotaErrorDetails {
	if db == nil || ownerID == uuid.Nil || organizationID == uuid.Nil {
		return nil
	}

	quotaCtx := ctx
	if actor, ok := dbauthz.ActorFromContext(ctx); !ok || actor.ID != ownerID.String() {
		ownerCtx, err := asOwner(ctx, db, ownerID)
		if err != nil {
			logger.Debug(ctx, "failed to load owner authorization for quota lookup",
				slog.F("owner_id", ownerID),
				slog.F("organization_id", organizationID),
				slog.Error(err),
			)
			return nil
		}
		quotaCtx = ownerCtx
	}

	consumed, err := db.GetQuotaConsumedForUser(quotaCtx, database.GetQuotaConsumedForUserParams{
		OwnerID:        ownerID,
		OrganizationID: organizationID,
	})
	if err != nil {
		logger.Debug(ctx, "failed to load consumed workspace quota",
			slog.F("owner_id", ownerID),
			slog.F("organization_id", organizationID),
			slog.Error(err),
		)
		return nil
	}
	budget, err := db.GetQuotaAllowanceForUser(quotaCtx, database.GetQuotaAllowanceForUserParams{
		UserID:         ownerID,
		OrganizationID: organizationID,
	})
	if err != nil {
		logger.Debug(ctx, "failed to load workspace quota allowance",
			slog.F("owner_id", ownerID),
			slog.F("organization_id", organizationID),
			slog.Error(err),
		)
		return nil
	}
	return &quotaErrorDetails{
		CreditsConsumed: consumed,
		Budget:          budget,
	}
}

func quotaErrorToolResponse(
	ctx context.Context,
	logger slog.Logger,
	db database.Store,
	ownerID uuid.UUID,
	organizationID uuid.UUID,
	msg string,
	buildID uuid.UUID,
	action buildFailureAction,
) fantasy.ToolResponse {
	quota := workspaceQuotaDetails(ctx, logger, db, ownerID, organizationID)
	return marshalToolResponse(newQuotaError(msg, buildID, action, quota))
}

// buildFailureToolResponse keeps build failures as JSON carried in a normal
// text tool response. The chatprompt pipeline flattens IsError responses into
// a single string and drops structured fields, so quota and generic build
// failures both keep IsError false and let the frontend detect failures via
// the "error" key.
func buildFailureToolResponse(
	ctx context.Context,
	logger slog.Logger,
	db database.Store,
	ownerID uuid.UUID,
	organizationID uuid.UUID,
	action buildFailureAction,
	buildID uuid.UUID,
	err error,
) fantasy.ToolResponse {
	msg := err.Error()
	if optimus-ide-collabsdk.JobIsInsufficientQuotaErrorCode(buildErrorCode(err)) {
		return quotaErrorToolResponse(
			ctx,
			logger,
			db,
			ownerID,
			organizationID,
			msg,
			buildID,
			action,
		)
	}
	return buildToolResponse(newBuildError(msg, buildID))
}
