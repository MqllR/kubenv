package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

const (
	version = "1.2.2"
)

var output string

// versionCmd cobra command for version
func versionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the CLI version",
		RunE: func(cmd *cobra.Command, args []string) error {
			switch output {
			case "json":
				v, _ := json.Marshal(map[string]string{"version": version})
				fmt.Fprintf(cmd.OutOrStdout(), "%s", string(v))
			case "":
				fmt.Fprintf(cmd.OutOrStdout(), "%v kubenv v%s\n", promptui.IconGood, version)
			default:
				fmt.Fprintf(cmd.OutOrStdout(), "Unknown output")
			}
			return nil
		},
		Aliases: []string{"v"},
	}

	f := cmd.Flags()
	f.StringVarP(&output, "output", "o", "", "Output: json")

	return cmd
}
