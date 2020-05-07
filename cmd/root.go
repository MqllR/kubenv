package cmd

import (
	goflag "flag"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/klog"
)

var (
	// Used for flags.
	cfgFile string
	profile string

	rootCmd = &cobra.Command{
		Use:   "kubenv",
		Short: "A tool to manage authentication on Kube",
	}
)

// Execute executes the root command.
func Execute() error {
	goflag.Set("logtostderr", "true")
	goflag.CommandLine.Parse([]string{})

	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	klog.InitFlags(nil)
	rootCmd.PersistentFlags().AddGoFlagSet(goflag.CommandLine)

	rootCmd.Flags().SortFlags = false
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kubenv.yaml)")

	// root cmd
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(k8sCmd)
	rootCmd.AddCommand(dependencyCmd)

	// k8s sub cmd
	k8sCmd.AddCommand(syncCmd)
	k8sCmd.AddCommand(useContextCmd)

	// dependency sub cmd
	dependencyCmd.AddCommand(depCheckCmd)
	dependencyCmd.AddCommand(depInstallCmd)
}

func initConfig() {
	viper.SetEnvPrefix("kubenv")
	viper.BindEnv("config")

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else if viper.GetString("config") != "" {
		viper.SetConfigFile(viper.GetString("config"))
	} else {
		viper.SetConfigName("kubenv")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("$HOME")
		viper.AddConfigPath(".")
	}

	viper.SetDefault("kubeConfig", os.Getenv("HOME")+"/.kube/config")

	err := viper.ReadInConfig()
	if err != nil {
		klog.Fatalf("Using config file %s: %s", viper.ConfigFileUsed(), err)
	}

	klog.V(5).Infof("Config file content: %s", viper.AllSettings())
}
