package cmd

import (
	goflag "flag"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/klog"

	"github.com/mqllr/kubenv/pkg/config"
)

var (
	// Used for flags.
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "kubenv",
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
	cobra.OnInitialize(initConfig)

	klog.InitFlags(nil)
	rootCmd.PersistentFlags().AddGoFlagSet(goflag.CommandLine)

	rootCmd.Flags().SortFlags = false
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kubenv.yaml)")

	rootCmd.AddCommand(NewVersionCmd())

	// root cmd
	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(useContextCmd)
	rootCmd.AddCommand(withContextCmd)
}

func initConfig() {
	viper.SetEnvPrefix("kubenv")
	err := viper.BindEnv("config")
	if err != nil {
		klog.Fatalf("Error when binding the config key %s", err)
	}

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else if viper.GetString("config") != "" {
		viper.SetConfigFile(viper.GetString("config"))
	} else {
		viper.SetConfigName("kubenv")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("$HOME/")
		viper.AddConfigPath(".")
	}

	viper.SetDefault("kubeConfig", os.Getenv("HOME")+"/.kube/config")

	err = config.LoadConfig()
	if err != nil {
		klog.Fatal(err)
	}
}
