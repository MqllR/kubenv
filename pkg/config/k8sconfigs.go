package config

// K8SConfig represents a K8SConfig, with only a sync description
type K8SConfig struct {
	Sync *K8SSync `yaml:"sync"`
}

// K8SSync represents a sync configuration
type K8SSync struct {
	Mode    string   `yaml:"mode"`
	Path    string   `yaml:"path,omitempty"`
	Command []string `yaml:"command,omitempty"`
}

// FindK8SConfig get a configuration for a single K8SConfig
func (c *Config) FindK8SConfig(config string) *K8SConfig {
	if conf, ok := c.K8SConfigs[config]; ok {
		return conf
	}

	return &K8SConfig{}
}

// ListK8SConfigsNames just list the names of all the K8SConfigs
func (c *Config) ListK8SConfigsNames() []string {
	var configs []string
	for config := range c.K8SConfigs {
		configs = append(configs, config)
	}

	return configs
}

/*
func (k *K8SConfigs) Validate() error {
	for _, conf := range k.Configs {
		for _, syncMode := range availableK8SSyncMode {
			if conf.Sync.Mode == syncMode {
				return nil
			}
		}

		return fmt.Errorf("K8S Sync mode %s not implemented", conf.Sync.Mode)
	}

	return nil
}
*/
