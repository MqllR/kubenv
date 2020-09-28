package k8s

import (
	"sort"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/klog"

	"github.com/mqllr/kubenv/pkg/helper"
	"github.com/mqllr/kubenv/pkg/k8s"
	"github.com/mqllr/kubenv/pkg/utils"
)

var (
	context   string
	autoLogin bool
)

var UseContextCmd = &cobra.Command{
	Use:   "use-context",
	Short: "Switch from k8s context",
	Args:  cobra.MaximumNArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		helper.IsConfigExist(
			[]string{
				"kubeconfig",
			},
		)
	},
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
	kubeConfig := viper.GetString("kubeconfig")
	kubeconfig, err := k8s.NewKubeConfigFromFile(kubeConfig)

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

		selectedContext, err = utils.Prompt("Select the context", contexts)
		if err != nil {
			klog.Fatalf("%s", err)
		}
	}

	kubeconfig.CurrentContext = selectedContext
	kubeconfig.WriteFile(kubeConfig)
}
