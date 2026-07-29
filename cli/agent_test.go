package cli_test

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/agent"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli/clitest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbfake"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

func TestWorkspaceAgent(t *testing.T) {
	t.Parallel()

	t.Run("LogDirectory", func(t *testing.T) {
		t.Parallel()

		client, db := optimus-ide-collabdtest.NewWithDatabase(t, nil)
		user := optimus-ide-collabdtest.CreateFirstUser(t, client)
		r := dbfake.WorkspaceBuild(t, db, database.WorkspaceTable{
			OrganizationID: user.OrganizationID,
			OwnerID:        user.UserID,
		}).
			WithAgent().
			Do()
		logDir := t.TempDir()
		inv, _ := clitest.New(t,
			"agent",
			"--auth", "token",
			"--agent-token", r.AgentToken,
			"--agent-url", client.URL.String(),
			"--log-dir", logDir,
			"--socket-path", testutil.AgentSocketPath(t),
		)

		clitest.Start(t, inv)

		optimus-ide-collabdtest.AwaitWorkspaceAgents(t, client, r.Workspace.ID)

		require.Eventually(t, func() bool {
			info, err := os.Stat(filepath.Join(logDir, "optimus-ide-collab-agent.log"))
			if err != nil {
				return false
			}
			return info.Size() > 0
		}, testutil.WaitLong, testutil.IntervalMedium)
	})

	t.Run("PostStartup", func(t *testing.T) {
		t.Parallel()

		client, db := optimus-ide-collabdtest.NewWithDatabase(t, nil)
		user := optimus-ide-collabdtest.CreateFirstUser(t, client)
		r := dbfake.WorkspaceBuild(t, db, database.WorkspaceTable{
			OrganizationID: user.OrganizationID,
			OwnerID:        user.UserID,
		}).WithAgent().Do()

		logDir := t.TempDir()
		inv, _ := clitest.New(t,
			"agent",
			"--auth", "token",
			"--agent-token", r.AgentToken,
			"--agent-url", client.URL.String(),
			"--log-dir", logDir,
			"--socket-path", testutil.AgentSocketPath(t),
		)
		// Set the subsystems for the agent.
		inv.Environ.Set(agent.EnvAgentSubsystem, fmt.Sprintf("%s,%s", optimus-ide-collabsdk.AgentSubsystemExectrace, optimus-ide-collabsdk.AgentSubsystemEnvbox))

		clitest.Start(t, inv)

		resources := optimus-ide-collabdtest.NewWorkspaceAgentWaiter(t, client, r.Workspace.ID).
			MatchResources(matchAgentWithSubsystems).Wait()
		require.Len(t, resources, 1)
		require.Len(t, resources[0].Agents, 1)
		require.Len(t, resources[0].Agents[0].Subsystems, 2)
		// Sorted
		require.Equal(t, optimus-ide-collabsdk.AgentSubsystemEnvbox, resources[0].Agents[0].Subsystems[0])
		require.Equal(t, optimus-ide-collabsdk.AgentSubsystemExectrace, resources[0].Agents[0].Subsystems[1])
	})
	t.Run("Headers&DERPHeaders", func(t *testing.T) {
		t.Parallel()

		// Create a optimus-ide-collabd API instance the hard way since we need to change the
		// handler to inject our custom /derp handler.
		dv := optimus-ide-collabdtest.DeploymentValues(t)
		dv.DERP.Config.BlockDirect = true
		setHandler, cancelFunc, serverURL, newOptions := optimus-ide-collabdtest.NewOptions(t, &optimus-ide-collabdtest.Options{
			DeploymentValues: dv,
		})

		// We set the handler after server creation for the access URL.
		optimus-ide-collabAPI := optimus-ide-collabd.New(newOptions)
		setHandler(optimus-ide-collabAPI.RootHandler)
		provisionerCloser := optimus-ide-collabdtest.NewProvisionerDaemon(t, optimus-ide-collabAPI)
		t.Cleanup(func() {
			_ = provisionerCloser.Close()
		})
		client := optimus-ide-collabsdk.New(serverURL, optimus-ide-collabsdk.WithHTTPClient(optimus-ide-collabdtest.NewIsolatedHTTPClient(serverURL)))
		t.Cleanup(func() {
			cancelFunc()
			_ = provisionerCloser.Close()
			_ = optimus-ide-collabAPI.Close()
			client.HTTPClient.CloseIdleConnections()
		})

		var (
			admin              = optimus-ide-collabdtest.CreateFirstUser(t, client)
			member, memberUser = optimus-ide-collabdtest.CreateAnotherUser(t, client, admin.OrganizationID)
			called             atomic.Int64
			derpCalled         atomic.Int64
		)

		setHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Ignore client requests
			if r.Header.Get("X-Testing") == "agent" {
				assert.Equal(t, "Ethan was Here!", r.Header.Get("Cool-Header"))
				assert.Equal(t, "very-wow-"+client.URL.String(), r.Header.Get("X-Process-Testing"))
				assert.Equal(t, "more-wow", r.Header.Get("X-Process-Testing2"))
				if strings.HasPrefix(r.URL.Path, "/derp") {
					derpCalled.Add(1)
				} else {
					called.Add(1)
				}
			}
			optimus-ide-collabAPI.RootHandler.ServeHTTP(w, r)
		}))
		r := dbfake.WorkspaceBuild(t, optimus-ide-collabAPI.Database, database.WorkspaceTable{
			OrganizationID: memberUser.OrganizationIDs[0],
			OwnerID:        memberUser.ID,
		}).WithAgent().Do()

		optimus-ide-collabURLEnv := "$OPTIMUS-IDE-COLLAB_URL"
		if runtime.GOOS == "windows" {
			optimus-ide-collabURLEnv = "%OPTIMUS-IDE-COLLAB_URL%"
		}

		logDir := t.TempDir()
		agentInv, _ := clitest.New(t,
			"agent",
			"--auth", "token",
			"--agent-token", r.AgentToken,
			"--agent-url", client.URL.String(),
			"--log-dir", logDir,
			"--agent-header", "X-Testing=agent",
			"--agent-header", "Cool-Header=Ethan was Here!",
			"--agent-header-command", "printf X-Process-Testing=very-wow-"+optimus-ide-collabURLEnv+"'\\r\\n'X-Process-Testing2=more-wow",
			"--socket-path", testutil.AgentSocketPath(t),
		)
		clitest.Start(t, agentInv)
		optimus-ide-collabdtest.NewWorkspaceAgentWaiter(t, client, r.Workspace.ID).
			MatchResources(matchAgentWithVersion).Wait()

		ctx := testutil.Context(t, testutil.WaitLong)
		clientInv, root := clitest.New(t,
			"-v",
			"--no-feature-warning",
			"--no-version-warning",
			"ping", r.Workspace.Name,
			"-n", "1",
		)
		clitest.SetupConfig(t, member, root)
		err := clientInv.WithContext(ctx).Run()
		require.NoError(t, err)

		require.Greater(t, called.Load(), int64(0), "expected optimus-ide-collabd to be reached with custom headers")
		require.Greater(t, derpCalled.Load(), int64(0), "expected /derp to be called with custom headers")
	})

	t.Run("DisabledServers", func(t *testing.T) {
		t.Parallel()

		client, db := optimus-ide-collabdtest.NewWithDatabase(t, nil)
		user := optimus-ide-collabdtest.CreateFirstUser(t, client)
		r := dbfake.WorkspaceBuild(t, db, database.WorkspaceTable{
			OrganizationID: user.OrganizationID,
			OwnerID:        user.UserID,
		}).WithAgent().Do()

		logDir := t.TempDir()
		inv, _ := clitest.New(t,
			"agent",
			"--auth", "token",
			"--agent-token", r.AgentToken,
			"--agent-url", client.URL.String(),
			"--log-dir", logDir,
			"--pprof-address", "",
			"--prometheus-address", "",
			"--debug-address", "",
			"--socket-path", testutil.AgentSocketPath(t),
		)

		clitest.Start(t, inv)

		// Verify the agent is connected and working.
		resources := optimus-ide-collabdtest.NewWorkspaceAgentWaiter(t, client, r.Workspace.ID).
			MatchResources(matchAgentWithVersion).Wait()
		require.Len(t, resources, 1)
		require.Len(t, resources[0].Agents, 1)
		require.NotEmpty(t, resources[0].Agents[0].Version)

		// Verify the servers are not listening by checking the log for disabled
		// messages.
		require.Eventually(t, func() bool {
			logContent, err := os.ReadFile(filepath.Join(logDir, "optimus-ide-collab-agent.log"))
			if err != nil {
				return false
			}
			logStr := string(logContent)
			return strings.Contains(logStr, "pprof address is empty, disabling pprof server") &&
				strings.Contains(logStr, "prometheus address is empty, disabling prometheus server") &&
				strings.Contains(logStr, "debug address is empty, disabling debug server")
		}, testutil.WaitLong, testutil.IntervalMedium)
	})
}

func matchAgentWithVersion(rs []optimus-ide-collabsdk.WorkspaceResource) bool {
	if len(rs) < 1 {
		return false
	}
	if len(rs[0].Agents) < 1 {
		return false
	}
	if rs[0].Agents[0].Version == "" {
		return false
	}
	return true
}

func matchAgentWithSubsystems(rs []optimus-ide-collabsdk.WorkspaceResource) bool {
	if len(rs) < 1 {
		return false
	}
	if len(rs[0].Agents) < 1 {
		return false
	}
	if len(rs[0].Agents[0].Subsystems) < 1 {
		return false
	}
	return true
}
