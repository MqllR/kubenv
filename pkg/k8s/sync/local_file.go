package sync

import "github.com/mqllr/kubenv/pkg/k8s"

type LocalFile struct {
	path string
}

func NewLocalFile(path string) *LocalFile {
	return &LocalFile{
		path: path,
	}
}

func (local *LocalFile) GetKubeConfig() (*k8s.KubeConfig, error) {
	return k8s.NewKubeConfigFromFile(local.path)
}
