package sync_test

import (
	"strings"
	"testing"

	"github.com/mqllr/kubenv/pkg/sync"
)

func TestLocalFile(t *testing.T) {
	r := strings.NewReader(kubeconfig1)

	localFile := sync.NewLocalFile(r)
	kubeconfig, err := localFile.GetKubeConfig()

	if err != nil {
		t.Errorf("Unexpeted err: %s", err)
	}

	if kubeconfig.GetContextNames()[0] != "fakecontext1" {
		t.Errorf("Got context %s but was expecting to get %s", kubeconfig.GetContextNames()[0], "fakecontext1")
	}
}
