package dep

import (
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"k8s.io/klog"

	awsgoogleauth "github.com/mqllr/kubenv/pkg/aws-google-auth"
	awsiamauthenticator "github.com/mqllr/kubenv/pkg/aws-iam-authenticator"
	"github.com/mqllr/kubenv/pkg/dep"
)

var InstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install a tools",
	Run: func(cmd *cobra.Command, args []string) {
		install(args)
	},
}

func install(args []string) {
	availableTools := map[string]dep.Dependency{
		"aws-google-auth":       &awsgoogleauth.AWSGoogleAuthExec{},
		"aws-iam-authenticator": &awsiamauthenticator.AWSIAMAuthExec{},
	}

	tools := make([]string, len(availableTools))
	i := 0
	for name := range availableTools {
		tools[i] = name
		i++
	}

	prompt := promptui.Select{
		Label: "Select an environment",
		Items: tools,
	}

	_, selectedTool, err := prompt.Run()

	if err != nil {
		klog.Fatalf("Prompt failed %v\n", err)
	}

	err = dep.Install(availableTools[selectedTool])
	if err != nil {
		klog.Errorf("Error when installing tool %s: %s", selectedTool, err)
	}

}
