package main

import (
	"github.com/mqllr/kubenv/cmd"
	"k8s.io/klog"
)

func main() {
	if err := cmd.Execute(); err != nil {
		klog.Fatal(err)
	}
}
