// Package db2sdk provides common conversion routines from database types to optimus-ide-collabsdk types
package db2sdk

import (
	"cmp"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/hcl/v2"
	"github.com/sqlc-dev/pqtype"
	"golang.org/x/xerrors"
	"tailscale.com/tailcfg"

	agentproto "github.com/optimus-ide-collab/optimus-ide-collab/v2/agent/proto"
	aibridgeutils "github.com/optimus-ide-collab/optimus-ide-collab/v2/aibridge/utils"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbtime"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/externalauth/gitprovider"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac/policy"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/render"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/util/ptr"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/util/slice"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/workspaceapps/appurl"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/x/chatd/chatprompt"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/provisionersdk/proto"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/tailnet"
	previewtypes "github.com/optimus-ide-collab/preview/types"
)

func APIAllowListTarget(entry rbac.AllowListElement) optimus-ide-collabsdk.APIAllowListTarget {
	return optimus-ide-collabsdk.APIAllowListTarget{
		Type: optimus-ide-collabsdk.RBACResource(entry.Type),
		ID:   entry.ID,
	}
}

// AIProvider converts a database row plus its API keys into the
// optimus-ide-collabsdk shape. The caller is responsible for ensuring the row and
// keys have been decrypted (i.e. fetched through the dbcrypt-wrapped
// store). Each api_key is masked via aibridge utils.MaskSecret and
// write-only fields on Settings are stripped, so the result is safe
// to echo back in API responses.
func AIProvider(row database.AIProvider, keys []database.AIProviderKey) (optimus-ide-collabsdk.AIProvider, error) {
	display := row.Name
	if row.DisplayName.Valid && row.DisplayName.String != "" {
		display = row.DisplayName.String
	}
	out := optimus-ide-collabsdk.AIProvider{
		ID:          row.ID,
		Type:        optimus-ide-collabsdk.AIProviderType(row.Type),
		Name:        row.Name,
		DisplayName: display,
		Icon:        row.Icon,
		Enabled:     row.Enabled,
		BaseURL:     row.BaseUrl,
		APIKeys:     maskAIProviderKeys(keys),
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
	s, err := AIProviderSettings(row.Settings)
	if err != nil {
		return optimus-ide-collabsdk.AIProvider{}, xerrors.Errorf("decode settings: %w", err)
	}
	out.Settings = redactAIProviderSettings(s)
	return out, nil
}

// AIProviderSettings parses the on-disk JSON form back into a optimus-ide-collabsdk
// settings value. SQL NULL and the empty string decode to the zero
// value.
func AIProviderSettings(col sql.NullString) (optimus-ide-collabsdk.AIProviderSettings, error) {
	if !col.Valid || col.String == "" {
		return optimus-ide-collabsdk.AIProviderSettings{}, nil
	}
	var s optimus-ide-collabsdk.AIProviderSettings
	if err := json.Unmarshal([]byte(col.String), &s); err != nil {
		return optimus-ide-collabsdk.AIProviderSettings{}, err
	}
	return s, nil
}

// maskAIProviderKeys converts the supplied database rows into the
// public-facing AIProviderKey shape, preserving order. Plaintext is
// replaced by a non-reversible mask (see aibridgeutils.MaskSecret) so
// the result is safe to embed in API responses.
func maskAIProviderKeys(keys []database.AIProviderKey) []optimus-ide-collabsdk.AIProviderKey {
	out := make([]optimus-ide-collabsdk.AIProviderKey, 0, len(keys))
	for _, k := range keys {
		out = append(out, optimus-ide-collabsdk.AIProviderKey{
			ID:        k.ID,
			Masked:    aibridgeutils.MaskSecret(k.APIKey),
			CreatedAt: k.CreatedAt,
		})
	}
	return out
}

// redactAIProviderSettings strips write-only fields from a settings
// value so it can be safely echoed back in API responses.
func redactAIProviderSettings(s optimus-ide-collabsdk.AIProviderSettings) optimus-ide-collabsdk.AIProviderSettings {
	out := s
	if out.Bedrock != nil {
		// Deep-copy so we don't mutate the caller's struct.
		b := *out.Bedrock
		b.AccessKey = nil
		b.AccessKeySecret = nil
		out.Bedrock = &b
	}
	return out
}

type ExternalAuthMeta struct {
	Authenticated bool
	ValidateError string
}

func ExternalAuths(auths []database.ExternalAuthLink, meta map[string]ExternalAuthMeta) []optimus-ide-collabsdk.ExternalAuthLink {
	out := make([]optimus-ide-collabsdk.ExternalAuthLink, 0, len(auths))
	for _, auth := range auths {
		out = append(out, ExternalAuth(auth, meta[auth.ProviderID]))
	}
	return out
}

func ExternalAuth(auth database.ExternalAuthLink, meta ExternalAuthMeta) optimus-ide-collabsdk.ExternalAuthLink {
	return optimus-ide-collabsdk.ExternalAuthLink{
		ProviderID:      auth.ProviderID,
		CreatedAt:       auth.CreatedAt,
		UpdatedAt:       auth.UpdatedAt,
		HasRefreshToken: auth.OAuthRefreshToken != "",
		Expires:         auth.OAuthExpiry,
		Authenticated:   meta.Authenticated,
		ValidateError:   meta.ValidateError,
	}
}

func WorkspaceBuildParameter(p database.WorkspaceBuildParameter) optimus-ide-collabsdk.WorkspaceBuildParameter {
	return optimus-ide-collabsdk.WorkspaceBuildParameter{
		Name:  p.Name,
		Value: p.Value,
	}
}

func WorkspaceBuildParameters(params []database.WorkspaceBuildParameter) []optimus-ide-collabsdk.WorkspaceBuildParameter {
	return slice.List(params, WorkspaceBuildParameter)
}

func TemplateVersionParameters(params []database.TemplateVersionParameter) ([]optimus-ide-collabsdk.TemplateVersionParameter, error) {
	out := make([]optimus-ide-collabsdk.TemplateVersionParameter, 0, len(params))
	for _, p := range params {
		np, err := TemplateVersionParameter(p)
		if err != nil {
			return nil, xerrors.Errorf("convert template version parameter %q: %w", p.Name, err)
		}
		out = append(out, np)
	}

	return out, nil
}

func TemplateVersionParameterFromPreview(param previewtypes.Parameter) (optimus-ide-collabsdk.TemplateVersionParameter, error) {
	descriptionPlaintext, err := render.PlaintextFromMarkdown(param.Description)
	if err != nil {
		return optimus-ide-collabsdk.TemplateVersionParameter{}, err
	}

	sdkParam := optimus-ide-collabsdk.TemplateVersionParameter{
		Name:                 param.Name,
		DisplayName:          param.DisplayName,
		Description:          param.Description,
		DescriptionPlaintext: descriptionPlaintext,
		Type:                 string(param.Type),
		FormType:             string(param.FormType),
		Mutable:              param.Mutable,
		DefaultValue:         param.DefaultValue.AsString(),
		Icon:                 param.Icon,
		Required:             param.Required,
		Ephemeral:            param.Ephemeral,
		Options:              slice.List(param.Options, TemplateVersionParameterOptionFromPreview),
		// Validation set after
	}
	if len(param.Validations) > 0 {
		validation := param.Validations[0]
		sdkParam.ValidationError = validation.Error
		if validation.Monotonic != nil {
			sdkParam.ValidationMonotonic = optimus-ide-collabsdk.ValidationMonotonicOrder(*validation.Monotonic)
		}
		if validation.Regex != nil {
			sdkParam.ValidationRegex = *validation.Regex
		}
		if validation.Min != nil {
			//nolint:gosec // No other choice
			sdkParam.ValidationMin = ptr.Ref(int32(*validation.Min))
		}
		if validation.Max != nil {
			//nolint:gosec // No other choice
			sdkParam.ValidationMax = ptr.Ref(int32(*validation.Max))
		}
	}

	return sdkParam, nil
}

func TemplateVersionParameter(param database.TemplateVersionParameter) (optimus-ide-collabsdk.TemplateVersionParameter, error) {
	options, err := templateVersionParameterOptions(param.Options)
	if err != nil {
		return optimus-ide-collabsdk.TemplateVersionParameter{}, err
	}

	descriptionPlaintext, err := render.PlaintextFromMarkdown(param.Description)
	if err != nil {
		return optimus-ide-collabsdk.TemplateVersionParameter{}, err
	}

	var validationMin *int32
	if param.ValidationMin.Valid {
		validationMin = &param.ValidationMin.Int32
	}

	var validationMax *int32
	if param.ValidationMax.Valid {
		validationMax = &param.ValidationMax.Int32
	}

	return optimus-ide-collabsdk.TemplateVersionParameter{
		Name:                 param.Name,
		DisplayName:          param.DisplayName,
		Description:          param.Description,
		DescriptionPlaintext: descriptionPlaintext,
		Type:                 param.Type,
		FormType:             string(param.FormType),
		Mutable:              param.Mutable,
		DefaultValue:         param.DefaultValue,
		Icon:                 param.Icon,
		Options:              options,
		ValidationRegex:      param.ValidationRegex,
		ValidationMin:        validationMin,
		ValidationMax:        validationMax,
		ValidationError:      param.ValidationError,
		ValidationMonotonic:  optimus-ide-collabsdk.ValidationMonotonicOrder(param.ValidationMonotonic),
		Required:             param.Required,
		Ephemeral:            param.Ephemeral,
	}, nil
}

func MinimalUser(user database.User) optimus-ide-collabsdk.MinimalUser {
	return optimus-ide-collabsdk.MinimalUser{
		ID:        user.ID,
		Username:  user.Username,
		Name:      user.Name,
		AvatarURL: user.AvatarURL,
	}
}

func MinimalUserFromVisibleUser(user database.VisibleUser) optimus-ide-collabsdk.MinimalUser {
	return optimus-ide-collabsdk.MinimalUser{
		ID:        user.ID,
		Username:  user.Username,
		Name:      user.Name,
		AvatarURL: user.AvatarURL,
	}
}

func ReducedUser(user database.User) optimus-ide-collabsdk.ReducedUser {
	return optimus-ide-collabsdk.ReducedUser{
		MinimalUser:      MinimalUser(user),
		Email:            user.Email,
		CreatedAt:        user.CreatedAt,
		UpdatedAt:        user.UpdatedAt,
		LastSeenAt:       user.LastSeenAt,
		Status:           optimus-ide-collabsdk.UserStatus(user.Status),
		LoginType:        optimus-ide-collabsdk.LoginType(user.LoginType),
		IsServiceAccount: user.IsServiceAccount,
	}
}

func UserFromGroupMember(member database.GroupMember) database.User {
	return database.User{
		ID:                 member.UserID,
		Email:              member.UserEmail,
		Username:           member.UserUsername,
		HashedPassword:     member.UserHashedPassword,
		CreatedAt:          member.UserCreatedAt,
		UpdatedAt:          member.UserUpdatedAt,
		Status:             member.UserStatus,
		RBACRoles:          member.UserRbacRoles,
		LoginType:          member.UserLoginType,
		AvatarURL:          member.UserAvatarUrl,
		Deleted:            member.UserDeleted,
		LastSeenAt:         member.UserLastSeenAt,
		QuietHoursSchedule: member.UserQuietHoursSchedule,
		Name:               member.UserName,
		GithubComUserID:    member.UserGithubComUserID,
		IsServiceAccount:   member.UserIsServiceAccount,
	}
}

func ReducedUserFromGroupMember(member database.GroupMember) optimus-ide-collabsdk.ReducedUser {
	return ReducedUser(UserFromGroupMember(member))
}

func ReducedUsersFromGroupMembers(members []database.GroupMember) []optimus-ide-collabsdk.ReducedUser {
	return slice.List(members, ReducedUserFromGroupMember)
}

func UserFromGroupMemberRow(member database.GetGroupMembersByGroupIDPaginatedRow) database.User {
	return database.User{
		ID:                 member.UserID,
		Email:              member.UserEmail,
		Username:           member.UserUsername,
		HashedPassword:     member.UserHashedPassword,
		CreatedAt:          member.UserCreatedAt,
		UpdatedAt:          member.UserUpdatedAt,
		Status:             member.UserStatus,
		RBACRoles:          member.UserRbacRoles,
		LoginType:          member.UserLoginType,
		AvatarURL:          member.UserAvatarUrl,
		Deleted:            member.UserDeleted,
		LastSeenAt:         member.UserLastSeenAt,
		QuietHoursSchedule: member.UserQuietHoursSchedule,
		Name:               member.UserName,
		GithubComUserID:    member.UserGithubComUserID,
		IsServiceAccount:   member.UserIsServiceAccount,
	}
}

func ReducedUserFromGroupMemberRow(member database.GetGroupMembersByGroupIDPaginatedRow) optimus-ide-collabsdk.ReducedUser {
	return ReducedUser(UserFromGroupMemberRow(member))
}

func ReducedUsersFromGroupMemberRows(members []database.GetGroupMembersByGroupIDPaginatedRow) []optimus-ide-collabsdk.ReducedUser {
	return slice.List(members, ReducedUserFromGroupMemberRow)
}

func ReducedUsers(users []database.User) []optimus-ide-collabsdk.ReducedUser {
	return slice.List(users, ReducedUser)
}

func User(user database.User, organizationIDs []uuid.UUID) optimus-ide-collabsdk.User {
	convertedUser := optimus-ide-collabsdk.User{
		ReducedUser:     ReducedUser(user),
		OrganizationIDs: organizationIDs,
		Roles:           SlimRolesFromNames(user.RBACRoles),
	}

	return convertedUser
}

func Users(users []database.User, organizationIDs map[uuid.UUID][]uuid.UUID) []optimus-ide-collabsdk.User {
	return slice.List(users, func(user database.User) optimus-ide-collabsdk.User {
		return User(user, organizationIDs[user.ID])
	})
}

func Group(row database.GetGroupsRow, members []database.GroupMember, totalMemberCount int) optimus-ide-collabsdk.Group {
	return optimus-ide-collabsdk.Group{
		ID:                      row.Group.ID,
		Name:                    row.Group.Name,
		DisplayName:             row.Group.DisplayName,
		OrganizationID:          row.Group.OrganizationID,
		AvatarURL:               row.Group.AvatarURL,
		Members:                 ReducedUsersFromGroupMembers(members),
		TotalMemberCount:        totalMemberCount,
		QuotaAllowance:          int(row.Group.QuotaAllowance),
		Source:                  optimus-ide-collabsdk.GroupSource(row.Group.Source),
		OrganizationName:        row.OrganizationName,
		OrganizationDisplayName: row.OrganizationDisplayName,
	}
}

func TemplateInsightsParameters(parameterRows []database.GetTemplateParameterInsightsRow) ([]optimus-ide-collabsdk.TemplateParameterUsage, error) {
	// Use a stable sort, similarly to how we would sort in the query, note that
	// we don't sort in the query because order varies depending on the table
	// collation.
	//
	// ORDER BY utp.name, utp.type, utp.display_name, utp.description, utp.options, wbp.value
	slices.SortFunc(parameterRows, func(a, b database.GetTemplateParameterInsightsRow) int {
		if a.Name != b.Name {
			return strings.Compare(a.Name, b.Name)
		}
		if a.Type != b.Type {
			return strings.Compare(a.Type, b.Type)
		}
		if a.DisplayName != b.DisplayName {
			return strings.Compare(a.DisplayName, b.DisplayName)
		}
		if a.Description != b.Description {
			return strings.Compare(a.Description, b.Description)
		}
		if string(a.Options) != string(b.Options) {
			return strings.Compare(string(a.Options), string(b.Options))
		}
		return strings.Compare(a.Value, b.Value)
	})

	parametersUsage := []optimus-ide-collabsdk.TemplateParameterUsage{}
	indexByNum := make(map[int64]int)
	for _, param := range parameterRows {
		if _, ok := indexByNum[param.Num]; !ok {
			var opts []optimus-ide-collabsdk.TemplateVersionParameterOption
			err := json.Unmarshal(param.Options, &opts)
			if err != nil {
				return nil, err
			}

			plaintextDescription, err := render.PlaintextFromMarkdown(param.Description)
			if err != nil {
				return nil, err
			}

			parametersUsage = append(parametersUsage, optimus-ide-collabsdk.TemplateParameterUsage{
				TemplateIDs: param.TemplateIDs,
				Name:        param.Name,
				Type:        param.Type,
				DisplayName: param.DisplayName,
				Description: plaintextDescription,
				Options:     opts,
			})
			indexByNum[param.Num] = len(parametersUsage) - 1
		}

		i := indexByNum[param.Num]
		parametersUsage[i].Values = append(parametersUsage[i].Values, optimus-ide-collabsdk.TemplateParameterValue{
			Value: param.Value,
			Count: param.Count,
		})
	}

	return parametersUsage, nil
}

func templateVersionParameterOptions(rawOptions json.RawMessage) ([]optimus-ide-collabsdk.TemplateVersionParameterOption, error) {
	var protoOptions []*proto.RichParameterOption
	err := json.Unmarshal(rawOptions, &protoOptions)
	if err != nil {
		return nil, err
	}

	options := make([]optimus-ide-collabsdk.TemplateVersionParameterOption, 0)
	for _, option := range protoOptions {
		options = append(options, optimus-ide-collabsdk.TemplateVersionParameterOption{
			Name:        option.Name,
			Description: option.Description,
			Value:       option.Value,
			Icon:        option.Icon,
		})
	}
	return options, nil
}

func TemplateVersionParameterOptionFromPreview(option *previewtypes.ParameterOption) optimus-ide-collabsdk.TemplateVersionParameterOption {
	return optimus-ide-collabsdk.TemplateVersionParameterOption{
		Name:        option.Name,
		Description: option.Description,
		Value:       option.Value.AsString(),
		Icon:        option.Icon,
	}
}

func OAuth2ProviderApp(accessURL *url.URL, dbApp database.OAuth2ProviderApp) optimus-ide-collabsdk.OAuth2ProviderApp {
	return optimus-ide-collabsdk.OAuth2ProviderApp{
		ID:          dbApp.ID,
		Name:        dbApp.Name,
		CallbackURL: dbApp.CallbackURL,
		Icon:        dbApp.Icon,
		Endpoints: optimus-ide-collabsdk.OAuth2AppEndpoints{
			Authorization: accessURL.ResolveReference(&url.URL{
				Path: "/oauth2/authorize",
			}).String(),
			Token: accessURL.ResolveReference(&url.URL{
				Path: "/oauth2/tokens",
			}).String(),
			// We do not currently support DeviceAuth.
			DeviceAuth: "",
			TokenRevoke: accessURL.ResolveReference(&url.URL{
				Path: "/oauth2/revoke",
			}).String(),
		},
	}
}

func OAuth2ProviderApps(accessURL *url.URL, dbApps []database.OAuth2ProviderApp) []optimus-ide-collabsdk.OAuth2ProviderApp {
	return slice.List(dbApps, func(dbApp database.OAuth2ProviderApp) optimus-ide-collabsdk.OAuth2ProviderApp {
		return OAuth2ProviderApp(accessURL, dbApp)
	})
}

func convertDisplayApps(apps []database.DisplayApp) []optimus-ide-collabsdk.DisplayApp {
	dapps := make([]optimus-ide-collabsdk.DisplayApp, 0, len(apps))
	for _, app := range apps {
		switch optimus-ide-collabsdk.DisplayApp(app) {
		case optimus-ide-collabsdk.DisplayAppVSCodeDesktop, optimus-ide-collabsdk.DisplayAppVSCodeInsiders, optimus-ide-collabsdk.DisplayAppPortForward, optimus-ide-collabsdk.DisplayAppWebTerminal, optimus-ide-collabsdk.DisplayAppSSH:
			dapps = append(dapps, optimus-ide-collabsdk.DisplayApp(app))
		}
	}

	return dapps
}

func WorkspaceAgentEnvironment(workspaceAgent database.WorkspaceAgent) (map[string]string, error) {
	var envs map[string]string
	if workspaceAgent.EnvironmentVariables.Valid {
		err := json.Unmarshal(workspaceAgent.EnvironmentVariables.RawMessage, &envs)
		if err != nil {
			return nil, xerrors.Errorf("unmarshal environment variables: %w", err)
		}
	}

	return envs, nil
}

func WorkspaceAgent(derpMap *tailcfg.DERPMap, coordinator tailnet.Coordinator,
	dbAgent database.WorkspaceAgent, apps []optimus-ide-collabsdk.WorkspaceApp, scripts []optimus-ide-collabsdk.WorkspaceAgentScript, logSources []optimus-ide-collabsdk.WorkspaceAgentLogSource,
	agentInactiveDisconnectTimeout time.Duration, agentFallbackTroubleshootingURL string,
) (optimus-ide-collabsdk.WorkspaceAgent, error) {
	envs, err := WorkspaceAgentEnvironment(dbAgent)
	if err != nil {
		return optimus-ide-collabsdk.WorkspaceAgent{}, err
	}
	troubleshootingURL := agentFallbackTroubleshootingURL
	if dbAgent.TroubleshootingURL != "" {
		troubleshootingURL = dbAgent.TroubleshootingURL
	}
	subsystems := make([]optimus-ide-collabsdk.AgentSubsystem, len(dbAgent.Subsystems))
	for i, subsystem := range dbAgent.Subsystems {
		subsystems[i] = optimus-ide-collabsdk.AgentSubsystem(subsystem)
	}

	legacyStartupScriptBehavior := optimus-ide-collabsdk.WorkspaceAgentStartupScriptBehaviorNonBlocking
	for _, script := range scripts {
		if !script.RunOnStart {
			continue
		}
		if !script.StartBlocksLogin {
			continue
		}
		legacyStartupScriptBehavior = optimus-ide-collabsdk.WorkspaceAgentStartupScriptBehaviorBlocking
	}

	workspaceAgent := optimus-ide-collabsdk.WorkspaceAgent{
		ID:                       dbAgent.ID,
		ParentID:                 dbAgent.ParentID,
		CreatedAt:                dbAgent.CreatedAt,
		UpdatedAt:                dbAgent.UpdatedAt,
		ResourceID:               dbAgent.ResourceID,
		InstanceID:               dbAgent.AuthInstanceID.String,
		Name:                     dbAgent.Name,
		Architecture:             dbAgent.Architecture,
		OperatingSystem:          dbAgent.OperatingSystem,
		Scripts:                  scripts,
		StartupScriptBehavior:    legacyStartupScriptBehavior,
		LogsLength:               dbAgent.LogsLength,
		LogsOverflowed:           dbAgent.LogsOverflowed,
		LogSources:               logSources,
		Version:                  dbAgent.Version,
		APIVersion:               dbAgent.APIVersion,
		EnvironmentVariables:     envs,
		Directory:                dbAgent.Directory,
		ExpandedDirectory:        dbAgent.ExpandedDirectory,
		Apps:                     apps,
		ConnectionTimeoutSeconds: dbAgent.ConnectionTimeoutSeconds,
		TroubleshootingURL:       troubleshootingURL,
		LifecycleState:           optimus-ide-collabsdk.WorkspaceAgentLifecycle(dbAgent.LifecycleState),
		Subsystems:               subsystems,
		DisplayApps:              convertDisplayApps(dbAgent.DisplayApps),
	}
	node := coordinator.Node(dbAgent.ID)
	if node != nil {
		workspaceAgent.DERPLatency = map[string]optimus-ide-collabsdk.DERPRegion{}
		for rawRegion, latency := range node.DERPLatency {
			regionIDStr, _, _ := strings.Cut(rawRegion, "-")
			regionID, err := strconv.Atoi(regionIDStr)
			if err != nil {
				return optimus-ide-collabsdk.WorkspaceAgent{}, xerrors.Errorf("convert derp region id %q: %w", rawRegion, err)
			}
			region, found := derpMap.Regions[regionID]
			if !found {
				// It's possible that a workspace agent is using an old DERPMap
				// and reports regions that do not exist. If that's the case,
				// report the region as unknown!
				region = &tailcfg.DERPRegion{
					RegionID:   regionID,
					RegionName: fmt.Sprintf("Unnamed %d", regionID),
				}
			}
			workspaceAgent.DERPLatency[region.RegionName] = optimus-ide-collabsdk.DERPRegion{
				Preferred:           node.PreferredDERP == regionID,
				LatencyMilliseconds: latency * 1000,
			}
		}
	}

	status := dbAgent.Status(dbtime.Now(), agentInactiveDisconnectTimeout)
	workspaceAgent.Status = optimus-ide-collabsdk.WorkspaceAgentStatus(status.Status)
	workspaceAgent.FirstConnectedAt = status.FirstConnectedAt
	workspaceAgent.LastConnectedAt = status.LastConnectedAt
	workspaceAgent.DisconnectedAt = status.DisconnectedAt

	if dbAgent.StartedAt.Valid {
		workspaceAgent.StartedAt = &dbAgent.StartedAt.Time
	}
	if dbAgent.ReadyAt.Valid {
		workspaceAgent.ReadyAt = &dbAgent.ReadyAt.Time
	}

	switch {
	case workspaceAgent.Status != optimus-ide-collabsdk.WorkspaceAgentConnected && workspaceAgent.LifecycleState == optimus-ide-collabsdk.WorkspaceAgentLifecycleOff:
		workspaceAgent.Health.Reason = "agent is not running"
	case workspaceAgent.Status == optimus-ide-collabsdk.WorkspaceAgentConnecting:
		// Note: the case above catches connecting+off as "not running".
		// This case handles connecting agents with a non-off lifecycle
		// (e.g. "created" or "starting"), where the agent binary has
		// not yet established a connection to optimus-ide-collabd.
		workspaceAgent.Health.Reason = "agent has not yet connected"
	case workspaceAgent.Status == optimus-ide-collabsdk.WorkspaceAgentTimeout:
		workspaceAgent.Health.Reason = "agent is taking too long to connect"
	case workspaceAgent.Status == optimus-ide-collabsdk.WorkspaceAgentDisconnected:
		workspaceAgent.Health.Reason = "agent has lost connection"
	// Note: We could also handle optimus-ide-collabsdk.WorkspaceAgentLifecycleStartTimeout
	// here, but it's more of a soft issue, so we don't want to mark the agent
	// as unhealthy.
	case workspaceAgent.LifecycleState == optimus-ide-collabsdk.WorkspaceAgentLifecycleStartError:
		workspaceAgent.Health.Reason = "agent startup script exited with an error"
	case workspaceAgent.LifecycleState.ShuttingDown():
		workspaceAgent.Health.Reason = "agent is shutting down"
	default:
		workspaceAgent.Health.Healthy = true
	}

	return workspaceAgent, nil
}

func AppSubdomain(dbApp database.WorkspaceApp, agentName, workspaceName, ownerName string) string {
	if !dbApp.Subdomain || agentName == "" || ownerName == "" || workspaceName == "" {
		return ""
	}

	appSlug := dbApp.Slug
	if appSlug == "" {
		appSlug = dbApp.DisplayName
	}

	// Agent name is optional when app slug is present
	normalizedAgentName := agentName
	if !appurl.PortRegex.MatchString(appSlug) {
		normalizedAgentName = ""
	}

	return appurl.ApplicationURL{
		// We never generate URLs with a prefix. We only allow prefixes when
		// parsing URLs from the hostname. Users that want this feature can
		// write out their own URLs.
		Prefix:        "",
		AppSlugOrPort: appSlug,
		AgentName:     normalizedAgentName,
		WorkspaceName: workspaceName,
		Username:      ownerName,
	}.String()
}

func Apps(dbApps []database.WorkspaceApp, statuses []database.WorkspaceAppStatus, agent database.WorkspaceAgent, ownerName string, workspace database.WorkspaceTable) []optimus-ide-collabsdk.WorkspaceApp {
	slices.SortFunc(dbApps, func(a, b database.WorkspaceApp) int {
		return cmp.Or(
			cmp.Compare(a.DisplayOrder, b.DisplayOrder),
			cmp.Compare(a.DisplayName, b.DisplayName),
			cmp.Compare(a.Slug, b.Slug),
		)
	})

	statusesByAppID := map[uuid.UUID][]database.WorkspaceAppStatus{}
	for _, status := range statuses {
		statusesByAppID[status.AppID] = append(statusesByAppID[status.AppID], status)
	}

	apps := make([]optimus-ide-collabsdk.WorkspaceApp, 0)
	for _, dbApp := range dbApps {
		statuses := statusesByAppID[dbApp.ID]
		apps = append(apps, optimus-ide-collabsdk.WorkspaceApp{
			ID:            dbApp.ID,
			URL:           dbApp.Url.String,
			External:      dbApp.External,
			Slug:          dbApp.Slug,
			DisplayName:   dbApp.DisplayName,
			Command:       dbApp.Command.String,
			Icon:          dbApp.Icon,
			Subdomain:     dbApp.Subdomain,
			SubdomainName: AppSubdomain(dbApp, agent.Name, workspace.Name, ownerName),
			SharingLevel:  optimus-ide-collabsdk.WorkspaceAppSharingLevel(dbApp.SharingLevel),
			Healthcheck: optimus-ide-collabsdk.Healthcheck{
				URL:       dbApp.HealthcheckUrl,
				Interval:  dbApp.HealthcheckInterval,
				Threshold: dbApp.HealthcheckThreshold,
			},
			Health:   optimus-ide-collabsdk.WorkspaceAppHealth(dbApp.Health),
			Group:    dbApp.DisplayGroup.String,
			Hidden:   dbApp.Hidden,
			OpenIn:   optimus-ide-collabsdk.WorkspaceAppOpenIn(dbApp.OpenIn),
			Tooltip:  dbApp.Tooltip,
			Statuses: WorkspaceAppStatuses(statuses),
		})
	}
	return apps
}

func WorkspaceAppStatuses(statuses []database.WorkspaceAppStatus) []optimus-ide-collabsdk.WorkspaceAppStatus {
	return slice.List(statuses, WorkspaceAppStatus)
}

func WorkspaceAppStatus(status database.WorkspaceAppStatus) optimus-ide-collabsdk.WorkspaceAppStatus {
	return optimus-ide-collabsdk.WorkspaceAppStatus{
		ID:          status.ID,
		CreatedAt:   status.CreatedAt,
		WorkspaceID: status.WorkspaceID,
		AgentID:     status.AgentID,
		AppID:       status.AppID,
		URI:         status.Uri.String,
		Message:     status.Message,
		State:       optimus-ide-collabsdk.WorkspaceAppStatusState(status.State),
	}
}

func ProvisionerJobLog(log database.ProvisionerJobLog) optimus-ide-collabsdk.ProvisionerJobLog {
	return optimus-ide-collabsdk.ProvisionerJobLog{
		ID:        log.ID,
		CreatedAt: log.CreatedAt,
		Source:    optimus-ide-collabsdk.LogSource(log.Source),
		Level:     optimus-ide-collabsdk.LogLevel(log.Level),
		Stage:     log.Stage,
		Output:    log.Output,
	}
}

func WorkspaceAgentLog(log database.WorkspaceAgentLog) optimus-ide-collabsdk.WorkspaceAgentLog {
	return optimus-ide-collabsdk.WorkspaceAgentLog{
		ID:        log.ID,
		CreatedAt: log.CreatedAt,
		Output:    log.Output,
		Level:     optimus-ide-collabsdk.LogLevel(log.Level),
		SourceID:  log.LogSourceID,
	}
}

func WorkspaceAgentScript(dbScript database.GetWorkspaceAgentScriptsByAgentIDsRow) optimus-ide-collabsdk.WorkspaceAgentScript {
	script := optimus-ide-collabsdk.WorkspaceAgentScript{
		ID:               dbScript.ID,
		LogPath:          dbScript.LogPath,
		LogSourceID:      dbScript.LogSourceID,
		Script:           dbScript.Script,
		Cron:             dbScript.Cron,
		RunOnStart:       dbScript.RunOnStart,
		RunOnStop:        dbScript.RunOnStop,
		StartBlocksLogin: dbScript.StartBlocksLogin,
		Timeout:          time.Duration(dbScript.TimeoutSeconds) * time.Second,
		DisplayName:      dbScript.DisplayName,
		ExitCode:         nullInt32Ptr(dbScript.ExitCode),
	}
	if dbScript.Status.Valid {
		status := optimus-ide-collabsdk.WorkspaceAgentScriptStatus(dbScript.Status.WorkspaceAgentScriptTimingStatus)
		script.Status = &status
	}
	return script
}

func ProvisionerDaemon(dbDaemon database.ProvisionerDaemon) optimus-ide-collabsdk.ProvisionerDaemon {
	result := optimus-ide-collabsdk.ProvisionerDaemon{
		ID:             dbDaemon.ID,
		OrganizationID: dbDaemon.OrganizationID,
		CreatedAt:      dbDaemon.CreatedAt,
		LastSeenAt:     optimus-ide-collabsdk.NullTime{NullTime: dbDaemon.LastSeenAt},
		Name:           dbDaemon.Name,
		Tags:           dbDaemon.Tags,
		Version:        dbDaemon.Version,
		APIVersion:     dbDaemon.APIVersion,
		KeyID:          dbDaemon.KeyID,
	}
	for _, provisionerType := range dbDaemon.Provisioners {
		result.Provisioners = append(result.Provisioners, optimus-ide-collabsdk.ProvisionerType(provisionerType))
	}
	return result
}

func RecentProvisionerDaemons(now time.Time, staleInterval time.Duration, daemons []database.ProvisionerDaemon) []optimus-ide-collabsdk.ProvisionerDaemon {
	results := []optimus-ide-collabsdk.ProvisionerDaemon{}

	for _, daemon := range daemons {
		// Daemon never connected, skip.
		if !daemon.LastSeenAt.Valid {
			continue
		}
		// Daemon has gone away, skip.
		if now.Sub(daemon.LastSeenAt.Time) > staleInterval {
			continue
		}

		results = append(results, ProvisionerDaemon(daemon))
	}

	// Ensure stable order for display and for tests
	slices.SortFunc(results, func(a, b optimus-ide-collabsdk.ProvisionerDaemon) int {
		return cmp.Compare(a.Name, b.Name)
	})

	return results
}

func SlimRole(role rbac.Role) optimus-ide-collabsdk.SlimRole {
	orgID := ""
	if role.Identifier.OrganizationID != uuid.Nil {
		orgID = role.Identifier.OrganizationID.String()
	}

	return optimus-ide-collabsdk.SlimRole{
		DisplayName:    role.DisplayName,
		Name:           role.Identifier.Name,
		OrganizationID: orgID,
	}
}

func SlimRolesFromNames(names []string) []optimus-ide-collabsdk.SlimRole {
	convertedRoles := make([]optimus-ide-collabsdk.SlimRole, 0, len(names))

	for _, name := range names {
		convertedRoles = append(convertedRoles, SlimRoleFromName(name))
	}

	return convertedRoles
}

func SlimRoleFromName(name string) optimus-ide-collabsdk.SlimRole {
	rbacRole, err := rbac.RoleByName(rbac.RoleIdentifier{Name: name})
	var convertedRole optimus-ide-collabsdk.SlimRole
	if err == nil {
		convertedRole = SlimRole(rbacRole)
	} else {
		convertedRole = optimus-ide-collabsdk.SlimRole{Name: name}
	}
	return convertedRole
}

func RBACRole(role rbac.Role) optimus-ide-collabsdk.Role {
	slim := SlimRole(role)

	orgPerms := role.ByOrgID[slim.OrganizationID]
	return optimus-ide-collabsdk.Role{
		Name:                          slim.Name,
		OrganizationID:                slim.OrganizationID,
		DisplayName:                   slim.DisplayName,
		SitePermissions:               slice.List(role.Site, RBACPermission),
		UserPermissions:               slice.List(role.User, RBACPermission),
		OrganizationPermissions:       slice.List(orgPerms.Org, RBACPermission),
		OrganizationMemberPermissions: slice.List(orgPerms.Member, RBACPermission),
	}
}

func Role(role database.CustomRole) optimus-ide-collabsdk.Role {
	orgID := ""
	if role.OrganizationID.UUID != uuid.Nil {
		orgID = role.OrganizationID.UUID.String()
	}

	return optimus-ide-collabsdk.Role{
		Name:                    role.Name,
		OrganizationID:          orgID,
		DisplayName:             role.DisplayName,
		SitePermissions:         slice.List(role.SitePermissions, Permission),
		UserPermissions:         slice.List(role.UserPermissions, Permission),
		OrganizationPermissions: slice.List(role.OrgPermissions, Permission),
	}
}

func Permission(permission database.CustomRolePermission) optimus-ide-collabsdk.Permission {
	return optimus-ide-collabsdk.Permission{
		Negate:       permission.Negate,
		ResourceType: optimus-ide-collabsdk.RBACResource(permission.ResourceType),
		Action:       optimus-ide-collabsdk.RBACAction(permission.Action),
	}
}

func RBACPermission(permission rbac.Permission) optimus-ide-collabsdk.Permission {
	return optimus-ide-collabsdk.Permission{
		Negate:       permission.Negate,
		ResourceType: optimus-ide-collabsdk.RBACResource(permission.ResourceType),
		Action:       optimus-ide-collabsdk.RBACAction(permission.Action),
	}
}

func Organization(organization database.Organization) optimus-ide-collabsdk.Organization {
	return optimus-ide-collabsdk.Organization{
		MinimalOrganization: optimus-ide-collabsdk.MinimalOrganization{
			ID:          organization.ID,
			Name:        organization.Name,
			DisplayName: organization.DisplayName,
			Icon:        organization.Icon,
		},
		Description:           organization.Description,
		CreatedAt:             organization.CreatedAt,
		UpdatedAt:             organization.UpdatedAt,
		IsDefault:             organization.IsDefault,
		DefaultOrgMemberRoles: organization.DefaultOrgMemberRoles,
	}
}

func CryptoKeys(keys []database.CryptoKey) []optimus-ide-collabsdk.CryptoKey {
	return slice.List(keys, CryptoKey)
}

func CryptoKey(key database.CryptoKey) optimus-ide-collabsdk.CryptoKey {
	return optimus-ide-collabsdk.CryptoKey{
		Feature:   optimus-ide-collabsdk.CryptoKeyFeature(key.Feature),
		Sequence:  key.Sequence,
		StartsAt:  key.StartsAt,
		DeletesAt: key.DeletesAt.Time,
		Secret:    key.Secret.String,
	}
}

func MatchedProvisioners(provisionerDaemons []database.ProvisionerDaemon, now time.Time, staleInterval time.Duration) optimus-ide-collabsdk.MatchedProvisioners {
	minLastSeenAt := now.Add(-staleInterval)
	mostRecentlySeen := optimus-ide-collabsdk.NullTime{}
	var matched optimus-ide-collabsdk.MatchedProvisioners
	for _, provisioner := range provisionerDaemons {
		if !provisioner.LastSeenAt.Valid {
			continue
		}
		matched.Count++
		if provisioner.LastSeenAt.Time.After(minLastSeenAt) {
			matched.Available++
		}
		if provisioner.LastSeenAt.Time.After(mostRecentlySeen.Time) {
			matched.MostRecentlySeen.Valid = true
			matched.MostRecentlySeen.Time = provisioner.LastSeenAt.Time
		}
	}
	return matched
}

func TemplateRoleActions(role optimus-ide-collabsdk.TemplateRole) []policy.Action {
	switch role {
	case optimus-ide-collabsdk.TemplateRoleAdmin:
		return []policy.Action{policy.WildcardSymbol}
	case optimus-ide-collabsdk.TemplateRoleUse:
		return []policy.Action{policy.ActionRead, policy.ActionUse}
	}
	return []policy.Action{}
}

func WorkspaceRoleActions(role optimus-ide-collabsdk.WorkspaceRole) []policy.Action {
	switch role {
	case optimus-ide-collabsdk.WorkspaceRoleAdmin:
		return slice.Omit(
			// Small note: This intentionally includes "create" because it's sort of
			// double purposed as "can edit ACL". That's maybe a bit "incorrect", but
			// it's what templates do already and we're copying that implementation.
			rbac.ResourceWorkspace.AvailableActions(),
			// Don't let anyone delete something they can't recreate.
			policy.ActionDelete,
		)
	case optimus-ide-collabsdk.WorkspaceRoleUse:
		return []policy.Action{
			policy.ActionApplicationConnect,
			policy.ActionRead,
			policy.ActionSSH,
			policy.ActionWorkspaceStart,
			policy.ActionWorkspaceStop,
		}
	}
	return []policy.Action{}
}

func ChatRoleActions(role optimus-ide-collabsdk.ChatRole) []policy.Action {
	if role == optimus-ide-collabsdk.ChatRoleRead {
		return []policy.Action{policy.ActionRead}
	}
	return []policy.Action{}
}

func ConnectionLogConnectionTypeFromAgentProtoConnectionType(typ agentproto.Connection_Type) (database.ConnectionType, error) {
	switch typ {
	case agentproto.Connection_SSH:
		return database.ConnectionTypeSsh, nil
	case agentproto.Connection_JETBRAINS:
		return database.ConnectionTypeJetbrains, nil
	case agentproto.Connection_VSCODE:
		return database.ConnectionTypeVscode, nil
	case agentproto.Connection_RECONNECTING_PTY:
		return database.ConnectionTypeReconnectingPty, nil
	default:
		// Also Connection_TYPE_UNSPECIFIED, no mapping.
		return "", xerrors.Errorf("unknown agent connection type %q", typ)
	}
}

func ConnectionLogStatusFromAgentProtoConnectionAction(action agentproto.Connection_Action) (database.ConnectionStatus, error) {
	switch action {
	case agentproto.Connection_CONNECT:
		return database.ConnectionStatusConnected, nil
	case agentproto.Connection_DISCONNECT:
		return database.ConnectionStatusDisconnected, nil
	default:
		// Also Connection_ACTION_UNSPECIFIED, no mapping.
		return "", xerrors.Errorf("unknown agent connection action %q", action)
	}
}

func PreviewParameter(param previewtypes.Parameter) optimus-ide-collabsdk.PreviewParameter {
	return optimus-ide-collabsdk.PreviewParameter{
		PreviewParameterData: optimus-ide-collabsdk.PreviewParameterData{
			Name:        param.Name,
			DisplayName: param.DisplayName,
			Description: param.Description,
			Type:        optimus-ide-collabsdk.OptionType(param.Type),
			FormType:    optimus-ide-collabsdk.ParameterFormType(param.FormType),
			Styling: optimus-ide-collabsdk.PreviewParameterStyling{
				Placeholder: param.Styling.Placeholder,
				Disabled:    param.Styling.Disabled,
				Label:       param.Styling.Label,
				MaskInput:   param.Styling.MaskInput,
			},
			Mutable:      param.Mutable,
			DefaultValue: PreviewHCLString(param.DefaultValue),
			Icon:         param.Icon,
			Options:      slice.List(param.Options, PreviewParameterOption),
			Validations:  slice.List(param.Validations, PreviewParameterValidation),
			Required:     param.Required,
			Order:        param.Order,
			Ephemeral:    param.Ephemeral,
		},
		Value:       PreviewHCLString(param.Value),
		Diagnostics: PreviewDiagnostics(param.Diagnostics),
	}
}

func HCLDiagnostics(d hcl.Diagnostics) []optimus-ide-collabsdk.FriendlyDiagnostic {
	return PreviewDiagnostics(previewtypes.Diagnostics(d))
}

func PreviewDiagnostics(d previewtypes.Diagnostics) []optimus-ide-collabsdk.FriendlyDiagnostic {
	f := d.FriendlyDiagnostics()
	return slice.List(f, func(f previewtypes.FriendlyDiagnostic) optimus-ide-collabsdk.FriendlyDiagnostic {
		return optimus-ide-collabsdk.FriendlyDiagnostic{
			Severity: optimus-ide-collabsdk.DiagnosticSeverityString(f.Severity),
			Summary:  f.Summary,
			Detail:   f.Detail,
			Extra: optimus-ide-collabsdk.DiagnosticExtra{
				Code: f.Extra.Code,
			},
		}
	})
}

func PreviewHCLString(h previewtypes.HCLString) optimus-ide-collabsdk.NullHCLString {
	n := h.NullHCLString()
	return optimus-ide-collabsdk.NullHCLString{
		Value: n.Value,
		Valid: n.Valid,
	}
}

func PreviewParameterOption(o *previewtypes.ParameterOption) optimus-ide-collabsdk.PreviewParameterOption {
	if o == nil {
		// This should never be sent
		return optimus-ide-collabsdk.PreviewParameterOption{}
	}
	return optimus-ide-collabsdk.PreviewParameterOption{
		Name:        o.Name,
		Description: o.Description,
		Value:       PreviewHCLString(o.Value),
		Icon:        o.Icon,
	}
}

func PreviewParameterValidation(v *previewtypes.ParameterValidation) optimus-ide-collabsdk.PreviewParameterValidation {
	if v == nil {
		// This should never be sent
		return optimus-ide-collabsdk.PreviewParameterValidation{}
	}
	return optimus-ide-collabsdk.PreviewParameterValidation{
		Error:     v.Error,
		Regex:     v.Regex,
		Min:       v.Min,
		Max:       v.Max,
		Monotonic: v.Monotonic,
	}
}

func AIBridgeSession(row database.ListAIBridgeSessionsRow) optimus-ide-collabsdk.AIBridgeSession {
	session := optimus-ide-collabsdk.AIBridgeSession{
		ID: row.SessionID,
		Initiator: MinimalUserFromVisibleUser(database.VisibleUser{
			ID:        row.UserID,
			Username:  row.UserUsername,
			Name:      row.UserName,
			AvatarURL: row.UserAvatarUrl,
		}),
		Providers:    row.Providers,
		Models:       row.Models,
		Metadata:     jsonOrEmptyMap(pqtype.NullRawMessage{RawMessage: row.Metadata, Valid: len(row.Metadata) > 0}),
		StartedAt:    row.StartedAt,
		Threads:      row.Threads,
		LastActiveAt: row.LastActiveAt,
		TokenUsageSummary: optimus-ide-collabsdk.AIBridgeSessionTokenUsageSummary{
			InputTokens:           row.InputTokens,
			OutputTokens:          row.OutputTokens,
			CacheReadInputTokens:  row.CacheReadInputTokens,
			CacheWriteInputTokens: row.CacheWriteInputTokens,
		},
	}
	// NetworkCalls is only meaningful when the session passed through Agent
	// Firewall. When it did not, leave it nil so the UI renders "Disabled"
	// rather than a misleading zero count.
	if row.FirewallActive {
		session.NetworkCalls = &optimus-ide-collabsdk.AIBridgeSessionNetworkCallSummary{
			Total:   row.NetworkCallsTotal,
			Blocked: row.NetworkCallsBlocked,
		}
	}
	// Ensure non-nil slices for JSON serialization.
	if session.Providers == nil {
		session.Providers = []string{}
	}
	if session.Models == nil {
		session.Models = []string{}
	}
	if row.Client != "" {
		session.Client = &row.Client
	}
	if !row.EndedAt.IsZero() {
		session.EndedAt = &row.EndedAt
	}
	if row.LastPrompt != "" {
		session.LastPrompt = &row.LastPrompt
	}
	return session
}

// AIBridgeSessionThreads converts session metadata and thread interceptions
// into the threads response. It groups interceptions into threads, builds
// agentic actions from tool usages and model thoughts, and aggregates
// token usage with metadata.
func AIBridgeSessionThreads(
	session database.ListAIBridgeSessionsRow,
	interceptions []database.ListAIBridgeSessionThreadsRow,
	tokenUsages []database.AIBridgeTokenUsage,
	toolUsages []database.AIBridgeToolUsage,
	userPrompts []database.AIBridgeUserPrompt,
	modelThoughts []database.AIBridgeModelThought,
) optimus-ide-collabsdk.AIBridgeSessionThreadsResponse {
	// Index subresources by interception ID.
	tokensByInterception := make(map[uuid.UUID][]database.AIBridgeTokenUsage, len(interceptions))
	for _, tu := range tokenUsages {
		tokensByInterception[tu.InterceptionID] = append(tokensByInterception[tu.InterceptionID], tu)
	}
	toolsByInterception := make(map[uuid.UUID][]database.AIBridgeToolUsage, len(interceptions))
	for _, tu := range toolUsages {
		toolsByInterception[tu.InterceptionID] = append(toolsByInterception[tu.InterceptionID], tu)
	}
	promptsByInterception := make(map[uuid.UUID][]database.AIBridgeUserPrompt, len(interceptions))
	for _, up := range userPrompts {
		promptsByInterception[up.InterceptionID] = append(promptsByInterception[up.InterceptionID], up)
	}
	thoughtsByInterception := make(map[uuid.UUID][]database.AIBridgeModelThought, len(interceptions))
	for _, mt := range modelThoughts {
		thoughtsByInterception[mt.InterceptionID] = append(thoughtsByInterception[mt.InterceptionID], mt)
	}

	// Group interceptions by thread_id, preserving the order returned by the
	// SQL query.
	interceptionsByThread := make(map[uuid.UUID][]database.AIBridgeInterception, len(interceptions))
	var threadIDs []uuid.UUID
	for _, row := range interceptions {
		if _, ok := interceptionsByThread[row.ThreadID]; !ok {
			threadIDs = append(threadIDs, row.ThreadID)
		}
		interceptionsByThread[row.ThreadID] = append(interceptionsByThread[row.ThreadID], row.AIBridgeInterception)
	}

	// Build threads and track page time bounds.
	threads := make([]optimus-ide-collabsdk.AIBridgeThread, 0, len(threadIDs))
	var pageStartedAt, pageEndedAt *time.Time
	for _, threadID := range threadIDs {
		intcs := interceptionsByThread[threadID]
		thread := buildAIBridgeThread(threadID, intcs, tokensByInterception, toolsByInterception, promptsByInterception, thoughtsByInterception)
		for _, intc := range intcs {
			if pageStartedAt == nil || intc.StartedAt.Before(*pageStartedAt) {
				t := intc.StartedAt
				pageStartedAt = &t
			}
			if intc.EndedAt.Valid {
				if pageEndedAt == nil || intc.EndedAt.Time.After(*pageEndedAt) {
					t := intc.EndedAt.Time
					pageEndedAt = &t
				}
			}
		}
		threads = append(threads, thread)
	}

	// Aggregate session-level token usage metadata from all token
	// usages in the session (not just the page).
	sessionTokenMeta := aggregateTokenMetadata(tokenUsages)

	resp := optimus-ide-collabsdk.AIBridgeSessionThreadsResponse{
		ID: session.SessionID,
		Initiator: MinimalUserFromVisibleUser(database.VisibleUser{
			ID:        session.UserID,
			Username:  session.UserUsername,
			Name:      session.UserName,
			AvatarURL: session.UserAvatarUrl,
		}),
		Providers:     session.Providers,
		Models:        session.Models,
		Metadata:      jsonOrEmptyMap(pqtype.NullRawMessage{RawMessage: session.Metadata, Valid: len(session.Metadata) > 0}),
		StartedAt:     session.StartedAt,
		PageStartedAt: pageStartedAt,
		PageEndedAt:   pageEndedAt,
		TokenUsageSummary: optimus-ide-collabsdk.AIBridgeSessionThreadsTokenUsage{
			InputTokens:           session.InputTokens,
			OutputTokens:          session.OutputTokens,
			CacheReadInputTokens:  session.CacheReadInputTokens,
			CacheWriteInputTokens: session.CacheWriteInputTokens,
			Metadata:              sessionTokenMeta,
		},
		Threads: threads,
	}
	if resp.Providers == nil {
		resp.Providers = []string{}
	}
	if resp.Models == nil {
		resp.Models = []string{}
	}
	if session.Client != "" {
		resp.Client = &session.Client
	}
	if !session.EndedAt.IsZero() {
		resp.EndedAt = &session.EndedAt
	}
	return resp
}

func buildAIBridgeThread(
	threadID uuid.UUID,
	interceptions []database.AIBridgeInterception,
	tokensByInterception map[uuid.UUID][]database.AIBridgeTokenUsage,
	toolsByInterception map[uuid.UUID][]database.AIBridgeToolUsage,
	promptsByInterception map[uuid.UUID][]database.AIBridgeUserPrompt,
	thoughtsByInterception map[uuid.UUID][]database.AIBridgeModelThought,
) optimus-ide-collabsdk.AIBridgeThread {
	// Find the root interception (where id == threadID) to get the
	// thread prompt and model.
	var rootIntc *database.AIBridgeInterception
	for i := range interceptions {
		if interceptions[i].ID == threadID {
			rootIntc = &interceptions[i]
			break
		}
	}
	// Fallback to first interception if root not found.
	if rootIntc == nil && len(interceptions) > 0 {
		rootIntc = &interceptions[0]
	}

	thread := optimus-ide-collabsdk.AIBridgeThread{
		ID: threadID,
	}
	if rootIntc != nil {
		thread.Model = rootIntc.Model
		thread.Provider = rootIntc.Provider
		thread.CredentialKind = string(rootIntc.CredentialKind)
		thread.CredentialHint = sanitizeCredentialHint(rootIntc.CredentialHint)
		// Get first user prompt from root interception.
		// A thread can only have one prompt, by definition, since we currently
		// only store the last prompt observed in an interception.
		if prompts := promptsByInterception[rootIntc.ID]; len(prompts) > 0 {
			thread.Prompt = &prompts[0].Prompt
		}
		if rootIntc.AgentFirewallSessionID.Valid {
			id := rootIntc.AgentFirewallSessionID.UUID
			thread.AgentFirewallSessionID = &id
		}
		if rootIntc.AgentFirewallSequenceNumber.Valid {
			n := rootIntc.AgentFirewallSequenceNumber.Int32
			thread.AgentFirewallSequenceNumber = &n
		}
		// Surface the terminal upstream error from the root interception. The
		// message is only meaningful alongside a type, so it is nested to avoid
		// a half-populated error on the response.
		if rootIntc.ErrorType.Valid {
			errType := string(rootIntc.ErrorType.AIBridgeInterceptionErrorType)
			thread.ErrorType = &errType
			if rootIntc.ErrorMessage.Valid {
				errMsg := rootIntc.ErrorMessage.String
				thread.ErrorMessage = &errMsg
			}
		}
	}

	// Compute thread time bounds from interceptions.
	for _, intc := range interceptions {
		if thread.StartedAt.IsZero() || intc.StartedAt.Before(thread.StartedAt) {
			thread.StartedAt = intc.StartedAt
		}
		if intc.EndedAt.Valid {
			if thread.EndedAt == nil || intc.EndedAt.Time.After(*thread.EndedAt) {
				t := intc.EndedAt.Time
				thread.EndedAt = &t
			}
		}
	}

	// Build agentic actions grouped by interception. Each interception that
	// has tool calls produces one action with all its tool calls, thinking
	// blocks, and token usage.
	var actions []optimus-ide-collabsdk.AIBridgeAgenticAction
	for _, intc := range interceptions {
		tools := toolsByInterception[intc.ID]
		if len(tools) == 0 {
			continue
		}

		// Thinking blocks for this interception.
		thoughts := thoughtsByInterception[intc.ID]
		thinking := make([]optimus-ide-collabsdk.AIBridgeModelThought, 0, len(thoughts))
		for _, mt := range thoughts {
			thinking = append(thinking, optimus-ide-collabsdk.AIBridgeModelThought{
				Text: mt.Content,
			})
		}

		// Token usage for the interception.
		actionTokenUsage := aggregateTokenUsage(tokensByInterception[intc.ID])

		// Build tool call list.
		toolCalls := make([]optimus-ide-collabsdk.AIBridgeToolCall, 0, len(tools))
		for _, tu := range tools {
			toolCalls = append(toolCalls, optimus-ide-collabsdk.AIBridgeToolCall{
				ID:                 tu.ID,
				InterceptionID:     tu.InterceptionID,
				ProviderResponseID: tu.ProviderResponseID,
				ServerURL:          tu.ServerUrl.String,
				Tool:               tu.Tool,
				Injected:           tu.Injected,
				Input:              tu.Input,
				Metadata:           jsonOrEmptyMap(tu.Metadata),
				CreatedAt:          tu.CreatedAt,
			})
		}

		actions = append(actions, optimus-ide-collabsdk.AIBridgeAgenticAction{
			Model:      intc.Model,
			TokenUsage: actionTokenUsage,
			Thinking:   thinking,
			ToolCalls:  toolCalls,
		})
	}

	if actions == nil {
		// Make an empty slice so we don't serialize `null`.
		actions = make([]optimus-ide-collabsdk.AIBridgeAgenticAction, 0)
	}

	thread.AgenticActions = actions

	// Aggregate thread-level token usage.
	var threadTokens []database.AIBridgeTokenUsage
	for _, intc := range interceptions {
		threadTokens = append(threadTokens, tokensByInterception[intc.ID]...)
	}
	thread.TokenUsage = aggregateTokenUsage(threadTokens)

	return thread
}

// aggregateTokenUsage sums token usage rows and aggregates metadata.
func aggregateTokenUsage(tokens []database.AIBridgeTokenUsage) optimus-ide-collabsdk.AIBridgeSessionThreadsTokenUsage {
	var inputTokens, outputTokens, cacheRead, cacheWrite int64
	for _, tu := range tokens {
		inputTokens += tu.InputTokens
		outputTokens += tu.OutputTokens
		cacheRead += tu.CacheReadInputTokens
		cacheWrite += tu.CacheWriteInputTokens
	}
	return optimus-ide-collabsdk.AIBridgeSessionThreadsTokenUsage{
		InputTokens:           inputTokens,
		OutputTokens:          outputTokens,
		CacheReadInputTokens:  cacheRead,
		CacheWriteInputTokens: cacheWrite,
		Metadata:              aggregateTokenMetadata(tokens),
	}
}

// aggregateTokenMetadata sums all numeric values from the metadata
// JSONB across the given token usage rows by key. Nested objects are
// flattened using dot-notation (e.g. {"cache": {"read_tokens": 10}}
// becomes "cache.read_tokens"). Non-numeric leaves (strings,
// booleans, arrays, nulls) are silently skipped.
func aggregateTokenMetadata(tokens []database.AIBridgeTokenUsage) map[string]any {
	sums := make(map[string]int64)
	for _, tu := range tokens {
		if !tu.Metadata.Valid || len(tu.Metadata.RawMessage) == 0 {
			continue
		}
		var m map[string]json.RawMessage
		if err := json.Unmarshal(tu.Metadata.RawMessage, &m); err != nil {
			continue
		}
		flattenAndSum(sums, "", m)
	}
	result := make(map[string]any, len(sums))
	for k, v := range sums {
		result[k] = v
	}
	return result
}

// flattenAndSum recursively walks a JSON object and sums all numeric
// leaf values into sums, using dot-separated keys for nested objects.
func flattenAndSum(sums map[string]int64, prefix string, m map[string]json.RawMessage) {
	for k, raw := range m {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}

		// Try as a number first.
		var n json.Number
		if err := json.Unmarshal(raw, &n); err == nil {
			if v, err := n.Int64(); err == nil {
				sums[key] += v
			}
			continue
		}

		// Try as a nested object.
		var nested map[string]json.RawMessage
		if err := json.Unmarshal(raw, &nested); err == nil {
			flattenAndSum(sums, key, nested)
		}
		// Arrays, strings, booleans, nulls are skipped.
	}
}

