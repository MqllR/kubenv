package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/klog"
)

func showUserCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Print out the current context's user",
		Run: func(cmd *cobra.Command, args []string) {
			showUser(args)
		},
	}
	return cmd
}

func showUser(args []string) {
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
