package optimus-ide-collabd

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/audit"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/db2sdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpapi"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpmw"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac/policy"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/util/slice"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

// postOrgRoles will allow creating a custom organization role
//
// @Summary Insert a custom organization role
// @ID insert-a-custom-organization-role
// @Security Optimus-IDE-CollabSessionToken
// @Accept json
// @Produce json
// @Param organization path string true "Organization ID" format(uuid)
// @Param request body optimus-ide-collabsdk.CustomRoleRequest true "Insert role request"
// @Tags Members
// @Success 200 {array} optimus-ide-collabsdk.Role
// @Router /api/v2/organizations/{organization}/members/roles [post]
func (api *API) postOrgRoles(rw http.ResponseWriter, r *http.Request) {
	var (
		ctx               = r.Context()
		db                = api.Database
		auditor           = api.AGPL.Auditor.Load()
		organization      = httpmw.OrganizationParam(r)
		aReq, commitAudit = audit.InitRequest[database.CustomRole](rw, &audit.RequestParams{
			Audit:          *auditor,
			Log:            api.Logger,
			Request:        r,
			Action:         database.AuditActionCreate,
			OrganizationID: organization.ID,
		})
	)
	defer commitAudit()

	var req optimus-ide-collabsdk.CustomRoleRequest
	if !httpapi.Read(ctx, rw, r, &req) {
		return
	}

	if !validOrganizationRoleRequest(ctx, req, rw) {
		return
	}

	inserted, err := db.InsertCustomRole(ctx, database.InsertCustomRoleParams{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		OrganizationID: uuid.NullUUID{
			UUID:  organization.ID,
			Valid: true,
		},
		SitePermissions: slice.List(req.SitePermissions, sdkPermissionToDB),
		OrgPermissions:  slice.List(req.OrganizationPermissions, sdkPermissionToDB),
		UserPermissions: slice.List(req.UserPermissions, sdkPermissionToDB),
		// Satisfy the linter (we don't support member permissions in non-system roles).
		MemberPermissions: database.CustomRolePermissions{},
		IsSystem:          false,
	})
	if httpapi.Is404Error(err) {
		httpapi.ResourceNotFound(rw)
		return
	}
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message: "Failed to update role permissions",
			Detail:  err.Error(),
		})
		return
	}
	aReq.New = inserted

	httpapi.Write(ctx, rw, http.StatusOK, db2sdk.Role(inserted))
}

// putOrgRoles will allow updating a custom organization role
//
// @Summary Update a custom organization role
// @ID update-a-custom-organization-role
// @Security Optimus-IDE-CollabSessionToken
// @Accept json
// @Produce json
// @Param organization path string true "Organization ID" format(uuid)
// @Param request body optimus-ide-collabsdk.CustomRoleRequest true "Update role request"
// @Tags Members
// @Success 200 {array} optimus-ide-collabsdk.Role
// @Router /api/v2/organizations/{organization}/members/roles [put]
func (api *API) putOrgRoles(rw http.ResponseWriter, r *http.Request) {
	var (
		ctx               = r.Context()
		db                = api.Database
		auditor           = api.AGPL.Auditor.Load()
		organization      = httpmw.OrganizationParam(r)
		aReq, commitAudit = audit.InitRequest[database.CustomRole](rw, &audit.RequestParams{
			Audit:          *auditor,
			Log:            api.Logger,
			Request:        r,
			Action:         database.AuditActionWrite,
			OrganizationID: organization.ID,
		})
	)
	defer commitAudit()

	var req optimus-ide-collabsdk.CustomRoleRequest
	if !httpapi.Read(ctx, rw, r, &req) {
		return
	}

	if !validOrganizationRoleRequest(ctx, req, rw) {
		return
	}

	originalRoles, err := db.CustomRoles(ctx, database.CustomRolesParams{
		LookupRoles: []database.NameOrganizationPair{
			{
				Name:           req.Name,
				OrganizationID: organization.ID,
			},
		},
		ExcludeOrgRoles:    false,
		OrganizationID:     organization.ID,
		IncludeSystemRoles: false,
	})
	// If it is a 404 (not found) error, ignore it.
	if err != nil && !httpapi.Is404Error(err) {
		httpapi.InternalServerError(rw, err)
		return
	}
	if len(originalRoles) == 1 {
		// For auditing changes to a role.
		aReq.Old = originalRoles[0]
	}

	updated, err := db.UpdateCustomRole(ctx, database.UpdateCustomRoleParams{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		OrganizationID: uuid.NullUUID{
			UUID:  organization.ID,
			Valid: true,
		},
		// Invalid permissions are filtered out. If this is changed
		// to throw an error, then the story of a previously valid role
		// now being invalid has to be addressed. Optimus-IDE-Collab can change permissions,
		// objects, and actions at any time.
		SitePermissions: slice.List(filterInvalidPermissions(req.SitePermissions), sdkPermissionToDB),
		OrgPermissions:  slice.List(filterInvalidPermissions(req.OrganizationPermissions), sdkPermissionToDB),
		UserPermissions: slice.List(filterInvalidPermissions(req.UserPermissions), sdkPermissionToDB),
		// Satisfy the linter (we don't support member permissions in non-system roles).
		MemberPermissions: database.CustomRolePermissions{},
	})
	if httpapi.Is404Error(err) {
		httpapi.ResourceNotFound(rw)
		return
	}
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message: "Failed to update role permissions",
			Detail:  err.Error(),
		})
		return
	}
	aReq.New = updated

	httpapi.Write(ctx, rw, http.StatusOK, db2sdk.Role(updated))
}

