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
	KubeConfig    string                   `yaml:"kubeConfig"`
	K8SConfigs    map[string]*K8SConfig    `mapstructure:"k8sConfigs"`
	AuthProviders map[string]*AuthProvider `mapstructure:"authProviders"`
	AuthAccounts  map[string]*AuthAccount  `mapstructure:"authAccounts"`
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

func (c *Config) FindAuthAccount(account string) *AuthAccount {
	if acc, ok := c.AuthAccounts[account]; ok {
		return acc
	}

	return &AuthAccount{}
}

func (c *Config) ListAuthAccountNames() []string {
	var accounts []string
	for account := range c.AuthAccounts {
		accounts = append(accounts, account)
	}

	return accounts
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

func (c *Config) FindAuthProvider(provider string) *AuthProvider {
	if acc, ok := c.AuthProviders[provider]; ok {
		return acc
	}

	return &AuthProvider{}
}