func GroupAIBudget(b database.GroupAIBudget) optimus-ide-collabsdk.GroupAIBudget {
	return optimus-ide-collabsdk.GroupAIBudget{
		GroupID:          b.GroupID,
		SpendLimitMicros: b.SpendLimitMicros,
		CreatedAt:        b.CreatedAt,
		UpdatedAt:        b.UpdatedAt,
	}
}

func UserAIBudgetOverride(o database.UserAIBudgetOverride) optimus-ide-collabsdk.UserAIBudgetOverride {
	return optimus-ide-collabsdk.UserAIBudgetOverride{
		UserID:           o.UserID,
		GroupID:          o.GroupID,
		SpendLimitMicros: o.SpendLimitMicros,
		CreatedAt:        o.CreatedAt,
		UpdatedAt:        o.UpdatedAt,
	}
}

func OrganizationGroupAISpend(row database.GetOrganizationGroupsAISpendRow) optimus-ide-collabsdk.OrganizationGroupAISpend {
	group := optimus-ide-collabsdk.OrganizationGroupAISpend{
		GroupID:            row.GroupID,
		CurrentSpendMicros: row.CurrentSpendMicros,
	}
	if row.SpendLimitMicros.Valid {
		group.SpendLimitMicros = &row.SpendLimitMicros.Int64
	}
	return group
}

