package history

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const historyPrefixFile = "config"

type IGenerator interface {
	TimestampedFile() string
}

type Generator struct {
	baseConfigPath string
}

var _ IGenerator = &Generator{}

func NewKubeHistory(baseConfigPath string) (*Generator, error) {
	generator := &Generator{}

	if baseConfigPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("cannot get the home dir: %w", err)
		}

		generator.baseConfigPath = filepath.Join(home, ".kube")
	}

	return generator, nil
}

func (g *Generator) TimestampedFile() string {
	now := time.Now().Unix()

	return fmt.Sprintf("%s/%s-%s", g.baseConfigPath, historyPrefixFile, strconv.FormatInt(now, 10))
}
