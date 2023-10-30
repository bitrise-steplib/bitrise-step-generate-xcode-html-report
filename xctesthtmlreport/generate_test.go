package xctesthtmlreport

import (
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-steplib/bitrise-step-generate-xcode-html-report/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReportGenerator_Run(t *testing.T) {
	testingToolPath := "my/awesome/tool"
	testDeployDir := t.TempDir()
	tempDir := t.TempDir()
	cmd := new(mocks.Command)
	cmd.On("RunAndReturnTrimmedCombinedOutput").Return("", nil).Once()

	commandFactory := new(mocks.Factory)
	runFunc := func(args mock.Arguments) {
		arguments, ok := args.Get(1).([]string)
		if !ok {
			t.FailNow()
		}

		require.Equal(t, "--output", arguments[0])
		require.Equal(t, tempDir, arguments[1])
		require.Equal(t, testDeployDir, arguments[2])
	}
	commandFactory.On("Create", testingToolPath, mock.Anything, mock.Anything).Run(runFunc).Return(cmd).Once()

	bitriseXchtmlGenerator := BitriseXchtmlGenerator{
		commandFactory: commandFactory,
		logger:         log.NewLogger(),
		toolPath:       testingToolPath,
	}
	err := bitriseXchtmlGenerator.Generate(tempDir, testDeployDir)
	require.NoError(t, err)

	cmd.AssertExpectations(t)
	commandFactory.AssertExpectations(t)
}
