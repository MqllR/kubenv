package aws

import (
	"fmt"
	"os"

	"gopkg.in/ini.v1"
)

type ConfigFile struct {
	Path string
}

func NewConfigFile() (*ConfigFile, error) {
	if file, exist := os.LookupEnv("AWS_CONFIG_FILE"); exist {
		return &ConfigFile{
			Path: file,
		}, nil
	}

	homeFile := os.Getenv("HOME") + "/.aws/config"
	if _, err := os.Stat(homeFile); err == nil {
		return &ConfigFile{
			Path: homeFile,
		}, nil
	}

	return nil, fmt.Errorf("The AWS config file is not found.")
}

func (c *ConfigFile) EnsureIniSection(section string, params map[string]string) error {
	cfg, err := ini.Load(c.Path)
	if err != nil {
		return nil
	}

	s, _ := cfg.GetSection(section)

	if s == nil {
		s, err = cfg.NewSection(section)
		if err != nil {
			return err
		}
	}

	for key, value := range params {
		s.Key(key).SetValue(value)
	}

	err = cfg.SaveTo(c.Path)
	if err != nil {
		return err
	}

	return nil
}
