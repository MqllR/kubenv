package cmd

import (
	goflag "flag"
	"fmt"

	"github.com/mqllr/kubenv/cmd/helpers"
	"github.com/mqllr/kubenv/pkg/k8s"
	"github.com/spf13/cobra"
	"k8s.io/klog"
)

var (
	kubeconfig *k8s.KubeConfig

	rootCmd = &cobra.Command{
		Use: "kubenv",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			var err error
			kubeconfig, err = helpers.NewKubeConfig()
			if err != nil {
				klog.Fatalf("Cannot load the kubeconfig: %s", err)
			}
		},
		Short: "A tool to manage multiple Kube cluster",
	}
)

// Execute executes the root command.
func Execute() error {
	err := goflag.Set("logtostderr", "true")
	if err != nil {
		return fmt.Errorf("Error when setting the value to logtostderr %s", err)
	}
	err = goflag.CommandLine.Parse([]string{})
	if err != nil {
		return fmt.Errorf("Error when parsing params %s", err)
	}

	return rootCmd.Execute()
}

func init() {
	klog.InitFlags(nil)
	rootCmd.PersistentFlags().AddGoFlagSet(goflag.CommandLine)

	rootCmd.Flags().SortFlags = false

	rootCmd.AddCommand(versionCmd())

	// root cmd
	rootCmd.AddCommand(syncCommand())
	rootCmd.AddCommand(useContextCmd())
	rootCmd.AddCommand(withContextCmd())

	// show cmd
	s := showCmd()
	s.AddCommand(showUserCmd())
	s.AddCommand(showClusterCmd())

	rootCmd.AddCommand(s)

	// edit cmd
	rootCmd.AddCommand(editCmd())
}
