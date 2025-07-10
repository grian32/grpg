package main

import (
	"fmt"
	"os"
)

func main() {
	usageStr := "Usage: grpgscript [file]"

	if len(os.Args) > 2 {
		fmt.Println(usageStr)
		os.Exit(64)
	} else if len(os.Args) == 2 {
		RunFile(os.Args[1])
	} else {
		fmt.Println(usageStr)
	}
}
