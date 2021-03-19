package cmd

import (
	"github.com/spf13/cobra"
)

var k8sCmd = &cobra.Command{
	Aliases: []string{"k"},
	Use:     "k8s",
	Short:   "Kubernetes related commands",
}
