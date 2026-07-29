package optimus-ide-collabd

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/xerrors"

	"cdr.dev/slog/v3"
	agpl "github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/audit"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/db2sdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbauthz"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpapi"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpmw"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac/acl"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac/policy"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/searchquery"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/util/slice"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

// @Summary Get template available acl users/groups
// @ID get-template-available-acl-usersgroups
// @Security Optimus-IDE-CollabSessionToken
// @Produce json
// @Tags Enterprise
// @Param template path string true "Template ID" format(uuid)
// @Success 200 {array} optimus-ide-collabsdk.ACLAvailable
// @Router /api/v2/templates/{template}/acl/available [get]
func (api *API) templateAvailablePermissions(rw http.ResponseWriter, r *http.Request) {
	var (
		ctx      = r.Context()
		template = httpmw.TemplateParam(r)
	)

	// Requires update permission on the template to list all avail users/groups
	// for assignment.
	if !api.Authorize(r, policy.ActionUpdate, template) {
		httpapi.ResourceNotFound(rw)
		return
	}

	// We have to use the system restricted context here because the caller
	// might not have permission to read all users.
	// nolint:gocritic
	users, _, ok := api.AGPL.GetUsers(rw, r.WithContext(dbauthz.AsSystemRestricted(ctx)))
	if !ok {
		return
	}

	// Apply the same q/limit semantics to groups as the users half of this response.
	// The query semantics are defined for the users, which is awkward. But we can
	// just reuse the search part of the query which is a fuzzy match.
	userFilter, verr := searchquery.Users(r.URL.Query().Get("q"))
	if len(verr) > 0 {
		httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message:     "Invalid user search query.",
			Validations: verr,
		})
		return
	}
	groupPagination, ok := agpl.ParsePagination(rw, r)
	if !ok {
		return
	}

	// Perm check is the template update check.
	// nolint:gocritic
	groups, err := api.Database.GetGroups(dbauthz.AsSystemRestricted(ctx), database.GetGroupsParams{
		OrganizationID: template.OrganizationID,
		Search:         userFilter.Search,
		// #nosec G115 - Pagination limits are small and fit in int32
		LimitOpt: int32(groupPagination.Limit),
	})
	if err != nil {
		httpapi.InternalServerError(rw, err)
		return
	}

	// Fetch member counts for all groups in a single query to avoid an
	// N+1 lookup pattern that was making this endpoint extremely slow on
	// deployments with many groups. The per-group member lists are
	// intentionally not populated here: callers of this endpoint only
	// surface total_member_count (see Group.TotalMemberCount, which is
	// already documented as the canonical value).
	groupIDs := make([]uuid.UUID, len(groups))
	for i, g := range groups {
		groupIDs[i] = g.Group.ID
	}

	// nolint:gocritic // Same justification as the GetGroups call above.
	countRows, err := api.Database.GetGroupMembersCountByGroupIDs(dbauthz.AsSystemRestricted(ctx), database.GetGroupMembersCountByGroupIDsParams{
		GroupIds:      groupIDs,
		IncludeSystem: false,
	})
	if err != nil {
		httpapi.InternalServerError(rw, err)
		return
	}
	countByGroup := make(map[uuid.UUID]int64, len(countRows))
	for _, row := range countRows {
		countByGroup[row.GroupID] = row.MemberCount
	}

	sdkGroups := make([]optimus-ide-collabsdk.Group, 0, len(groups))
	for _, group := range groups {
		sdkGroups = append(sdkGroups, db2sdk.Group(group, nil, int(countByGroup[group.Group.ID])))
	}

	httpapi.Write(ctx, rw, http.StatusOK, optimus-ide-collabsdk.ACLAvailable{
		// TODO: @emyrk we should return a MinimalUser here instead of a full user.
		// The FE requires the `email` field, so this cannot be done without
		// a UI change.
		Users:  db2sdk.ReducedUsers(users),
		Groups: sdkGroups,
	})
}

