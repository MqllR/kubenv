package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"github.com/mqllr/kubenv/pkg/config"
	"github.com/mqllr/kubenv/pkg/k8s"
	"github.com/mqllr/kubenv/pkg/sync"
)

type SyncOptions struct {
	AppendTo bool
	Mode     string
	Path     string
	Command  string
	Glob     string
}

func syncCommand() *cobra.Command {
	opts := SyncOptions{}

	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Synchronize the kubernetes config files",
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

func runSync(opts *SyncOptions) error {
	fmt.Printf("%v Start to synchronize the kubeconfig file into %s ...\n", promptui.IconSelect, config.GetKubeConfig())

	baseKubeConfig := k8s.NewKubeConfig()
	var err error

	if opts.AppendTo {
		baseKubeConfig = kubeconfig
	}

	var s sync.Syncer

	switch opts.Mode {
	case "local":
		f, err := os.Open(opts.Path)
		if err != nil {
			return nil
		}
		s = sync.NewLocalFile(f)
	case "exec":
		cmd := strings.Split(opts.Command, " ")
		s = sync.NewCommandExec(cmd)
	case "glob":
		fs := os.DirFS("/")
		s = sync.NewGlob(fs, opts.Glob)
	default:
		return fmt.Errorf("Mode %s not supported", opts.Mode)
	}

	remoteConfig, err := s.GetKubeConfig()
	if err != nil {
		return fmt.Errorf("cannot retrieve the kubeconfig: %s", err)
	}

	baseKubeConfig.Append(remoteConfig)

	backupAndSave(baseKubeConfig)

	return nil
}
