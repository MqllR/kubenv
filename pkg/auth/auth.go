package auth

import (
	"github.com/mqllr/kubenv/pkg/config"
	executil "k8s.io/utils/exec"

	awsgoogleauth "github.com/mqllr/kubenv/pkg/aws-google-auth"
)

// Interface is an injectable interface for running aws-azure-login commands
type Interface interface {
	Authenticate() error
	GetVersion() (string, error)
	GetRemoteVersion() (string, error)
	Install() error
}

type runner struct {
	exec executil.Interface
}

// New returns a new Interface which will exec aws-google-auth
func NewRunner(exec executil.Interface) Interface {
	return &runner{
		exec: exec,
	}
}

type Auth struct {
	Account  *config.AuthAccount
	Provider *config.AuthProvider
}

func NewAuth(account *config.AuthAccount, provider *config.AuthProvider) *Auth {
	return &Auth{
		Account:  account,
		Provider: provider,
	}
}

func (auth *Auth) Authenticate() error {
	switch auth.Account.AuthProvider {
	case "aws-google-auth":
		err := auth.authWithGoogleAuth()
		return err
	case "aws-azure-login":
		err := auth.authWithAzureLogin()
		return err
	case "aws-sts":
		auth.authWithAWSSTS()
		err := auth.authWithAzureLogin()
		return err
	}
	return nil
}

func (auth *Auth) authWithGoogleAuth() error {
	a := awsgoogleauth.NewAWSGoogleAuth(
		auth.Provider.IDP,
		auth.Provider.SP,
		auth.Provider.UserName,
	)

	a.AWSRole = auth.Account.AWSRole
	a.AWSProfile = auth.Account.AWSProfile
	a.Region = auth.Account.Region

	execer := executil.New()
	runner := awsgoogleauth.New(execer)

	err := runner.Authenticate(a)
	if err != nil {
		return err
	}

	return nil
}

func (auth *Auth) authWithAzureLogin() error {
	a := awsazurelogin.NewAWSAzureLogin(
		auth.Provider.TenantID,
		auth.Provider.AppIDUri,
		auth.Provider.UserName,
	)

	a.AWSRole = auth.Account.AWSRole
	a.AWSProfile = auth.Account.AWSProfile
	a.Duration = auth.Account.Duration

	err := a.Configure()
	if err != nil {
		return err
	}

	execer := executil.New()
	runner := awsazurelogin.New(execer)

	err = runner.Authenticate(a)
	if err != nil {
		return err
	}

	return nil
}

func (auth *Auth) authWithAWSSTS() error {
	/*
		var sess *aws.SharedSession
		var err error
		if account.DependsOn != "" {
			auth := authAccountsConfig.FindAuthAccount(account.DependsOn)

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
		var sess *aws.SharedSession
		var err error
		if account.DependsOn != "" {
			auth := authAccountsConfig.FindAuthAccount(account.DependsOn)

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
	*/
	return nil
}
