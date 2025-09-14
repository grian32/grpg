package main

import (
	grpgscriptlsp "grpgscript-lsp"

	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopmentConfig().Build()

	grpgscriptlsp.StartServer(logger)
}
