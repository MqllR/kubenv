package cmd

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/manifoldco/promptui"
)

func TestVersionOutput(t *testing.T) {
	options := map[string]struct {
		arg   string
		regex string
	}{
		"empty":         {"", fmt.Sprintf("%s kubenv vdev\n", promptui.IconGood)},
		"json":          {"json", `{"version":"dev"}`},
		"non supported": {"foo", `Unknown output`},
	}

	for name, option := range options {
		t.Run(name, func(t *testing.T) {
			b := bytes.NewBufferString("")

			cmd := versionCmd()
			cmd.SetOut(b)
			cmd.SetArgs([]string{"-o", option.arg})

			err := cmd.Execute()
			if err != nil {
				t.Errorf("Cmd executing failed with arg %s, err: %s", option.arg, err)
			}

			output := b.String()

			if option.regex != output {
				t.Errorf("Expected %s, got: %s", option.regex, output)
			}
		})
	}
}
