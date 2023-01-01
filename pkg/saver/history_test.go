package saver_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/mqllr/kubenv/pkg/history"
	"github.com/mqllr/kubenv/pkg/saver"
)

func TestHistorySaveConfig(t *testing.T) {
	bR := strings.NewReader("foobar")
	var bW, hW bytes.Buffer

	b := history.NewBackup(bR, &bW)

	s := saver.NewHistorySave(&hW, b)

	s.SaveConfig([]byte("john"))

	if bW.String() != "foobar" {
		t.Errorf("Expect %s but got: %s", "foobar", bW.String())
	}

	if hW.String() != "john" {
		t.Errorf("Expect %s but got: %s", "john", bW.String())
	}
}
