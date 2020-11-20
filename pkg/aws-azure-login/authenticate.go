package awsazurelogin

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/mqllr/kubenv/pkg/aws"
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
func NewAWSAzureLogin(tenantId string, appIDUri string, userName string, awsProfile string, awsRole string) *AWSAzureLogin {

	return &AWSAzureLogin{
		TenantID:   tenantId,
		AppIDUri:   appIDUri,
		UserName:   userName,
		AWSProfile: awsProfile,
		AWSRole:    awsRole,
	}
}

func (a *AWSAzureLogin) SetDuration(duration int) {
	a.Duration = duration
}

func (a *AWSAzureLogin) SetRemerberMe(remeberMe bool) {
	a.RemeberMe = remeberMe
}

// SetDefaults inject default value if not set
func (a *AWSAzureLogin) SetDefaults() {
	if a.Duration == 0 {
		a.Duration = DefaultDuration
	}
}

// Validate ensure every fields are correctly defined
func (a *AWSAzureLogin) Validate() bool {
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

func (a *AWSAzureLogin) configure() error {
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

func (a *AWSAzureLogin) Authenticate() error {
	err := a.configure()
	if err != nil {
		return fmt.Errorf("Failed to configure the profile: %s", err)
	}

	args := []string{
		"--no-prompt",
		"--profile",
		a.AWSProfile,
	}

	cmd := exec.Command(AWSAzureLoginCmd, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("Error when running cmd: %s", err)
	}

	return nil
}
