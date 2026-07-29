package optimus-ide-collabd

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"golang.org/x/xerrors"

	"cdr.dev/slog/v3"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/audit"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/db2sdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbauthz"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbtime"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/gitsshkey"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpapi"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpmw"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/notifications"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac/policy"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/searchquery"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/telemetry"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/userpassword"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/util/ptr"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/util/slice"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

// userDebugOIDC returns the OIDC debug context for the user.
// Not going to expose this via swagger as the return payload is not guaranteed
// to be consistent between releases.
//
// @Summary Debug OIDC context for a user
// @ID debug-oidc-context-for-a-user
// @Security Optimus-IDE-CollabSessionToken
// @Tags Agents
// @Success 200 "Success"
// @Param user path string true "User ID, name, or me"
// @Router /api/v2/debug/{user}/debug-link [get]
// @x-apidocgen {"skip": true}
func (api *API) userDebugOIDC(rw http.ResponseWriter, r *http.Request) {
	var (
		ctx  = r.Context()
		user = httpmw.UserParam(r)
	)

	if user.LoginType != database.LoginTypeOIDC {
		httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message: "User is not an OIDC user.",
		})
		return
	}

	link, err := api.Database.GetUserLinkByUserIDLoginType(ctx, database.GetUserLinkByUserIDLoginTypeParams{
		UserID:    user.ID,
		LoginType: database.LoginTypeOIDC,
	})
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Failed to get user links.",
			Detail:  err.Error(),
		})
		return
	}

	httpapi.Write(ctx, rw, http.StatusOK, link.Claims)
}

// Returns the merged OIDC claims for the authenticated user.
//
// @Summary Get OIDC claims for the authenticated user
// @ID get-oidc-claims-for-the-authenticated-user
// @Security Optimus-IDE-CollabSessionToken
// @Produce json
// @Tags Users
// @Success 200 {object} optimus-ide-collabsdk.OIDCClaimsResponse
// @Router /api/v2/users/oidc-claims [get]
func (api *API) userOIDCClaims(rw http.ResponseWriter, r *http.Request) {
	var (
		ctx    = r.Context()
		apiKey = httpmw.APIKey(r)
	)

	user, err := api.Database.GetUserByID(ctx, apiKey.UserID)
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Failed to get user.",
			Detail:  err.Error(),
		})
		return
	}

	if user.LoginType != database.LoginTypeOIDC {
		httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message: "User is not an OIDC user.",
		})
		return
	}

	//nolint:gocritic // GetUserLinkByUserIDLoginType requires reading
	// rbac.ResourceSystem. The endpoint is scoped to the authenticated
	// user's own identity via apiKey, so this is safe.
	link, err := api.Database.GetUserLinkByUserIDLoginType(
		dbauthz.AsSystemRestricted(ctx),
		database.GetUserLinkByUserIDLoginTypeParams{
			UserID:    user.ID,
			LoginType: database.LoginTypeOIDC,
		},
	)
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Failed to get user link.",
			Detail:  err.Error(),
		})
		return
	}

	claims := link.Claims.MergedClaims
	if claims == nil {
		claims = map[string]interface{}{}
	}
	httpapi.Write(ctx, rw, http.StatusOK, optimus-ide-collabsdk.OIDCClaimsResponse{
		Claims: claims,
	})
}

// Returns whether the initial user has been created or not.
//
// @Summary Check initial user created
// @ID check-initial-user-created
// @Security Optimus-IDE-CollabSessionToken
// @Produce json
// @Tags Users
// @Success 200 {object} optimus-ide-collabsdk.Response
// @Router /api/v2/users/first [get]
func (api *API) firstUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// nolint:gocritic // Getting user count is a system function.
	userCount, err := api.Database.GetUserCount(dbauthz.AsSystemRestricted(ctx), false)
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Internal error fetching user count.",
			Detail:  err.Error(),
		})
		return
	}

	if userCount == 0 {
		httpapi.Write(ctx, rw, http.StatusNotFound, optimus-ide-collabsdk.Response{
			Message: "The initial user has not been created!",
		})
		return
	}

	httpapi.Write(ctx, rw, http.StatusOK, optimus-ide-collabsdk.Response{
		Message: "The initial user has already been created!",
	})
}

// Creates the initial user for a Optimus-IDE-Collab deployment.
//
// @Summary Create initial user
// @ID create-initial-user
// @Security Optimus-IDE-CollabSessionToken
// @Accept json
// @Produce json
// @Tags Users
// @Param request body optimus-ide-collabsdk.CreateFirstUserRequest true "First user request"
// @Success 201 {object} optimus-ide-collabsdk.CreateFirstUserResponse
// @Router /api/v2/users/first [post]
func (api *API) postFirstUser(rw http.ResponseWriter, r *http.Request) {
	// The first user can also be created via oidc, so if making changes to the flow,
	// ensure that the oidc flow is also updated.
	ctx := r.Context()
	var createUser optimus-ide-collabsdk.CreateFirstUserRequest
	if !httpapi.Read(ctx, rw, r, &createUser) {
		return
	}

	// This should only function for the first user.
	// nolint:gocritic // Getting user count is a system function.
	userCount, err := api.Database.GetUserCount(dbauthz.AsSystemRestricted(ctx), false)
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Internal error fetching user count.",
			Detail:  err.Error(),
		})
		return
	}

	// If a user already exists, the initial admin user no longer can be created.
	if userCount != 0 {
		httpapi.Write(ctx, rw, http.StatusConflict, optimus-ide-collabsdk.Response{
			Message: "The initial user has already been created.",
		})
		return
	}

	err = userpassword.Validate(createUser.Password)
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message: "Password is invalid",
			Validations: []optimus-ide-collabsdk.ValidationError{{
				Field:  "password",
				Detail: err.Error(),
			}},
		})
		return
	}

	if createUser.Trial && api.TrialGenerator != nil {
		err = api.TrialGenerator(ctx, optimus-ide-collabsdk.LicensorTrialRequest{
			Email:       createUser.Email,
			FirstName:   createUser.TrialInfo.FirstName,
			LastName:    createUser.TrialInfo.LastName,
			PhoneNumber: createUser.TrialInfo.PhoneNumber,
			JobTitle:    createUser.TrialInfo.JobTitle,
			CompanyName: createUser.TrialInfo.CompanyName,
			Country:     createUser.TrialInfo.Country,
			Developers:  createUser.TrialInfo.Developers,
		})
		if err != nil {
			httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
				Message: "Failed to generate trial",
				Detail:  err.Error(),
			})
			return
		}
	}

	//nolint:gocritic // needed to create first user
	defaultOrg, err := api.Database.GetDefaultOrganization(dbauthz.AsSystemRestricted(ctx))
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Internal error fetching default organization. If you are encountering this error, you will have to restart the Optimus-IDE-Collab deployment.",
			Detail:  err.Error(),
		})
		return
	}

	//nolint:gocritic // needed to create first user
	user, err := api.CreateUser(dbauthz.AsSystemRestricted(ctx), api.Database, CreateUserRequest{
		CreateUserRequestWithOrgs: optimus-ide-collabsdk.CreateUserRequestWithOrgs{
			Email:    createUser.Email,
			Username: createUser.Username,
			Name:     createUser.Name,
			Password: createUser.Password,
			// There's no reason to create the first user as dormant, since you have
			// to login immediately anyways.
			UserStatus:      ptr.Ref(optimus-ide-collabsdk.UserStatusActive),
			OrganizationIDs: []uuid.UUID{defaultOrg.ID},
		},
		LoginType:          database.LoginTypePassword,
		RBACRoles:          []string{rbac.RoleOwner().String()},
		accountCreatorName: "optimus-ide-collab",
	})
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Internal error creating user.",
			Detail:  err.Error(),
		})
		return
	}

	if api.RefreshEntitlements != nil {
		err = api.RefreshEntitlements(ctx)
		if err != nil {
			api.Logger.Error(ctx, "failed to refresh entitlements after generating trial license")
			return
		}
	} else {
		api.Logger.Debug(ctx, "entitlements will not be refreshed")
	}

	telemetryUser := telemetry.ConvertUser(user)
	// Send the initial users email address!
	telemetryUser.Email = &user.Email
	// Only populate onboarding data when the client actually sent it. A nil
	// OnboardingInfo means the request came from an older client, the CLI, or
	// the OIDC flow — not from a user who answered "no" to every question.
	var onboarding *telemetry.FirstUserOnboarding
	if createUser.OnboardingInfo != nil {
		onboarding = &telemetry.FirstUserOnboarding{
			NewsletterMarketing: createUser.OnboardingInfo.NewsletterMarketing,
			NewsletterReleases:  createUser.OnboardingInfo.NewsletterReleases,
		}
	}
	api.Telemetry.Report(&telemetry.Snapshot{
		Users:               []telemetry.User{telemetryUser},
		FirstUserOnboarding: onboarding,
	})

	httpapi.Write(ctx, rw, http.StatusCreated, optimus-ide-collabsdk.CreateFirstUserResponse{
		UserID:         user.ID,
		OrganizationID: defaultOrg.ID,
	})
}