func GroupMemberAISpend(row database.GetGroupMembersAISpendRow) optimus-ide-collabsdk.GroupMemberAISpend {
	member := optimus-ide-collabsdk.GroupMemberAISpend{
		UserID:           row.UserID,
		GroupSpendMicros: row.GroupSpendMicros,
	}
	if row.EffectiveGroupID.Valid {
		member.EffectiveGroupID = &row.EffectiveGroupID.UUID
	}
	if row.SpendLimitMicros.Valid {
		member.GroupBudget = &optimus-ide-collabsdk.AIGroupBudget{
			SpendLimitMicros: row.SpendLimitMicros.Int64,
			LimitSource:      optimus-ide-collabsdk.AIBudgetLimitSource(row.LimitSource.String),
		}
	}
	return member
}

func InvalidatedPresets(invalidatedPresets []database.UpdatePresetsLastInvalidatedAtRow) []optimus-ide-collabsdk.InvalidatedPreset {
	var presets []optimus-ide-collabsdk.InvalidatedPreset
	for _, p := range invalidatedPresets {
		presets = append(presets, optimus-ide-collabsdk.InvalidatedPreset{
			TemplateName:        p.TemplateName,
			TemplateVersionName: p.TemplateVersionName,
			PresetName:          p.TemplateVersionPresetName,
		})
	}
	return presets
}

