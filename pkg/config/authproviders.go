package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type AuthProvider struct {
	UserName string `yaml:"UserName"`
	IDP      string `yaml:"IDP,omitempty"`
	SP       string `yaml:"SP,omitempty"`
	AppIDUri string `yaml:"AppIDUri,omitempty"`
	TenantID string `yaml:"TenantID,omitempty"`
}

type AuthProviders struct {
	Providers map[string]*AuthProvider `mapstructure:"authProviders"`
}

func NewAuthProvidersConfig() (*AuthProviders, error) {
	var auth AuthProviders
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

func (a *AuthProviders) Validate() error {
	for provider, auth := range a.Providers {
		if auth.UserName == "" {
			return fmt.Errorf("UserName undefined for provider %s", provider)
		}

		switch provider {
		case "aws-google-auth":
			if auth.IDP == "" {
				return fmt.Errorf("IDP undefined for provider %s", provider)
			}
			if auth.SP == "" {
				return fmt.Errorf("SP undefined for provider %s", provider)
			}
		case "aws-azure-login":
			if auth.TenantID == "" {
				return fmt.Errorf("TenantID undefined for provider %s", provider)
			}
			if auth.AppIDUri == "" {
				return fmt.Errorf("AppIDUri undefined for provider %s", provider)
			}
		case "aws-sts":
			continue
		default:
			return fmt.Errorf("Provider %s not implemented", provider)
		}
	}

	return nil
}

func (a *AuthProviders) FindAuthProvider(provider string) *AuthProvider {
	if acc, ok := a.Providers[provider]; ok {
		return acc
	}

	return &AuthProvider{}
}