// @Summary Get users
// @ID get-users
// @Security Optimus-IDE-CollabSessionToken
// @Produce json
// @Tags Users
// @Param q query string false "Search query"
// @Param after_id query string false "After ID" format(uuid)
// @Param limit query int false "Page limit"
// @Param offset query int false "Page offset"
// @Success 200 {object} optimus-ide-collabsdk.GetUsersResponse
// @Router /api/v2/users [get]
func (api *API) users(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	users, userCount, ok := api.GetUsers(rw, r)
	if !ok {
		return
	}

	userIDs := make([]uuid.UUID, 0, len(users))
	for _, user := range users {
		userIDs = append(userIDs, user.ID)
	}
	organizationIDsByMemberIDsRows, err := api.Database.GetOrganizationIDsByMemberIDs(ctx, userIDs)
	if xerrors.Is(err, sql.ErrNoRows) {
		err = nil
	}
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Internal error fetching user's organizations.",
			Detail:  err.Error(),
		})
		return
	}
	organizationIDsByUserID := map[uuid.UUID][]uuid.UUID{}
	for _, organizationIDsByMemberIDsRow := range organizationIDsByMemberIDsRows {
		organizationIDsByUserID[organizationIDsByMemberIDsRow.UserID] = organizationIDsByMemberIDsRow.OrganizationIDs
	}

	var aiSeatSet map[uuid.UUID]struct{}
	if api.Entitlements.Enabled(optimus-ide-collabsdk.FeatureAIGovernanceUserLimit) {
		var aiSeatUserIDs []uuid.UUID
		//nolint:gocritic // AI seat state is a system-level read gated by entitlement.
		aiSeatUserIDs, err = api.Database.GetUserAISeatStates(dbauthz.AsSystemRestricted(ctx), userIDs)
		if err != nil {
			if !xerrors.Is(err, sql.ErrNoRows) {
				api.Logger.Warn(
					ctx,
					"failed to fetch AI seat states for users",
					slog.F("user_count", len(userIDs)),
					slog.Error(err),
				)
			}
			aiSeatUserIDs = nil
		}

		aiSeatSet = make(map[uuid.UUID]struct{}, len(aiSeatUserIDs))
		for _, uid := range aiSeatUserIDs {
			aiSeatSet[uid] = struct{}{}
		}
	}

	httpapi.Write(ctx, rw, http.StatusOK, optimus-ide-collabsdk.GetUsersResponse{
		Users: convertUsers(users, organizationIDsByUserID, aiSeatSet),
		Count: int(userCount),
	})
}

func (api *API) GetUsers(rw http.ResponseWriter, r *http.Request) ([]database.User, int64, bool) {
	ctx := r.Context()
	query := r.URL.Query().Get("q")
	params, errs := searchquery.Users(query)
	if len(errs) > 0 {
		httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message:     "Invalid user search query.",
			Validations: errs,
		})
		return nil, -1, false
	}

	paginationParams, ok := ParsePagination(rw, r)
	if !ok {
		return nil, -1, false
	}

	userRows, err := api.Database.GetUsers(ctx, database.GetUsersParams{
		AfterID:          paginationParams.AfterID,
		Search:           params.Search,
		Name:             params.Name,
		Status:           params.Status,
		IsServiceAccount: params.IsServiceAccount,
		RbacRole:         params.RbacRole,
		LastSeenBefore:   params.LastSeenBefore,
		LastSeenAfter:    params.LastSeenAfter,
		CreatedAfter:     params.CreatedAfter,
		CreatedBefore:    params.CreatedBefore,
		GithubComUserID:  params.GithubComUserID,
		LoginType:        params.LoginType,
		// #nosec G115 - Pagination offsets are small and fit in int32
		OffsetOpt: int32(paginationParams.Offset),
		// #nosec G115 - Pagination limits are small and fit in int32
		LimitOpt: int32(paginationParams.Limit),
	})
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Internal error fetching users.",
			Detail:  err.Error(),
		})
		return nil, -1, false
	}

	// GetUsers does not return ErrNoRows because it uses a window function to get the count.
	// So we need to check if the userRows is empty and return an empty array if so.
	if len(userRows) == 0 {
		return []database.User{}, 0, true
	}

	users := database.ConvertUserRows(userRows)
	return users, userRows[0].Count, true
}

// Creates a new user.
//
// @Summary Create new user
// @ID create-new-user
// @Security Optimus-IDE-CollabSessionToken
// @Accept json
// @Produce json
// @Tags Users
// @Param request body optimus-ide-collabsdk.CreateUserRequestWithOrgs true "Create user request"
// @Success 201 {object} optimus-ide-collabsdk.User
// @Router /api/v2/users [post]
func (api *API) postUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	auditor := *api.Auditor.Load()
	aReq, commitAudit := audit.InitRequest[database.User](rw, &audit.RequestParams{
		Audit:   auditor,
		Log:     api.Logger,
		Request: r,
		Action:  database.AuditActionCreate,
	})
	defer commitAudit()

	var req optimus-ide-collabsdk.CreateUserRequestWithOrgs
	if !httpapi.Read(ctx, rw, r, &req) {
		return
	}

	// Service accounts must use login_type 'none' and have no password
	// or email.
	if req.ServiceAccount {
		// The client can omit login type for a service account and it will be
		// set for them below. But if they request the wrong one, we have to let
		// them know.
		if req.UserLoginType != "" && req.UserLoginType != optimus-ide-collabsdk.LoginTypeNone {
			httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
				Message: "Service accounts must use login type 'none'.",
			})
			return
		}
		if req.Password != "" {
			httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
				Message: "Password cannot be set for service accounts.",
			})
			return
		}
		if req.Email != "" {
			httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
				Message: "Email cannot be set for service accounts.",
			})
			return
		}

		req.UserLoginType = optimus-ide-collabsdk.LoginTypeNone

		// Service accounts are a Premium feature.
		if !api.Entitlements.Enabled(optimus-ide-collabsdk.FeatureServiceAccounts) {
			httpapi.Write(ctx, rw, http.StatusForbidden, optimus-ide-collabsdk.Response{
				Message: fmt.Sprintf("%s is a Premium feature. Contact sales!", optimus-ide-collabsdk.FeatureServiceAccounts.Humanize()),
			})
			return
		}
	} else if req.UserLoginType == "" {
		// Default to password auth
		req.UserLoginType = optimus-ide-collabsdk.LoginTypePassword
	}

	if !req.ServiceAccount && req.UserLoginType == optimus-ide-collabsdk.LoginTypeNone {
		httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message: "Login type 'none' requires a service account.",
		})
		return
	}

	if req.UserLoginType != optimus-ide-collabsdk.LoginTypePassword && req.Password != "" {
		httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message: fmt.Sprintf("Password cannot be set for non-password (%q) authentication.", req.UserLoginType),
		})
		return
	}

	// If password auth is disabled, don't allow new users to be
	// created with a password!
	if api.DeploymentValues.DisablePasswordAuth && req.UserLoginType == optimus-ide-collabsdk.LoginTypePassword {
		httpapi.Write(ctx, rw, http.StatusForbidden, optimus-ide-collabsdk.Response{
			Message: "Password based authentication is disabled! Unable to provision new users with password authentication.",
		})
		return
	}

	if len(req.OrganizationIDs) == 0 {
		httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message: "No organization specified to place the user as a member of. It is required to specify at least one organization id to place the user in.",
			Detail:  "required at least 1 value for the array 'organization_ids'",
			Validations: []optimus-ide-collabsdk.ValidationError{
				{
					Field:  "organization_ids",
					Detail: "Missing values, this cannot be empty",
				},
			},
		})
		return
	}

	// TODO: @emyrk Authorize the organization create if the createUser will do that.

	_, err := api.Database.GetUserByEmailOrUsername(ctx, database.GetUserByEmailOrUsernameParams{
		Username: req.Username,
		Email:    req.Email,
	})
	if err == nil {
		httpapi.Write(ctx, rw, http.StatusConflict, optimus-ide-collabsdk.Response{
			Message: "User already exists.",
		})
		return
	}
	if !errors.Is(err, sql.ErrNoRows) {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Internal error fetching user.",
			Detail:  err.Error(),
		})
		return
	}

	// If an organization was provided, make sure it exists.
	for i, orgID := range req.OrganizationIDs {
		var orgErr error
		if orgID != uuid.Nil {
			_, orgErr = api.Database.GetOrganizationByID(ctx, orgID)
		} else {
			var defaultOrg database.Organization
			defaultOrg, orgErr = api.Database.GetDefaultOrganization(ctx)
			if orgErr == nil {
				// converts uuid.Nil --> default org.ID
				req.OrganizationIDs[i] = defaultOrg.ID
			}
		}
		if orgErr != nil {
			if httpapi.Is404Error(orgErr) {
				httpapi.Write(ctx, rw, http.StatusNotFound, optimus-ide-collabsdk.Response{
					Message: fmt.Sprintf("Organization does not exist with the provided id %q.", orgID),
				})
				return
			}

			httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
				Message: "Internal error fetching organization.",
				Detail:  orgErr.Error(),
			})
			return
		}
	}

	var loginType database.LoginType
	switch req.UserLoginType {
	case optimus-ide-collabsdk.LoginTypeNone:
		loginType = database.LoginTypeNone
	case optimus-ide-collabsdk.LoginTypePassword:
		err = userpassword.Validate(req.Password)
		if err != nil {
			httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
				Message: "Password is invalid",
				Validations: []optimus-ide-collabsdk.ValidationError{{
					Field:  "password",
					Detail: err.Error(),
				}},
			})
			return
		}
		loginType = database.LoginTypePassword
	case optimus-ide-collabsdk.LoginTypeOIDC:
		if api.OIDCConfig == nil {
			httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
				Message: "You must configure OIDC before creating OIDC users.",
			})
			return
		}
		loginType = database.LoginTypeOIDC
	case optimus-ide-collabsdk.LoginTypeGithub:
		loginType = database.LoginTypeGithub
	default:
		httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message: fmt.Sprintf("Unsupported login type %q for manually creating new users.", req.UserLoginType),
		})
		return
	}

	apiKey := httpmw.APIKey(r)

	accountCreator, err := api.Database.GetUserByID(ctx, apiKey.UserID)
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Unable to determine the details of the actor creating the account.",
		})
		return
	}

	user, err := api.CreateUser(ctx, api.Database, CreateUserRequest{
		CreateUserRequestWithOrgs: req,
		LoginType:                 loginType,
		accountCreatorName:        accountCreator.Name,
		RBACRoles:                 req.Roles,
	})

	if dbauthz.IsNotAuthorizedError(err) {
		httpapi.Write(ctx, rw, http.StatusForbidden, optimus-ide-collabsdk.Response{
			Message: "You are not authorized to create users.",
		})
		return
	}
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Internal error creating user.",
			Detail:  err.Error(),
		})
		return
	}

	aReq.New = user

	// Report when users are added!
	api.Telemetry.Report(&telemetry.Snapshot{
		Users: []telemetry.User{telemetry.ConvertUser(user)},
	})

	sdkUser := db2sdk.User(user, req.OrganizationIDs)
	api.enrichUserAISeat(ctx, &sdkUser)
	httpapi.Write(ctx, rw, http.StatusCreated, sdkUser)
}

