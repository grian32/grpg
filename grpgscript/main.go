package main

import (
	"fmt"
	"grpgscript/run"
	"os"
)

func main() {
	usageStr := "Usage: grpgscript [file]"

	if len(os.Args) > 2 {
		fmt.Println(usageStr)
		os.Exit(64)
	} else if len(os.Args) == 2 {
		if os.Args[1] == "repl" {
			fmt.Printf("Welcome to GRPGScript REPL:\n")
			run.Start(os.Stdin, os.Stdout)
		} else {
			run.RunFile(os.Args[1])
		}
	} else {
		fmt.Println(usageStr)
	}
}
