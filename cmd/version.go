package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

const (
	Version = "0.0.2"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the CLI version",
	Run: func(cmd *cobra.Command, args []string) {
		version(io.Writer(os.Stdout))
	},
}

func version(w io.Writer) {
	fmt.Fprintf(w, "%v kubenv v%s\n", promptui.IconGood, Version)
}
