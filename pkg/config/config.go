package config

import (
	"os"
	"path"
)

func GetKubeConfig() string {
	kubeconfig := os.Getenv("KUBECONFIG")

	if kubeconfig == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}

		return path.Join(home, ".kube/config")
	}

	return kubeconfig
}