// @Summary Get template ACLs
// @ID get-template-acls
// @Security Optimus-IDE-CollabSessionToken
// @Produce json
// @Tags Enterprise
// @Param template path string true "Template ID" format(uuid)
// @Success 200 {object} optimus-ide-collabsdk.TemplateACL
// @Router /api/v2/templates/{template}/acl [get]
func (api *API) templateACL(rw http.ResponseWriter, r *http.Request) {
	var (
		ctx      = r.Context()
		template = httpmw.TemplateParam(r)
	)

	users, err := api.Database.GetTemplateUserRoles(ctx, template.ID)
	if err != nil {
		httpapi.InternalServerError(rw, err)
		return
	}

	dbGroups, err := api.Database.GetTemplateGroupRoles(ctx, template.ID)
	if err != nil {
		httpapi.InternalServerError(rw, err)
		return
	}

	userIDs := make([]uuid.UUID, 0, len(users))
	for _, user := range users {
		userIDs = append(userIDs, user.ID)
	}

	orgIDsByMemberIDsRows, err := api.Database.GetOrganizationIDsByMemberIDs(r.Context(), userIDs)
	if err != nil && !xerrors.Is(err, sql.ErrNoRows) {
		httpapi.InternalServerError(rw, err)
		return
	}

	organizationIDsByUserID := map[uuid.UUID][]uuid.UUID{}
	for _, organizationIDsByMemberIDsRow := range orgIDsByMemberIDsRows {
		organizationIDsByUserID[organizationIDsByMemberIDsRow.UserID] = organizationIDsByMemberIDsRow.OrganizationIDs
	}

	groups := make([]optimus-ide-collabsdk.TemplateGroup, 0, len(dbGroups))
	for _, group := range dbGroups {
		var members []database.GroupMember

		// This is a bit of a hack. The caller might not have permission to do this,
		// but they can read the acl list if the function got this far. So we let
		// them read the group members.
		// We should probably at least return more truncated user data here.
		// nolint:gocritic
		members, err = api.Database.GetGroupMembersByGroupID(dbauthz.AsSystemRestricted(ctx), database.GetGroupMembersByGroupIDParams{
			GroupID:       group.Group.ID,
			IncludeSystem: false,
		})
		if err != nil {
			httpapi.InternalServerError(rw, err)
			return
		}
		// nolint:gocritic
		memberCount, err := api.Database.GetGroupMembersCountByGroupID(dbauthz.AsSystemRestricted(ctx), database.GetGroupMembersCountByGroupIDParams{
			GroupID:       group.Group.ID,
			IncludeSystem: false,
		})
		if err != nil {
			httpapi.InternalServerError(rw, err)
			return
		}
		groups = append(groups, optimus-ide-collabsdk.TemplateGroup{
			Group: db2sdk.Group(database.GetGroupsRow{
				Group:                   group.Group,
				OrganizationName:        template.OrganizationName,
				OrganizationDisplayName: template.OrganizationDisplayName,
			}, members, int(memberCount)),
			Role: convertToTemplateRole(group.Actions),
		})
	}

	httpapi.Write(ctx, rw, http.StatusOK, optimus-ide-collabsdk.TemplateACL{
		Users:  convertTemplateUsers(users, organizationIDsByUserID),
		Groups: groups,
	})
}

// @Summary Update template ACL
// @ID update-template-acl
// @Security Optimus-IDE-CollabSessionToken
// @Accept json
// @Produce json
// @Tags Enterprise
// @Param template path string true "Template ID" format(uuid)
// @Param request body optimus-ide-collabsdk.UpdateTemplateACL true "Update template ACL request"
// @Success 200 {object} optimus-ide-collabsdk.Response
// @Router /api/v2/templates/{template}/acl [patch]
func (api *API) patchTemplateACL(rw http.ResponseWriter, r *http.Request) {
	var (
		ctx               = r.Context()
		template          = httpmw.TemplateParam(r)
		auditor           = api.AGPL.Auditor.Load()
		aReq, commitAudit = audit.InitRequest[database.Template](rw, &audit.RequestParams{
			Audit:          *auditor,
			Log:            api.Logger,
			Request:        r,
			Action:         database.AuditActionWrite,
			OrganizationID: template.OrganizationID,
		})
	)
	defer commitAudit()
	aReq.Old = template

	var req optimus-ide-collabsdk.UpdateTemplateACL
	if !httpapi.Read(ctx, rw, r, &req) {
		return
	}

	validErrs := acl.Validate(ctx, api.Database, TemplateACLUpdateValidator(req))
	if len(validErrs) > 0 {
		httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message:     "Invalid request to update template ACL",
			Validations: validErrs,
		})
		return
	}

	err := api.Database.InTx(func(tx database.Store) error {
		var err error
		template, err = tx.GetTemplateByID(ctx, template.ID)
		if err != nil {
			return xerrors.Errorf("get template by ID: %w", err)
		}

		for id, role := range req.UserPerms {
			if role == optimus-ide-collabsdk.TemplateRoleDeleted {
				delete(template.UserACL, id)
				continue
			}
			template.UserACL[id] = db2sdk.TemplateRoleActions(role)
		}

		for id, role := range req.GroupPerms {
			if role == optimus-ide-collabsdk.TemplateRoleDeleted {
				delete(template.GroupACL, id)
				continue
			}
			template.GroupACL[id] = db2sdk.TemplateRoleActions(role)
		}

		err = tx.UpdateTemplateACLByID(ctx, database.UpdateTemplateACLByIDParams{
			ID:       template.ID,
			UserACL:  template.UserACL,
			GroupACL: template.GroupACL,
		})
		if err != nil {
			return xerrors.Errorf("update template ACL by ID: %w", err)
		}
		template, err = tx.GetTemplateByID(ctx, template.ID)
		if err != nil {
			return xerrors.Errorf("get updated template by ID: %w", err)
		}
		return nil
	}, nil)
	if err != nil {
		httpapi.InternalServerError(rw, err)
		return
	}

	aReq.New = template

	httpapi.Write(ctx, rw, http.StatusOK, optimus-ide-collabsdk.Response{
		Message: "Successfully updated template ACL list.",
	})
}

