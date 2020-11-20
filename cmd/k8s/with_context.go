package k8s

import (
	"io/ioutil"
	"os"
	"sort"

	"github.com/spf13/cobra"
	"k8s.io/klog"

	"github.com/mqllr/kubenv/pkg/config"
	"github.com/mqllr/kubenv/pkg/k8s"
	"github.com/mqllr/kubenv/pkg/prompt"
)

var WithContextCmd = &cobra.Command{
	Use:   "with-context [context] command",
	Short: "Execute a command with a k8s context",
	Args:  cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		withContext(args)
	},
}

// with-context command
func withContext(args []string) {
	var (
		tempKubeConfig string
		kubeConfig     = config.Conf.KubeConfig
		_              = args[len(args)-1]
	)

	config, err := ioutil.ReadFile(kubeConfig)
	if err != nil {
		klog.Fatalf("Cannot read kubeconfig file: %s", err)
	}

	temp, err := ioutil.TempFile("/tmp", "kubeconfig-*")
	if err != nil {
		klog.Fatalf("Cannot create a temporary file %s", err)
	}

	tempKubeConfig = temp.Name()

	defer func() {
		temp.Close()
		os.Remove(tempKubeConfig)
	}()

	_, err = temp.Write(config)
	if err != nil {
		klog.Fatalf("")
	}
	klog.V(2).Infof("Original kubeconfig copied to %s", tempKubeConfig)

	kubeconfig := k8s.NewKubeConfig()
	if err = kubeconfig.Unmarshal(config); err != nil {
		klog.Fatalf("Cannot unmarshal config: %s", err)
	}

	contexts := kubeconfig.GetContextNames()

	var selectedContext string

	sort.Strings(contexts)

	selectedContext, err = prompt.Prompt("Select the context", contexts)
	if err != nil {
		klog.Fatalf("%s", err)
	}

	klog.Info(selectedContext)

	//	kubeconfig.CurrentContext = selectedContext
	//	kubeconfig.WriteFile(kubeConfig)
}
