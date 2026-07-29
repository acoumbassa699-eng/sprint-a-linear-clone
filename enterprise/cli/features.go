package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/xerrors"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cli/cliui"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/serpent"
)

func (r *RootCmd) features() *serpent.Command {
	cmd := &serpent.Command{
		Short:   "List Enterprise features",
		Use:     "features",
		Aliases: []string{"feature"},
		Handler: func(inv *serpent.Invocation) error {
			return inv.Command.HelpHandler(inv)
		},
		Children: []*serpent.Command{
			r.featuresList(),
		},
	}
	return cmd
}

func (r *RootCmd) featuresList() *serpent.Command {
	var (
		featureColumns = []string{"name", "entitlement", "enabled", "limit", "actual"}
		columns        []string
		outputFormat   string
	)

	cmd := &serpent.Command{
		Use:        "list",
		Aliases:    []string{"ls"},
		Middleware: serpent.Chain(),
		Handler: func(inv *serpent.Invocation) error {
			client, err := r.InitClient(inv)
			if err != nil {
				return err
			}
			entitlements, err := client.Entitlements(inv.Context())
			var apiError *optimus-ide-collabsdk.Error
			if errors.As(err, &apiError) && apiError.StatusCode() == http.StatusNotFound {
				return xerrors.New("You are on the AGPL licensed version of Optimus-IDE-Collab that does not have Enterprise functionality!")
			}
			if err != nil {
				return err
			}

			// This uses custom formatting as the JSON output outputs an object
			// as opposed to a list from the table output.
			out := ""
			switch outputFormat {
			case "table", "":
				out, err = displayFeatures(columns, entitlements.Features)
				if err != nil {
					return xerrors.Errorf("render table: %w", err)
				}
			case "json":
				buf := new(bytes.Buffer)
				enc := json.NewEnoptimus-ide-collab(buf)
				enc.SetIndent("", "  ")
				err = enc.Encode(entitlements)
				if err != nil {
					return xerrors.Errorf("marshal features to JSON: %w", err)
				}
				out = buf.String()
			default:
				return xerrors.Errorf(`unknown output format %q, only "table" and "json" are supported`, outputFormat)
			}

			_, err = fmt.Fprintln(inv.Stdout, out)
			return err
		},
	}

	cmd.Options = serpent.OptionSet{
		{
			Flag:          "column",
			FlagShorthand: "c",
			Description:   "Specify columns to filter in the table.",
			Default:       strings.Join(featureColumns, ","),
			Value:         serpent.EnumArrayOf(&columns, featureColumns...),
		},
		{
			Flag:          "output",
			FlagShorthand: "o",
			Description:   "Output format.",
			Default:       "table",
			Value:         serpent.EnumOf(&outputFormat, "table", "json"),
		},
	}

	return cmd
}

type featureRow struct {
	Name        optimus-ide-collabsdk.FeatureName `table:"name,default_sort"`
	Entitlement string               `table:"entitlement"`
	Enabled     bool                 `table:"enabled"`
	Limit       *int64               `table:"limit"`
	Actual      *int64               `table:"actual"`
}

// displayFeatures will return a table displaying all features passed in.
// filterColumns must be a subset of the feature fields and will determine which
// columns to display
func displayFeatures(filterColumns []string, features map[optimus-ide-collabsdk.FeatureName]optimus-ide-collabsdk.Feature) (string, error) {
	rows := make([]featureRow, 0, len(features))
	for name, feat := range features {
		rows = append(rows, featureRow{
			Name:        name,
			Entitlement: string(feat.Entitlement),
			Enabled:     feat.Enabled,
			Limit:       feat.Limit,
			Actual:      feat.Actual,
		})
	}

	return cliui.DisplayTable(rows, "name", filterColumns)
}