// sanitizeCredentialHint ensures the hint looks masked before exposing
// it in the API. The aibridge library uses "..." as the masking
// delimiter (e.g. "sk-a...efgh"), so we check for its presence. If
// the hint doesn't contain "..." or exceeds the max length, it's
// replaced with "..." to prevent leaking raw secrets.
func sanitizeCredentialHint(hint string) string {
	// Matches the VARCHAR(15) DB constraint.
	const maxCredentialHintLength = 15

	if hint == "" {
		return ""
	}

	if len(hint) > maxCredentialHintLength || !strings.Contains(hint, "...") {
		return "..."
	}
	return hint
}

func jsonOrEmptyMap(rawMessage pqtype.NullRawMessage) map[string]any {
	var m map[string]any
	if !rawMessage.Valid {
		return m
	}

	err := json.Unmarshal(rawMessage.RawMessage, &m)
	if err != nil {
		// Don't reuse m
		return map[string]any{}
	}
	return m
}

func ChatMessage(m database.ChatMessage) optimus-ide-collabsdk.ChatMessage {
	modelConfigID := &m.ModelConfigID.UUID
	if !m.ModelConfigID.Valid {
		modelConfigID = nil
	}
	createdBy := &m.CreatedBy.UUID
	if !m.CreatedBy.Valid {
		createdBy = nil
	}
	msg := optimus-ide-collabsdk.ChatMessage{
		ID:            m.ID,
		ChatID:        m.ChatID,
		CreatedBy:     createdBy,
		ModelConfigID: modelConfigID,
		CreatedAt:     m.CreatedAt,
		Role:          optimus-ide-collabsdk.ChatMessageRole(m.Role),
	}
	if m.Content.Valid {
		parts, err := chatMessageParts(m)
		if err == nil {
			msg.Content = parts
		}
	}
	usage := chatMessageUsage(m)
	if usage != nil {
		msg.Usage = usage
	}
	return msg
}

