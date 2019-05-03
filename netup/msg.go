package netup

import (
	"fmt"
	"strings"
)

func msgPrintln(v ...interface{}) {
	fmt.Println(v...)
}

func msgPrintf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func msgPrintProgress(percent int) {
	if percent < 0 {
		percent = 0
	} else if percent > 100 {
		percent = 100
	}
	const col = 68
	n := percent * col / 100
	bar := strings.Repeat("=", n) + strings.Repeat(" ", col-n)
	msgPrintf("\r    %3d%% |%s|", percent, bar)
}
