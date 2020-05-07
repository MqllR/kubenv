package awsiamauthenticator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	utilexec "k8s.io/utils/exec"
)

// Interface is an injectable interface for running aws-google-auth commands
type Interface interface {
	// GetVersion returns the version string for aws-iam-authenticator
	GetVersion() (string, error)
	GetRemoteVersion() (string, error)
}

const (
	AWSIAMAuthCmd     = "aws-iam-authenticator"
	GithubAPIEndpoint = "https://api.github.com/repos/kubernetes-sigs/aws-iam-authenticator/releases/latest"
)

type runner struct {
	exec utilexec.Interface
}

// New returns a new Interface which will exec aws-iam-authenticator
func New(exec utilexec.Interface) Interface {
	return &runner{
		exec: exec,
	}
}

// GetVersion return the version number of the aws-iam-authenticator tools
func (runner *runner) GetVersion() (string, error) {
	args := []string{"version"}

	bytes, err := runner.exec.Command(AWSIAMAuthCmd, args...).CombinedOutput()

	if err != nil {
		return "", fmt.Errorf("Error when getting version, error: %v", err)
	}

	type version struct {
		Version string
		Commit  string
	}

	v := &version{}
	if err := json.Unmarshal(bytes, v); err != nil {
		return "", fmt.Errorf("Cannot unmarshal json after command %s: %s", AWSIAMAuthCmd+" "+strings.Join(args, " "), err)
	}

	return v.Version, nil
}

// GetRemoteVersion return the version number of the aws-iam-authenticator from github
func (runner *runner) GetRemoteVersion() (string, error) {
	resp, err := http.Get(GithubAPIEndpoint)
	if err != nil {
		return "", fmt.Errorf("Error on GET request: %s", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Cannot ready body: %s", err)
	}

	type githubRelease struct {
		TagName string                 `json:"tag_name"`
		X       map[string]interface{} `json:"-"`
	}

	release := &githubRelease{}

	err = json.Unmarshal(body, release)
	if err != nil {
		return "", fmt.Errorf("Cannot unmarshal body: %s", err)
	}

	if release.TagName == "" {
		return "", fmt.Errorf("Empty tag_name")
	}

	return release.TagName, nil
}