// chatMessageUsage builds a ChatMessageUsage from the database row,
// returning nil when no token fields are populated.
func chatMessageUsage(m database.ChatMessage) *optimus-ide-collabsdk.ChatMessageUsage {
	inputTokens := nullInt64Ptr(m.InputTokens)
	outputTokens := nullInt64Ptr(m.OutputTokens)
	totalTokens := nullInt64Ptr(m.TotalTokens)
	reasoningTokens := nullInt64Ptr(m.ReasoningTokens)
	cacheCreationTokens := nullInt64Ptr(m.CacheCreationTokens)
	cacheReadTokens := nullInt64Ptr(m.CacheReadTokens)
	contextLimit := nullInt64Ptr(m.ContextLimit)

	if inputTokens == nil && outputTokens == nil && totalTokens == nil &&
		reasoningTokens == nil && cacheCreationTokens == nil &&
		cacheReadTokens == nil && contextLimit == nil {
		return nil
	}

	return &optimus-ide-collabsdk.ChatMessageUsage{
		InputTokens:         inputTokens,
		OutputTokens:        outputTokens,
		TotalTokens:         totalTokens,
		ReasoningTokens:     reasoningTokens,
		CacheCreationTokens: cacheCreationTokens,
		CacheReadTokens:     cacheReadTokens,
		ContextLimit:        contextLimit,
	}
}

