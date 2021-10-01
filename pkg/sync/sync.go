package sync

import (
	"os"
	"strings"

	"github.com/mqllr/kubenv/pkg/k8s"
)

type SyncOptions struct {
	AppendTo bool
	Mode     string
	Path     string
	Command  string
	Glob     string
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
	// TODO handle errors
	sync := &Service{}

	switch opts.Mode {
	case "local":
		f, err := os.Open(opts.Path)
		if err != nil {
			return nil
		}
		sync.s = NewLocalFile(f)
	case "exec":
		cmd := strings.Split(opts.Command, " ")
		sync.s = NewCommandExec(cmd)
	case "glob":
		fs := os.DirFS("/")
		sync.s = NewGlob(fs, opts.Glob)
	default:
		return nil
	}

	return sync
}

// GetKubeConfig retrieve a kubeconfig using a sync mode
func (s *Service) GetKubeConfig() (*k8s.KubeConfig, error) {
	return s.s.GetKubeConfig()
}
