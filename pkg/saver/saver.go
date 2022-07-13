package saver

import (
	"fmt"
	"io"
	"os"
)

type Saver interface {
	BackupHistory() error
	SaveConfig() error
}

type Save struct {
	generator IGenerator
	source    string
}

func NewSave(gen IGenerator, filename string) *Save {
	return &Save{
		generator: gen,
		source:    filename,
	}
}

func (s *Save) BackupHistory() error {
	src, err := os.Open(s.source)
	if err != nil {
		return fmt.Errorf("Cannot open the filename %s: %s", s.source, err)
	}
	defer src.Close()

	filename := s.generator.GenerateHistoryFilename()
	dest, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("Cannot open the filename %s: %s", filename, err)
	}
	defer dest.Close()

	_, err = io.Copy(dest, src)
	if err != nil {
		return fmt.Errorf("Cannot copy the file: %s", err)
	}

	return nil
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
