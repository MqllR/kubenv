package awsazurelogin

import (
	"fmt"
	"os"
	"strconv"

	"github.com/mqllr/kubenv/pkg/auth"
	"github.com/mqllr/kubenv/pkg/aws"
)

// Interface is an injectable interface for running aws-azure-login commands
type Interface interface {
	// Authenticate try to authenticate using aws-azure-login
	Authenticate(auth *AWSAzureLogin) error
}

const (
	AWSAzureLoginCmd = "aws-azure-login"
	npmCmd           = "sudo npm -g"
	DefaultDuration  = 43200
)

type AWSAzureLogin struct {
	TenantID   string
	AppIDUri   string
	UserName   string
	RemeberMe  bool
	Duration   int
	AWSProfile string
	AWSRole    string
}

// NewAWSGoogleAuth create an AWSGoogleAuth struct
func NewAWSAzureLogin(tid string, addid string, username string) *AWSAzureLogin {
	return &AWSAzureLogin{
		TenantID: tid,
		AppIDUri: addid,
		UserName: username,
	}
}

// SetDefaults inject default value if not set
func (a *AWSAzureLogin) SetDefaults() {
	if a.Duration == 0 {
		a.Duration = DefaultDuration
	}
}

// Validate ensure every fields are correctly defined
func (a *AWSAzureLogin) ConfigureValidate() bool {
	if a.TenantID == "" {
		return false
	}
	if a.AppIDUri == "" {
		return false
	}
	if a.AWSProfile == "" {
		return false
	}
	if a.AWSRole == "" {
		return false
	}
	if a.UserName == "" {
		return false
	}

	return true
}

// Authenticate a username with aws-google-a
func (a *AWSAzureLogin) Configure() error {
	a.SetDefaults()

	valid := a.ConfigureValidate()
	if !valid {
		return fmt.Errorf("Error aws-azure-login profile is invalid")
	}

	ini, err := aws.NewConfigFile()
	if err != nil {
		return err
	}

	err = ini.EnsureSectionAndSave("profile "+a.AWSProfile, map[string]string{
		"azure_tenant_id":              a.TenantID,
		"azure_app_id_uri":             a.AppIDUri,
		"azure_default_username":       a.UserName,
		"azure_default_role_arn":       a.AWSRole,
		"azure_default_duration_hours": strconv.Itoa(a.Duration / 60 / 60), // Seconds to hours
		"azure_default_remember_me":    strconv.FormatBool(a.RemeberMe),
	})

	if err != nil {
		return err
	}

	return nil
}

func (runner *auth.Runner) Authenticate(a *AWSAzureLogin) error {
	args := []string{
		"--no-prompt",
		"--profile",
		a.AWSProfile,
	}

	cmd := runner.exec.Command(AWSAzureLoginCmd, args...)
	cmd.SetStdin(os.Stdin)
	cmd.SetStdout(os.Stdout)
	cmd.SetStderr(os.Stderr)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Error when running cmd: %s", err)
	}

	return nil
}
