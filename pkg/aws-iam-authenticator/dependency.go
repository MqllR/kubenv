package awsiamauthenticator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
)

type AWSIAMAuthExec struct{}

// GetVersion return the version number of the aws-iam-authenticator tools
func (e *AWSIAMAuthExec) GetLocalVersion() (string, error) {
	args := []string{"version"}

	bytes, err := exec.Command(AWSIAMAuthCmd, args...).CombinedOutput()

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
func (e *AWSIAMAuthExec) GetRemoteVersion() (string, error) {
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

func (e *AWSIAMAuthExec) Install() error {
	// TODO
	return nil
}
