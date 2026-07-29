package optimus-ide-collabd

import (
	"database/sql"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/xerrors"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/audit"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/db2sdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbauthz"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbtime"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpapi"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpmw"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac/rolestore"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

// @Summary Update organization
// @ID update-organization
// @Security Optimus-IDE-CollabSessionToken
// @Accept json
// @Produce json
// @Tags Organizations
// @Param organization path string true "Organization ID or name"
// @Param request body optimus-ide-collabsdk.UpdateOrganizationRequest true "Patch organization request"
// @Success 200 {object} optimus-ide-collabsdk.Organization
// @Router /api/v2/organizations/{organization} [patch]
func (api *API) patchOrganization(rw http.ResponseWriter, r *http.Request) {
	var (
		ctx               = r.Context()
		organization      = httpmw.OrganizationParam(r)
		auditor           = api.AGPL.Auditor.Load()
		aReq, commitAudit = audit.InitRequest[database.Organization](rw, &audit.RequestParams{
			Audit:          *auditor,
			Log:            api.Logger,
			Request:        r,
			Action:         database.AuditActionWrite,
			OrganizationID: organization.ID,
		})
	)
	aReq.Old = organization
	defer commitAudit()

	var req optimus-ide-collabsdk.UpdateOrganizationRequest
	if !httpapi.Read(ctx, rw, r, &req) {
		return
	}

	// "default" is a reserved name that always refers to the default org (much like the way we
	// use "me" for users).
	if req.Name == optimus-ide-collabsdk.DefaultOrganization {
		httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message: fmt.Sprintf("Organization name %q is reserved.", optimus-ide-collabsdk.DefaultOrganization),
		})
		return
	}

	// Deviations from rbac.DefaultOrgMemberRoles require the
	// minimum-implicit-member experiment.
	if req.DefaultOrgMemberRoles != nil &&
		!slices.Equal(*req.DefaultOrgMemberRoles, rbac.DefaultOrgMemberRoles()) &&
		!api.AGPL.Experiments.Enabled(optimus-ide-collabsdk.ExperimentMinimumImplicitMember) {
		httpapi.Write(ctx, rw, http.StatusForbidden, optimus-ide-collabsdk.Response{
			Message: "Changing default organization roles is not enabled on this deployment.",
			Detail:  fmt.Sprintf("Setting default_org_member_roles to anything other than %v requires the %q experiment.", rbac.DefaultOrgMemberRoles(), optimus-ide-collabsdk.ExperimentMinimumImplicitMember),
		})
		return
	}

	// default_org_member_roles currently accepts built-in role names only.
	// Custom (DB-stored) roles are intentionally rejected here so the
	// caller cannot land a malformed name that would break role expansion
	// for every member of the org. A future change can extend this to
	// custom org roles by routing through canAssignRoles in dbauthz.
	if req.DefaultOrgMemberRoles != nil {
		for _, name := range *req.DefaultOrgMemberRoles {
			if _, err := rbac.RoleByName(rbac.RoleIdentifier{Name: name, OrganizationID: organization.ID}); err != nil {
				httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
					Message: "Invalid default_org_member_roles entry.",
					Detail:  fmt.Sprintf("%q is not a built-in role; default_org_member_roles currently accepts built-in role names only.", name),
					Validations: []optimus-ide-collabsdk.ValidationError{{
						Field:  "default_org_member_roles",
						Detail: fmt.Sprintf("%q is not a built-in role.", name),
					}},
				})
				return
			}
		}
	}

	err := database.ReadModifyUpdate(api.Database, func(tx database.Store) error {
		var err error
		organization, err = tx.GetOrganizationByID(ctx, organization.ID)
		if err != nil {
			return err
		}

		updateOrgParams := database.UpdateOrganizationParams{
			UpdatedAt:             dbtime.Now(),
			ID:                    organization.ID,
			Name:                  organization.Name,
			DisplayName:           organization.DisplayName,
			Description:           organization.Description,
			Icon:                  organization.Icon,
			DefaultOrgMemberRoles: organization.DefaultOrgMemberRoles,
		}

		if req.Name != "" {
			updateOrgParams.Name = req.Name
		}
		if req.DisplayName != "" {
			updateOrgParams.DisplayName = req.DisplayName
		}
		if req.Description != nil {
			updateOrgParams.Description = *req.Description
		}
		if req.Icon != nil {
			updateOrgParams.Icon = *req.Icon
		}
		if req.DefaultOrgMemberRoles != nil {
			updateOrgParams.DefaultOrgMemberRoles = *req.DefaultOrgMemberRoles
		}

		organization, err = tx.UpdateOrganization(ctx, updateOrgParams)
		if err != nil {
			return err
		}
		return nil
	})

	if httpapi.Is404Error(err) {
		httpapi.ResourceNotFound(rw)
		return
	}
	if database.IsUniqueViolation(err) {
		httpapi.Write(ctx, rw, http.StatusConflict, optimus-ide-collabsdk.Response{
			Message: fmt.Sprintf("Organization already exists with the name %q.", req.Name),
			Validations: []optimus-ide-collabsdk.ValidationError{{
				Field:  "name",
				Detail: "This value is already in use and should be unique.",
			}},
		})
		return
	}
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Internal error updating organization.",
			Detail:  fmt.Sprintf("update organization: %s", err.Error()),
		})
		return
	}

	aReq.New = organization
	httpapi.Write(ctx, rw, http.StatusOK, db2sdk.Organization(organization))
}

