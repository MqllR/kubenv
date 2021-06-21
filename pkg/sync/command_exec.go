package sync

import (
	"fmt"
	"os/exec"

	"github.com/mqllr/kubenv/pkg/k8s"
)

type CommandExec struct {
	cmd []string
}

// NewCommandExec just creates a CommandExec
func NewCommandExec(cmd []string) *CommandExec {
	return &CommandExec{
		cmd: cmd,
	}
}

// GetKubeConfig executes the command
// returns a KubeConfig
func (cmd *CommandExec) GetKubeConfig() (*k8s.KubeConfig, error) {
	com := exec.Command(cmd.cmd[0], cmd.cmd[1:]...)

	output, err := com.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("Error on sync command: %s", err)
	}

	kubeconfig := k8s.NewKubeConfig()
	err = kubeconfig.Unmarshal(output)
	if err != nil {
		return nil, fmt.Errorf("Bad kubeconfig file: %s", err)
	}

	return kubeconfig, nil
}
