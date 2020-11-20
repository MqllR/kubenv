package dep

import (
	"fmt"
	"sync"

	"github.com/manifoldco/promptui"
	"k8s.io/klog"
)

type Dependency interface {
	GetLocalVersion() (string, error)
	GetRemoteVersion() (string, error)
	Install() error
}

func Install(exec Dependency) error {
	return exec.Install()
}

func GetVersions(tool string, exec Dependency, wg *sync.WaitGroup) {
	defer wg.Done()

	localVersion, err := exec.GetLocalVersion()
	if err != nil {
		klog.Errorf("Error when getting local version %s", err)
	}

	remoteVersion, err := exec.GetRemoteVersion()
	if err != nil {
		klog.Errorf("Error when getting local version %s", err)
	}

	var icon string
	if localVersion == remoteVersion {
		icon = promptui.IconGood
	} else {
		icon = promptui.IconBad
	}

	fmt.Printf("%s %s: local: %s\tremote: %s\n", icon, tool, localVersion, remoteVersion)
}
