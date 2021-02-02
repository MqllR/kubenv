package sync

import (
	"fmt"

	"github.com/mqllr/kubenv/pkg/config"
	"github.com/mqllr/kubenv/pkg/k8s"
)

type Sync interface {
	GetKubeConfig() (*k8s.KubeConfig, error)
}

type SyncService struct {
	s      Sync
	config config.K8SSync
}

func NewSyncService(conf config.K8SSync) (*SyncService, error) {
	sync := &SyncService{
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

// AppendKubeConfig merge 2 kubeConfig file in 1
func (s *SyncService) AppendKubeConfig(kubeConfig *k8s.KubeConfig) error {
	k, err := s.s.GetKubeConfig()
	if err != nil {
		return fmt.Errorf("Cannot get the kubeconfig: %s", err)
	}

	kubeConfig.Append(k)

	return nil
}
