package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/klog"
)

func showClusterCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "Print out the current context's cluster",
		Run: func(cmd *cobra.Command, args []string) {
			showCluster(args)
		},
	}

	return cmd
}

func showCluster(args []string) {
	cluster, err := kubeconfig.GetClusterByContextName(kubeconfig.CurrentContext)
	if err != nil {
		klog.Fatalf("Cannot get the cluster: %s", err)
	}

	clusterJ, err := json.MarshalIndent(cluster, "", "  ")
	if err != nil {
		klog.Fatalf("Cannot unmarshal the cluster: %s", err)
	}
	fmt.Println(string(clusterJ))
}
