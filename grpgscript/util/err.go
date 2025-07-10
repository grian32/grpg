package util

import (
	"fmt"
	"os"
)

func ReportErr(line int32, msg string) {
	fmt.Printf("[line %d] Error: %s", line, msg)
	os.Exit(65)
}
