package dep

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"k8s.io/klog"
	executil "k8s.io/utils/exec"

	awsgoogleauth "github.com/mqllr/kubenv/pkg/aws-google-auth"
	awsiamauthenticator "github.com/mqllr/kubenv/pkg/aws-iam-authenticator"
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
	execer := executil.New()

	gAuth := awsgoogleauth.New(execer)
	version, err := gAuth.GetVersion()
	if err != nil {
		klog.Errorf("Error when getting local version %s", err)
	}

	remoteVersion, err := gAuth.GetRemoteVersion()
	if err != nil {
		klog.Errorf("Error when getting remote version %s", err)
	}

	fmt.Printf("%s:\n", awsgoogleauth.AWSGoogleAuthCmd)
	compareVersion(version, remoteVersion)

	aAuth := awsiamauthenticator.New(execer)
	version, err = aAuth.GetVersion()

	if err != nil {
		klog.Errorf("Error when getting version %s", err)
	}

	remoteVersion, err = aAuth.GetRemoteVersion()
	if err != nil {
		klog.Errorf("Error when remote version %s", err)
	}

	fmt.Printf("%s:\n", awsiamauthenticator.AWSIAMAuthCmd)
	compareVersion(version, remoteVersion)
}

func compareVersion(localVersion string, remoteVersion string) {
	var icon string
	if localVersion == remoteVersion {
		icon = promptui.IconGood
	} else {
		icon = promptui.IconBad
	}

	fmt.Printf("\t%v local: %s\tremote: %s\n", icon, localVersion, remoteVersion)
}
