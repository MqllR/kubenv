package cmd

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"k8s.io/klog"

	"github.com/mqllr/kubenv/pkg/config"
	"github.com/mqllr/kubenv/pkg/k8s"
	"github.com/mqllr/kubenv/pkg/sync"
)

func syncCommand() *cobra.Command {
	opts := sync.SyncOptions{}

	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Synchronize the kubernetes config files",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return validateFlags(&opts)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSync(&opts)
		},
	}

	f := cmd.Flags()
	f.BoolVarP(&opts.AppendTo, "append", "a", true, "Append the new kubeconfig files to ~/.kube/config")
	f.StringVarP(&opts.Mode, "mode", "m", "local", "Mode to read a kubeconfig file. Either local, exec or glob")
	f.StringVar(&opts.Path, "path", "", "A path to the local kubeconfig file")
	f.StringVar(&opts.Command, "command", "", "A command to execute to retrieve the kubeconfig file")
	f.StringVar(&opts.Glob, "glob", "", "A glob pattern to retrieve the kubeconfig files. Relative path to / (ex: home/foo/bar/*)")

	return cmd
}

func validateFlags(opts *sync.SyncOptions) error {
	existInSyncMode := func(mode string) bool {
		for _, m := range config.SyncMode {
			if mode == m {
				return true
			}
		}
		return false
	}

	if !existInSyncMode(opts.Mode) {
		return fmt.Errorf("Mode %s not supported", opts.Mode)
	}

	return nil
}

func runSync(opts *sync.SyncOptions) error {
	fmt.Printf("%v Start to synchronize the kubeconfig file into %s ...\n", promptui.IconSelect, config.GetKubeConfig())

	baseKubeConfig := k8s.NewKubeConfig()
	var err error

	if opts.AppendTo {
		f, err := os.Open(config.GetKubeConfig())
		if err != nil {
			klog.Fatalf("Cannot open the kube config: %s", err)
		}

		baseKubeConfig, err = k8s.NewKubeConfigFromReader(f)
		if err != nil {
			return fmt.Errorf("Something went wrong: %s", err)
		}
	}

	svc := sync.NewService(opts)
	kubeconfig, err := svc.GetKubeConfig()
	if err != nil {
		return fmt.Errorf("Cannot retrieve the kubeconfig: %s", err)
	}

	baseKubeConfig.Append(kubeconfig)

	fh, err := os.OpenFile(config.GetKubeConfig(), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return fmt.Errorf("Cannot open the kubeconfig: %s", err)
	}
	defer fh.Close()

	err = baseKubeConfig.Save(fh)
	if err != nil {
		fmt.Printf("%v Failed to write the kubeconfig file: %s", promptui.IconBad, err)
	}

	return nil
}
