package cmd

import (
	"fmt"
	"os"
	"strings"

	"k8s.io/klog"

	"github.com/mqllr/kubenv/pkg/config"
	"github.com/mqllr/kubenv/pkg/history"
	"github.com/mqllr/kubenv/pkg/k8s"
	"github.com/mqllr/kubenv/pkg/saver"
)

func newKubeConfig() (*k8s.KubeConfig, error) {
	f, err := os.Open(config.GetKubeConfig())
	if err != nil {
		return nil, fmt.Errorf("Cannot open the kube config: %s", err)
	}

	kubeconfig, err := k8s.NewKubeConfigFromReader(f)
	if err != nil {
		return nil, fmt.Errorf("Cannot load the kubeconfig file: %s", err)
	}

	return kubeconfig, nil
}

func openKubeConfigWriter() (*os.File, error) {
	return os.OpenFile(config.GetKubeConfig(), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
}

func backupAndSave(kubeConfig *k8s.KubeConfig) {
	gen, _ := history.NewKubeHistory("")
	historyFile := gen.TimestampedFile()

	historyWriter, err := os.Create(historyFile)
	if err != nil {
		klog.Fatalf("cannot create a writer for the history: %s", err)
	}
	defer historyWriter.Close()

	myk, err := kubeconfig.ToString()
	if err != nil {
		klog.Fatalf("cannot kubeconfig convert to string: %s", err)
	}

	b := history.NewBackup(strings.NewReader(myk), historyWriter)

	c, _ := openKubeConfigWriter()
	s := saver.NewHistorySave(c, b)

	err = kubeConfig.Save(s)
	if err != nil {
		klog.Fatalf("cannot save kubeconfig: %s", err)
	}
}
