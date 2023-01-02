package sync

import (
	"github.com/mqllr/kubenv/pkg/k8s"
)

// Syncer implements a way to pick up a kubeconfig
type Syncer interface {
	GetKubeConfig() (*k8s.KubeConfig, error)
}
