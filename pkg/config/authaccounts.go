package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type AuthAccount struct {
	AuthProvider string `yaml:"AuthProvider"`
	AWSProfile   string `yaml:"AWSProfile,omitempty"`
	AWSRole      string `yaml:"AWSRole,omitempty"`
	Region       string `yaml:"Region,omitempty"`
}

type AuthAccounts struct {
	Env map[string]*AuthAccount `mapstructure:"authAccounts"`
}

var AuthProviders = []string{
	"aws-google-auth",
}

func NewAuthAccountsConfig() (*AuthAccounts, error) {
	var auth AuthAccounts
	err := viper.Unmarshal(&auth)
	if err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %v", err)
	}

	err = auth.Validate()
	if err != nil {
		return nil, err
	}

	return &auth, nil
}

func (a *AuthAccounts) Validate() error {
	for env, auth := range a.Env {
		if auth.AuthProvider == "" {
			return fmt.Errorf("Nil AuthProvider for environment %s", env)
		}

		for _, provider := range AuthProviders {
			if auth.AuthProvider == provider {
				return nil
			}
		}
		return fmt.Errorf("The AuthProvider %s doesn't exist", auth.AuthProvider)
	}
	return nil
}
