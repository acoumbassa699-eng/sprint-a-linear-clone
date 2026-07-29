package agentcontainers

import (
	"context"
	"os"
	"path/filepath"

	"github.com/google/uuid"

	"cdr.dev/slog/v3"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

const (
	// DevcontainerLocalFolderLabel is the label that contains the path to
	// the local workspace folder for a devcontainer.
	DevcontainerLocalFolderLabel = "devcontainer.local_folder"
	// DevcontainerConfigFileLabel is the label that contains the path to
	// the devcontainer.json configuration file.
	DevcontainerConfigFileLabel = "devcontainer.config_file"
	// DevcontainerIsTestRunLabel is set if the devcontainer is part of a test
	// and should be excluded.
	DevcontainerIsTestRunLabel = "devcontainer.is_test_run"
	// The default workspace folder inside the devcontainer.
	DevcontainerDefaultContainerWorkspaceFolder = "/workspaces"
)

func ExtractDevcontainerScripts(
	devcontainers []optimus-ide-collabsdk.WorkspaceAgentDevcontainer,
	scripts []optimus-ide-collabsdk.WorkspaceAgentScript,
) (filteredScripts []optimus-ide-collabsdk.WorkspaceAgentScript, devcontainerScripts map[uuid.UUID]optimus-ide-collabsdk.WorkspaceAgentScript) {
	devcontainerScripts = make(map[uuid.UUID]optimus-ide-collabsdk.WorkspaceAgentScript)
ScriptLoop:
	for _, script := range scripts {
		for _, dc := range devcontainers {
			// The devcontainer scripts match the devcontainer ID for
			// identification.
			if script.ID == dc.ID {
				devcontainerScripts[dc.ID] = script
				continue ScriptLoop
			}
		}

		filteredScripts = append(filteredScripts, script)
	}

	return filteredScripts, devcontainerScripts
}

// ExpandAllDevcontainerPaths expands all devcontainer paths in the given
// devcontainers. This is required by the devcontainer CLI, which requires
// absolute paths for the workspace folder and config path.
func ExpandAllDevcontainerPaths(logger slog.Logger, expandPath func(string) (string, error), devcontainers []optimus-ide-collabsdk.WorkspaceAgentDevcontainer) []optimus-ide-collabsdk.WorkspaceAgentDevcontainer {
	expanded := make([]optimus-ide-collabsdk.WorkspaceAgentDevcontainer, 0, len(devcontainers))
	for _, dc := range devcontainers {
		expanded = append(expanded, expandDevcontainerPaths(logger, expandPath, dc))
	}
	return expanded
}

func expandDevcontainerPaths(logger slog.Logger, expandPath func(string) (string, error), dc optimus-ide-collabsdk.WorkspaceAgentDevcontainer) optimus-ide-collabsdk.WorkspaceAgentDevcontainer {
	logger = logger.With(slog.F("devcontainer", dc.Name), slog.F("workspace_folder", dc.WorkspaceFolder), slog.F("config_path", dc.ConfigPath))

	if wf, err := expandPath(dc.WorkspaceFolder); err != nil {
		logger.Warn(context.Background(), "expand devcontainer workspace folder failed", slog.Error(err))
	} else {
		dc.WorkspaceFolder = wf
	}
	if dc.ConfigPath != "" {
		// Let expandPath handle home directory, otherwise assume relative to
		// workspace folder or absolute.
		if dc.ConfigPath[0] == '~' {
			if cp, err := expandPath(dc.ConfigPath); err != nil {
				logger.Warn(context.Background(), "expand devcontainer config path failed", slog.Error(err))
			} else {
				dc.ConfigPath = cp
			}
		} else {
			dc.ConfigPath = relativePathToAbs(dc.WorkspaceFolder, dc.ConfigPath)
		}
	}
	return dc
}

func relativePathToAbs(workdir, path string) string {
	path = os.ExpandEnv(path)
	if !filepath.IsAbs(path) {
		path = filepath.Join(workdir, path)
	}
	return path
}
