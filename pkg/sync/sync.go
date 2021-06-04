package sync

import (
	"fmt"

	"github.com/mqllr/kubenv/pkg/config"
	"github.com/mqllr/kubenv/pkg/k8s"
)

// Sync implements a way to pick up a kubeconfig
type Sync interface {
	GetKubeConfig() (*k8s.KubeConfig, error)
}

// Service represents the required information to
// pick a kubeconfig according to the config
type Service struct {
	s      Sync
	config config.K8SSync
}

// NewService creates a SyncService according to the
// sync type
func NewService(conf config.K8SSync) (*Service, error) {
	sync := &Service{
		config: conf,
	}

	switch conf.Mode {
	case "local":
		sync.s = NewLocalFile(conf.Path)
	case "exec":
		sync.s = NewCommandExec(conf.Command)
	default:
		return nil, fmt.Errorf("Sync mode not implemented")
	}

	return sync, nil
}

// AppendKubeConfig merges the kubeconfig synchronised into the
// kubeConfig in argument
func (s *Service) AppendKubeConfig(kubeConfig *k8s.KubeConfig) error {
	k, err := s.s.GetKubeConfig()
	if err != nil {
		return fmt.Errorf("Cannot get the kubeconfig: %s", err)
	}

	kubeConfig.Append(k)

	return nil
}