// @Summary Delete user
// @ID delete-user
// @Security Optimus-IDE-CollabSessionToken
// @Tags Users
// @Param user path string true "User ID, name, or me"
// @Success 200
// @Router /api/v2/users/{user} [delete]
func (api *API) deleteUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	auditor := *api.Auditor.Load()
	user := httpmw.UserParam(r)
	auth := httpmw.UserAuthorization(r.Context())
	aReq, commitAudit := audit.InitRequest[database.User](rw, &audit.RequestParams{
		Audit:   auditor,
		Log:     api.Logger,
		Request: r,
		Action:  database.AuditActionDelete,
	})
	aReq.Old = user
	defer commitAudit()

	if auth.ID == user.ID.String() {
		httpapi.Write(ctx, rw, http.StatusForbidden, optimus-ide-collabsdk.Response{
			Message: "You cannot delete yourself!",
		})
		return
	}

	// This query is ONLY done to get the workspace count, so we use a system
	// context to return ALL workspaces. Not just workspaces the user can view.
	// nolint:gocritic
	workspaces, err := api.Database.GetWorkspaces(dbauthz.AsSystemRestricted(ctx), database.GetWorkspacesParams{
		OwnerID: user.ID,
	})
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Internal error fetching workspaces.",
			Detail:  err.Error(),
		})
		return
	}
	if len(workspaces) > 0 {
		httpapi.Write(ctx, rw, http.StatusExpectationFailed, optimus-ide-collabsdk.Response{
			Message: "You cannot delete a user that has workspaces. Delete their workspaces and try again!",
		})
		return
	}

	err = api.Database.UpdateUserDeletedByID(ctx, user.ID)
	if dbauthz.IsNotAuthorizedError(err) {
		httpapi.Forbidden(rw)
		return
	}
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Internal error deleting user.",
			Detail:  err.Error(),
		})
		return
	}
	user.Deleted = true
	aReq.New = user

	userAdmins, err := findUserAdmins(ctx, api.Database)
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Internal error fetching user admins.",
			Detail:  err.Error(),
		})
		return
	}

	apiKey := httpmw.APIKey(r)

	accountDeleter, err := api.Database.GetUserByID(ctx, apiKey.UserID)
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Unable to determine the details of the actor deleting the account.",
		})
		return
	}

	for _, u := range userAdmins {
		// nolint: gocritic // Need notifier actor to enqueue notifications
		if _, err := api.NotificationsEnqueuer.Enqueue(dbauthz.AsNotifier(ctx), u.ID, notifications.TemplateUserAccountDeleted,
			map[string]string{
				"deleted_account_name":      user.Username,
				"deleted_account_user_name": user.Name,
				"initiator":                 accountDeleter.Name,
			},
			"api-users-delete",
			user.ID,
		); err != nil {
			api.Logger.Warn(ctx, "unable to notify about deleted user", slog.F("deleted_user", user.Username), slog.Error(err))
		}
	}

	httpapi.Write(ctx, rw, http.StatusOK, optimus-ide-collabsdk.Response{
		Message: "User has been deleted!",
	})
}

// Returns the parameterized user requested. All validation
// is completed in the middleware for this route.
//
// @Summary Get user by name
// @ID get-user-by-name
// @Security Optimus-IDE-CollabSessionToken
// @Produce json
// @Tags Users
// @Param user path string true "User ID, username, or me"
// @Success 200 {object} optimus-ide-collabsdk.User
// @Router /api/v2/users/{user} [get]
func (api *API) userByName(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := httpmw.UserParam(r)
	organizationIDs, err := userOrganizationIDs(ctx, api, user)
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Internal error fetching user's organizations.",
			Detail:  err.Error(),
		})
		return
	}

	sdkUser := db2sdk.User(user, organizationIDs)
	api.enrichUserAISeat(ctx, &sdkUser)
	httpapi.Write(ctx, rw, http.StatusOK, sdkUser)
}

// Returns recent build parameters for the signed-in user.
//
// @Summary Get autofill build parameters for user
// @ID get-autofill-build-parameters-for-user
// @Security Optimus-IDE-CollabSessionToken
// @Produce json
// @Tags Users
// @Param user path string true "User ID, username, or me"
// @Param template_id query string true "Template ID"
// @Success 200 {array} optimus-ide-collabsdk.UserParameter
// @Router /api/v2/users/{user}/autofill-parameters [get]
func (api *API) userAutofillParameters(rw http.ResponseWriter, r *http.Request) {
	user := httpmw.UserParam(r)

	p := httpapi.NewQueryParamParser().RequiredNotEmpty("template_id")
	templateID := p.UUID(r.URL.Query(), uuid.UUID{}, "template_id")
	p.ErrorExcessParams(r.URL.Query())
	if len(p.Errors) > 0 {
		httpapi.Write(r.Context(), rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message:     "Invalid query parameters.",
			Validations: p.Errors,
		})
		return
	}

	params, err := api.Database.GetUserWorkspaceBuildParameters(
		r.Context(),
		database.GetUserWorkspaceBuildParametersParams{
			OwnerID:    user.ID,
			TemplateID: templateID,
		},
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		httpapi.Write(r.Context(), rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Internal error fetching user's parameters.",
			Detail:  err.Error(),
		})
		return
	}

	sdkParams := []optimus-ide-collabsdk.UserParameter{}
	for _, param := range params {
		sdkParams = append(sdkParams, optimus-ide-collabsdk.UserParameter{
			Name:  param.Name,
			Value: param.Value,
		})
	}

	httpapi.Write(r.Context(), rw, http.StatusOK, sdkParams)
}

// Returns the user's login type. This only works if the api key for authorization
// and the requested user match. Eg: 'me'
//
// @Summary Get user login type
// @ID get-user-login-type
// @Security Optimus-IDE-CollabSessionToken
// @Produce json
// @Tags Users
// @Param user path string true "User ID, name, or me"
// @Success 200 {object} optimus-ide-collabsdk.UserLoginType
// @Router /api/v2/users/{user}/login-type [get]
func (*API) userLoginType(rw http.ResponseWriter, r *http.Request) {
	var (
		ctx  = r.Context()
		user = httpmw.UserParam(r)
		key  = httpmw.APIKey(r)
	)

	if key.UserID != user.ID {
		// Currently this is only valid for querying yourself.
		httpapi.Write(ctx, rw, http.StatusForbidden, optimus-ide-collabsdk.Response{
			Message: "You are not authorized to view this user's login type.",
		})
		return
	}

	httpapi.Write(ctx, rw, http.StatusOK, optimus-ide-collabsdk.UserLoginType{
		LoginType: optimus-ide-collabsdk.LoginType(user.LoginType),
	})
}

