package sync_test

import (
	"strings"
	"testing"

	"github.com/mqllr/kubenv/pkg/sync"
)

func TestLocalFile(t *testing.T) {
	r := strings.NewReader(`
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: FAKEVALUE
    server: https://fakeurl.com
  name: fakecluster1
contexts:
- context:
    cluster: fakecluster1
    namespace: fakens1
    user: fakeuser1
  name: fakecontext1
kind: Config
preferences: {}
users:
- name: fakeuser1
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1alpha1
      args:
      - token
      - -i
      - fakecluster
      command: aws-iam-authenticator
`)

	localFile := sync.NewLocalFile(r)
	kubeconfig, err := localFile.GetKubeConfig()

	if err != nil {
		t.Errorf("Unexpeted err: %s", err)
	}

	if kubeconfig.GetContextNames()[0] != "fakecontext1" {
		t.Errorf("Got context %s but was expecting to get %s", kubeconfig.GetContextNames()[0], "fakecontext1")
	}
}
