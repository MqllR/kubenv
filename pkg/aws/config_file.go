package aws

import (
	"fmt"
	"os"
)

func NewConfigFile() (*IniFile, error) {
	if file, exist := os.LookupEnv("AWS_CONFIG_FILE"); exist {
		return NewIniFile(file)
	}

	homeFile := os.Getenv("HOME") + "/.aws/config"
	if _, err := os.Stat(homeFile); err == nil {
		return NewIniFile(homeFile)
	}

	return nil, fmt.Errorf("The AWS config file is not found.")
}
