package portsharing

import (
	"golang.org/x/xerrors"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

type PortSharer interface {
	AuthorizedLevel(template database.Template, level optimus-ide-collabsdk.WorkspaceAgentPortShareLevel) error
	ValidateTemplateMaxLevel(level optimus-ide-collabsdk.WorkspaceAgentPortShareLevel) error
	ConvertMaxLevel(level database.AppSharingLevel) optimus-ide-collabsdk.WorkspaceAgentPortShareLevel
}

type AGPLPortSharer struct{}

func (AGPLPortSharer) AuthorizedLevel(_ database.Template, _ optimus-ide-collabsdk.WorkspaceAgentPortShareLevel) error {
	return nil
}

func (AGPLPortSharer) ValidateTemplateMaxLevel(_ optimus-ide-collabsdk.WorkspaceAgentPortShareLevel) error {
	return xerrors.New("Restricting port sharing level is an enterprise feature that is not enabled.")
}

func (AGPLPortSharer) ConvertMaxLevel(_ database.AppSharingLevel) optimus-ide-collabsdk.WorkspaceAgentPortShareLevel {
	return optimus-ide-collabsdk.WorkspaceAgentPortShareLevelPublic
}

var DefaultPortSharer PortSharer = AGPLPortSharer{}
