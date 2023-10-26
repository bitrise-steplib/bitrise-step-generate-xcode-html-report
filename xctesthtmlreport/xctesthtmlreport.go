package xctesthtmlreport

import (
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
)

const toolCmd = "xchtmlreport-bitrise"
const defaultRemoteVersion = "1.0.0"
const BitriseXcHTMLReportVersionEnvKey = "BITRISE_XCHMTL_REPORT_VERSION"

type BitriseXchtmlGenerator struct {
	Logger         log.Logger
	CommandFactory command.Factory
	EnvRepository  env.Repository
	toolPath       string
}

type Generator interface {
	Install() error
	Generate(outputPath, xcresultPath string) error
}
