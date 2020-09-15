package helper

import (
	"k8s.io/klog"

	"github.com/mqllr/kubenv/pkg/aws"
	"github.com/mqllr/kubenv/pkg/config"
)

func IsExpired(authAccount *config.AuthAccount) bool {
	session, err := aws.NewSharedSession(authAccount.AWSProfile)
	if err != nil {
		klog.Fatalf("Cannot create an AWS session: %s", err)
	}

	if !session.IsExpired() {
		klog.V(2).Infof("Token is still valid for AWS profile %s", authAccount.AWSProfile)
		return false
	}

	return true
}
