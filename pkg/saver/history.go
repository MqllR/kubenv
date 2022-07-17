package saver

import (
	"fmt"
	"io"
	"os"
)

type History struct {
	generator IGenerator
	source    string
}

func NewHistory(gen IGenerator, filename string) *History {
	return &History{
		generator: gen,
		source:    filename,
	}
}

func (h *History) SaveConfig(data []byte) error {
	err := h.backupHistory()
	if err != nil {
		return fmt.Errorf("Cannot save the history: %s", err)
	}

	fh, err := os.OpenFile(h.source, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
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

func (h *History) backupHistory() error {
	src, err := os.Open(h.source)
	if err != nil {
		return fmt.Errorf("Cannot open the filename %s: %s", h.source, err)
	}
	defer src.Close()

	filename := h.generator.GenerateHistoryFilename()
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
