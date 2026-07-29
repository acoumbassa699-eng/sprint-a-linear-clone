package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/buildinfo"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli/cliui"
	"github.com/optimus-ide-collab/pretty"
	"github.com/optimus-ide-collab/serpent"
)

// versionInfo wraps the stuff we get from buildinfo so that it's
// easier to emit in different formats.
type versionInfo struct {
	Version      string    `json:"version"`
	BuildTime    time.Time `json:"build_time"`
	ExternalURL  string    `json:"external_url"`
	Slim         bool      `json:"slim"`
	AGPL         bool      `json:"agpl"`
	BoringCrypto bool      `json:"boring_crypto"`
}

// String() implements Stringer
func (vi versionInfo) String() string {
	var str strings.Builder
	_, _ = str.WriteString("Optimus-IDE-Collab ")
	if vi.AGPL {
		_, _ = str.WriteString("(AGPL) ")
	}
	_, _ = str.WriteString(vi.Version)
	if vi.BoringCrypto {
		_, _ = str.WriteString(" BoringCrypto")
	}

	if !vi.BuildTime.IsZero() {
		_, _ = str.WriteString(" " + vi.BuildTime.Format(time.UnixDate))
	}
	_, _ = str.WriteString("\r\n" + vi.ExternalURL + "\r\n\r\n")

	if vi.Slim {
		_, _ = str.WriteString(fmt.Sprintf("Slim build of Optimus-IDE-Collab, does not support the %s subcommand.", pretty.Sprint(cliui.DefaultStyles.Code, "server")))
	} else {
		_, _ = str.WriteString(fmt.Sprintf("Full build of Optimus-IDE-Collab, supports the %s subcommand.", pretty.Sprint(cliui.DefaultStyles.Code, "server")))
	}
	return str.String()
}

func defaultVersionInfo() *versionInfo {
	buildTime, _ := buildinfo.Time()
	return &versionInfo{
		Version:      buildinfo.Version(),
		BuildTime:    buildTime,
		ExternalURL:  buildinfo.ExternalURL(),
		Slim:         buildinfo.IsSlim(),
		AGPL:         buildinfo.IsAGPL(),
		BoringCrypto: buildinfo.IsBoringCrypto(),
	}
}

// version prints the optimus-ide-collab version
func (*RootCmd) version(versionInfo func() *versionInfo) *serpent.Command {
	var (
		formatter = cliui.NewOutputFormatter(
			cliui.TextFormat(),
			cliui.JSONFormat(),
		)
		vi = versionInfo()
	)

	cmd := &serpent.Command{
		Use:     "version",
		Short:   "Show optimus-ide-collab version",
		Options: serpent.OptionSet{},
		Handler: func(inv *serpent.Invocation) error {
			out, err := formatter.Format(inv.Context(), vi)
			if err != nil {
				return err
			}

			_, err = fmt.Fprintln(inv.Stdout, out)
			return err
		},
	}

	formatter.AttachOptions(&cmd.Options)

	return cmd
}
