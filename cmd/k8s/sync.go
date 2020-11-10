package k8s

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"k8s.io/klog"

	"github.com/mqllr/kubenv/pkg/config"
	"github.com/mqllr/kubenv/pkg/helper"
	"github.com/mqllr/kubenv/pkg/k8s"
)

const k8sConfigFile = "config"

var SyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Synchronize the kubernetes config files",
	PreRun: func(cmd *cobra.Command, args []string) {
		helper.IsConfigExist(
			[]string{
				"kubeconfig",
				"k8sconfigs",
			},
		)
	},
	Run: func(cmd *cobra.Command, args []string) {
		sync(args)
	},
}

func sync(args []string) {
	fmt.Printf("%v Start the synchronization of kubeconfig file into %s ...\n", promptui.IconSelect, config.Conf.KubeConfig)

	var err error
	fullConfig := k8s.NewKubeConfig()

	for name, conf := range config.Conf.K8SConfigs {
		fmt.Printf("Sync kubeconfig %s", name)

		var k *k8s.KubeConfig

		if conf.Sync.Mode == "local" {
			klog.V(2).Info("Sync start in local mode")

			k, err = k8s.NewKubeConfigFromFile(conf.Sync.Path)
			if err != nil {
				fmt.Printf(" %v Error when loading the kubeconfig file %s: %s\n", promptui.IconBad, name, err)
				continue
			}
		}

		if conf.AuthAccount != "" {
			account := config.Conf.FindAuthAccount(conf.AuthAccount)

			for _, user := range k.Users {
				env := &k8s.Env{Name: "AWS_PROFILE", Value: account.AWSProfile}
				user.User.Exec.Env = []*k8s.Env{env}
			}
		}

		fullConfig.Append(k)
		fmt.Printf(" %v\n", promptui.IconGood)
	}

	fullConfig.WriteFile(config.Conf.KubeConfig)
}
