package portsharing

import (
	"golang.org/x/xerrors"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

type EnterprisePortSharer struct{}

func NewEnterprisePortSharer() *EnterprisePortSharer {
	return &EnterprisePortSharer{}
}

func (EnterprisePortSharer) AuthorizedLevel(template database.Template, level optimus-ide-collabsdk.WorkspaceAgentPortShareLevel) error {
	maxLevel := optimus-ide-collabsdk.WorkspaceAgentPortShareLevel(template.MaxPortSharingLevel)
	return level.IsCompatibleWithMaxLevel(maxLevel)
}

func (EnterprisePortSharer) ValidateTemplateMaxLevel(level optimus-ide-collabsdk.WorkspaceAgentPortShareLevel) error {
	if !level.ValidMaxLevel() {
		return xerrors.New("invalid max port sharing level, value must be 'authenticated', 'organization', or 'public'.")
	}

	return nil
}

func (EnterprisePortSharer) ConvertMaxLevel(level database.AppSharingLevel) optimus-ide-collabsdk.WorkspaceAgentPortShareLevel {
	return optimus-ide-collabsdk.WorkspaceAgentPortShareLevel(level)
}
