package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var Version = "dev"

type versionOptions struct {
	output string
}

// versionCmd cobra command for version
func versionCmd() *cobra.Command {
	opts := versionOptions{}

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the CLI version",
		RunE: func(cmd *cobra.Command, args []string) error {
			switch opts.output {
			case "json":
				v, _ := json.Marshal(map[string]string{"version": Version})
				fmt.Fprintf(cmd.OutOrStdout(), "%s", string(v))
			case "":
				fmt.Fprintf(cmd.OutOrStdout(), "%v kubenv v%s\n", promptui.IconGood, Version)
			default:
				fmt.Fprintf(cmd.OutOrStdout(), "Unknown output")
			}
			return nil
		},
		Aliases: []string{"v"},
	}

	f := cmd.Flags()
	f.StringVarP(&opts.output, "output", "o", "", "Output: json")

	return cmd
}