// @Summary Delete organization
// @ID delete-organization
// @Security Optimus-IDE-CollabSessionToken
// @Produce json
// @Tags Organizations
// @Param organization path string true "Organization ID or name"
// @Success 200 {object} optimus-ide-collabsdk.Response
// @Router /api/v2/organizations/{organization} [delete]
func (api *API) deleteOrganization(rw http.ResponseWriter, r *http.Request) {
	var (
		ctx               = r.Context()
		organization      = httpmw.OrganizationParam(r)
		auditor           = api.AGPL.Auditor.Load()
		aReq, commitAudit = audit.InitRequest[database.Organization](rw, &audit.RequestParams{
			Audit:          *auditor,
			Log:            api.Logger,
			Request:        r,
			Action:         database.AuditActionDelete,
			OrganizationID: organization.ID,
		})
	)
	aReq.Old = organization
	defer commitAudit()

	if organization.IsDefault {
		httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message: "Default organization cannot be deleted.",
		})
		return
	}

	err := api.Database.InTx(func(tx database.Store) error {
		err := tx.UpdateOrganizationDeletedByID(ctx, database.UpdateOrganizationDeletedByIDParams{
			ID:        organization.ID,
			UpdatedAt: dbtime.Now(),
		})
		if err != nil {
			return xerrors.Errorf("delete organization: %w", err)
		}
		return nil
	}, nil)
	if err != nil {
		orgResourcesRow, queryErr := api.Database.GetOrganizationResourceCountByID(ctx, organization.ID)
		if queryErr != nil {
			httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
				Message: "Internal error deleting organization.",
				Detail:  fmt.Sprintf("delete organization: %s", err.Error()),
			})

			return
		}

		detailParts := make([]string, 0)

		addDetailPart := func(resource string, count int64) {
			if count == 1 {
				detailParts = append(detailParts, fmt.Sprintf("1 %s", resource))
			} else if count > 1 {
				detailParts = append(detailParts, fmt.Sprintf("%d %ss", count, resource))
			}
		}

		addDetailPart("workspace", orgResourcesRow.WorkspaceCount)
		addDetailPart("template", orgResourcesRow.TemplateCount)

		// There will always be one member and group so instead we need to check that
		// the count is greater than one.
		addDetailPart("member", orgResourcesRow.MemberCount-1)
		addDetailPart("group", orgResourcesRow.GroupCount-1)

		addDetailPart("provisioner key", orgResourcesRow.ProvisionerKeyCount)

		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Error deleting organization.",
			Detail:  fmt.Sprintf("This organization has %s that must be deleted first.", strings.Join(detailParts, ", ")),
		})

		return
	}

	aReq.New = database.Organization{}
	httpapi.Write(ctx, rw, http.StatusOK, optimus-ide-collabsdk.Response{
		Message: "Organization has been deleted.",
	})
}

