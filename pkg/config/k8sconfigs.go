package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type K8SSync struct {
	Mode string `yaml:"mode"`
	Path string `yaml:"path,omitempty"`
}

type K8SConfig struct {
	Sync        *K8SSync `yaml:"sync"`
	AuthAccount string   `yaml:"authAccount,omitempty"`
}

type K8SConfigs struct {
	Configs map[string]*K8SConfig `mapstructure:"k8sConfigs"`
}

func NewK8SConfigs() (*K8SConfigs, error) {
	var k8s K8SConfigs
	err := viper.Unmarshal(&k8s)
	if err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %v", err)
	}

	err = k8s.Validate()
	if err != nil {
		return nil, err
	}

	return &k8s, nil
}

func (k *K8SConfigs) Validate() error {
	for _, conf := range k.Configs {
		for _, syncMode := range AvailableK8SSyncMode {
			if conf.Sync.Mode == syncMode {
				return nil
			}
		}

		return fmt.Errorf("K8S Sync mode %s not implemented", conf.Sync.Mode)
	}

	return nil
}

func (k *K8SConfigs) FindK8SConfig(config string) *K8SConfig {
	if conf, ok := k.Configs[config]; ok {
		return conf
	}

	return &K8SConfig{}
}

func (k *K8SConfigs) ListK8SConfigsNames() []string {
	var configs []string
	for config := range k.Configs {
		configs = append(configs, config)
	}

	return configs
}