// @Summary Update user profile
// @ID update-user-profile
// @Security Optimus-IDE-CollabSessionToken
// @Accept json
// @Produce json
// @Tags Users
// @Param user path string true "User ID, name, or me"
// @Param request body optimus-ide-collabsdk.UpdateUserProfileRequest true "Updated profile"
// @Success 200 {object} optimus-ide-collabsdk.User
// @Router /api/v2/users/{user}/profile [put]
func (api *API) putUserProfile(rw http.ResponseWriter, r *http.Request) {
	var (
		ctx               = r.Context()
		user              = httpmw.UserParam(r)
		auditor           = *api.Auditor.Load()
		aReq, commitAudit = audit.InitRequest[database.User](rw, &audit.RequestParams{
			Audit:   auditor,
			Log:     api.Logger,
			Request: r,
			Action:  database.AuditActionWrite,
		})
	)
	defer commitAudit()
	aReq.Old = user

	var params optimus-ide-collabsdk.UpdateUserProfileRequest
	if !httpapi.Read(ctx, rw, r, &params) {
		return
	}

	// If caller wants to update user's username, they need "update_users" permission.
	// This is restricted to user admins only.
	if params.Username != user.Username && !api.Authorize(r, policy.ActionUpdate, user) {
		httpapi.ResourceNotFound(rw)
		return
	}

	existentUser, err := api.Database.GetUserByEmailOrUsername(ctx, database.GetUserByEmailOrUsernameParams{
		Username: params.Username,
	})
	isDifferentUser := existentUser.ID != user.ID

	if err == nil && isDifferentUser {
		responseErrors := []optimus-ide-collabsdk.ValidationError{{
			Field:  "username",
			Detail: "This username is already in use.",
		}}
		httpapi.Write(ctx, rw, http.StatusConflict, optimus-ide-collabsdk.Response{
			Message:     "A user with this username already exists.",
			Validations: responseErrors,
		})
		return
	}
	if !errors.Is(err, sql.ErrNoRows) && isDifferentUser {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Internal error fetching user.",
			Detail:  err.Error(),
		})
		return
	}

	// Avatars for password and none login types are managed manually. For
	// other login types the avatar is synced from the identity provider on
	// login, so we preserve the existing value and ignore any submitted one.
	avatarURL := user.AvatarURL
	if user.LoginType == database.LoginTypePassword || user.LoginType == database.LoginTypeNone {
		avatarURL = params.AvatarURL
	}

	updatedUserProfile, err := api.Database.UpdateUserProfile(ctx, database.UpdateUserProfileParams{
		ID:        user.ID,
		Email:     user.Email,
		Name:      params.Name,
		AvatarURL: avatarURL,
		Username:  params.Username,
		UpdatedAt: dbtime.Now(),
	})
	aReq.New = updatedUserProfile

	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Internal error updating user.",
			Detail:  err.Error(),
		})
		return
	}

	organizationIDs, err := userOrganizationIDs(ctx, api, user)
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Internal error fetching user's organizations.",
			Detail:  err.Error(),
		})
		return
	}

	sdkUser := db2sdk.User(updatedUserProfile, organizationIDs)
	api.enrichUserAISeat(ctx, &sdkUser)
	httpapi.Write(ctx, rw, http.StatusOK, sdkUser)
}

// @Summary Suspend user account
// @ID suspend-user-account
// @Security Optimus-IDE-CollabSessionToken
// @Produce json
// @Tags Users
// @Param user path string true "User ID, name, or me"
// @Success 200 {object} optimus-ide-collabsdk.User
// @Router /api/v2/users/{user}/status/suspend [put]
func (api *API) putSuspendUserAccount() func(rw http.ResponseWriter, r *http.Request) {
	return api.putUserStatus(database.UserStatusSuspended)
}

// @Summary Activate user account
// @ID activate-user-account
// @Security Optimus-IDE-CollabSessionToken
// @Produce json
// @Tags Users
// @Param user path string true "User ID, name, or me"
// @Success 200 {object} optimus-ide-collabsdk.User
// @Router /api/v2/users/{user}/status/activate [put]
func (api *API) putActivateUserAccount() func(rw http.ResponseWriter, r *http.Request) {
	return api.putUserStatus(database.UserStatusActive)
}

func (api *API) putUserStatus(status database.UserStatus) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		var (
			ctx               = r.Context()
			user              = httpmw.UserParam(r)
			apiKey            = httpmw.APIKey(r)
			auditor           = *api.Auditor.Load()
			aReq, commitAudit = audit.InitRequest[database.User](rw, &audit.RequestParams{
				Audit:   auditor,
				Log:     api.Logger,
				Request: r,
				Action:  database.AuditActionWrite,
			})
		)
		defer commitAudit()
		aReq.Old = user

		if status == database.UserStatusSuspended {
			// There are some manual protections when suspending a user to
			// prevent certain situations.
			switch {
			case user.ID == apiKey.UserID:
				// Suspending yourself is not allowed, as you can lock yourself
				// out of the system.
				httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
					Message: "You cannot suspend yourself.",
				})
				return
			case slice.Contains(user.RBACRoles, rbac.RoleOwner().String()):
				// You may not suspend an owner
				httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
					Message: fmt.Sprintf("You cannot suspend a user with the %q role. You must remove the role first.", rbac.RoleOwner()),
				})
				return
			}
		}

		actingUser, err := api.Database.GetUserByID(ctx, apiKey.UserID)
		if err != nil {
			httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
				Message: "Unable to determine the details of the actor creating the account.",
			})
			return
		}

		targetUser, err := api.Database.UpdateUserStatus(ctx, database.UpdateUserStatusParams{
			ID:         user.ID,
			Status:     status,
			UpdatedAt:  dbtime.Now(),
			UserIsSeen: false,
		})
		if err != nil {
			httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
				Message: fmt.Sprintf("Internal error updating user's status to %q.", status),
				Detail:  err.Error(),
			})
			return
		}
		aReq.New = targetUser

		err = api.notifyUserStatusChanged(ctx, actingUser.Name, user, status)
		if err != nil {
			api.Logger.Warn(ctx, "unable to notify about changed user's status", slog.F("affected_user", user.Username), slog.Error(err))
		}

		organizations, err := userOrganizationIDs(ctx, api, user)
		if err != nil {
			httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
				Message: "Internal error fetching user's organizations.",
				Detail:  err.Error(),
			})
			return
		}

		sdkUser := db2sdk.User(targetUser, organizations)
		api.enrichUserAISeat(ctx, &sdkUser)
		httpapi.Write(ctx, rw, http.StatusOK, sdkUser)
	}
}

func (api *API) notifyUserStatusChanged(ctx context.Context, actingUserName string, targetUser database.User, status database.UserStatus) error {
	var labels map[string]string
	var data map[string]any
	var adminTemplateID, personalTemplateID uuid.UUID
	switch status {
	case database.UserStatusSuspended:
		labels = map[string]string{
			"suspended_account_name":      targetUser.Username,
			"suspended_account_user_name": targetUser.Name,
			"initiator":                   actingUserName,
		}
		data = map[string]any{
			"user": map[string]any{"id": targetUser.ID, "name": targetUser.Name, "email": targetUser.Email},
		}
		adminTemplateID = notifications.TemplateUserAccountSuspended
		personalTemplateID = notifications.TemplateYourAccountSuspended
	case database.UserStatusActive:
		labels = map[string]string{
			"activated_account_name":      targetUser.Username,
			"activated_account_user_name": targetUser.Name,
			"initiator":                   actingUserName,
		}
		data = map[string]any{
			"user": map[string]any{"id": targetUser.ID, "name": targetUser.Name, "email": targetUser.Email},
		}
		adminTemplateID = notifications.TemplateUserAccountActivated
		personalTemplateID = notifications.TemplateYourAccountActivated
	default:
		api.Logger.Error(ctx, "user status is not supported", slog.F("username", targetUser.Username), slog.F("user_status", string(status)))
		return xerrors.Errorf("unable to notify admins as the user's status is unsupported")
	}

	userAdmins, err := findUserAdmins(ctx, api.Database)
	if err != nil {
		api.Logger.Error(ctx, "unable to find user admins", slog.Error(err))
	}

	// Send notifications to user admins and affected user
	for _, u := range userAdmins {
		// nolint:gocritic // Need notifier actor to enqueue notifications
		if _, err := api.NotificationsEnqueuer.EnqueueWithData(dbauthz.AsNotifier(ctx), u.ID, adminTemplateID,
			labels, data, "api-put-user-status",
			targetUser.ID,
		); err != nil {
			api.Logger.Warn(ctx, "unable to notify about changed user's status", slog.F("affected_user", targetUser.Username), slog.Error(err))
		}
	}
	// nolint:gocritic // Need notifier actor to enqueue notifications
	if _, err := api.NotificationsEnqueuer.EnqueueWithData(dbauthz.AsNotifier(ctx), targetUser.ID, personalTemplateID,
		labels, data, "api-put-user-status",
		targetUser.ID,
	); err != nil {
		api.Logger.Warn(ctx, "unable to notify user about status change of their account", slog.F("affected_user", targetUser.Username), slog.Error(err))
	}
	return nil
}

