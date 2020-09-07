package awsgoogleauth

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"k8s.io/klog"
	utilexec "k8s.io/utils/exec"
)

// Interface is an injectable interface for running aws-google-auth commands
type Interface interface {
	// Authenticate try to authenticate using aws-google-auth
	Authenticate(auth *AWSGoogleAuth) error
	// GetVersion returns the "X.Y" version string for aws-google-auth.
	GetVersion() (string, error)
	// GetRemoteVersion returns the "X.Y" version string for aws-google-auth from pypi
	GetRemoteVersion() (string, error)
	// Install will install the latest version aws-google-auth
	Install() error
}

const (
	AWSGoogleAuthCmd = "aws-google-auth"
	pipCmd           = "sudo pip"
	DefaultDuration  = 28800
)

var VersionPattern = `(\d+\.){1,2}\d+`

type AWSGoogleAuth struct {
	UserName   string
	IDP        string
	SP         string
	Region     string
	Duration   int
	AWSProfile string
	AWSRole    string
}

// NewAWSGoogleAuth create an AWSGoogleAuth struct
func NewAWSGoogleAuth(idp string, sp string) *AWSGoogleAuth {
	return &AWSGoogleAuth{
		IDP: idp,
		SP:  sp,
	}
}

// SetDefaults inject default value if not set
func (auth *AWSGoogleAuth) SetDefaults() {
	if auth.Duration == 0 {
		auth.Duration = DefaultDuration
	}
}

// Validate ensure every fields are correctly defined
func (auth *AWSGoogleAuth) Validate() bool {
	if auth.IDP == "" {
		return false
	}
	if auth.SP == "" {
		return false
	}
	if auth.Region == "" {
		return false
	}
	if auth.AWSProfile == "" {
		return false
	}
	if auth.AWSRole == "" {
		return false
	}

	return true
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

// Authenticate a username with aws-google-auth
func (runner *runner) Authenticate(auth *AWSGoogleAuth) error {
	auth.SetDefaults()

	valid := auth.Validate()
	if !valid {
		return fmt.Errorf("Error AWSGoogleAuth is invalid")
	}

	args := []string{
		"-k",
		"-I",
		auth.IDP,
		"-S",
		auth.SP,
		"-d",
		strconv.Itoa(auth.Duration),
		"-p",
		auth.AWSProfile,
		"-r",
		auth.AWSRole,
		"-R",
		auth.Region,
	}

	if auth.UserName != "" {
		args = append(args, []string{"-u", auth.UserName}...)
	}

	if klog.V(5) {
		args = append(args, []string{"-l", "debug"}...)
	}

	klog.V(2).Infof("Running cmd: %s %s", AWSGoogleAuthCmd, strings.Join(args, " "))

	cmd := runner.exec.Command(AWSGoogleAuthCmd, args...)
	cmd.SetStdin(os.Stdin)
	cmd.SetStdout(os.Stdout)
	cmd.SetStderr(os.Stderr)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Error when running cmd: %s", err)
	}

	return nil
}

// GetVersion return the version number of the aws-google-auth tools
func (runner *runner) GetVersion() (string, error) {
	args := []string{"-V"}

	bytes, err := runner.exec.Command(AWSGoogleAuthCmd, args...).CombinedOutput()

	if err != nil {
		return "", fmt.Errorf("Error when getting version, error: %v", err)
	}

	versionMatcher := regexp.MustCompile(VersionPattern)
	match := versionMatcher.FindStringSubmatch(string(bytes))

	if match == nil {
		return "", fmt.Errorf("No aws-google-auth version found in string: %s", bytes)
	}

	return match[0], nil
}

// GetRemoteVersion return the version number of the aws-google-auth from the pypi repository
func (runner *runner) GetRemoteVersion() (string, error) {
	searchPattern := `^aws\-google\-auth \((?P<Version>` + VersionPattern + `)\)`

	bytes, err := runner.exec.Command(pipCmd, []string{"search", AWSGoogleAuthCmd}...).CombinedOutput()

	if err != nil {
		return "", fmt.Errorf("Error when executing pip search, error: %v", err)
	}

	versionMatcher := regexp.MustCompile(searchPattern)
	match := versionMatcher.FindStringSubmatch(string(bytes))

	if match == nil {
		return "", fmt.Errorf("No aws-google-auth version found in pip search: %s", bytes)
	}

	return match[1], nil
}

// Install will install the latest version aws-google-auth
func (runner *runner) Install() error {
	cmd := runner.exec.Command(pipCmd, []string{"install", AWSGoogleAuthCmd + "[u2f]"}...)

	cmd.SetStdin(os.Stdin)
	cmd.SetStdout(os.Stdout)
	cmd.SetStderr(os.Stderr)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Error when installing %s: %s", AWSGoogleAuthCmd, err)
	}

	return nil
}
