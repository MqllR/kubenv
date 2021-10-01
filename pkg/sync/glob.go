package sync

import (
	"fmt"
	"io/fs"

	"github.com/mqllr/kubenv/pkg/k8s"
)

type Glob struct {
	fileSystem fs.FS
	pattern    string
}

func NewGlob(fileSystem fs.FS, pattern string) *Glob {
	return &Glob{
		fileSystem,
		pattern,
	}
}

func (g *Glob) GetKubeConfig() (*k8s.KubeConfig, error) {
	matches, err := fs.Glob(g.fileSystem, g.pattern)
	if err != nil {
		return nil, fmt.Errorf("Cannot get the files matching the pattern: %s", err)
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("No matches found for the pattern %s", g.pattern)
	}

	kubeconfig := k8s.NewKubeConfig()
	for _, file := range matches {
		f, err := g.fileSystem.Open(file)
		if err != nil {
			return nil, fmt.Errorf("Cannot open the file: %s", err)
		}
		k, err := k8s.NewKubeConfigFromReader(f)
		if err != nil {
			return nil, fmt.Errorf("Cannot get the kubeconfig: %s", err)
		}

		kubeconfig.Append(k)
	}

	return kubeconfig, nil
}
