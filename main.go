package main

import (
	"conductor/cmd"

	"github.com/golang/glog"
)

func main() {
	if err := cmd.ConductorCmd.Execute(); err != nil {
		glog.Fatal(err)
	}
}
