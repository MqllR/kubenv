package awsgoogleauth

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"k8s.io/klog"
)

type AWSGoogleAuth struct {
	UserName   string
	IDP        string
	SP         string
	Duration   int
	AWSRegion  string
	AWSProfile string
	AWSRole    string
}

// NewAWSGoogleAuth create an AWSGoogleAuth struct
func NewAWSGoogleAuth(idp string, sp string, userName string, awsProfile string, awsRole string) *AWSGoogleAuth {

	return &AWSGoogleAuth{
		IDP:        idp,
		SP:         sp,
		UserName:   userName,
		AWSProfile: awsProfile,
		AWSRole:    awsRole,
	}
}

func (a *AWSGoogleAuth) SetDuration(duration int) {
	a.Duration = duration
}

func (a *AWSGoogleAuth) SetAWSRegion(awsRegion string) {
	a.AWSRegion = awsRegion
}

// SetDefaults inject default value if not set
func (a *AWSGoogleAuth) SetDefaults() {
	if a.Duration == 0 {
		a.Duration = DefaultDuration
	}

	if a.AWSRegion == "" {
		a.AWSRegion = DefaultAWSRegion
	}
}

// Validate ensure every fields are correctly defined
func (a *AWSGoogleAuth) Validate() bool {
	if a.IDP == "" {
		return false
	}
	if a.SP == "" {
		return false
	}
	if a.AWSRegion == "" {
		return false
	}
	if a.AWSProfile == "" {
		return false
	}
	if a.AWSRole == "" {
		return false
	}

	return true
}

// Authenticate a username with aws-google-auth
func (a *AWSGoogleAuth) Authenticate() error {
	args := []string{
		"-k",
		"-I",
		a.IDP,
		"-S",
		a.SP,
		"-d",
		strconv.Itoa(a.Duration),
		"-p",
		a.AWSProfile,
		"-r",
		a.AWSRole,
		"-R",
		a.AWSRegion,
		"-u",
		a.UserName,
	}

	if bool(klog.V(5)) {
		args = append(args, []string{"-l", "debug"}...)
	}

	klog.V(2).Infof("Running cmd: %s %s", AWSGoogleAuthCmd, strings.Join(args, " "))

	cmd := exec.Command(AWSGoogleAuthCmd, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Error when running cmd: %s", err)
	}

	return nil
}
