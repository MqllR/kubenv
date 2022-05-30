package cmd

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"
)

const releaseVersionRe = `\d+(\.\d+){2}`

func TestVersionOutput(t *testing.T) {
	options := map[string]struct {
		arg   string
		regex string
	}{
		"empty":         {"", fmt.Sprintf(`^.* kubenv v%s\n`, releaseVersionRe)},
		"json":          {"json", fmt.Sprintf(`\{"version":"%s"\}`, releaseVersionRe)},
		"non supported": {"foo", `Unknown output`},
	}

	for name, option := range options {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			b := bytes.NewBufferString("")

			cmd := versionCmd()
			cmd.SetOut(b)
			cmd.SetArgs([]string{"-o", option.arg})

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
		})
	}
}

func TestVersionPublished(t *testing.T) {
	// refs/tags/0.3.0
	tag := os.Getenv("GITHUB_REF")
	switch ok, _ := regexp.MatchString(releaseVersionRe, tag); {
	case !ok:
		t.Log("Not within a Github workflow")
	case ok:
		a := strings.Split(tag, "/")
		if a[len(a)-1] != version {
			t.Errorf("Version mismatch, expected %s but got %s", version, tag)
		}
	}
}
