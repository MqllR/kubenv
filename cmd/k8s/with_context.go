package k8s

import (
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/klog"

	"github.com/mqllr/kubenv/pkg/config"
	"github.com/mqllr/kubenv/pkg/k8s"
	"github.com/mqllr/kubenv/pkg/prompt"
)

var WithContextCmd = &cobra.Command{
	Use:   "with-context command ...",
	Short: "Execute a command with a k8s context",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		withContext(args)
	},
}

// with-context command
func withContext(args []string) {
	var (
		kubeConfig     = config.Conf.KubeConfig
		tempKubeConfig string
	)

	klog.V(2).Infof("Read the kubeconfig file from %s", kubeConfig)

	c, err := k8s.NewKubeConfigFromFile(kubeConfig)
	contexts := c.GetContextNames()
	sort.Strings(contexts)

	klog.V(5).Infof("List of contexts: %v", contexts)

	selectedContext, err := prompt.Prompt("Select the context to use", contexts)
	if err != nil {
		klog.Fatalf("%s", err)
	}

	c.SetCurrentContext(selectedContext)

	klog.V(2).Info("Create a temporary kubeconfig file")

	temp, err := ioutil.TempFile("/tmp", "kubeconfig-*")
	if err != nil {
		klog.Fatalf("Cannot create a temporary file %s", err)
	}

	tempKubeConfig = temp.Name()
	defer func() {
		temp.Close()
		os.Remove(tempKubeConfig)
	}()

	data, err := c.Marshal()
	if err != nil {
		klog.Fatalf("Unable to marshal kubeconfig: %s", err)
	}

	_, err = temp.Write(data)
	if err != nil {
		klog.Fatalf("Error when writting the temporary kubeconfig: %s", err)
	}

	klog.V(2).Infof("Original kubeconfig copied to %s using context %s", tempKubeConfig, selectedContext)

	cmd := exec.Command("/bin/sh", "-c", strings.Join(args, " "))
	cmd.Env = []string{
		"KUBECONFIG=" + tempKubeConfig,
	}
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	klog.V(2).Infof("Running command %s with environment var KUBECONFIG=%s", strings.Join(args, " "), tempKubeConfig)

	cmd.Run()
}
