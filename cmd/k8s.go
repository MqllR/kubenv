package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/klog"

	"github.com/mqllr/kubenv/pkg/k8s"
	"github.com/mqllr/kubenv/pkg/utils"
)

const k8sConfigFile = "config"

var kubeConfigsPath string

var k8sCmd = &cobra.Command{
	Use:   "k8s",
	Short: "Kubernetes related commands",
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Synchronize the kubernetes config files",
	PreRun: func(cmd *cobra.Command, args []string) {
		isConfigExist(
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

var useContextCmd = &cobra.Command{
	Use:   "use-context [context]",
	Short: "Select a k8s context",
	Args:  cobra.MaximumNArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		isConfigExist(
			[]string{
				"kubeconfig",
			},
		)
	},
	Run: func(cmd *cobra.Command, args []string) {
		useContext(args)
	},
}

func init() {
	syncCmd.Flags().StringVarP(&kubeConfigsPath, "kubeconfigs-path", "p", "", "Path to the directory where all kubeconfigs are stored (required)")
	syncCmd.MarkFlagRequired("kubeconfigs-path")
}

// sync command
func sync(args []string) {
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

	for k8sconfig, _ := range k8sConfigs {
		configPath := path.Join([]string{kubeConfigsPath, k8sconfig, k8sConfigFile}...)
		fmt.Printf("Sync kubeconfig %s", k8sconfig)

		err := existAndDirectory(configPath)
		if err != nil {
			fmt.Printf(" %v %s\n", promptui.IconBad, err)
			continue
		}

		config, err := ioutil.ReadFile(configPath)
		if err != nil {
			fmt.Printf(" %v %s\n", promptui.IconBad, err)
			continue
		}

		kubeconfig := k8s.NewKubeConfig()

		if err = kubeconfig.Unmarshal(config); err != nil {
			fmt.Printf(" %v Cannot unmarshal config %s: %s\n", promptui.IconBad, configPath, err)
			continue
		}

		// Retrieve account's profile
		accountKey := "k8sconfigs." + k8sconfig + ".authAccount"
		if viper.IsSet(accountKey) {
			account := viper.GetString(accountKey)

			if !viper.IsSet("authaccounts." + account) {
				fmt.Printf(" %v account %s doesn't exist\n", promptui.IconWarn, account)
			} else {
				sub := viper.Sub("authaccounts." + account)
				if sub.GetString("authprovider") == "aws-google-auth" {
					awsProfile := sub.GetString("awsprofile")
					klog.V(2).Infof("Using AWS profile %s", awsProfile)

					for _, user := range kubeconfig.Users {
						env := &k8s.Env{Name: "AWS_PROFILE", Value: awsProfile}
						user.User.Exec.Env = []*k8s.Env{env}
					}
				}
			}
		}

		fullConfig.Append(kubeconfig)
		fmt.Printf(" %v\n", promptui.IconGood)
	}

	fullConfig.WriteFile(kubeConfig)
}

// use-context command
func useContext(args []string) {
	var (
		kubeConfig = viper.GetString("kubeconfig")
	)

	kubeconfig := k8s.NewKubeConfig()

	config, err := ioutil.ReadFile(kubeConfig)
	if err != nil {
		klog.Fatalf("Cannot read file %s: %s", kubeConfig, err)
	}

	if err = kubeconfig.Unmarshal(config); err != nil {
		klog.Fatalf("Cannot unmarshal config: %s", err)
	}

	contexts := kubeconfig.GetContextNames()

	var selectedContext string

	if len(args) == 1 {
		exist := func(slice []string, item string) bool {
			for _, s := range slice {
				if item == s {
					return true
				}
			}
			return false
		}

		if !exist(contexts, args[0]) {
			klog.Fatalf("Context %s doesn't exist", args[0])
		}

		selectedContext = args[0]
	} else {
		sort.Strings(contexts)

		selectedContext, err = utils.Prompt("Select the context", contexts)
		if err != nil {
			klog.Fatalf("%s", err)
		}
	}

	kubeconfig.CurrentContext = selectedContext
	kubeconfig.WriteFile(kubeConfig)
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
