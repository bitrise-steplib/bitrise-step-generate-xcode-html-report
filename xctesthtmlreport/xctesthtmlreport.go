package xctesthtmlreport

import (
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
)

const toolCmd = "xchtmlreport-bitrise"
const defaultRemoteVersion = "1.0.0"
const BitriseXcHTMLReportVersionEnvKey = "BITRISE_XCHTML_REPORT_VERSION"

type BitriseXchtmlGenerator struct {
	logger         log.Logger
	commandFactory command.Factory
	envRepository  env.Repository
	downloader     Downloader
	toolPath       string
}

type Downloader interface {
	Get(destination, source string) error
}

func NewBitriseXchtmlGenerator(
	logger log.Logger,
	commandFactory command.Factory,
	envRepository env.Repository,
	downloader Downloader,
) *BitriseXchtmlGenerator {
	return &BitriseXchtmlGenerator{
		logger:         logger,
		commandFactory: commandFactory,
		envRepository:  envRepository,
		downloader:     downloader,
	}
}

type Generator interface {
	Install() error
	Generate(outputPath, xcresultPath string) error
}
