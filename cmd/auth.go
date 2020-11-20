package cmd

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"k8s.io/klog"

	authenticate "github.com/mqllr/kubenv/pkg/auth"
	"github.com/mqllr/kubenv/pkg/aws"
	awsazurelogin "github.com/mqllr/kubenv/pkg/aws-azure-login"
	awsgoogleauth "github.com/mqllr/kubenv/pkg/aws-google-auth"
	awssts "github.com/mqllr/kubenv/pkg/aws-sts"
	"github.com/mqllr/kubenv/pkg/config"
	"github.com/mqllr/kubenv/pkg/helper"
	"github.com/mqllr/kubenv/pkg/prompt"
)

var (
	account string
	all     bool
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication related tasks",
	Run: func(cmd *cobra.Command, args []string) {
		auth()
	},
}

func init() {
	authCmd.Flags().StringVarP(&account, "account", "a", "", "Account name to authenticate")
	authCmd.Flags().BoolVarP(&all, "all", "", false, "Authenticate all account")
}

func auth() {
	switch {
	case all:
		for _, auth := range config.Conf.ListAuthAccountNames() {
			authAccount(config.Conf.FindAuthAccount(auth))
		}
		break
	case account != "":
		authAccount(config.Conf.FindAuthAccount(account))
		break
	default:
		item, err := prompt.Prompt("Select an account", config.Conf.ListAuthAccountNames())
		if err != nil {
			klog.Fatalf("%s", err)
		}

		authAccount(config.Conf.FindAuthAccount(item))
	}
}

func authAccount(account *config.AuthAccount) {
	fmt.Printf("%v Authentication using %s...\n", promptui.IconSelect, account.AuthProvider)

	if !helper.IsExpired(account) {
		fmt.Printf("%v Token already active. Skipping.\n", promptui.IconGood)
		return
	}

	var acc authenticate.Auth
	switch account.AuthProvider {
	case "aws-google-auth":
		acc = authWithGoogleAuth(config.Conf.FindAuthProvider("aws-google-auth"), account)
	case "aws-azure-login":
		acc = authWithAzureLogin(config.Conf.FindAuthProvider("aws-azure-login"), account)
	case "aws-sts":
		acc = authWithAWSSTS(config.Conf.FindAuthProvider("aws-sts"), account)
	}

	authenticate.Authenticate(acc)
}

func authWithGoogleAuth(provider *config.AuthProvider, account *config.AuthAccount) authenticate.Auth {
	idp := provider.IDP
	sp := provider.SP
	username := provider.UserName
	awsRole := account.AWSRole
	awsProfile := account.AWSProfile

	a := awsgoogleauth.NewAWSGoogleAuth(
		idp,
		sp,
		username,
		awsProfile,
		awsRole,
	)

	a.SetAWSRegion(account.AWSRegion)
	a.SetDuration(account.Duration)

	klog.V(2).Infof("Authenticate using aws-google-auth IDP: %s", idp)
	klog.V(2).Infof("Authenticate using aws-google-auth SP: %s", sp)
	klog.V(2).Infof("Authenticate using aws-google-auth UserName: %s", username)
	klog.V(2).Infof("Authenticate using aws-google-auth AWSRole: %s", awsRole)
	klog.V(2).Infof("Authenticate using aws-google-auth AWSProfile: %s", awsProfile)
	klog.V(2).Infof("Authenticate using aws-google-auth Region: %s", a.AWSRegion)
	klog.V(2).Infof("Authenticate using aws-azure-login Duration: %d", a.Duration)

	return a
}

func authWithAzureLogin(provider *config.AuthProvider, account *config.AuthAccount) authenticate.Auth {
	tid := provider.TenantID
	appid := provider.AppIDUri
	username := provider.UserName
	awsRole := account.AWSRole
	awsProfile := account.AWSProfile

	a := awsazurelogin.NewAWSAzureLogin(
		tid,
		appid,
		username,
		awsProfile,
		awsRole,
	)

	a.SetDuration(account.Duration)

	klog.V(2).Infof("Authenticate using aws-azure-login TenantID: %s", tid)
	klog.V(2).Infof("Authenticate using aws-azure-login AppIDUri: %s", appid)
	klog.V(2).Infof("Authenticate using aws-azure-login UserName: %s", username)
	klog.V(2).Infof("Authenticate using aws-azure-login AWSRole: %s", awsRole)
	klog.V(2).Infof("Authenticate using aws-azure-login AWSProfile: %s", awsProfile)
	klog.V(2).Infof("Authenticate using aws-azure-login Duration: %d", a.Duration)

	return a
}

func authWithAWSSTS(provider *config.AuthProvider, account *config.AuthAccount) authenticate.Auth {
	var sess *aws.SharedSession
	var err error

	if account.DependsOn == "" {
		klog.Fatal("DependsOn is required when using the aws-sts account provider")
	}

	auth := config.Conf.FindAuthAccount(account.DependsOn)
	fmt.Printf("%v Depends on %s\n", promptui.IconWarn, account.DependsOn)
	authAccount(auth)

	sess, err = aws.NewSharedSession(auth.AWSProfile)
	if err != nil {
		klog.Fatalf("Error when creating a new session: %s", err)
	}

	a := awssts.NewAssumeRole(
		account.AWSRole,
		provider.UserName,
		sess,
		account.AWSProfile,
	)

	d := int64(account.Duration)
	a.SetDuration(&d)
	a.SetAWSRegion(account.AWSRegion)

	return a
}
