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
	DependsOn    string `yaml:"DependsOn,omitempty"`
}

type AuthAccounts struct {
	Envs map[string]*AuthAccount `mapstructure:"authAccounts"`
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

// TODO we should test if the provider is declared
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

func (a *AuthAccounts) FindAuthAccount(account string) *AuthAccount {
	if acc, ok := a.Envs[account]; ok {
		return acc
	}

	return &AuthAccount{}
}

func (a *AuthAccounts) ListAuthAccountNames() []string {
	var accounts []string
	for account := range a.Envs {
		accounts = append(accounts, account)
	}

	return accounts
}
