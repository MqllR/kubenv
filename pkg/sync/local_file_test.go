package sync_test

import (
	"os"
	"testing"

	"github.com/mqllr/kubenv/pkg/sync"
)

func TestLocalFile(t *testing.T) {
	kubeconfig1, err := os.Open("testdata/kubeconfig1.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer kubeconfig1.Close()

	localFile := sync.NewLocalFile(kubeconfig1)
	kubeconfig, err := localFile.GetKubeConfig()

	if err != nil {
		t.Errorf("Unexpeted err: %s", err)
	}

	if kubeconfig.GetContextNames()[0] != "fakecontext1" {
		t.Errorf("Got context %s but was expecting to get %s", kubeconfig.GetContextNames()[0], "fakecontext1")
	}
}
