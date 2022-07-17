package saver

import (
	"fmt"
	"os"
)

type Save struct {
	source string
}

func NewSave(filename string) *Save {
	return &Save{
		source: filename,
	}
}

func (s *Save) SaveConfig(data []byte) error {
	fh, err := os.OpenFile(s.source, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("Cannot open the kubeconfig: %s", err)
	}
	defer fh.Close()

	_, err = fh.Write(data)
	if err != nil {
		return fmt.Errorf("Cannot write the file: %s", err)
	}

	return nil
}
