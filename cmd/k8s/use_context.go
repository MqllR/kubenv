package k8s

import (
	"sort"

	"github.com/spf13/cobra"
	"k8s.io/klog"

	"github.com/mqllr/kubenv/pkg/config"
	"github.com/mqllr/kubenv/pkg/k8s"
	"github.com/mqllr/kubenv/pkg/prompt"
)

var (
	context   string
	autoLogin bool
)

var UseContextCmd = &cobra.Command{
	Use:   "use-context",
	Short: "Switch from k8s context",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		useContext(args)
	},
}

func init() {
	UseContextCmd.Flags().StringVarP(&context, "context", "c", "", "Kubernetes context to switch")
	UseContextCmd.Flags().BoolVarP(&autoLogin, "auto-login", "a", true, "Auto-login if authAccount is set")
}

// use-context command
func useContext(args []string) {
	kubeconfig, err := k8s.NewKubeConfigFromFile(config.Conf.KubeConfig)
	if err != nil {
		klog.Fatalf("Cannot load kubeconfig file: %s", err)
	}

	contexts := kubeconfig.GetContextNames()

	var selectedContext string

	if context != "" {
		if !kubeconfig.IsContextExist(context) {
			klog.Fatalf("Context %s doesn't exist", context)
		}
		selectedContext = context
	} else {
		sort.Strings(contexts)

		selectedContext, err = prompt.Prompt("Select the context", contexts)
		if err != nil {
			klog.Fatalf("%s", err)
		}
	}

	err = kubeconfig.SetCurrentContext(selectedContext)
	if err != nil {
		klog.Fatalf("Cannot set the current context %s: %s", selectedContext, err)
	}

	kubeconfig.WriteFile(config.Conf.KubeConfig)
}
