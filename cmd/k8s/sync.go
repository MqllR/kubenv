package k8s

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/klog"

	"github.com/mqllr/kubenv/pkg/config"
	"github.com/mqllr/kubenv/pkg/helper"
	"github.com/mqllr/kubenv/pkg/k8s"
)

const k8sConfigFile = "config"

var kubeConfigsPath string

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

func init() {
	SyncCmd.Flags().StringVarP(&kubeConfigsPath, "kubeconfigs-path", "p", "", "Path to the directory where all kubeconfigs are stored (required)")
	SyncCmd.MarkFlagRequired("kubeconfigs-path")
}

// sync command
func sync(args []string) {
	// TODO having a config object for k8sconfigs
	var (
		k8sConfigs = viper.GetStringMap("k8sconfigs")
		kubeConfig = viper.GetString("kubeconfig")
	)

	err := existAndDirectory(kubeConfigsPath)
	if err != nil {
		klog.Fatalf("%s", err)
	}

	fmt.Printf("%v Start the synchronization of kubeconfigs into %s ...\n", promptui.IconSelect, kubeConfig)

	fullConfig := k8s.NewKubeConfig()

	authAccounts, err := config.NewAuthAccountsConfig()
	if err != nil {
		klog.Fatal("Error on loading the authAccounts")
	}

	for k8sconfig := range k8sConfigs {
		configPath := path.Join([]string{kubeConfigsPath, k8sconfig, k8sConfigFile}...)
		fmt.Printf("Sync kubeconfig %s", k8sconfig)

		err := existAndDirectory(configPath)
		if err != nil {
			fmt.Printf(" %v %s\n", promptui.IconBad, err)
			continue
		}

		kconfig, err := ioutil.ReadFile(configPath)
		if err != nil {
			fmt.Printf(" %v %s\n", promptui.IconBad, err)
			continue
		}

		kubeconfig := k8s.NewKubeConfig()

		if err = kubeconfig.Unmarshal(kconfig); err != nil {
			fmt.Printf(" %v Cannot unmarshal config %s: %s\n", promptui.IconBad, configPath, err)
			continue
		}

		// Retrieve account's profile
		accountKey := "k8sconfigs." + k8sconfig + ".authAccount"
		if viper.IsSet(accountKey) {
			account := authAccounts.FindAuthAccount(viper.GetString(accountKey))

			for _, user := range kubeconfig.Users {
				env := &k8s.Env{Name: "AWS_PROFILE", Value: account.AWSProfile}
				user.User.Exec.Env = []*k8s.Env{env}
			}
		}

		fullConfig.Append(kubeconfig)
		fmt.Printf(" %v\n", promptui.IconGood)
	}

	fullConfig.WriteFile(kubeConfig)
}

func existAndDirectory(path string) error {
	info, err := os.Stat(kubeConfigsPath)

	if os.IsNotExist(err) {
		return fmt.Errorf("File %s doesn't exist: %s", path, err)
	}

	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", path)
	}

	return nil
}