// @Summary Get user appearance settings
// @ID get-user-appearance-settings
// @Security Optimus-IDE-CollabSessionToken
// @Produce json
// @Tags Users
// @Param user path string true "User ID, name, or me"
// @Success 200 {object} optimus-ide-collabsdk.UserAppearanceSettings
// @Router /api/v2/users/{user}/appearance [get]
func (api *API) userAppearanceSettings(rw http.ResponseWriter, r *http.Request) {
	var (
		ctx  = r.Context()
		user = httpmw.UserParam(r)
	)

	settings, err := api.Database.GetUserAppearanceSettings(ctx, user.ID)
	if err != nil {
		writeUserSettingsReadError(ctx, rw, err)
		return
	}

	httpapi.Write(ctx, rw, http.StatusOK, userAppearanceSettingsFromRow(settings))
}

func userAppearanceSettingsFromRow(settings database.GetUserAppearanceSettingsRow) optimus-ide-collabsdk.UserAppearanceSettings {
	return optimus-ide-collabsdk.UserAppearanceSettings{
		ThemePreference: settings.ThemePreference,
		ThemeMode:       optimus-ide-collabsdk.ThemeMode(settings.ThemeMode),
		ThemeLight:      settings.ThemeLight,
		ThemeDark:       settings.ThemeDark,
		TerminalFont:    optimus-ide-collabsdk.TerminalFontName(settings.TerminalFont),
	}
}

func isLegacyAutoThemePreference(themePreference string) bool {
	switch themePreference {
	case "auto", "auto-protan-deuter", "auto-tritan":
		return true
	default:
		return false
	}
}

func writeUserSettingsReadError(ctx context.Context, rw http.ResponseWriter, err error) {
	httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
		Message: "Error reading user settings.",
		Detail:  err.Error(),
	})
}

// @Summary Update user appearance settings
// @ID update-user-appearance-settings
// @Security Optimus-IDE-CollabSessionToken
// @Accept json
// @Produce json
// @Tags Users
// @Param user path string true "User ID, name, or me"
// @Param request body optimus-ide-collabsdk.UpdateUserAppearanceSettingsRequest true "New appearance settings"
// @Success 200 {object} optimus-ide-collabsdk.UserAppearanceSettings
// @Router /api/v2/users/{user}/appearance [put]
func (api *API) putUserAppearanceSettings(rw http.ResponseWriter, r *http.Request) {
	var (
		ctx  = r.Context()
		user = httpmw.UserParam(r)
	)

	var params optimus-ide-collabsdk.UpdateUserAppearanceSettingsRequest
	if !httpapi.Read(ctx, rw, r, &params) {
		return
	}

	if !isValidFontName(params.TerminalFont) {
		httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message: "Unsupported font family.",
		})
		return
	}

	// theme_mode is optional for backward compatibility. Older CLI
	// clients do not know about theme_mode or the sync slots, so an
	// omitted mode must leave those fields untouched instead of replacing
	// them with single-mode defaults. Legacy auto values are the exception:
	// the old UI used them to mean sync-with-system, so clearing theme_mode
	// lets modern clients migrate them on read.
	themeModeProvided := params.ThemeMode != optimus-ide-collabsdk.ThemeModeUnset
	updateThemeMode := themeModeProvided
	isSyncMode := params.ThemeMode == optimus-ide-collabsdk.ThemeModeSync
	isSingleMode := params.ThemeMode == optimus-ide-collabsdk.ThemeModeSingle
	updateThemeLight := isSyncMode || (isSingleMode && params.ThemeLight != "")
	updateThemeDark := isSyncMode || (isSingleMode && params.ThemeDark != "")
	themeMode := params.ThemeMode
	if !updateThemeMode && isLegacyAutoThemePreference(params.ThemePreference) {
		updateThemeMode = true
		themeMode = optimus-ide-collabsdk.ThemeModeUnset
	}

	var updatedSettings database.GetUserAppearanceSettingsRow

	err := api.Database.InTx(func(tx database.Store) error {
		_, err := tx.UpdateUserThemePreference(ctx, database.UpdateUserThemePreferenceParams{
			UserID:          user.ID,
			ThemePreference: params.ThemePreference,
		})
		if err != nil {
			return xerrors.Errorf("update user theme preference: %w", err)
		}

		if updateThemeMode {
			_, err = tx.UpdateUserThemeMode(ctx, database.UpdateUserThemeModeParams{
				UserID:    user.ID,
				ThemeMode: string(themeMode),
			})
			if err != nil {
				return xerrors.Errorf("update user theme mode: %w", err)
			}
		}

		if updateThemeLight {
			_, err = tx.UpdateUserThemeLight(ctx, database.UpdateUserThemeLightParams{
				UserID:     user.ID,
				ThemeLight: params.ThemeLight,
			})
			if err != nil {
				return xerrors.Errorf("update user theme light: %w", err)
			}
		}

		if updateThemeDark {
			_, err = tx.UpdateUserThemeDark(ctx, database.UpdateUserThemeDarkParams{
				UserID:    user.ID,
				ThemeDark: params.ThemeDark,
			})
			if err != nil {
				return xerrors.Errorf("update user theme dark: %w", err)
			}
		}

		_, err = tx.UpdateUserTerminalFont(ctx, database.UpdateUserTerminalFontParams{
			UserID:       user.ID,
			TerminalFont: string(params.TerminalFont),
		})
		if err != nil {
			return xerrors.Errorf("update user terminal font: %w", err)
		}

		updatedSettings, err = tx.GetUserAppearanceSettings(ctx, user.ID)
		if err != nil {
			return xerrors.Errorf("get updated user appearance settings: %w", err)
		}

		return nil
	}, nil)
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Internal error updating user appearance settings.",
			Detail:  err.Error(),
		})
		return
	}

	httpapi.Write(ctx, rw, http.StatusOK, userAppearanceSettingsFromRow(updatedSettings))
}

// @Summary Get user preference settings
// @ID get-user-preference-settings
// @Security Optimus-IDE-CollabSessionToken
// @Produce json
// @Tags Users
// @Param user path string true "User ID, name, or me"
// @Success 200 {object} optimus-ide-collabsdk.UserPreferenceSettings
// @Router /api/v2/users/{user}/preferences [get]
func (api *API) userPreferenceSettings(rw http.ResponseWriter, r *http.Request) {
	var (
		ctx  = r.Context()
		user = httpmw.UserParam(r)
	)

	taskAlertDismissed, err := api.Database.GetUserTaskNotificationAlertDismissed(ctx, user.ID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
				Message: "Error reading user preference settings.",
				Detail:  err.Error(),
			})
			return
		}
	}

	thinkingMode, err := api.Database.GetUserThinkingDisplayMode(ctx, user.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Error reading user preference settings.",
			Detail:  err.Error(),
		})
		return
	}

	shellToolMode, err := api.Database.GetUserShellToolDisplayMode(ctx, user.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Error reading user preference settings.",
			Detail:  err.Error(),
		})
		return
	}

	codeDiffMode, err := api.Database.GetUserCodeDiffDisplayMode(ctx, user.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Error reading user preference settings.",
			Detail:  err.Error(),
		})
		return
	}

	agentChatSendShortcut, err := api.Database.GetUserAgentChatSendShortcut(ctx, user.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Error reading user preference settings.",
			Detail:  err.Error(),
		})
		return
	}

	httpapi.Write(ctx, rw, http.StatusOK, optimus-ide-collabsdk.UserPreferenceSettings{
		TaskNotificationAlertDismissed: taskAlertDismissed,
		ThinkingDisplayMode:            sanitizeThinkingDisplayMode(thinkingMode),
		ShellToolDisplayMode:           sanitizeShellToolDisplayMode(shellToolMode),
		CodeDiffDisplayMode:            sanitizeAgentDisplayMode(codeDiffMode),
		AgentChatSendShortcut:          sanitizeAgentChatSendShortcut(agentChatSendShortcut),
	})
}

