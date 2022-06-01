package cmd

import (
	"os"
	"sort"

	"github.com/spf13/cobra"
	"k8s.io/klog"

	"github.com/mqllr/kubenv/pkg/config"
	"github.com/mqllr/kubenv/pkg/k8s"
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
	f, err := os.Open(config.GetKubeConfig())
	if err != nil {
		klog.Fatalf("Cannot open the kube config: %s", err)
	}

	kubeconfig, err := k8s.NewKubeConfigFromReader(f)
	if err != nil {
		klog.Fatalf("Cannot load kubeconfig file: %s", err)
	}

	contexts := kubeconfig.GetContextNames()

	var selectedContext string

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

	fh, err := os.OpenFile(config.GetKubeConfig(), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		klog.Fatalf("Cannot open the kubeconfig: %s", err)
	}
	defer fh.Close()

	err = kubeconfig.Save(fh)
	if err != nil {
		klog.Fatalf("Cannot write the kubeconfig file: %s", err)
	}
}