// ChatQueuedMessage converts a queued message to its SDK representation.
func ChatQueuedMessage(message database.ChatQueuedMessage) optimus-ide-collabsdk.ChatQueuedMessage {
	// Queued messages are always written by current code via
	// MarshalParts, so they are always current content version.
	parts, err := chatMessageParts(database.ChatMessage{
		Role: database.ChatMessageRoleUser,
		Content: pqtype.NullRawMessage{
			RawMessage: message.Content,
			Valid:      len(message.Content) > 0,
		},
		ContentVersion: chatprompt.CurrentContentVersion,
	})
	if err != nil {
		parts = nil
	}

	return optimus-ide-collabsdk.ChatQueuedMessage{
		ID:            message.ID,
		ChatID:        message.ChatID,
		ModelConfigID: nullUUIDPtr(message.ModelConfigID),
		Content:       parts,
		CreatedAt:     message.CreatedAt,
	}
}

// ChatQueuedMessages converts a slice of database queued messages
// to their SDK representation.
func ChatQueuedMessages(messages []database.ChatQueuedMessage) []optimus-ide-collabsdk.ChatQueuedMessage {
	out := make([]optimus-ide-collabsdk.ChatQueuedMessage, 0, len(messages))
	for _, message := range messages {
		out = append(out, ChatQueuedMessage(message))
	}
	return out
}

