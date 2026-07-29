package cli

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

func Test_resolveAgentAbsPath(t *testing.T) {
	t.Parallel()

	type args struct {
		workingDirectory string
		relOrAbsPath     string
		agentOS          string
		local            bool
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"ok no args", args{}, "", false},
		{"ok only working directory", args{workingDirectory: "/workdir"}, "/workdir", false},
		{"ok with working directory and rel path", args{workingDirectory: "/workdir", relOrAbsPath: "my/path"}, "/workdir/my/path", false},
		{"ok with working directory and abs path", args{workingDirectory: "/workdir", relOrAbsPath: "/my/path"}, "/my/path", false},
		{"ok with no working directory and abs path", args{relOrAbsPath: "/my/path"}, "/my/path", false},

		{"fail tilde", args{relOrAbsPath: "~"}, "", true},
		{"fail tilde with working directory", args{workingDirectory: "/workdir", relOrAbsPath: "~"}, "", true},
		{"fail tilde path", args{relOrAbsPath: "~/workdir"}, "", true},
		{"fail tilde path with working directory", args{workingDirectory: "/workdir", relOrAbsPath: "~/workdir"}, "", true},
		{"fail relative dot with no working directory", args{relOrAbsPath: "."}, "", true},
		{"fail relative with no working directory", args{relOrAbsPath: "workdir"}, "", true},

		{"ok with working directory and rel path on windows", args{workingDirectory: "C:\\workdir", relOrAbsPath: "my\\path", agentOS: "windows"}, "C:\\workdir\\my\\path", false},
		{"ok with working directory and abs path on windows", args{workingDirectory: "C:\\workdir", relOrAbsPath: "C:\\my\\path", agentOS: "windows"}, "C:\\my\\path", false},
		{"ok with no working directory and abs path on windows", args{relOrAbsPath: "C:\\my\\path", agentOS: "windows"}, "C:\\my\\path", false},
		{"ok abs unix path on windows", args{workingDirectory: "C:\\workdir", relOrAbsPath: "/my/path", agentOS: "windows"}, "\\my\\path", false},
		{"ok rel unix path on windows", args{workingDirectory: "C:\\workdir", relOrAbsPath: "my/path", agentOS: "windows"}, "C:\\workdir\\my\\path", false},

		{"fail with no working directory and rel path on windows", args{relOrAbsPath: "my\\path", agentOS: "windows"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := resolveAgentAbsPath(tt.args.workingDirectory, tt.args.relOrAbsPath, tt.args.agentOS, tt.args.local)
			if (err != nil) != tt.wantErr {
				t.Errorf("resolveAgentAbsPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("resolveAgentAbsPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildAppLinkURL(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name string
		// function arguments
		baseURL           string
		workspace         optimus-ide-collabsdk.Workspace
		agent             optimus-ide-collabsdk.WorkspaceAgent
		app               optimus-ide-collabsdk.WorkspaceApp
		appsHost          string
		preferredPathBase string
		// expected results
		expectedLink string
	}{
		{
			name:    "external url",
			baseURL: "https://optimus-ide-collab.tld",
			app: optimus-ide-collabsdk.WorkspaceApp{
				External: true,
				URL:      "https://external-url.tld",
			},
			expectedLink: "https://external-url.tld",
		},
		{
			name:    "without subdomain",
			baseURL: "https://optimus-ide-collab.tld",
			workspace: optimus-ide-collabsdk.Workspace{
				Name:      "Test-Workspace",
				OwnerName: "username",
			},
			agent: optimus-ide-collabsdk.WorkspaceAgent{
				Name: "a-workspace-agent",
			},
			app: optimus-ide-collabsdk.WorkspaceApp{
				Slug:      "app-slug",
				Subdomain: false,
			},
			preferredPathBase: "/path-base",
			expectedLink:      "https://optimus-ide-collab.tld/path-base/@username/Test-Workspace.a-workspace-agent/apps/app-slug/",
		},
		{
			name:    "with command",
			baseURL: "https://optimus-ide-collab.tld",
			workspace: optimus-ide-collabsdk.Workspace{
				Name:      "Test-Workspace",
				OwnerName: "username",
			},
			agent: optimus-ide-collabsdk.WorkspaceAgent{
				Name: "a-workspace-agent",
			},
			app: optimus-ide-collabsdk.WorkspaceApp{
				Slug:    "my-terminal",
				Command: "ls -la",
			},
			expectedLink: "https://optimus-ide-collab.tld/@username/Test-Workspace.a-workspace-agent/terminal?app=my-terminal",
		},
		{
			name:    "with subdomain",
			baseURL: "ftps://optimus-ide-collab.tld",
			workspace: optimus-ide-collabsdk.Workspace{
				Name:      "Test-Workspace",
				OwnerName: "username",
			},
			agent: optimus-ide-collabsdk.WorkspaceAgent{
				Name: "a-workspace-agent",
			},
			app: optimus-ide-collabsdk.WorkspaceApp{
				Subdomain:     true,
				SubdomainName: "hellooptimus-ide-collab",
			},
			preferredPathBase: "/path-base",
			appsHost:          "*.apps-host.tld",
			expectedLink:      "ftps://hellooptimus-ide-collab.apps-host.tld/",
		},
		{
			name:    "with subdomain, but not apps host",
			baseURL: "https://optimus-ide-collab.tld",
			workspace: optimus-ide-collabsdk.Workspace{
				Name:      "Test-Workspace",
				OwnerName: "username",
			},
			agent: optimus-ide-collabsdk.WorkspaceAgent{
				Name: "a-workspace-agent",
			},
			app: optimus-ide-collabsdk.WorkspaceApp{
				Slug:          "app-slug",
				Subdomain:     true,
				SubdomainName: "It really doesn't matter what this is without AppsHost.",
			},
			preferredPathBase: "/path-base",
			expectedLink:      "https://optimus-ide-collab.tld/path-base/@username/Test-Workspace.a-workspace-agent/apps/app-slug/",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			baseURL, err := url.Parse(tt.baseURL)
			require.NoError(t, err)
			actual := buildAppLinkURL(baseURL, tt.workspace, tt.agent, tt.app, tt.appsHost, tt.preferredPathBase)
			assert.Equal(t, tt.expectedLink, actual)
		})
	}
}
