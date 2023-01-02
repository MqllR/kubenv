package sync

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/mqllr/kubenv/pkg/k8s"
)

type CommandExec struct {
	cmd []string
}

var _ Syncer = &CommandExec{}

// NewCommandExec just creates a CommandExec
func NewCommandExec(cmd []string) *CommandExec {
	return &CommandExec{
		cmd: cmd,
	}
}

// GetKubeConfig executes the command
// returns a KubeConfig
func (cmd *CommandExec) GetKubeConfig() (*k8s.KubeConfig, error) {
	command := exec.Command(cmd.cmd[0], cmd.cmd[1:]...)

	output, err := command.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("Error on sync command: %s", err)
	}

	r := strings.NewReader(string(output))

	kubeconfig, err := k8s.NewKubeConfigFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("Bad kubeconfig file: %s", err)
	}

	return kubeconfig, nil
}
