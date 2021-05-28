package config

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

var (
	// Conf store the full config. It can be used to access
	// to any configuration information after calling LoadConfig()
	Conf   *Config
	confMu = &sync.Mutex{}
)

// Config global config description
type Config struct {
	KubeConfig string                `yaml:"kubeConfig"`
	K8SConfigs map[string]*K8SConfig `mapstructure:"k8sConfigs"`
}

// LoadConfig should be call first to load the configuration. It stores
// the configuration in Conf.
func LoadConfig() error {
	confMu.Lock()
	defer confMu.Unlock()

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("Error Using config file %s: %s", viper.ConfigFileUsed(), err)
	}

	Conf = &Config{}
	err = viper.Unmarshal(Conf)
	if err != nil {
		return fmt.Errorf("Error when unmarshaling the config file %s: %s", viper.ConfigFileUsed(), err)
	}

	return nil
}
