package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

const (
	Version = "0.1.0"
)

var output string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the CLI version",
	Run: func(cmd *cobra.Command, args []string) {
		version(io.Writer(os.Stdout))
	},
	Aliases: []string{"v"},
}

func init() {
	versionCmd.Flags().StringVarP(&output, "output", "o", "", "Output: json")
}

func version(w io.Writer) {
	switch output {
	case "json":
		v, _ := json.Marshal(map[string]string{"version": Version})
		fmt.Fprintf(w, string(v))
	case "":
		fmt.Fprintf(w, "%v kubenv v%s\n", promptui.IconGood, Version)
	default:
		fmt.Fprintf(w, "Unknown output")
	}
}
