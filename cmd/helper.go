package cmd

import (
	"github.com/spf13/viper"
	"k8s.io/klog"
)

// Ensure all keys exist in the configuration. Exist with a fatal if not.
func isConfigExist(keysToCheck []string) {
	for _, key := range keysToCheck {
		if !viper.IsSet(key) {
			klog.Fatalf("Missing config key %s", key)
		}
	}
}

// Look for the keywork "all" in a slice
func containsAll(search []string) bool {
	for _, item := range search {
		if item == "all" {
			return true
		}
	}

	return false
}
