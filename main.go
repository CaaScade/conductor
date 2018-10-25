package main

import (
	"conductor/cmd"

	//initialized api types
	_ "conductor/pkg/api/v1alpha1"
	_ "conductor/pkg/store"

	"github.com/golang/glog"
)

func main() {
	if err := cmd.ConductorCmd.Execute(); err != nil {
		glog.Fatal(err)
	}
}
