package saver

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type IGenerator interface {
	GenerateHistoryFilename() string
}

type Generator struct {
	baseConfigPath string
}

func NewGenerator(baseConfigPath string) (*Generator, error) {
	generator := &Generator{}

	if baseConfigPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("Cannot get the home dir: %s", err)
		}
		generator.baseConfigPath = filepath.Join(home, ".kube")
	}

	return generator, nil
}

func (g *Generator) GenerateHistoryFilename() string {
	now := time.Now().Unix()
	return fmt.Sprintf("%s/kubeconfig-%s", g.baseConfigPath, strconv.FormatInt(now, 10))
}
