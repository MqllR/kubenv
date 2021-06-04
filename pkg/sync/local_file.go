package sync

import "github.com/mqllr/kubenv/pkg/k8s"

// LocalFile is the path of a kubeconfig file
type LocalFile struct {
	path string
}

// NewLocalFile creates a LocalFile
func NewLocalFile(path string) *LocalFile {
	return &LocalFile{
		path: path,
	}
}

// GetKubeConfig just returns a KubeConfig
func (local *LocalFile) GetKubeConfig() (*k8s.KubeConfig, error) {
	return k8s.NewKubeConfigFromFile(local.path)
}
