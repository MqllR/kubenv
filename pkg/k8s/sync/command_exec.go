package sync

import "github.com/mqllr/kubenv/pkg/k8s"

type CommandExec struct{}

func NewCommandExec() *CommandExec {
	return &CommandExec{}
}

func (cmd *CommandExec) GetKubeConfig() (*k8s.KubeConfig, error) {
	return nil, nil
}