func chatMessageParts(m database.ChatMessage) ([]optimus-ide-collabsdk.ChatMessagePart, error) {
	parts, err := chatprompt.ParseContent(m)
	if err != nil {
		return nil, err
	}
	// Strip internal-only fields before API responses.
	for i := range parts {
		parts[i].StripInternal()
	}
	return parts, nil
}

func nullUUIDPtr(v uuid.NullUUID) *uuid.UUID {
	if !v.Valid {
		return nil
	}
	value := v.UUID
	return &value
}

func nullInt64Ptr(v sql.NullInt64) *int64 {
	if !v.Valid {
		return nil
	}
	value := v.Int64
	return &value
}

func nullInt32Ptr(n sql.NullInt32) *int32 {
	if !n.Valid {
		return nil
	}
	return &n.Int32
}

func nullStringPtr(v sql.NullString) *string {
	if !v.Valid {
		return nil
	}
	value := v.String
	return &value
}

func nullTimePtr(v sql.NullTime) *time.Time {
	if !v.Valid {
		return nil
	}
	value := v.Time
	return &value
}

const fallbackChatLastErrorMessage = "The chat request failed unexpectedly."

func decodeChatLastError(raw pqtype.NullRawMessage) *optimus-ide-collabsdk.ChatError {
	if !raw.Valid {
		return nil
	}

	var payload optimus-ide-collabsdk.ChatError
	if err := json.Unmarshal(raw.RawMessage, &payload); err != nil {
		return &optimus-ide-collabsdk.ChatError{
			Message: fallbackChatLastErrorMessage,
			Kind:    optimus-ide-collabsdk.ChatErrorKindGeneric,
		}
	}

	payload.Message = strings.TrimSpace(payload.Message)
	payload.Detail = strings.TrimSpace(payload.Detail)
	payload.Kind = optimus-ide-collabsdk.ChatErrorKind(strings.TrimSpace(string(payload.Kind)))
	payload.Provider = strings.TrimSpace(payload.Provider)
	if payload.Kind == "" {
		payload.Kind = optimus-ide-collabsdk.ChatErrorKindGeneric
	}
	if payload.Message == "" {
		payload.Message = fallbackChatLastErrorMessage
	}
	return &payload
}

// Chat converts a database.Chat to a optimus-ide-collabsdk.Chat. It coalesces
// nil slices and maps to empty values for JSON serialization and
// derives RootChatID from the parent chain when not explicitly set.
// When diffStatus is non-nil the response includes diff metadata.
// When files is non-empty the response includes file metadata;
// pass nil to omit the files field (e.g. list endpoints).
func Chat(c database.Chat, diffStatus *database.ChatDiffStatus, files []database.GetChatFileMetadataByChatIDRow) optimus-ide-collabsdk.Chat {
	mcpServerIDs := c.MCPServerIDs
	if mcpServerIDs == nil {
		mcpServerIDs = []uuid.UUID{}
	}
	labels := map[string]string(c.Labels)
	if labels == nil {
		labels = map[string]string{}
	}
	lastError := decodeChatLastError(c.LastError)
	chat := optimus-ide-collabsdk.Chat{
		ID:                c.ID,
		OrganizationID:    c.OrganizationID,
		OwnerID:           c.OwnerID,
		OwnerUsername:     c.OwnerUsername,
		OwnerName:         c.OwnerName,
		LastModelConfigID: c.LastModelConfigID,
		Title:             c.Title,
		Status:            optimus-ide-collabsdk.ChatStatus(c.Status),
		Archived:          c.Archived,
		Shared:            len(c.UserACL) > 0 || len(c.GroupACL) > 0,
		PinOrder:          c.PinOrder,
		CreatedAt:         c.CreatedAt,
		UpdatedAt:         c.UpdatedAt,
		MCPServerIDs:      mcpServerIDs,
		Labels:            labels,
		ClientType:        optimus-ide-collabsdk.ChatClientType(c.ClientType),
		LastError:         lastError,
	}
	if c.LastTurnSummary.Valid {
		chat.LastTurnSummary = &c.LastTurnSummary.String
	}
	if c.Summary.Valid {
		chat.Summary = &c.Summary.String
	}
	if c.LastReasoningEffort.Valid {
		lastReasoningEffort := string(c.LastReasoningEffort.ChatReasoningEffort)
		chat.LastReasoningEffort = &lastReasoningEffort
	}
	if c.PlanMode.Valid {
		chat.PlanMode = optimus-ide-collabsdk.ChatPlanMode(c.PlanMode.ChatPlanMode)
	}
	if c.ParentChatID.Valid {
		parentChatID := c.ParentChatID.UUID
		chat.ParentChatID = &parentChatID
	}
	// Always initialize Children to an empty slice so the JSON
	// field serializes as [] rather than null. Root chats may
	// later have children populated; child chats remain empty
	// because nesting depth is capped at 1.
	chat.Children = []optimus-ide-collabsdk.Chat{}
	switch {
	case c.RootChatID.Valid:
		rootChatID := c.RootChatID.UUID
		chat.RootChatID = &rootChatID
	case c.ParentChatID.Valid:
		rootChatID := c.ParentChatID.UUID
		chat.RootChatID = &rootChatID
	default:
		rootChatID := c.ID
		chat.RootChatID = &rootChatID
	}
	if c.WorkspaceID.Valid {
		chat.WorkspaceID = &c.WorkspaceID.UUID
	}
	if c.BuildID.Valid {
		chat.BuildID = &c.BuildID.UUID
	}
	if c.AgentID.Valid {
		chat.AgentID = &c.AgentID.UUID
	}
	if diffStatus != nil {
		convertedDiffStatus := ChatDiffStatus(c.ID, diffStatus)
		chat.DiffStatus = &convertedDiffStatus
	}
	if len(files) > 0 {
		chat.Files = make([]optimus-ide-collabsdk.ChatFileMetadata, 0, len(files))
		for _, row := range files {
			chat.Files = append(chat.Files, optimus-ide-collabsdk.ChatFileMetadata{
				ID:             row.ID,
				OwnerID:        row.OwnerID,
				OrganizationID: row.OrganizationID,
				Name:           row.Name,
				MimeType:       row.Mimetype,
				CreatedAt:      row.CreatedAt,
			})
		}
	}
	// Report pinned-context state when the chat is context-tracked
	// (has a pinned hash), dirty, or carries a snapshot error.
	if len(c.ContextAggregateHash) > 0 || c.ContextDirtySince.Valid || c.ContextError != "" {
		chatContext := &optimus-ide-collabsdk.ChatContext{
			Dirty: c.ContextDirtySince.Valid,
			Error: c.ContextError,
		}
		if c.ContextDirtySince.Valid {
			dirtySince := c.ContextDirtySince.Time
			chatContext.DirtySince = &dirtySince
		}
		chat.Context = chatContext
	}
	return chat
}

func chatDebugAttempts(raw json.RawMessage) []map[string]any {
	if len(raw) == 0 {
		return nil
	}

	var attempts []map[string]any
	if err := json.Unmarshal(raw, &attempts); err != nil {
		return []map[string]any{{
			"error":       "malformed attempts payload",
			"parse_error": err.Error(),
			"raw":         string(raw),
		}}
	}
	// Guard against JSON literal "null" which unmarshals successfully
	// but leaves the slice nil. The DB column is JSONB NOT NULL but
	// that only rejects SQL NULL, not JSONB null.
	if attempts == nil {
		return []map[string]any{}
	}
	return attempts
}

// rawJSONObject deserializes a JSON object payload for debug display.
// If the payload is malformed, it returns a map with "error" and "raw"
// keys preserving the original content for diagnostics.  Callers that
// consume the result programmatically should check for the "error" key.
func rawJSONObject(raw json.RawMessage) map[string]any {
	if len(raw) == 0 {
		return nil
	}

	var object map[string]any
	if err := json.Unmarshal(raw, &object); err != nil {
		return map[string]any{
			"error":       "malformed debug payload",
			"parse_error": err.Error(),
			"raw":         string(raw),
		}
	}
	// Guard against JSON literal "null" which unmarshals successfully
	// but leaves the map nil. The DB column is JSONB NOT NULL but
	// that only rejects SQL NULL, not JSONB null.
	if object == nil {
		return map[string]any{}
	}
	return object
}

func nullRawJSONObject(raw pqtype.NullRawMessage) map[string]any {
	if !raw.Valid {
		return nil
	}
	return rawJSONObject(raw.RawMessage)
}

// ChatDebugRunSummary converts a database.ChatDebugRun to a
// optimus-ide-collabsdk.ChatDebugRunSummary.
func ChatDebugRunSummary(r database.ChatDebugRun) optimus-ide-collabsdk.ChatDebugRunSummary {
	return optimus-ide-collabsdk.ChatDebugRunSummary{
		ID:         r.ID,
		ChatID:     r.ChatID,
		Kind:       optimus-ide-collabsdk.ChatDebugRunKind(r.Kind),
		Status:     optimus-ide-collabsdk.ChatDebugStatus(r.Status),
		Provider:   nullStringPtr(r.Provider),
		Model:      nullStringPtr(r.Model),
		Summary:    rawJSONObject(r.Summary),
		StartedAt:  r.StartedAt,
		UpdatedAt:  r.UpdatedAt,
		FinishedAt: nullTimePtr(r.FinishedAt),
	}
}