// @Summary Update user preference settings
// @ID update-user-preference-settings
// @Security Optimus-IDE-CollabSessionToken
// @Accept json
// @Produce json
// @Tags Users
// @Param user path string true "User ID, name, or me"
// @Param request body optimus-ide-collabsdk.UpdateUserPreferenceSettingsRequest true "New preference settings"
// @Success 200 {object} optimus-ide-collabsdk.UserPreferenceSettings
// @Router /api/v2/users/{user}/preferences [put]
func (api *API) putUserPreferenceSettings(rw http.ResponseWriter, r *http.Request) {
	var (
		ctx  = r.Context()
		user = httpmw.UserParam(r)
	)

	var params optimus-ide-collabsdk.UpdateUserPreferenceSettingsRequest
	if !httpapi.Read(ctx, rw, r, &params) {
		return
	}

	if params.ThinkingDisplayMode != "" &&
		!slices.Contains(optimus-ide-collabsdk.ValidThinkingDisplayModes, params.ThinkingDisplayMode) {
		httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message: "Invalid thinking display mode.",
			Validations: []optimus-ide-collabsdk.ValidationError{
				{Field: "thinking_display_mode", Detail: thinkingDisplayModeValidationDetail},
			},
		})
		return
	}
	if params.ShellToolDisplayMode != "" &&
		!slices.Contains(optimus-ide-collabsdk.ValidAgentDisplayModes, params.ShellToolDisplayMode) {
		httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message: "Invalid shell tool display mode.",
			Validations: []optimus-ide-collabsdk.ValidationError{
				{Field: "shell_tool_display_mode", Detail: agentDisplayModeValidationDetail},
			},
		})
		return
	}
	if params.CodeDiffDisplayMode != "" &&
		!slices.Contains(optimus-ide-collabsdk.ValidAgentDisplayModes, params.CodeDiffDisplayMode) {
		httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message: "Invalid code diff display mode.",
			Validations: []optimus-ide-collabsdk.ValidationError{
				{Field: "code_diff_display_mode", Detail: agentDisplayModeValidationDetail},
			},
		})
		return
	}
	if params.AgentChatSendShortcut != "" &&
		!slices.Contains(optimus-ide-collabsdk.ValidAgentChatSendShortcuts, params.AgentChatSendShortcut) {
		httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message: "Invalid agent chat send shortcut.",
			Validations: []optimus-ide-collabsdk.ValidationError{
				{Field: "agent_chat_send_shortcut", Detail: agentChatSendShortcutValidationDetail},
			},
		})
		return
	}
	var settings optimus-ide-collabsdk.UserPreferenceSettings
	err := api.Database.InTx(func(tx database.Store) error {
		var err error
		if params.TaskNotificationAlertDismissed != nil {
			settings.TaskNotificationAlertDismissed, err = tx.UpdateUserTaskNotificationAlertDismissed(ctx, database.UpdateUserTaskNotificationAlertDismissedParams{
				UserID:                         user.ID,
				TaskNotificationAlertDismissed: *params.TaskNotificationAlertDismissed,
			})
			if err != nil {
				return newUserPreferenceSettingsAPIError("Internal error updating user task notification alert dismissed.", err)
			}
		} else {
			settings.TaskNotificationAlertDismissed, err = tx.GetUserTaskNotificationAlertDismissed(ctx, user.ID)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				return newUserPreferenceSettingsAPIError("Error reading task notification alert dismissed.", err)
			}
		}

		if params.ThinkingDisplayMode != "" {
			updated, err := tx.UpdateUserThinkingDisplayMode(ctx, database.UpdateUserThinkingDisplayModeParams{
				UserID:              user.ID,
				ThinkingDisplayMode: string(params.ThinkingDisplayMode),
			})
			if err != nil {
				return newUserPreferenceSettingsAPIError("Internal error updating thinking display mode.", err)
			}
			settings.ThinkingDisplayMode = sanitizeThinkingDisplayMode(updated)
		} else {
			stored, err := tx.GetUserThinkingDisplayMode(ctx, user.ID)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				return newUserPreferenceSettingsAPIError("Error reading thinking display mode.", err)
			}
			settings.ThinkingDisplayMode = sanitizeThinkingDisplayMode(stored)
		}

		if params.ShellToolDisplayMode != "" {
			updated, err := tx.UpdateUserShellToolDisplayMode(ctx, database.UpdateUserShellToolDisplayModeParams{
				UserID:               user.ID,
				ShellToolDisplayMode: string(params.ShellToolDisplayMode),
			})
			if err != nil {
				return newUserPreferenceSettingsAPIError("Internal error updating shell tool display mode.", err)
			}
			settings.ShellToolDisplayMode = sanitizeShellToolDisplayMode(updated)
		} else {
			stored, err := tx.GetUserShellToolDisplayMode(ctx, user.ID)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				return newUserPreferenceSettingsAPIError("Error reading shell tool display mode.", err)
			}
			settings.ShellToolDisplayMode = sanitizeShellToolDisplayMode(stored)
		}

		if params.CodeDiffDisplayMode != "" {
			updated, err := tx.UpdateUserCodeDiffDisplayMode(ctx, database.UpdateUserCodeDiffDisplayModeParams{
				UserID:              user.ID,
				CodeDiffDisplayMode: string(params.CodeDiffDisplayMode),
			})
			if err != nil {
				return newUserPreferenceSettingsAPIError("Internal error updating code diff display mode.", err)
			}
			settings.CodeDiffDisplayMode = sanitizeAgentDisplayMode(updated)
		} else {
			stored, err := tx.GetUserCodeDiffDisplayMode(ctx, user.ID)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				return newUserPreferenceSettingsAPIError("Error reading code diff display mode.", err)
			}
			settings.CodeDiffDisplayMode = sanitizeAgentDisplayMode(stored)
		}

		if params.AgentChatSendShortcut != "" {
			updated, err := tx.UpdateUserAgentChatSendShortcut(ctx, database.UpdateUserAgentChatSendShortcutParams{
				UserID:                user.ID,
				AgentChatSendShortcut: string(params.AgentChatSendShortcut),
			})
			if err != nil {
				return newUserPreferenceSettingsAPIError("Internal error updating agent chat send shortcut.", err)
			}
			settings.AgentChatSendShortcut = sanitizeAgentChatSendShortcut(updated)
		} else {
			stored, err := tx.GetUserAgentChatSendShortcut(ctx, user.ID)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				return newUserPreferenceSettingsAPIError("Error reading agent chat send shortcut.", err)
			}
			settings.AgentChatSendShortcut = sanitizeAgentChatSendShortcut(stored)
		}
		return nil
	}, database.DefaultTXOptions().WithID("user_preference_settings"))
	if err != nil {
		var apiErr userPreferenceSettingsAPIError
		if errors.As(err, &apiErr) {
			httpapi.Write(ctx, rw, apiErr.statusCode, apiErr.response)
			return
		}
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Internal error updating user preference settings.",
			Detail:  err.Error(),
		})
		return
	}

	httpapi.Write(ctx, rw, http.StatusOK, settings)
}

type userPreferenceSettingsAPIError struct {
	statusCode int
	response   optimus-ide-collabsdk.Response
	err        error
}

func newUserPreferenceSettingsAPIError(message string, err error) userPreferenceSettingsAPIError {
	return userPreferenceSettingsAPIError{
		statusCode: http.StatusInternalServerError,
		response: optimus-ide-collabsdk.Response{
			Message: message,
			Detail:  err.Error(),
		},
		err: err,
	}
}

func (e userPreferenceSettingsAPIError) Error() string {
	return fmt.Sprintf("%s: %s", e.response.Message, e.err)
}

func (e userPreferenceSettingsAPIError) Unwrap() error {
	return e.err
}

const (
	thinkingDisplayModeValidationDetail   = "must be one of: auto, preview, always_expanded, always_collapsed"
	agentDisplayModeValidationDetail      = "must be one of: auto, always_expanded, always_collapsed"
	agentChatSendShortcutValidationDetail = "must be one of: enter, modifier_enter"
)

func sanitizeThinkingDisplayMode(raw string) optimus-ide-collabsdk.ThinkingDisplayMode {
	mode := optimus-ide-collabsdk.ThinkingDisplayMode(raw)
	if slices.Contains(optimus-ide-collabsdk.ValidThinkingDisplayModes, mode) {
		return mode
	}
	return optimus-ide-collabsdk.ThinkingDisplayModeAuto
}

func sanitizeShellToolDisplayMode(raw string) optimus-ide-collabsdk.AgentDisplayMode {
	mode := sanitizeAgentDisplayMode(raw)
	if mode == "" {
		return optimus-ide-collabsdk.AgentDisplayModeAlwaysCollapsed
	}
	return mode
}

func sanitizeAgentDisplayMode(raw string) optimus-ide-collabsdk.AgentDisplayMode {
	mode := optimus-ide-collabsdk.AgentDisplayMode(raw)
	if slices.Contains(optimus-ide-collabsdk.ValidAgentDisplayModes, mode) {
		return mode
	}
	return ""
}

func sanitizeAgentChatSendShortcut(raw string) optimus-ide-collabsdk.AgentChatSendShortcut {
	shortcut := optimus-ide-collabsdk.AgentChatSendShortcut(raw)
	if slices.Contains(optimus-ide-collabsdk.ValidAgentChatSendShortcuts, shortcut) {
		return shortcut
	}
	return optimus-ide-collabsdk.AgentChatSendShortcutEnter
}

func isValidFontName(font optimus-ide-collabsdk.TerminalFontName) bool {
	return slices.Contains(optimus-ide-collabsdk.TerminalFontNames, font)
}

