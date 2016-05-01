package main

import "fmt"

func msgPrint(v ...interface{}) {
	fmt.Print(v...)
}

func msgPrintln(v ...interface{}) {
	fmt.Println(v...)
}

func msgPrintf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}
