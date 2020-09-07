package awsazurelogin

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mqllr/kubenv/pkg/aws"
	"k8s.io/klog"
	utilexec "k8s.io/utils/exec"
)

// Interface is an injectable interface for running aws-azure-login commands
type Interface interface {
	// Authenticate try to authenticate using aws-azure-login
	Authenticate(auth *AWSAzureLogin) error
}

const (
	AWSAzureLoginCmd = "aws-azure-login"
	npmCmd           = "sudo npm -g"
	DefaultDuration  = 12
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
func (auth *AWSAzureLogin) SetDefaults() {
	if auth.Duration == 0 {
		auth.Duration = DefaultDuration
	}
}

// Validate ensure every fields are correctly defined
func (auth *AWSAzureLogin) ConfigureValidate() bool {
	if auth.TenantID == "" {
		return false
	}
	if auth.AppIDUri == "" {
		return false
	}
	if auth.AWSProfile == "" {
		return false
	}
	if auth.AWSRole == "" {
		return false
	}
	if auth.UserName == "" {
		return false
	}

	return true
}

// Authenticate a username with aws-google-auth
func (auth *AWSAzureLogin) Configure() error {
	auth.SetDefaults()

	valid := auth.ConfigureValidate()
	if !valid {
		return fmt.Errorf("Error aws-azure-login profile is invalid")
	}

	config, err := aws.NewConfigFile()
	if err != nil {
		return err
	}

	iniSection := map[string]string{
		"azure_tenant_id":              auth.TenantID,
		"azure_app_id_uri":             auth.AppIDUri,
		"azure_default_username":       auth.UserName,
		"azure_default_role_arn":       auth.AWSRole,
		"azure_default_duration_hours": strconv.Itoa(auth.Duration),
		"azure_default_remember_me":    strconv.FormatBool(auth.RemeberMe),
	}

	err = config.EnsureIniSection("profile "+auth.AWSProfile, iniSection)
	if err != nil {
		return err
	}

	return nil
}

type runner struct {
	exec utilexec.Interface
}

// New returns a new Interface which will exec aws-google-auth
func New(exec utilexec.Interface) Interface {
	return &runner{
		exec: exec,
	}
}

func (runner *runner) Authenticate(auth *AWSAzureLogin) error {
	args := []string{
		"--no-prompt",
		"--profile",
		auth.AWSProfile,
	}

	klog.V(2).Infof("Running cmd: %s %s", AWSAzureLoginCmd, strings.Join(args, " "))

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
