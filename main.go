package main

import (
	"fmt"
	"github.com/bitrise-io/go-steputils/v2/export"
	"os"

	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/errorutil"
	. "github.com/bitrise-io/go-utils/v2/exitcode"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-steplib/bitrise-step-generate-xcode-html-report/step"
)

func main() {
	exitCode := run()
	os.Exit(int(exitCode))
}

func run() ExitCode {
	logger := log.NewLogger()
	reportGenerator := createStep(logger)

	config, err := reportGenerator.ProcessConfig()
	if err != nil {
		logger.Println()
		logger.Errorf(errorutil.FormattedError(fmt.Errorf("Failed to process Step inputs: %w", err)))
		return Failure
	}

	if err := reportGenerator.InstallDependencies(); err != nil {
		logger.Println()
		logger.Errorf(errorutil.FormattedError(fmt.Errorf("Failed to install Step dependencies: %w", err)))
		return Failure
	}

	result, err := reportGenerator.Run(*config)
	if err != nil {
		logger.Println()
		logger.Errorf(errorutil.FormattedError(fmt.Errorf("Failed to execute Step: %w", err)))
		return Failure
	}

	if err := reportGenerator.Export(result); err != nil {
		logger.Println()
		logger.Errorf(errorutil.FormattedError(fmt.Errorf("Failed to export outputs: %w", err)))
		return Failure
	}

	return Success
}

func createStep(logger log.Logger) step.ReportGenerator {
	envRepository := env.NewRepository()
	inputParser := stepconf.NewInputParser(envRepository)
	commandFactory := command.NewFactory(envRepository)
	exporter := export.NewExporter(commandFactory)

	return step.NewReportGenerator(envRepository, inputParser, commandFactory, exporter, logger)
}
