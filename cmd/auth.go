package cmd

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/klog"
	executil "k8s.io/utils/exec"

	"github.com/mqllr/kubenv/pkg/aws"
	awsgoogleauth "github.com/mqllr/kubenv/pkg/aws-google-auth"
	"github.com/mqllr/kubenv/pkg/config"
	"github.com/mqllr/kubenv/pkg/helper"
	"github.com/mqllr/kubenv/pkg/utils"
)

var userName string

var authCmd = &cobra.Command{
	Use:   "auth [account]",
	Short: "Authentication related tasks",
	Args:  cobra.MaximumNArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		helper.IsConfigExist(
			[]string{
				"authProviders",
				"authAccounts",
			},
		)
	},
	Run: func(cmd *cobra.Command, args []string) {
		auth(args)
	},
}

func init() {
	authCmd.Flags().StringVarP(&userName, "username", "u", "", "The username to authenticate")
}

func auth(args []string) {
	authAccountsConfig, err := config.NewAuthAccountsConfig()
	if err != nil {
		klog.Fatalf("%s", err)
	}

	switch {
	case len(args) == 1 && args[0] == "all":
		for env, account := range authAccountsConfig.Env {
			authAccount(env, account)
		}
	case len(args) == 1 && args[0] != "all":
		authAccount(args[0], authAccountsConfig.Env[args[0]])
	default:
		var items []string
		for env, _ := range authAccountsConfig.Env {
			items = append(items, env)
		}
		item, err := utils.Prompt("Select an account", items)
		if err != nil {
			klog.Fatalf("%s", err)
		}

		authAccount(item, authAccountsConfig.Env[item])
	}
}

func authAccount(env string, account *config.AuthAccount) {
	fmt.Printf("%v Authentication using %s on %s...\n", promptui.IconSelect, account.AuthProvider, env)
	provider := getViperProvider(account.AuthProvider)
	if account.AuthProvider == "aws-google-auth" {
		authWithGoogleAuth(provider, account)
	}
}

func authWithGoogleAuth(authCfg *viper.Viper, account *config.AuthAccount) {
	idp := authCfg.GetString("IDP")
	sp := authCfg.GetString("SP")

	klog.V(2).Infof("Authenticate using aws-google-auth IDP: %s", idp)
	klog.V(2).Infof("Authenticate using aws-google-auth SP: %s", sp)

	auth := awsgoogleauth.NewAWSGoogleAuth(
		idp,
		sp,
	)

	if userName != "" {
		auth.UserName = userName
		klog.V(2).Infof("Authenticate using aws-google-auth UserName: %s", userName)
	}

	auth.AWSRole = account.AWSRole
	auth.AWSProfile = account.AWSProfile
	auth.Region = account.Region

	klog.V(2).Infof("Authenticate using aws-google-auth AWSRole: %s", auth.AWSRole)
	klog.V(2).Infof("Authenticate using aws-google-auth AWSProfile: %s", auth.AWSProfile)
	klog.V(2).Infof("Authenticate using aws-google-auth Region: %s", auth.Region)

	session, err := aws.NewSharedSession(auth.AWSProfile)
	if err != nil {
		klog.Fatalf("Cannot create an AWS session: %s", err)
	}

	if !session.IsExpired() {
		fmt.Printf("%v Your token is still valid. Would you like to renew it? [Y/n] ", promptui.IconInitial)
		var answer string
		fmt.Scanln(&answer)

		if answer != "Y" {
			return
		}
	}

	execer := executil.New()
	runner := awsgoogleauth.New(execer)

	err = runner.Authenticate(auth)
	if err != nil {
		klog.Fatalf("Error on authentication: %s", err)
	}
}

func getViperProvider(provider string) *viper.Viper {
	if !viper.IsSet("authproviders." + provider) {
		klog.Fatalf("Provider %s doesn't exist", provider)
	}

	return viper.Sub("authproviders." + provider)
}