type TemplateACLUpdateValidator optimus-ide-collabsdk.UpdateTemplateACL

var (
	templateACLUpdateUsersFieldName  = "user_perms"
	templateACLUpdateGroupsFieldName = "group_perms"
)

// TemplateACLUpdateValidator implements acl.UpdateValidator[optimus-ide-collabsdk.TemplateRole]
var _ acl.UpdateValidator[optimus-ide-collabsdk.TemplateRole] = TemplateACLUpdateValidator{}

func (w TemplateACLUpdateValidator) Users() (map[string]optimus-ide-collabsdk.TemplateRole, string) {
	return w.UserPerms, templateACLUpdateUsersFieldName
}

func (w TemplateACLUpdateValidator) Groups() (map[string]optimus-ide-collabsdk.TemplateRole, string) {
	return w.GroupPerms, templateACLUpdateGroupsFieldName
}

func (TemplateACLUpdateValidator) ValidateRole(role optimus-ide-collabsdk.TemplateRole) error {
	actions := db2sdk.TemplateRoleActions(role)
	if len(actions) == 0 && role != optimus-ide-collabsdk.TemplateRoleDeleted {
		return xerrors.Errorf("role %q is not a valid template role", role)
	}

	return nil
}

func convertTemplateUsers(tus []database.TemplateUser, orgIDsByUserIDs map[uuid.UUID][]uuid.UUID) []optimus-ide-collabsdk.TemplateUser {
	users := make([]optimus-ide-collabsdk.TemplateUser, 0, len(tus))

	for _, tu := range tus {
		users = append(users, optimus-ide-collabsdk.TemplateUser{
			User: db2sdk.User(tu.User, orgIDsByUserIDs[tu.User.ID]),
			Role: convertToTemplateRole(tu.Actions),
		})
	}

	return users
}

func convertToTemplateRole(actions []policy.Action) optimus-ide-collabsdk.TemplateRole {
	switch {
	case slice.SameElements(actions, db2sdk.TemplateRoleActions(optimus-ide-collabsdk.TemplateRoleAdmin)):
		return optimus-ide-collabsdk.TemplateRoleAdmin
	case slice.SameElements(actions, db2sdk.TemplateRoleActions(optimus-ide-collabsdk.TemplateRoleUse)):
		return optimus-ide-collabsdk.TemplateRoleUse
	}

	return optimus-ide-collabsdk.TemplateRoleDeleted
}

// TODO move to api.RequireFeatureMW when we are OK with changing the behavior.
func (api *API) templateRBACEnabledMW(next http.Handler) http.Handler {
	return api.RequireFeatureMW(optimus-ide-collabsdk.FeatureTemplateRBAC)(next)
}

func (api *API) RequireFeatureMW(feat optimus-ide-collabsdk.FeatureName) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			// Entitlement must be enabled.
			if !api.Entitlements.Enabled(feat) {
				// All feature warnings should be "Premium", not "Enterprise".
				httpapi.Write(r.Context(), rw, http.StatusForbidden, optimus-ide-collabsdk.Response{
					Message: fmt.Sprintf("%s is a Premium feature. Contact sales!", feat.Humanize()),
				})
				return
			}

			next.ServeHTTP(rw, r)
		})
	}
}

// @Summary Invalidate presets for template
// @ID invalidate-presets-for-template
// @Security Optimus-IDE-CollabSessionToken
// @Produce json
// @Tags Enterprise
// @Param template path string true "Template ID" format(uuid)
// @Success 200 {object} optimus-ide-collabsdk.InvalidatePresetsResponse
// @Router /api/v2/templates/{template}/prebuilds/invalidate [post]
func (api *API) postInvalidateTemplatePresets(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	template := httpmw.TemplateParam(r)

	// Authorization: user must be able to update the template
	if !api.Authorize(r, policy.ActionUpdate, template) {
		httpapi.ResourceNotFound(rw)
		return
	}

	// Update last_invalidated_at for all presets of the active template version
	invalidatedPresets, err := api.Database.UpdatePresetsLastInvalidatedAt(ctx, database.UpdatePresetsLastInvalidatedAtParams{
		TemplateID:        template.ID,
		LastInvalidatedAt: sql.NullTime{Time: api.Clock.Now(), Valid: true},
	})
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Failed to invalidate presets.",
			Detail:  err.Error(),
		})
		return
	}

	api.Logger.Info(ctx, "invalidated presets",
		slog.F("template_id", template.ID),
		slog.F("template_name", template.Name),
		slog.F("preset_count", len(invalidatedPresets)),
	)

	invalidated := db2sdk.InvalidatedPresets(invalidatedPresets)
	if invalidated == nil {
		invalidated = []optimus-ide-collabsdk.InvalidatedPreset{} // need to avoid nil value
	}

	httpapi.Write(ctx, rw, http.StatusOK, optimus-ide-collabsdk.InvalidatePresetsResponse{
		Invalidated: invalidated,
	})
}
