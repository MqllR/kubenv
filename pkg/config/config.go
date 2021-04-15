package config

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
)

var (
	Conf   *Config
	confMu = &sync.Mutex{}
)

type Config struct {
	KubeConfig string                `yaml:"kubeConfig"`
	K8SConfigs map[string]*K8SConfig `mapstructure:"k8sConfigs"`
}

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

func (c *Config) FindK8SConfig(config string) *K8SConfig {
	if conf, ok := c.K8SConfigs[config]; ok {
		return conf
	}

	return &K8SConfig{}
}

func (c *Config) ListK8SConfigsNames() []string {
	var configs []string
	for config := range c.K8SConfigs {
		configs = append(configs, config)
	}

	return configs
}
