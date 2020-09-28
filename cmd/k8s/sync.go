package k8s

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	kubeConfig := viper.GetString("kubeconfig")

	k8sConfigs, err := config.NewK8SConfigs()
	if err != nil {
		klog.Fatalf("Error when loading k8sConfigs: %s", err)
	}

	authAccounts, err := config.NewAuthAccountsConfig()
	if err != nil {
		klog.Fatal("Error when loading the authAccounts")
	}

	fmt.Printf("%v Start the synchronization of kubeconfig file into %s ...\n", promptui.IconSelect, kubeConfig)

	fullConfig := k8s.NewKubeConfig()

	for name, config := range k8sConfigs.Configs {
		fmt.Printf("Sync kubeconfig %s", name)

		var k *k8s.KubeConfig

		if config.Sync.Mode == "local" {
			klog.V(2).Info("Sync start in local mode")

			k, err = k8s.NewKubeConfigFromFile(config.Sync.Path)
			if err != nil {
				fmt.Printf(" %v Error when loading the kubeconfig file %s: %s\n", promptui.IconBad, name, err)
				continue
			}
		}

		if config.AuthAccount != "" {
			account := authAccounts.FindAuthAccount(config.AuthAccount)

			for _, user := range k.Users {
				env := &k8s.Env{Name: "AWS_PROFILE", Value: account.AWSProfile}
				user.User.Exec.Env = []*k8s.Env{env}
			}
		}

		fullConfig.Append(k)
		fmt.Printf(" %v\n", promptui.IconGood)
	}

	fullConfig.WriteFile(kubeConfig)
}
