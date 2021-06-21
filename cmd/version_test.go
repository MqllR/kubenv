package cmd

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"testing"
)

const releaseVersionRe = `\d+(\.\d+){2}`

func TestVersionOutput(t *testing.T) {
	os.Setenv("KUBENV_CONFIG", os.Getenv("PWD")+"/../example/kubenv_example.yaml")

	options := []struct {
		arg   string
		regex string
	}{
		{"", fmt.Sprintf(`^.* kubenv v%s\n`, releaseVersionRe)},
		{"json", fmt.Sprintf(`\{"version":"%s"\}`, releaseVersionRe)},
		{"foo", `Unknown output`},
	}

	cmd := NewVersionCmd()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)

	for _, option := range options {
		cmd.SetArgs([]string{"-o", option.arg})
		b.Reset()

		err := cmd.Execute()
		if err != nil {
			t.Errorf("Cmd executing failed with arg %s, err: %s", option.arg, err)
		}

		output := b.String()
		match, err := regexp.MatchString(option.regex, output)
		if err != nil {
			t.Errorf("Regex failed for arg %s: %s", option.arg, err)
		}

		if !match {
			t.Errorf("Regex doesn't match for arg %s, output: %s", option.arg, option)
		}
	}
}

func TestVersionPublished(t *testing.T) {
	tag := os.Getenv("GITHUB_REF")
	switch ok, _ := regexp.MatchString(releaseVersionRe, tag); {
	case !ok:
		t.Log("Not within a Github workflow")
	case ok:
		if tag != version {
			t.Errorf("Version mismatch, expected %s but got %s", version, tag)
		}
	}
}
