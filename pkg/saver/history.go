package saver

import (
	"fmt"
	"io"

	"github.com/mqllr/kubenv/pkg/history"
)

type History struct {
	w io.Writer
	b *history.Back
}

var _ Saver = &History{}

func NewHistorySave(w io.Writer, backup *history.Back) *History {
	return &History{w, backup}
}

func (h *History) SaveConfig(data []byte) error {
	if err := h.b.Backup(); err != nil {
		return fmt.Errorf("failed to backup: %w", err)
	}

	_, err := h.w.Write(data)
	if err != nil {
		return fmt.Errorf("cannot save config: %w", err)
	}

	return nil
}
