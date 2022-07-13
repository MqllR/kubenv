package cmd

import (
	"sort"

	"github.com/spf13/cobra"
	"k8s.io/klog"

	"github.com/mqllr/kubenv/pkg/prompt"
)

type useContextOptions struct {
	context string
}

func useContextCmd() *cobra.Command {
	opts := useContextOptions{}

	cmd := &cobra.Command{
		Use:     "use-context",
		Short:   "Switch to k8s context",
		Aliases: []string{"uc"},
		Run: func(cmd *cobra.Command, args []string) {
			useContext(&opts)
		},
	}

	f := cmd.Flags()
	f.StringVarP(&opts.context, "context", "c", "", "Kubernetes context to switch")

	return cmd
}

// use-context command
func useContext(opts *useContextOptions) {
	contexts := kubeconfig.GetContextNames()

	var selectedContext string
	var err error

	if opts.context != "" {
		if !kubeconfig.IsContextExist(opts.context) {
			klog.Fatalf("Context %s doesn't exist", opts.context)
		}
		selectedContext = opts.context
	} else {
		sort.Strings(contexts)

		p := prompt.NewPrompt("Select the context", contexts)
		selectedContext, err = p.PromptSelect()
		if err != nil {
			klog.Fatalf("Cannot get the answer from the prompt: %s", err)
		}
	}

	err = kubeconfig.SetCurrentContext(selectedContext)
	if err != nil {
		klog.Fatalf("Cannot set the current context %s: %s", selectedContext, err)
	}

	err = kubeconfig.Save()
	if err != nil {
		klog.Fatalf("Cannot write the kubeconfig file: %s", err)
	}
}
