package helper

import (
	"github.com/spf13/viper"
	"k8s.io/klog"
)

// Ensure all keys exist in the configuration. Exist with a fatal if not.
func IsConfigExist(keysToCheck []string) {
	for _, key := range keysToCheck {
		if !viper.IsSet(key) {
			klog.Fatalf("Missing config key %s", key)
		}
	}
}
