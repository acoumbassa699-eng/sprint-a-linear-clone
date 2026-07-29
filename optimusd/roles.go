package optimus-ide-collabd

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/db2sdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpapi"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpmw"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac/policy"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

// AssignableSiteRoles returns all site wide roles that can be assigned.
//
// @Summary Get site member roles
// @ID get-site-member-roles
// @Security Optimus-IDE-CollabSessionToken
// @Produce json
// @Tags Members
// @Success 200 {array} optimus-ide-collabsdk.AssignableRoles
// @Router /api/v2/users/roles [get]
func (api *API) AssignableSiteRoles(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	actorRoles := httpmw.UserAuthorization(r.Context())
	if !api.Authorize(r, policy.ActionRead, rbac.ResourceAssignRole) {
		httpapi.Forbidden(rw)
		return
	}

	dbCustomRoles, err := api.Database.CustomRoles(ctx, database.CustomRolesParams{
		LookupRoles: nil,
		// Only site wide custom roles to be included
		ExcludeOrgRoles:    true,
		OrganizationID:     uuid.Nil,
		IncludeSystemRoles: false,
	})
	if err != nil {
		httpapi.InternalServerError(rw, err)
		return
	}

	siteRoles := rbac.SiteBuiltInRoles()

	httpapi.Write(ctx, rw, http.StatusOK,
		assignableRoles(actorRoles.Roles, siteRoles, dbCustomRoles))
}

// assignableOrgRoles returns all org wide roles that can be assigned.
//
// @Summary Get member roles by organization
// @ID get-member-roles-by-organization
// @Security Optimus-IDE-CollabSessionToken
// @Produce json
// @Tags Members
// @Param organization path string true "Organization ID" format(uuid)
// @Success 200 {array} optimus-ide-collabsdk.AssignableRoles
// @Router /api/v2/organizations/{organization}/members/roles [get]
func (api *API) assignableOrgRoles(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	organization := httpmw.OrganizationParam(r)
	actorRoles := httpmw.UserAuthorization(r.Context())

	if !api.Authorize(r, policy.ActionRead, rbac.ResourceAssignOrgRole.InOrg(organization.ID)) {
		httpapi.ResourceNotFound(rw)
		return
	}

	roles := rbac.OrganizationRoles(organization.ID)
	dbCustomRoles, err := api.Database.CustomRoles(ctx, database.CustomRolesParams{
		LookupRoles:        nil,
		ExcludeOrgRoles:    false,
		OrganizationID:     organization.ID,
		IncludeSystemRoles: false,
	})
	if err != nil {
		httpapi.InternalServerError(rw, err)
		return
	}

	httpapi.Write(ctx, rw, http.StatusOK, assignableRoles(actorRoles.Roles, roles, dbCustomRoles))
}

func assignableRoles(actorRoles rbac.ExpandableRoles, roles []rbac.Role, customRoles []database.CustomRole) []optimus-ide-collabsdk.AssignableRoles {
	assignable := make([]optimus-ide-collabsdk.AssignableRoles, 0)
	for _, role := range roles {
		// The member role is implied, and not assignable.
		// If there is no display name, then the role is also unassigned.
		// This is not the ideal logic, but works for now.
		if role.Identifier == rbac.RoleMember() || (role.DisplayName == "") {
			continue
		}
		assignable = append(assignable, optimus-ide-collabsdk.AssignableRoles{
			Role:       db2sdk.RBACRole(role),
			Assignable: rbac.CanAssignRole(actorRoles, role.Identifier),
			BuiltIn:    true,
		})
	}

	for _, role := range customRoles {
		canAssign := rbac.CanAssignRole(actorRoles, rbac.CustomSiteRole())
		if role.RoleIdentifier().IsOrgRole() {
			canAssign = rbac.CanAssignRole(actorRoles, rbac.CustomOrganizationRole(role.OrganizationID.UUID))
		}

		assignable = append(assignable, optimus-ide-collabsdk.AssignableRoles{
			Role:       db2sdk.Role(role),
			Assignable: canAssign,
			BuiltIn:    false,
		})
	}
	return assignable
}
