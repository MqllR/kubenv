package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/mqllr/kubenv/pkg/config"
	"github.com/mqllr/kubenv/pkg/k8s"
	"github.com/spf13/cobra"
	"k8s.io/klog"
)

var showUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Print out the current context's user",
	Run: func(cmd *cobra.Command, args []string) {
		showUser(args)
	},
}

func showUser(args []string) {
	kubeconfig, err := k8s.NewKubeConfigFromFile(config.Conf.KubeConfig)
	if err != nil {
		klog.Fatalf("Cannot load the kubeconfig file: %s", err)
	}

	user, err := kubeconfig.GetUserByContextName(kubeconfig.CurrentContext)
	if err != nil {
		klog.Fatalf("Cannot get the user: %s", err)
	}

	userJ, err := json.MarshalIndent(user, "", "  ")
	if err != nil {
		klog.Fatalf("Cannot unmarshal the cluster: %s", err)
	}
	fmt.Println(string(userJ))
}
