package dep

import (
	"fmt"
	"sync"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	awsgoogleauth "github.com/mqllr/kubenv/pkg/aws-google-auth"
	awsiamauthenticator "github.com/mqllr/kubenv/pkg/aws-iam-authenticator"
	"github.com/mqllr/kubenv/pkg/dep"
)

var CheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check the local and remote version",
	Run: func(cmd *cobra.Command, args []string) {
		check(args)
	},
}

func check(args []string) {
	fmt.Printf("%v Checking dependencies ...\n", promptui.IconSelect)

	dependencies := map[string]dep.Dependency{
		"aws-google-auth":       &awsgoogleauth.AWSGoogleAuthExec{},
		"aws-iam-authenticator": &awsiamauthenticator.AWSIAMAuthExec{},
	}

	var wg sync.WaitGroup
	wg.Add(len(dependencies))

	for tool, depend := range dependencies {
		go dep.GetVersions(tool, depend, &wg)
	}

	wg.Wait()
}
