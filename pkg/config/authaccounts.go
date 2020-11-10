package config

type AuthAccount struct {
	AuthProvider string `yaml:"AuthProvider"`
	AWSProfile   string `yaml:"AWSProfile,omitempty"`
	AWSRole      string `yaml:"AWSRole,omitempty"`
	Region       string `yaml:"Region,omitempty"`
	DependsOn    string `yaml:"DependsOn,omitempty"`
	Duration     int    `yaml:"Duration,omitempty"`
}

// TODO we should test if the provider is declared
/*
func (a *AuthAccounts) Validate() error {
	for env, auth := range a.Envs {
		if auth.AuthProvider == "" {
			return fmt.Errorf("Nil AuthProvider for environment %s", env)
		}

		for _, provider := range AvailableAuthProviders {
			if auth.AuthProvider == provider {
				return nil
			}
		}

		return fmt.Errorf("AuthProvider %s not implemented", auth.AuthProvider)
	}

	return nil
}
*/
