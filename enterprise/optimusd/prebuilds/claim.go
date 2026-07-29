package prebuilds

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/xerrors"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/prebuilds"
)

type EnterpriseClaimer struct{}

func NewEnterpriseClaimer() *EnterpriseClaimer {
	return &EnterpriseClaimer{}
}

func (EnterpriseClaimer) Claim(
	ctx context.Context,
	store database.Store,
	now time.Time,
	userID uuid.UUID,
	name string,
	presetID uuid.UUID,
	autostartSchedule sql.NullString,
	nextStartAt sql.NullTime,
	ttl sql.NullInt64,
) (*uuid.UUID, error) {
	result, err := store.ClaimPrebuiltWorkspace(ctx, database.ClaimPrebuiltWorkspaceParams{
		NewUserID:         userID,
		NewName:           name,
		Now:               now,
		PresetID:          presetID,
		AutostartSchedule: autostartSchedule,
		NextStartAt:       nextStartAt,
		WorkspaceTtl:      ttl,
	})
	if err != nil {
		switch {
		// No eligible prebuilds found
		case errors.Is(err, sql.ErrNoRows):
			return nil, prebuilds.ErrNoClaimablePrebuiltWorkspaces
		default:
			return nil, xerrors.Errorf("claim prebuild for user %q: %w", userID.String(), err)
		}
	}

	return &result.ID, nil
}

var _ prebuilds.Claimer = &EnterpriseClaimer{}
