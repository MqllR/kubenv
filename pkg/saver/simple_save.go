package saver

import (
	"fmt"
	"io"
)

type SimpleSave struct {
	w io.Writer
}

var _ Saver = &SimpleSave{}

func NewSimpleSave(w io.Writer) *SimpleSave {
	return &SimpleSave{w}
}

func (s *SimpleSave) SaveConfig(data []byte) error {
	_, err := s.w.Write(data)
	if err != nil {
		return fmt.Errorf("cannot save config: %w", err)
	}

	return nil
}
