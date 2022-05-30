package cmd

import "github.com/spf13/cobra"

func showCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "show",
		Short:   "show information about contexts",
		Aliases: []string{"sh"},
	}
}