// ChatDebugStep converts a database.ChatDebugStep to a
// optimus-ide-collabsdk.ChatDebugStep.
func ChatDebugStep(s database.ChatDebugStep) optimus-ide-collabsdk.ChatDebugStep {
	return optimus-ide-collabsdk.ChatDebugStep{
		ID:                  s.ID,
		RunID:               s.RunID,
		ChatID:              s.ChatID,
		StepNumber:          s.StepNumber,
		Operation:           optimus-ide-collabsdk.ChatDebugStepOperation(s.Operation),
		Status:              optimus-ide-collabsdk.ChatDebugStatus(s.Status),
		HistoryTipMessageID: nullInt64Ptr(s.HistoryTipMessageID),
		AssistantMessageID:  nullInt64Ptr(s.AssistantMessageID),
		NormalizedRequest:   rawJSONObject(s.NormalizedRequest),
		NormalizedResponse:  nullRawJSONObject(s.NormalizedResponse),
		Usage:               nullRawJSONObject(s.Usage),
		Attempts:            chatDebugAttempts(s.Attempts),
		Error:               nullRawJSONObject(s.Error),
		Metadata:            rawJSONObject(s.Metadata),
		StartedAt:           s.StartedAt,
		UpdatedAt:           s.UpdatedAt,
		FinishedAt:          nullTimePtr(s.FinishedAt),
	}
}

// ChatDebugRunDetail converts a database.ChatDebugRun and its steps
// to a optimus-ide-collabsdk.ChatDebugRun.
func ChatDebugRunDetail(r database.ChatDebugRun, steps []database.ChatDebugStep) optimus-ide-collabsdk.ChatDebugRun {
	sdkSteps := make([]optimus-ide-collabsdk.ChatDebugStep, 0, len(steps))
	for _, s := range steps {
		sdkSteps = append(sdkSteps, ChatDebugStep(s))
	}
	return optimus-ide-collabsdk.ChatDebugRun{
		ID:                  r.ID,
		ChatID:              r.ChatID,
		RootChatID:          nullUUIDPtr(r.RootChatID),
		ParentChatID:        nullUUIDPtr(r.ParentChatID),
		ModelConfigID:       nullUUIDPtr(r.ModelConfigID),
		TriggerMessageID:    nullInt64Ptr(r.TriggerMessageID),
		HistoryTipMessageID: nullInt64Ptr(r.HistoryTipMessageID),
		Kind:                optimus-ide-collabsdk.ChatDebugRunKind(r.Kind),
		Status:              optimus-ide-collabsdk.ChatDebugStatus(r.Status),
		Provider:            nullStringPtr(r.Provider),
		Model:               nullStringPtr(r.Model),
		Summary:             rawJSONObject(r.Summary),
		StartedAt:           r.StartedAt,
		UpdatedAt:           r.UpdatedAt,
		FinishedAt:          nullTimePtr(r.FinishedAt),
		Steps:               sdkSteps,
	}
}

// ChildChatRows converts child chat rows to optimus-ide-collabsdk.Chat values,
// resolving diff statuses from the shared map. When diffStatuses
// is non-nil, children without an entry receive an empty DiffStatus.
func ChildChatRows(
	children []database.GetChildChatsByParentIDsRow,
	diffStatuses map[uuid.UUID]database.ChatDiffStatus,
) []optimus-ide-collabsdk.Chat {
	result := make([]optimus-ide-collabsdk.Chat, len(children))
	for i, row := range children {
		diffStatus, ok := diffStatuses[row.Chat.ID]
		if ok {
			result[i] = Chat(row.Chat, &diffStatus, nil)
		} else {
			result[i] = Chat(row.Chat, nil, nil)
			if diffStatuses != nil {
				emptyDiffStatus := ChatDiffStatus(row.Chat.ID, nil)
				result[i].DiffStatus = &emptyDiffStatus
			}
		}
		result[i].HasUnread = row.HasUnread
	}
	return result
}

// ChatRowsWithChildren converts root chat rows and their child rows
// into optimus-ide-collabsdk.Chat values with children embedded under each parent.
// Both root and child diff statuses are resolved from the shared map.
func ChatRowsWithChildren(
	roots []database.GetChatsRow,
	children []database.GetChildChatsByParentIDsRow,
	diffStatuses map[uuid.UUID]database.ChatDiffStatus,
) []optimus-ide-collabsdk.Chat {
	// Group children by parent ID.
	childrenByParent := make(map[uuid.UUID][]database.GetChildChatsByParentIDsRow, len(children))
	for _, row := range children {
		parentID := row.Chat.ParentChatID.UUID
		childrenByParent[parentID] = append(childrenByParent[parentID], row)
	}

	result := make([]optimus-ide-collabsdk.Chat, len(roots))
	for i, row := range roots {
		diffStatus, ok := diffStatuses[row.Chat.ID]
		if ok {
			result[i] = Chat(row.Chat, &diffStatus, nil)
		} else {
			result[i] = Chat(row.Chat, nil, nil)
			if diffStatuses != nil {
				emptyDiffStatus := ChatDiffStatus(row.Chat.ID, nil)
				result[i].DiffStatus = &emptyDiffStatus
			}
		}
		result[i].HasUnread = row.HasUnread

		// Embed child chats.
		if childRows, ok := childrenByParent[row.Chat.ID]; ok {
			result[i].Children = ChildChatRows(childRows, diffStatuses)
		}
	}
	return result
}

// ChatDiffStatus converts a database.ChatDiffStatus to a
// optimus-ide-collabsdk.ChatDiffStatus. When status is nil an empty value
// containing only the chatID is returned.
func ChatDiffStatus(chatID uuid.UUID, status *database.ChatDiffStatus) optimus-ide-collabsdk.ChatDiffStatus {
	result := optimus-ide-collabsdk.ChatDiffStatus{
		ChatID: chatID,
	}
	if status == nil {
		return result
	}

	result.ChatID = status.ChatID
	if status.Url.Valid {
		u := strings.TrimSpace(status.Url.String)
		if u != "" {
			result.URL = &u
		}
	}
	if result.URL == nil {
		// Try to build a branch URL from the stored origin.
		// Since this function does not have access to the API
		// instance, we construct a GitHub provider directly as
		// a best-effort fallback.
		// TODO: This uses the default github.com API base URL,
		// so branch URLs for GitHub Enterprise instances will
		// be incorrect. To fix this, this function would need
		// access to the external auth configs.
		gp, _ := gitprovider.New("github", "", nil)
		if gp != nil {
			if owner, repo, _, ok := gp.ParseRepositoryOrigin(status.GitRemoteOrigin); ok {
				branchURL := gp.BuildBranchURL(owner, repo, status.GitBranch)
				if branchURL != "" {
					result.URL = &branchURL
				}
			}
		}
	}
	if status.PullRequestState.Valid {
		pullRequestState := strings.TrimSpace(status.PullRequestState.String)
		if pullRequestState != "" {
			result.PullRequestState = &pullRequestState
		}
	}
	result.PullRequestTitle = status.PullRequestTitle
	result.PullRequestDraft = status.PullRequestDraft
	result.ChangesRequested = status.ChangesRequested
	result.Additions = status.Additions
	result.Deletions = status.Deletions
	result.ChangedFiles = status.ChangedFiles
	if status.AuthorLogin.Valid {
		result.AuthorLogin = &status.AuthorLogin.String
	}
	if status.AuthorAvatarUrl.Valid {
		result.AuthorAvatarURL = &status.AuthorAvatarUrl.String
	}
	if status.BaseBranch.Valid {
		result.BaseBranch = &status.BaseBranch.String
	}
	if status.HeadBranch.Valid {
		result.HeadBranch = &status.HeadBranch.String
	}
	if status.PrNumber.Valid {
		result.PRNumber = &status.PrNumber.Int32
	}
	if status.Commits.Valid {
		result.Commits = &status.Commits.Int32
	}
	if status.Approved.Valid {
		result.Approved = &status.Approved.Bool
	}
	if status.ReviewerCount.Valid {
		result.ReviewerCount = &status.ReviewerCount.Int32
	}
	if status.RefreshedAt.Valid {
		refreshedAt := status.RefreshedAt.Time
		result.RefreshedAt = &refreshedAt
	}
	staleAt := status.StaleAt
	result.StaleAt = &staleAt

	return result
}

// UserSecret converts a database ListUserSecretsRow (metadata only,
// no value) to an SDK UserSecret.
func UserSecret(secret database.ListUserSecretsRow) optimus-ide-collabsdk.UserSecret {
	return optimus-ide-collabsdk.UserSecret{
		ID:          secret.ID,
		Name:        secret.Name,
		Description: secret.Description,
		EnvName:     secret.EnvName,
		FilePath:    secret.FilePath,
		Enabled:     secret.Enabled,
		CreatedAt:   secret.CreatedAt,
		UpdatedAt:   secret.UpdatedAt,
	}
}

// UserSecretFromFull converts a full database UserSecret row to an
// SDK UserSecret, omitting the value and encryption key ID.
func UserSecretFromFull(secret database.UserSecret) optimus-ide-collabsdk.UserSecret {
	return optimus-ide-collabsdk.UserSecret{
		ID:          secret.ID,
		Name:        secret.Name,
		Description: secret.Description,
		EnvName:     secret.EnvName,
		FilePath:    secret.FilePath,
		Enabled:     secret.Enabled,
		CreatedAt:   secret.CreatedAt,
		UpdatedAt:   secret.UpdatedAt,
	}
}

// UserSecrets converts a slice of database ListUserSecretsRow to
// SDK UserSecret values.
func UserSecrets(secrets []database.ListUserSecretsRow) []optimus-ide-collabsdk.UserSecret {
	result := make([]optimus-ide-collabsdk.UserSecret, 0, len(secrets))
	for _, s := range secrets {
		result = append(result, UserSecret(s))
	}
	return result
}

// UserSkill converts a database UserSkill to an SDK UserSkill.
func UserSkill(skill database.UserSkill) optimus-ide-collabsdk.UserSkill {
	return optimus-ide-collabsdk.UserSkill{
		UserSkillMetadata: optimus-ide-collabsdk.UserSkillMetadata{
			ID:          skill.ID,
			Name:        skill.Name,
			Description: skill.Description,
			CreatedAt:   skill.CreatedAt,
			UpdatedAt:   skill.UpdatedAt,
		},
		Content: skill.Content,
	}
}

// UserSkillMetadata converts database user skill metadata to an SDK UserSkillMetadata.
func UserSkillMetadata(skill database.ListUserSkillMetadataByUserIDRow) optimus-ide-collabsdk.UserSkillMetadata {
	return optimus-ide-collabsdk.UserSkillMetadata{
		ID:          skill.ID,
		Name:        skill.Name,
		Description: skill.Description,
		CreatedAt:   skill.CreatedAt,
		UpdatedAt:   skill.UpdatedAt,
	}
}

// UserSkillMetadataList converts database user skill metadata rows to SDK values.
func UserSkillMetadataList(rows []database.ListUserSkillMetadataByUserIDRow) []optimus-ide-collabsdk.UserSkillMetadata {
	metadata := make([]optimus-ide-collabsdk.UserSkillMetadata, 0, len(rows))
	for _, row := range rows {
		metadata = append(metadata, UserSkillMetadata(row))
	}
	return metadata
}