// @Summary Update user password
// @ID update-user-password
// @Security Optimus-IDE-CollabSessionToken
// @Accept json
// @Tags Users
// @Param user path string true "User ID, name, or me"
// @Param request body optimus-ide-collabsdk.UpdateUserPasswordRequest true "Update password request"
// @Success 204
// @Router /api/v2/users/{user}/password [put]
func (api *API) putUserPassword(rw http.ResponseWriter, r *http.Request) {
	var (
		ctx               = r.Context()
		user              = httpmw.UserParam(r)
		params            optimus-ide-collabsdk.UpdateUserPasswordRequest
		apiKey            = httpmw.APIKey(r)
		auditor           = *api.Auditor.Load()
		aReq, commitAudit = audit.InitRequest[database.User](rw, &audit.RequestParams{
			Audit:   auditor,
			Log:     api.Logger,
			Request: r,
			Action:  database.AuditActionWrite,
		})
	)
	defer commitAudit()
	aReq.Old = user

	if !api.Authorize(r, policy.ActionUpdatePersonal, user) {
		httpapi.ResourceNotFound(rw)
		return
	}

	// Only owners can change the password of another owner.
	if apiKey.UserID != user.ID && slices.Contains(user.RBACRoles, rbac.RoleOwner().String()) {
		actingUser, err := api.Database.GetUserByID(ctx, apiKey.UserID)
		if err != nil {
			httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
				Message: "Internal error fetching acting user.",
				Detail:  err.Error(),
			})
			return
		}
		if !slices.Contains(actingUser.RBACRoles, rbac.RoleOwner().String()) {
			httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
				Message: "Only owners can change the password of an owner.",
			})
			return
		}
	}

	if !httpapi.Read(ctx, rw, r, &params) {
		return
	}

	if user.LoginType != database.LoginTypePassword {
		httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message: "Users without password login type cannot change their password.",
		})
		return
	}

	// A user need to put its own password to update it
	if apiKey.UserID == user.ID && params.OldPassword == "" {
		httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message: "Old password is required.",
		})
		return
	}

	err := userpassword.Validate(params.Password)
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message: "Invalid password.",
			Validations: []optimus-ide-collabsdk.ValidationError{
				{
					Field:  "password",
					Detail: err.Error(),
				},
			},
		})
		return
	}

	if params.OldPassword != "" {
		// if they send something let's validate it
		ok, err := userpassword.Compare(string(user.HashedPassword), params.OldPassword)
		if err != nil {
			httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
				Message: "Internal error with passwords.",
				Detail:  err.Error(),
			})
			return
		}
		if !ok {
			httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
				Message: "Old password is incorrect.",
				Validations: []optimus-ide-collabsdk.ValidationError{
					{
						Field:  "old_password",
						Detail: "Old password is incorrect.",
					},
				},
			})
			return
		}
	}

	// Prevent users reusing their old password.
	if match, _ := userpassword.Compare(string(user.HashedPassword), params.Password); match {
		httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message: "New password cannot match old password.",
		})
		return
	}

	hashedPassword, err := userpassword.Hash(params.Password)
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Internal error hashing new password.",
			Detail:  err.Error(),
		})
		return
	}

	err = api.Database.InTx(func(tx database.Store) error {
		err = tx.UpdateUserHashedPassword(ctx, database.UpdateUserHashedPasswordParams{
			ID:             user.ID,
			HashedPassword: []byte(hashedPassword),
		})
		if err != nil {
			return xerrors.Errorf("update user hashed password: %w", err)
		}

		//nolint:gocritic // Password resets must revoke all keys owned by the
		// target user, not just keys addressable by the caller's actor.
		err = tx.DeleteAPIKeysByUserID(dbauthz.AsAPIKeyRevoker(ctx, user.ID), user.ID)
		if err != nil {
			return xerrors.Errorf("delete api keys by user ID: %w", err)
		}

		return nil
	}, nil)
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Internal error updating user's password.",
			Detail:  err.Error(),
		})
		return
	}

	newUser := user
	newUser.HashedPassword = []byte(hashedPassword)
	aReq.New = newUser

	rw.WriteHeader(http.StatusNoContent)
}

// @Summary Get user roles
// @ID get-user-roles
// @Security Optimus-IDE-CollabSessionToken
// @Produce json
// @Tags Users
// @Param user path string true "User ID, name, or me"
// @Success 200 {object} optimus-ide-collabsdk.User
// @Router /api/v2/users/{user}/roles [get]
func (api *API) userRoles(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := httpmw.UserParam(r)

	if !api.Authorize(r, policy.ActionReadPersonal, user) {
		httpapi.ResourceNotFound(rw)
		return
	}

	// TODO: Replace this with "GetAuthorizationUserRoles"
	resp := optimus-ide-collabsdk.UserRoles{
		Roles:             user.RBACRoles,
		OrganizationRoles: make(map[uuid.UUID][]string),
	}

	memberships, err := api.Database.OrganizationMembers(ctx, database.OrganizationMembersParams{
		UserID:         user.ID,
		OrganizationID: uuid.Nil,
		IncludeSystem:  false,
		GithubUserID:   0,
	})
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Internal error fetching user's organization memberships.",
			Detail:  err.Error(),
		})
		return
	}

	for _, mem := range memberships {
		resp.OrganizationRoles[mem.OrganizationMember.OrganizationID] = mem.OrganizationMember.Roles
	}

	httpapi.Write(ctx, rw, http.StatusOK, resp)
}

// @Summary Assign role to user
// @ID assign-role-to-user
// @Security Optimus-IDE-CollabSessionToken
// @Accept json
// @Produce json
// @Tags Users
// @Param user path string true "User ID, name, or me"
// @Param request body optimus-ide-collabsdk.UpdateRoles true "Update roles request"
// @Success 200 {object} optimus-ide-collabsdk.User
// @Router /api/v2/users/{user}/roles [put]
func (api *API) putUserRoles(rw http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		// User is the user to modify.
		user              = httpmw.UserParam(r)
		apiKey            = httpmw.APIKey(r)
		auditor           = *api.Auditor.Load()
		aReq, commitAudit = audit.InitRequest[database.User](rw, &audit.RequestParams{
			Audit:   auditor,
			Log:     api.Logger,
			Request: r,
			Action:  database.AuditActionWrite,
		})
	)
	defer commitAudit()
	aReq.Old = user

	if user.LoginType == database.LoginTypeOIDC && api.IDPSync.SiteRoleSyncEnabled() {
		httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message: "Cannot modify roles for OIDC users when role sync is enabled.",
			Detail:  "'User Role Field' is set in the OIDC configuration. All role changes must come from the oidc identity provider.",
		})
		return
	}

	if apiKey.UserID == user.ID {
		httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message: "You cannot change your own roles.",
		})
		return
	}

	var params optimus-ide-collabsdk.UpdateRoles
	if !httpapi.Read(ctx, rw, r, &params) {
		return
	}

	updatedUser, err := api.Database.UpdateUserRoles(ctx, database.UpdateUserRolesParams{
		GrantedRoles: params.Roles,
		ID:           user.ID,
	})
	if dbauthz.IsNotAuthorizedError(err) {
		httpapi.Forbidden(rw)
		return
	}
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message: err.Error(),
		})
		return
	}
	aReq.New = updatedUser

	organizationIDs, err := userOrganizationIDs(ctx, api, user)
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Internal error fetching user's organizations.",
			Detail:  err.Error(),
		})
		return
	}

	sdkUser := db2sdk.User(updatedUser, organizationIDs)
	api.enrichUserAISeat(ctx, &sdkUser)
	httpapi.Write(ctx, rw, http.StatusOK, sdkUser)
}

// Returns organizations the parameterized user has access to.
//
// @Summary Get organizations by user
// @ID get-organizations-by-user
// @Security Optimus-IDE-CollabSessionToken
// @Produce json
// @Tags Users
// @Param user path string true "User ID, name, or me"
// @Success 200 {array} optimus-ide-collabsdk.Organization
// @Router /api/v2/users/{user}/organizations [get]
func (api *API) organizationsByUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := httpmw.UserParam(r)

	organizations, err := api.Database.GetOrganizationsByUserID(ctx, database.GetOrganizationsByUserIDParams{
		UserID:  user.ID,
		Deleted: sql.NullBool{Bool: false, Valid: true},
	})
	if errors.Is(err, sql.ErrNoRows) {
		err = nil
		organizations = []database.Organization{}
	}
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Internal error fetching user's organizations.",
			Detail:  err.Error(),
		})
		return
	}

	// Only return orgs the user can read.
	organizations, err = AuthorizeFilter(api.HTTPAuth, r, policy.ActionRead, organizations)
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Internal error fetching organizations.",
			Detail:  err.Error(),
		})
		return
	}

	httpapi.Write(ctx, rw, http.StatusOK, slice.List(organizations, db2sdk.Organization))
}

// @Summary Get organization by user and organization name
// @ID get-organization-by-user-and-organization-name
// @Security Optimus-IDE-CollabSessionToken
// @Produce json
// @Tags Users
// @Param user path string true "User ID, name, or me"
// @Param organizationname path string true "Organization name"
// @Success 200 {object} optimus-ide-collabsdk.Organization
// @Router /api/v2/users/{user}/organizations/{organizationname} [get]
func (api *API) organizationByUserAndName(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	organizationName := chi.URLParam(r, "organizationname")
	organization, err := api.Database.GetOrganizationByName(ctx, database.GetOrganizationByNameParams{
		Name:    organizationName,
		Deleted: false,
	})
	if httpapi.Is404Error(err) {
		httpapi.ResourceNotFound(rw)
		return
	}
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Internal error fetching organization.",
			Detail:  err.Error(),
		})
		return
	}

	httpapi.Write(ctx, rw, http.StatusOK, db2sdk.Organization(organization))
}

