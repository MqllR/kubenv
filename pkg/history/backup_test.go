package history_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/mqllr/kubenv/pkg/history"
)

func TestBackup(t *testing.T) {
	r := strings.NewReader("coucou")
	var w bytes.Buffer

	b := history.NewBackup(r, &w)
	b.Backup()

	if w.String() != "coucou" {
		t.Errorf("Expect %s but got: %s", "coucou", w.String())
	}
}
