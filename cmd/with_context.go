package cmd

import (
	"sort"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"k8s.io/klog"

	"github.com/mqllr/kubenv/pkg/prompt"
)

func withContextCmd() *cobra.Command {
	cmd := &cobra.Command{
		Aliases: []string{"wc"},
		Use:     "with-context command ...",
		Short:   "Execute a command with one or multiple k8s context",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			withContext(args)
		},
	}
	return cmd
}

// with-context command
func withContext(args []string) {
	contexts := kubeconfig.GetContextNames()
	sort.Strings(contexts)

	klog.V(5).Infof("List of contexts: %v", contexts)

	p := prompt.NewPrompt("Select the contexts:", contexts)
	selectedContexts, err := p.PromptMultipleSelect()
	if err != nil {
		klog.Fatalf("Cannot get the answer from the prompt: %s", err)
	}

	for _, context := range selectedContexts {
		color.Green("-> Context %s\n", context)
		err := kubeconfig.ExecCommand(context, args)
		if err != nil {
			klog.Errorf("Cmd error: %s", err)
		}
	}
}