type CreateUserRequest struct {
	optimus-ide-collabsdk.CreateUserRequestWithOrgs
	LoginType          database.LoginType
	SkipNotifications  bool
	accountCreatorName string
	RBACRoles          []string
}

func (api *API) CreateUser(ctx context.Context, store database.Store, req CreateUserRequest) (database.User, error) {
	// Ensure the username is valid. It's the caller's responsibility to ensure
	// the username is valid and unique.
	if usernameValid := optimus-ide-collabsdk.NameValid(req.Username); usernameValid != nil {
		return database.User{}, xerrors.Errorf("invalid username %q: %w", req.Username, usernameValid)
	}

	// If the caller didn't specify rbac roles, default to
	// a member of the site.
	rbacRoles := []string{}
	if req.RBACRoles != nil {
		rbacRoles = req.RBACRoles
	}

	var user database.User
	err := store.InTx(func(tx database.Store) error {
		orgRoles := make([]string, 0)

		status := ""
		if req.UserStatus != nil {
			status = string(*req.UserStatus)
		}
		params := database.InsertUserParams{
			ID:               uuid.New(),
			Email:            req.Email,
			Username:         req.Username,
			Name:             optimus-ide-collabsdk.NormalizeRealUsername(req.Name),
			CreatedAt:        dbtime.Now(),
			UpdatedAt:        dbtime.Now(),
			HashedPassword:   []byte{},
			RBACRoles:        rbacRoles,
			LoginType:        req.LoginType,
			Status:           status,
			IsServiceAccount: req.ServiceAccount,
		}
		// If a user signs up with OAuth, they can have no password!
		if req.Password != "" {
			hashedPassword, err := userpassword.Hash(req.Password)
			if err != nil {
				return xerrors.Errorf("hash password: %w", err)
			}
			params.HashedPassword = []byte(hashedPassword)
		}

		var err error
		user, err = tx.InsertUser(ctx, params)
		if err != nil {
			return xerrors.Errorf("create user: %w", err)
		}

		privateKey, publicKey, err := gitsshkey.Generate(api.SSHKeygenAlgorithm)
		if err != nil {
			return xerrors.Errorf("generate user gitsshkey: %w", err)
		}
		_, err = tx.InsertGitSSHKey(ctx, database.InsertGitSSHKeyParams{
			UserID:          user.ID,
			CreatedAt:       dbtime.Now(),
			UpdatedAt:       dbtime.Now(),
			PrivateKey:      privateKey,
			PrivateKeyKeyID: sql.NullString{}, // dbcrypt will set as required
			PublicKey:       publicKey,
		})
		if err != nil {
			return xerrors.Errorf("insert user gitsshkey: %w", err)
		}

		for _, orgID := range req.OrganizationIDs {
			_, err = tx.InsertOrganizationMember(ctx, database.InsertOrganizationMemberParams{
				OrganizationID: orgID,
				UserID:         user.ID,
				CreatedAt:      dbtime.Now(),
				UpdatedAt:      dbtime.Now(),
				// By default give them membership to the organization.
				Roles: orgRoles,
			})
			if err != nil {
				return xerrors.Errorf("create organization member for %q: %w", orgID.String(), err)
			}
		}

		return nil
	}, nil)
	if err != nil || req.SkipNotifications {
		return user, err
	}

	userAdmins, err := findUserAdmins(ctx, store)
	if err != nil {
		return user, xerrors.Errorf("find user admins: %w", err)
	}

	for _, u := range userAdmins {
		if u.ID == user.ID {
			// If the new user is an admin, don't notify them about themselves.
			continue
		}
		if _, err := api.NotificationsEnqueuer.EnqueueWithData(
			// nolint:gocritic // Need notifier actor to enqueue notifications
			dbauthz.AsNotifier(ctx),
			u.ID,
			notifications.TemplateUserAccountCreated,
			map[string]string{
				"created_account_name":      user.Username,
				"created_account_user_name": user.Name,
				"initiator":                 req.accountCreatorName,
			},
			map[string]any{
				"user": map[string]any{"id": user.ID, "name": user.Name, "email": user.Email},
			},
			"api-users-create",
			user.ID,
		); err != nil {
			api.Logger.Warn(ctx, "unable to notify about created user", slog.F("created_user", user.Username), slog.Error(err))
		}
	}

	return user, err
}

// findUserAdmins fetches all users with user admin permission including owners.
func findUserAdmins(ctx context.Context, store database.Store) ([]database.GetUsersRow, error) {
	userAdmins, err := store.GetUsers(ctx, database.GetUsersParams{
		RbacRole: []string{optimus-ide-collabsdk.RoleOwner, optimus-ide-collabsdk.RoleUserAdmin},
	})
	if err != nil {
		return nil, xerrors.Errorf("get owners: %w", err)
	}
	return userAdmins, nil
}

// enrichUserAISeat sets HasAISeat on the user when the feature is entitled.
func (api *API) enrichUserAISeat(ctx context.Context, user *optimus-ide-collabsdk.User) {
	if !api.Entitlements.Enabled(optimus-ide-collabsdk.FeatureAIGovernanceUserLimit) {
		return
	}

	//nolint:gocritic // AI seat state is a system-level read gated by entitlement.
	aiSeatUserIDs, err := api.Database.GetUserAISeatStates(
		dbauthz.AsSystemRestricted(ctx),
		[]uuid.UUID{user.ID},
	)
	if err != nil {
		if !xerrors.Is(err, sql.ErrNoRows) {
			api.Logger.Warn(
				ctx,
				"failed to fetch AI seat state for user",
				slog.F("user_id", user.ID),
				slog.Error(err),
			)
		}
		return
	}

	user.HasAISeat = len(aiSeatUserIDs) > 0
}

func convertUsers(users []database.User, organizationIDsByUserID map[uuid.UUID][]uuid.UUID, aiSeatSet map[uuid.UUID]struct{}) []optimus-ide-collabsdk.User {
	converted := make([]optimus-ide-collabsdk.User, 0, len(users))
	for _, u := range users {
		userOrganizationIDs := organizationIDsByUserID[u.ID]
		_, hasAISeat := aiSeatSet[u.ID]
		convertedUser := db2sdk.User(u, userOrganizationIDs)
		convertedUser.HasAISeat = hasAISeat
		converted = append(converted, convertedUser)
	}
	return converted
}

func userOrganizationIDs(ctx context.Context, api *API, user database.User) ([]uuid.UUID, error) {
	organizationIDsByMemberIDsRows, err := api.Database.GetOrganizationIDsByMemberIDs(ctx, []uuid.UUID{user.ID})
	if err != nil {
		return []uuid.UUID{}, err
	}

	// If you are in no orgs, then return an empty list.
	if len(organizationIDsByMemberIDsRows) == 0 {
		return []uuid.UUID{}, nil
	}

	member := organizationIDsByMemberIDsRows[0]
	return member.OrganizationIDs, nil
}

func convertAPIKey(k database.APIKey) optimus-ide-collabsdk.APIKey {
	// Derive a single legacy scope name for response compatibility.
	// Historically, the API exposed only two scope strings: "all" and
	// "application_connect". Continue to return those for clients even
	// though the database stores canonical values (e.g. "optimus-ide-collab:all")
	// and may include low-level scopes.
	var legacyScope optimus-ide-collabsdk.APIKeyScope
	if k.Scopes.Has(database.ApiKeyScopeOptimus-IDE-CollabApplicationConnect) {
		legacyScope = optimus-ide-collabsdk.APIKeyScopeApplicationConnect
	} else if k.Scopes.Has(database.ApiKeyScopeOptimus-IDE-CollabAll) {
		legacyScope = optimus-ide-collabsdk.APIKeyScopeAll
	}

	scopes := make([]optimus-ide-collabsdk.APIKeyScope, 0, len(k.Scopes))
	for _, s := range k.Scopes {
		scopes = append(scopes, optimus-ide-collabsdk.APIKeyScope(s))
	}

	return optimus-ide-collabsdk.APIKey{
		ID:              k.ID,
		UserID:          k.UserID,
		LastUsed:        k.LastUsed,
		ExpiresAt:       k.ExpiresAt,
		CreatedAt:       k.CreatedAt,
		UpdatedAt:       k.UpdatedAt,
		LoginType:       optimus-ide-collabsdk.LoginType(k.LoginType),
		Scope:           legacyScope,
		Scopes:          scopes,
		LifetimeSeconds: k.LifetimeSeconds,
		TokenName:       k.TokenName,
		AllowList:       slice.List(k.AllowList, db2sdk.APIAllowListTarget),
	}
}
