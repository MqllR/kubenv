package aws

import (
	"fmt"
	"os"
)

func NewCredFile() (*IniFile, error) {
	if file, exist := os.LookupEnv("AWS_SHARED_CREDENTIALS_FILE"); exist {
		return NewIniFile(file)
	}

	homeFile := os.Getenv("HOME") + "/.aws/credentials"
	if _, err := os.Stat(homeFile); err == nil {
		return NewIniFile(homeFile)
	}

	return nil, fmt.Errorf("The AWS credentials file is not found.")
}
