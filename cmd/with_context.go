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

var withContextCmd = &cobra.Command{
	Aliases: []string{"wc"},
	Use:     "with-context command ...",
	Short:   "Execute a command with one or multiple k8s context",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		withContext(args)
	},
}

// with-context command
func withContext(args []string) {
	c, err := k8s.NewKubeConfigFromFile(config.Conf.KubeConfig)
	if err != nil {
		klog.Fatalf("Error when loading kubeconfig file: %s", err)
	}
	contexts := c.GetContextNames()
	sort.Strings(contexts)

	klog.V(5).Infof("List of contexts: %v", contexts)

	p := prompt.NewPrompt("Select the contexts:", contexts)
	selectedContexts, err := p.PromptMultipleSelect()
	if err != nil {
		klog.Fatalf("Cannot get the answer from the prompt: %s", err)
	}

	for _, context := range selectedContexts {
		color.Green("-> Context %s\n", context)
		err := c.ExecCommand(context, args)
		if err != nil {
			klog.Errorf("Cmd error: %s", err)
		}
	}
}
