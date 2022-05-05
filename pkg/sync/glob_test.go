package sync_test

import (
	"testing"
	"testing/fstest"

	"github.com/mqllr/kubenv/pkg/sync"
)

func TestGlob(t *testing.T) {
	fs := fstest.MapFS{
		"folder/foo": {Data: []byte(kubeconfig1)},
		"folder/bar": {Data: []byte(kubeconfig2)},
	}

	glob := sync.NewGlob(fs, "folder/*")
	kubeconfig, err := glob.GetKubeConfig()

	if err != nil {
		t.Errorf("Unexpeted err: %s", err)
	}

	contexts := kubeconfig.GetContextNames()
	if len(contexts) != 2 {
		t.Error("Number of kubeconfig return is not 2")
	}

	if contexts[1] != "fakecontext1" {
		t.Errorf("Got context %s but was expecting to get %s", contexts[0], "fakecontext1")
	}

	if contexts[0] != "fakecontext2" {
		t.Errorf("Got context %s but was expecting to get %s", contexts[0], "fakecontext2")
	}
}
