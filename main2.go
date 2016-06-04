package main

import (
	"github.com/koron/netupvim/netup"
)

func main2() {
	var pkg netup.Package
	err := netup.Update(".", pkg)
	if err != nil {
		netup.LogFatal(err)
	}
}
