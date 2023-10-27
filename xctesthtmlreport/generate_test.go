package xctesthtmlreport

//
//import (
//	"github.com/bitrise-io/go-steputils/v2/export"
//	"github.com/bitrise-io/go-steputils/v2/stepconf"
//	"github.com/bitrise-io/go-utils/v2/command"
//	"github.com/bitrise-io/go-utils/v2/env"
//	"github.com/bitrise-io/go-utils/v2/log"
//	"github.com/bitrise-steplib/bitrise-step-generate-xcode-html-report/mocks"
//	"github.com/stretchr/testify/mock"
//	"github.com/stretchr/testify/require"
//	"testing"
//)
//
//func TestReportGenerator_Run(t *testing.T) {
//	testDeployDir := setupRunEnvironment(t)
//
//	cmd := new(mocks.Command)
//	cmd.On("RunAndReturnTrimmedCombinedOutput").Return("", nil).Once()
//
//	commandFactory := new(mocks.Factory)
//	runFunc := func(args mock.Arguments) {
//		arguments, ok := args.Get(1).([]string)
//		if !ok {
//			t.FailNow()
//		}
//
//		require.Equal(t, "run", arguments[0])
//		require.Equal(t, "bitrise-io/XCTestHTMLReport@speed-improvements", arguments[1])
//		require.Equal(t, "--output", arguments[2])
//	}
//	commandFactory.On("Create", "mint", mock.Anything, mock.Anything).Run(runFunc).Return(cmd).Once()
//
//	envRepository := env.NewRepository()
//	generator := ReportGenerator{
//		envRepository:  envRepository,
//		inputParser:    stepconf.NewInputParser(envRepository),
//		commandFactory: commandFactory,
//		exporter:       export.NewExporter(command.NewFactory(envRepository)),
//		logger:         log.NewLogger(),
//	}
//
//	result, err := generator.Run(Config{TestDeployDir: testDeployDir})
//	require.NoError(t, err)
//	require.NotEqual(t, "", result.HtmlReportDir)
//
//	cmd.AssertExpectations(t)
//	commandFactory.AssertExpectations(t)
//}
