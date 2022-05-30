package config_test

import (
	"os"
	"testing"

	"github.com/mqllr/kubenv/pkg/config"
)

func TestGetKubeConfig(t *testing.T) {
	err := os.Setenv("HOME", "/home/foo")
	if err != nil {
		t.Fatal("Cannot set the home")
	}

	testCases := map[string]struct {
		envValue string
		expected string
	}{
		"with env":    {"/foo/bar", "/foo/bar"},
		"without env": {"", "/home/foo"},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := os.Setenv("KUBECONFIG", test.envValue)
			if err != nil {
				t.Errorf("Cannot set env variablie KUBECONFIG: %s", err)
			}

			got := config.GetKubeConfig()
			if got != test.expected {
				t.Errorf("Got %s but expected %s", got, test.expected)
			}
		})
	}
}
