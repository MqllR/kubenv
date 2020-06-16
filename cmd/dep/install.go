package dep

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"k8s.io/klog"
	executil "k8s.io/utils/exec"

	awsgoogleauth "github.com/mqllr/kubenv/pkg/aws-google-auth"
)

var InstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install a tools",
	Run: func(cmd *cobra.Command, args []string) {
		install(args)
	},
}

func install(args []string) {
	availableTools := []string{
		"aws-google-auth",
		"aws-iam-authenticator",
	}

	prompt := promptui.Select{
		Label: "Select an environment",
		Items: availableTools,
	}

	_, selectedTool, err := prompt.Run()

	if err != nil {
		klog.Fatalf("Prompt failed %v\n", err)
	}

	execer := executil.New()

	switch selectedTool {
	case "aws-google-auth":
		runner := awsgoogleauth.New(execer)
		runner.Install()
	case "aws-iam-authenticator":
		fmt.Printf("%v Not supported yet\n", promptui.IconBad)
	default:
		fmt.Printf("%v unknow tool\n", promptui.IconBad)
	}
}
