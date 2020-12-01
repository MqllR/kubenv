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

	exe, err := exec.LookPath(args[0])
	if err != nil {
		klog.Fatal(err)
	}

	envs := os.Environ()
	isExist := func(envs []string, key string) (bool, int) {
		for i, env := range envs {
			if env == key {
				return true, i
			}
		}

		return false, 0
	}

	exist, i := isExist(envs, "KUBECONFIG")
	localKubeConfig := "KUBECONFIG=" + tempKubeConfig
	if exist {
		envs[i] = localKubeConfig
	} else {
		envs = append(envs, localKubeConfig)
	}

	cmd := exec.Cmd{
		Path:   exe,
		Args:   args[0:],
		Env:    envs,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	klog.V(2).Infof("Running command: %s", strings.Join(args, " "))
	klog.V(5).Infof("Running command: %s with environment variable %v", strings.Join(args, " "), envs)

	cmd.Run()
}