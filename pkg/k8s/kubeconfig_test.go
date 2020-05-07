package k8s

import (
	"testing"
)

var testingConfig = `
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: FAKEVALUE
    server: https://fakeurl.com
  name: fakecluster
contexts:
- context:
    cluster: fakecluster
    namespace: fakens
    user: fakeuser
  name: fakecontext
kind: Config
preferences: {}
users:
- name: fakeuser
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1alpha1
      args:
      - token
      - -i
      - fakecluster
      command: aws-iam-authenticator
`

func TestNewKubeConfig(t *testing.T) {
	kubeconfig := NewKubeConfig()

	if kubeconfig.APIVersion != "v1" {
		t.Errorf("Excepted value for ApiVersion %s but got %s", "v1", kubeconfig.APIVersion)
	}

	if kubeconfig.Kind != "Config" {
		t.Errorf("Excepted value for Kind %s but got %s", "Config", kubeconfig.Kind)
	}
}

func TestUnmarshal(t *testing.T) {
	kubeconfig := NewKubeConfig()
	kubeconfig.Unmarshal([]byte(testingConfig))

	if len(kubeconfig.Clusters) != 1 {
		t.Errorf("Excepted array length of clusters %d but got %d", 1, len(kubeconfig.Clusters))
	}

	if len(kubeconfig.Users) != 1 {
		t.Errorf("Excepted array length of users %d but got %d", 1, len(kubeconfig.Users))
	}

	if len(kubeconfig.Contexts) != 1 {
		t.Errorf("Excepted array length of contexts %d but got %d", 1, len(kubeconfig.Contexts))
	}

	if kubeconfig.Contexts[0].Name != "fakecontext" {
		t.Errorf("Excepted context name %s but got %s", "fakecontext", kubeconfig.Contexts[0].Name)
	}

	if kubeconfig.Users[0].Name != "fakeuser" {
		t.Errorf("Excepted user name %s but got %s", "fakeuser", kubeconfig.Users[0].Name)
	}

	if kubeconfig.Clusters[0].Name != "fakecluster" {
		t.Errorf("Excepted cluster name %s but got %s", "fakecluster", kubeconfig.Clusters[0].Name)
	}
}
