package cmd

import (
	"bytes"
	"io"
	"regexp"
	"testing"
)

func TestVersion(t *testing.T) {
	buf := &bytes.Buffer{}
	version(io.Writer(buf))

	output := buf.String()

	match, err := regexp.MatchString(`^.* kubenv v\d+(\.\d+){2}\n`, output)
	if err != nil {
		t.Errorf("Regex failed: %s", err)
	}

	if !match {
		t.Error("Regex doesn't match")
	}
}
