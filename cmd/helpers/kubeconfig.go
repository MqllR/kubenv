package helpers

import (
	"fmt"
	"os"

	"github.com/mqllr/kubenv/pkg/config"
	"github.com/mqllr/kubenv/pkg/k8s"
)

func NewKubeConfig() (*k8s.KubeConfig, error) {
	f, err := os.Open(config.GetKubeConfig())
	if err != nil {
		return nil, fmt.Errorf("cannot open the kube config: %s", err)
	}

	kubeconfig, err := k8s.NewKubeConfigFromReader(f)
	if err != nil {
		return nil, fmt.Errorf("cannot load the kubeconfig file: %s", err)
	}

	return kubeconfig, nil
}

func SaveKubeConfig(kubeconfig *k8s.KubeConfig) (err error) {
	fh, err := os.OpenFile(config.GetKubeConfig(), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("cannot open the kubeconfig: %s", err)
	}
	defer func() {
		if cerr := fh.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("cannot close the kubeconfig: %s", cerr)
		}
	}()

	return kubeconfig.Save(fh)
}
