package config

import (
	"os"
	"testing"
)

func TestGetKubeConfig(t *testing.T) {
	tests := map[string]struct {
		setEnv bool
		expect string
	}{
		"setenv": {setEnv: true, expect: "/tmp/kubeconfig"},
		"noenv":  {setEnv: false, expect: "/home/toto/.kube/config"},
	}

	err := os.Setenv("HOME", "/home/toto")
	if err != nil {
		t.Error("Cannot set the home env variable")
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if test.setEnv {
				t.Setenv("KUBECONFIG", "/tmp/kubeconfig")
			}

			kubeconfig := GetKubeConfig()
			if kubeconfig != test.expect {
				t.Errorf("kubeconfig path doesn't match, got %s but expected %s", kubeconfig, test.expect)
			}
		})
	}
}
