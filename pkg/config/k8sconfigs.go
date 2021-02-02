package config

type K8SConfig struct {
	Sync        *K8SSync `yaml:"sync"`
	AuthAccount string   `yaml:"authAccount,omitempty"`
}

type K8SSync struct {
	Mode    string   `yaml:"mode"`
	Path    string   `yaml:"path,omitempty"`
	Command []string `yaml:"command,omitempty"`
}

/*
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
*/
