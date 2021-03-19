package publish

/*
* Expect a tag ref from github
**/

import (
	"encoding/json"
	"os"
	"os/exec"
	"testing"
)

func TestPublished(t *testing.T) {
	cmd := exec.Command(
		"../../kubenv-linux-amd64",
		[]string{"version", "-o", "json"}...,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Error when running version command: %s", err)
	}

	version := &struct {
		V string `json:"version"`
	}{}

	err = json.Unmarshal(output, &version)
	if err != nil {
		t.Fatalf("Cannot unmarshal json output: %s", err)
	}

	tag := os.Getenv("TAG_REF")
	if tag == "" {
		t.Fatal("Missing TAG_REF env variable")
	}

	if version.V != tag {
		t.Errorf("Version mismatch, expected %s but got %s", version.V, tag)
	}
}