// @Summary Create organization
// @ID create-organization
// @Security Optimus-IDE-CollabSessionToken
// @Accept json
// @Produce json
// @Tags Organizations
// @Param request body optimus-ide-collabsdk.CreateOrganizationRequest true "Create organization request"
// @Success 201 {object} optimus-ide-collabsdk.Organization
// @Router /api/v2/organizations [post]
func (api *API) postOrganizations(rw http.ResponseWriter, r *http.Request) {
	var (
		// organizationID is required before the audit log entry is created.
		organizationID    = uuid.New()
		ctx               = r.Context()
		apiKey            = httpmw.APIKey(r)
		auditor           = api.AGPL.Auditor.Load()
		aReq, commitAudit = audit.InitRequest[database.Organization](rw, &audit.RequestParams{
			Audit:          *auditor,
			Log:            api.Logger,
			Request:        r,
			Action:         database.AuditActionCreate,
			OrganizationID: organizationID,
		})
	)
	aReq.Old = database.Organization{}
	defer commitAudit()

	var req optimus-ide-collabsdk.CreateOrganizationRequest
	if !httpapi.Read(ctx, rw, r, &req) {
		return
	}

	if req.Name == optimus-ide-collabsdk.DefaultOrganization {
		httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message: fmt.Sprintf("Organization name %q is reserved.", optimus-ide-collabsdk.DefaultOrganization),
		})
		return
	}

	_, err := api.Database.GetOrganizationByName(ctx, database.GetOrganizationByNameParams{
		Name:    req.Name,
		Deleted: false,
	})
	if err == nil {
		httpapi.Write(ctx, rw, http.StatusConflict, optimus-ide-collabsdk.Response{
			Message: "Organization already exists with that name.",
		})
		return
	}
	if !xerrors.Is(err, sql.ErrNoRows) {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: fmt.Sprintf("Internal error fetching organization %q.", req.Name),
			Detail:  err.Error(),
		})
		return
	}

	var organization database.Organization
	err = api.Database.InTx(func(tx database.Store) error {
		// Serialize creation and reconciliation of the org-member
		// system role across optimus-ide-collabd instances (e.g. during rolling
		// restarts).
		err := tx.AcquireLock(ctx, database.LockIDReconcileSystemRoles)
		if err != nil {
			return xerrors.Errorf("acquire system roles reconciliation lock: %w", err)
		}

		if req.DisplayName == "" {
			req.DisplayName = req.Name
		}

		organization, err = tx.InsertOrganization(ctx, database.InsertOrganizationParams{
			ID:                    organizationID,
			Name:                  req.Name,
			DisplayName:           req.DisplayName,
			Description:           req.Description,
			Icon:                  req.Icon,
			CreatedAt:             dbtime.Now(),
			UpdatedAt:             dbtime.Now(),
			DefaultOrgMemberRoles: rbac.DefaultOrgMemberRoles(),
		})
		if err != nil {
			return xerrors.Errorf("create organization: %w", err)
		}

		// Populate the placeholder system role(s) that the DB trigger
		// created for us.
		//nolint:gocritic // ReconcileOrgMemberRole needs the system:update
		// permission that user doesn't have.
		sysCtx := dbauthz.AsSystemRestricted(ctx)
		for roleName := range rolestore.SystemRoleNames {
			_, _, err = rolestore.ReconcileSystemRole(sysCtx, tx, database.CustomRole{
				Name:           roleName,
				OrganizationID: uuid.NullUUID{UUID: organizationID, Valid: true},
			}, organization)
			if err != nil {
				return xerrors.Errorf("reconcile %s role for organization %s: %w",
					roleName, organizationID, err)
			}
		}

		_, err = tx.InsertOrganizationMember(ctx, database.InsertOrganizationMemberParams{
			OrganizationID: organization.ID,
			UserID:         apiKey.UserID,
			CreatedAt:      dbtime.Now(),
			UpdatedAt:      dbtime.Now(),
			Roles:          []string{
				// TODO: When organizations are allowed to be created, we should
				// come back to determining the default role of the person who
				// creates the org. Until that happens, all users in an organization
				// should be just regular members.
			},
		})
		if err != nil {
			return xerrors.Errorf("create organization admin: %w", err)
		}

		_, err = tx.InsertAllUsersGroup(ctx, organization.ID)
		if err != nil {
			return xerrors.Errorf("create %q group: %w", database.EveryoneGroup, err)
		}
		return nil
	}, nil)
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Internal error inserting organization member.",
			Detail:  err.Error(),
		})
		return
	}

	aReq.New = organization
	httpapi.Write(ctx, rw, http.StatusCreated, db2sdk.Organization(organization))
}
