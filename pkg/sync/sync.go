package sync

import (
	"strings"

	"github.com/mqllr/kubenv/pkg/k8s"
)

type SyncOptions struct {
	AppendTo bool
	Mode     string
	Path     string
	Command  string
}

// Sync implements a way to pick up a kubeconfig
type Sync interface {
	GetKubeConfig() (*k8s.KubeConfig, error)
}

// Service represents the required information to
// pick a kubeconfig according to the config
type Service struct {
	s Sync
}

// NewService creates a SyncService according to the
// sync type
func NewService(opts *SyncOptions) *Service {
	sync := &Service{}

	switch opts.Mode {
	case "local":
		sync.s = NewLocalFile(opts.Path)
	case "exec":
		cmd := strings.Split(opts.Command, " ")
		sync.s = NewCommandExec(cmd)
	}

	return sync
}

// GetKubeConfig retrieve a kubeconfig using a sync mode
func (s *Service) GetKubeConfig() (*k8s.KubeConfig, error) {
	return s.s.GetKubeConfig()
}
