package xctesthtmlreport

import (
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/log"
)

const toolCmd = "xchtmlreport-bitrise"

type BitriseXchtmlGenerator struct {
	Logger         log.Logger
	CommandFactory command.Factory
	toolPath       string
}

type Generator interface {
	Install() error
	Generate(outputPath, xcresultPath string) error
}
