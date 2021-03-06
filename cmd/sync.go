package cmd

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"k8s.io/klog"

	"github.com/mqllr/kubenv/pkg/config"
	"github.com/mqllr/kubenv/pkg/k8s"
	k8ssync "github.com/mqllr/kubenv/pkg/sync"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Synchronize the kubernetes config files",
	Run: func(cmd *cobra.Command, args []string) {
		sync(args)
	},
}

func sync(args []string) {
	fmt.Printf("%v Start the synchronization of kubeconfig file into %s ...\n", promptui.IconSelect, config.Conf.KubeConfig)

	fullConfig := k8s.NewKubeConfig()

	for name, conf := range config.Conf.K8SConfigs {
		fmt.Printf("Sync kubeconfig %s", name)

		s, err := k8ssync.NewService(*conf.Sync)
		if err != nil {
			fmt.Printf(" %v\n", promptui.IconBad)
			klog.V(2).Infof("Cannot sync: %s", err)
			continue
		}

		err = s.AppendKubeConfig(fullConfig)
		if err != nil {
			fmt.Printf(" %v\n", promptui.IconBad)
			klog.V(2).Infof("Error when getting the config back: %s", err)
		} else {
			fmt.Printf(" %v\n", promptui.IconGood)
		}
	}

	err := fullConfig.WriteFile(config.Conf.KubeConfig)
	if err != nil {
		fmt.Printf("%v Failed to write the kubeconfig file: %s", promptui.IconBad, err)
	}
}
