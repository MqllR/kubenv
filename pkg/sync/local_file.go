package sync

import (
	"io"

	"github.com/mqllr/kubenv/pkg/k8s"
)

// LocalFile is the path of a kubeconfig file
type LocalFile struct {
	reader io.Reader
}

var _ Syncer = &LocalFile{}

// NewLocalFile creates a LocalFile
func NewLocalFile(reader io.Reader) *LocalFile {
	return &LocalFile{
		reader,
	}
}

// GetKubeConfig just returns a KubeConfig
func (local *LocalFile) GetKubeConfig() (*k8s.KubeConfig, error) {
	return k8s.NewKubeConfigFromReader(local.reader)
}
