package cmd

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"k8s.io/klog"
	executil "k8s.io/utils/exec"

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

	switch account.AuthProvider {
	case "aws-google-auth":
		authWithGoogleAuth(config.Conf.FindAuthProvider("aws-google-auth"), account)
	case "aws-azure-login":
		authWithAzureLogin(config.Conf.FindAuthProvider("aws-azure-login"), account)
	case "aws-sts":
		authWithAWSSTS(config.Conf.FindAuthProvider("aws-sts"), account)
	}
}

func authWithGoogleAuth(provider *config.AuthProvider, account *config.AuthAccount) {
	idp := provider.IDP
	sp := provider.SP
	username := provider.UserName

	klog.V(2).Infof("Authenticate using aws-google-auth IDP: %s", idp)
	klog.V(2).Infof("Authenticate using aws-google-auth SP: %s", sp)
	klog.V(2).Infof("Authenticate using aws-google-auth UserName: %s", username)

	auth := awsgoogleauth.NewAWSGoogleAuth(
		idp,
		sp,
		username,
	)

	auth.AWSRole = account.AWSRole
	auth.AWSProfile = account.AWSProfile
	auth.Region = account.Region

	klog.V(2).Infof("Authenticate using aws-google-auth AWSRole: %s", auth.AWSRole)
	klog.V(2).Infof("Authenticate using aws-google-auth AWSProfile: %s", auth.AWSProfile)
	klog.V(2).Infof("Authenticate using aws-google-auth Region: %s", auth.Region)
	klog.V(2).Infof("Authenticate using aws-azure-login Duration: %d", auth.Duration)

	execer := executil.New()
	runner := awsgoogleauth.New(execer)

	err := runner.Authenticate(auth)
	if err != nil {
		klog.Fatalf("Error on authentication: %s", err)
	}
}

func authWithAzureLogin(provider *config.AuthProvider, account *config.AuthAccount) {
	tid := provider.TenantID
	appid := provider.AppIDUri
	username := provider.UserName

	klog.V(2).Infof("Authenticate using aws-azure-login TenantID: %s", tid)
	klog.V(2).Infof("Authenticate using aws-azure-login AppIDUri: %s", appid)
	klog.V(2).Infof("Authenticate using aws-azure-login UserName: %s", username)

	auth := awsazurelogin.NewAWSAzureLogin(
		tid,
		appid,
		username,
	)

	auth.AWSRole = account.AWSRole
	auth.AWSProfile = account.AWSProfile
	auth.Duration = account.Duration

	klog.V(2).Infof("Authenticate using aws-azure-login AWSRole: %s", auth.AWSRole)
	klog.V(2).Infof("Authenticate using aws-azure-login AWSProfile: %s", auth.AWSProfile)
	klog.V(2).Infof("Authenticate using aws-azure-login Duration: %d", auth.Duration)

	execer := executil.New()
	runner := awsazurelogin.New(execer)

	err := auth.Configure()
	if err != nil {
		klog.Fatalf("Error when configuring the AWS profile: %s", err)
	}

	klog.V(2).Info("Profile configured for aws-azure-login tool")

	err = runner.Authenticate(auth)
	if err != nil {
		klog.Fatalf("Error on authentication: %s", err)
	}
}

func authWithAWSSTS(provider *config.AuthProvider, account *config.AuthAccount) {
	var sess *aws.SharedSession
	var err error
	if account.DependsOn != "" {
		auth := config.Conf.FindAuthAccount(account.DependsOn)

		fmt.Printf("%v Depends on %s\n", promptui.IconWarn, account.DependsOn)

		authAccount(auth)
		sess, err = aws.NewSharedSession(auth.AWSProfile)
		if err != nil {
			klog.Fatalf("Error when creating a new session: %s", err)
		}
	} else {
		sess, err = aws.NewSharedSession("")
		if err != nil {
			klog.Fatalf("Error when creating a new session: %s", err)
		}
	}

	a := awssts.NewAssumeRole(
		account.AWSRole,
		provider.UserName,
		sess,
		account.AWSProfile,
		account.Region,
	)

	d := int64(account.Duration)
	a.Duration = &d

	err = a.Authenticate()
	if err != nil {
		klog.Fatalf("Error when trying the get a STS session: %s", err)
	}

	fmt.Printf("%v Authenticated on %s\n", promptui.IconGood, account.AuthProvider)
}
