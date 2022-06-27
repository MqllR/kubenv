package cmd

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"k8s.io/klog"

	"github.com/mqllr/kubenv/cmd/helpers"
	"github.com/mqllr/kubenv/pkg/k8s"
)

func editCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "edit",
		Short: "manually edit a context",
		Run: func(cmd *cobra.Command, args []string) {
			editContext(args)
		},
	}
}

func editContext(args []string) {
	k, err := kubeconfig.GetKubeConfigByContextName(kubeconfig.CurrentContext)
	if err != nil {
		klog.Fatalf("Cannot get the current context: %s", err)
	}

	tempfile, err := k.WriteTempFile()
	if err != nil {
		klog.Fatalf("Cannot create a temporary file: %s", err)
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		klog.Fatal("$EDITOR is not set")
	}

	cmd := exec.Command(editor, tempfile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()

	if err != nil {
		klog.Fatalf("Cannot run the editor cmd: %s", err)
	}

	f, err := os.Open(tempfile)
	if err != nil {
		klog.Fatalf("Cannot run the editor cmd: %s", err)
	}

	k, err = k8s.NewKubeConfigFromReader(f)
	if err != nil {
		klog.Fatalf("Cannot load the updated kubeconfig: %s", err)
	}

	newKubeconfig := k8s.NewKubeConfig()

	for _, context := range kubeconfig.Contexts {
		if context.Name != kubeconfig.CurrentContext {
			c, err := kubeconfig.GetKubeConfigByContextName(context.Name)
			if err != nil {
				klog.Fatalf("Cannot load the context while rebuilding the global config: %s", err)
			}
			newKubeconfig.Append(c)
		}
	}

	newKubeconfig.Append(k)

	err = helpers.SaveKubeConfig(newKubeconfig)
	if err != nil {
		klog.Fatalf("Cannot load the updated kubeconfig: %s", err)
	}
}
