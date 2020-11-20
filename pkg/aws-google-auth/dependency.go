package awsgoogleauth

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

var VersionPattern = `(\d+\.){1,2}\d+`

type AWSGoogleAuthExec struct{}

// GetVersion return the version number of the aws-google-auth tools
func (e *AWSGoogleAuthExec) GetLocalVersion() (string, error) {
	args := []string{"-V"}

	bytes, err := exec.Command(AWSGoogleAuthCmd, args...).CombinedOutput()

	if err != nil {
		return "", fmt.Errorf("Error when getting local version, error: %v", err)
	}

	versionMatcher := regexp.MustCompile(VersionPattern)
	match := versionMatcher.FindStringSubmatch(string(bytes))

	if match == nil {
		return "", fmt.Errorf("No aws-google-auth version found in string: %s", bytes)
	}

	return match[0], nil
}

// GetRemoteVersion return the version number of the aws-google-auth from the pypi repository
func (e *AWSGoogleAuthExec) GetRemoteVersion() (string, error) {
	searchPattern := `^aws\-google\-auth \((?P<Version>` + VersionPattern + `)\)`

	bytes, err := exec.Command(pipCmd, []string{"search", AWSGoogleAuthCmd}...).CombinedOutput()

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
func (e *AWSGoogleAuthExec) Install() error {
	cmd := exec.Command(pipCmd, []string{"install", AWSGoogleAuthCmd + "[u2f]"}...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Error when installing %s: %s", AWSGoogleAuthCmd, err)
	}

	return nil
}