// deleteOrgRole will remove a custom role from an organization
//
// @Summary Delete a custom organization role
// @ID delete-a-custom-organization-role
// @Security Optimus-IDE-CollabSessionToken
// @Produce json
// @Param organization path string true "Organization ID" format(uuid)
// @Param roleName path string true "Role name"
// @Tags Members
// @Success 200 {array} optimus-ide-collabsdk.Role
// @Router /api/v2/organizations/{organization}/members/roles/{roleName} [delete]
func (api *API) deleteOrgRole(rw http.ResponseWriter, r *http.Request) {
	var (
		ctx               = r.Context()
		auditor           = api.AGPL.Auditor.Load()
		organization      = httpmw.OrganizationParam(r)
		aReq, commitAudit = audit.InitRequest[database.CustomRole](rw, &audit.RequestParams{
			Audit:          *auditor,
			Log:            api.Logger,
			Request:        r,
			Action:         database.AuditActionDelete,
			OrganizationID: organization.ID,
		})
	)
	defer commitAudit()

	rolename := chi.URLParam(r, "roleName")

	// Catch requests that try to delete system roles.
	if !validOrganizationRoleRequest(ctx, optimus-ide-collabsdk.CustomRoleRequest{Name: rolename}, rw) {
		return
	}

	roles, err := api.Database.CustomRoles(ctx, database.CustomRolesParams{
		LookupRoles: []database.NameOrganizationPair{
			{
				Name:           rolename,
				OrganizationID: organization.ID,
			},
		},
		ExcludeOrgRoles:    false,
		IncludeSystemRoles: false,
		// Linter requires all fields to be set. This field is not actually required.
		OrganizationID: organization.ID,
	})
	if err != nil {
		httpapi.InternalServerError(rw, err)
		return
	}
	if len(roles) == 0 {
		httpapi.Write(ctx, rw, http.StatusNotFound, optimus-ide-collabsdk.Response{
			Message:     fmt.Sprintf("No custom role with the name %s found", rolename),
			Detail:      "no role found",
			Validations: nil,
		})
		return
	}
	if len(roles) > 1 {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message:     fmt.Sprintf("Multiple roles with the name %s found", rolename),
			Detail:      "multiple roles found, this should never happen",
			Validations: nil,
		})
		return
	}
	aReq.Old = roles[0]

	err = api.Database.DeleteCustomRole(ctx, database.DeleteCustomRoleParams{
		Name: rolename,
		OrganizationID: uuid.NullUUID{
			UUID:  organization.ID,
			Valid: true,
		},
	})
	if httpapi.IsUnauthorizedError(err) {
		httpapi.Forbidden(rw)
		return
	}
	if err != nil {
		httpapi.InternalServerError(rw, err)
		return
	}
	aReq.New = database.CustomRole{}

	httpapi.Write(ctx, rw, http.StatusNoContent, nil)
}

func filterInvalidPermissions(permissions []optimus-ide-collabsdk.Permission) []optimus-ide-collabsdk.Permission {
	// Filter out any invalid permissions
	var validPermissions []optimus-ide-collabsdk.Permission
	for _, permission := range permissions {
		err := rbac.Permission{
			Negate:       permission.Negate,
			ResourceType: string(permission.ResourceType),
			Action:       policy.Action(permission.Action),
		}.Valid()
		if err != nil {
			continue
		}
		validPermissions = append(validPermissions, permission)
	}
	return validPermissions
}

func sdkPermissionToDB(p optimus-ide-collabsdk.Permission) database.CustomRolePermission {
	return database.CustomRolePermission{
		Negate:       p.Negate,
		ResourceType: string(p.ResourceType),
		Action:       policy.Action(p.Action),
	}
}

func validOrganizationRoleRequest(ctx context.Context, req optimus-ide-collabsdk.CustomRoleRequest, rw http.ResponseWriter) bool {
	// This check is not ideal, but we cannot enforce a unique role name in the db against
	// the built-in role names.
	if rbac.ReservedRoleName(req.Name) {
		httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message: "Reserved role name",
			Detail:  fmt.Sprintf("%q is a reserved role name, and not allowed to be used", req.Name),
		})
		return false
	}

	if err := optimus-ide-collabsdk.NameValid(req.Name); err != nil {
		httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message: "Invalid role name",
			Detail:  err.Error(),
		})
		return false
	}

	// Only organization permissions are allowed to be granted
	if len(req.SitePermissions) > 0 {
		httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message: "Invalid request, not allowed to assign site wide permissions for an organization role.",
			Detail:  "organization scoped roles may not contain site wide permissions",
		})
		return false
	}

	if len(req.UserPermissions) > 0 {
		httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message: "Invalid request, not allowed to assign user permissions for an organization role.",
			Detail:  "organization scoped roles may not contain user permissions",
		})
		return false
	}

	if len(req.OrganizationMemberPermissions) > 0 {
		httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message: "Invalid request, not allowed to assign organization member permissions for an organization role.",
			Detail:  "organization scoped roles may not contain organization member permissions",
		})
		return false
	}

	return true
}
