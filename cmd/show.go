package cmd

import "github.com/spf13/cobra"

var showCmd = &cobra.Command{
	Use:     "show",
	Short:   "show different information",
	Aliases: []string{"sh"},
}
