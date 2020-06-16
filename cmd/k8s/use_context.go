package k8s

import (
	"io/ioutil"
	"sort"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/klog"

	"github.com/mqllr/kubenv/pkg/helper"
	"github.com/mqllr/kubenv/pkg/k8s"
	"github.com/mqllr/kubenv/pkg/utils"
)

var UseContextCmd = &cobra.Command{
	Use:   "use-context [context]",
	Short: "Select a k8s context",
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

// use-context command
func useContext(args []string) {
	var (
		kubeConfig = viper.GetString("kubeconfig")
	)

	kubeconfig := k8s.NewKubeConfig()

	config, err := ioutil.ReadFile(kubeConfig)
	if err != nil {
		klog.Fatalf("Cannot read file %s: %s", kubeConfig, err)
	}

	if err = kubeconfig.Unmarshal(config); err != nil {
		klog.Fatalf("Cannot unmarshal config: %s", err)
	}

	contexts := kubeconfig.GetContextNames()

	var selectedContext string

	if len(args) == 1 {
		exist := func(slice []string, item string) bool {
			for _, s := range slice {
				if item == s {
					return true
				}
			}
			return false
		}

		if !exist(contexts, args[0]) {
			klog.Fatalf("Context %s doesn't exist", args[0])
		}

		selectedContext = args[0]
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
