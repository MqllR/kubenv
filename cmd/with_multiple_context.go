package cmd

import (
	"sort"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"k8s.io/klog"

	"github.com/mqllr/kubenv/pkg/config"
	"github.com/mqllr/kubenv/pkg/k8s"
	"github.com/mqllr/kubenv/pkg/prompt"
)

var withMultipleContextsCmd = &cobra.Command{
	Aliases: []string{"wmc"},
	Use:     "with-multiple-contexts command ...",
	Short:   "Execute a command with a k8s context",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		withMultipleContexts(args)
	},
}

// with-context command
func withMultipleContexts(args []string) {
	c, err := k8s.NewKubeConfigFromFile(config.Conf.KubeConfig)
	if err != nil {
		klog.Fatalf("Error when loading kubeconfig file: %s", err)
	}
	contexts := c.GetContextNames()
	klog.V(5).Infof("List of contexts: %v", contexts)

	sort.Strings(contexts)

	p := prompt.NewMultipleSelectPrompt("Select the contexts:", contexts)
	selectedContexts := p.Prompt()

	for _, context := range selectedContexts {
		color.Green("-> Context %s\n", context)
		c.ExecCommand(context, args)
	}
}
