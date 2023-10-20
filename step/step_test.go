package step

import (
	"fmt"
	"github.com/bitrise-steplib/bitrise-step-generate-xcode-html-report/mocks"
	"github.com/stretchr/testify/mock"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/bitrise-io/go-steputils/v2/export"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/stretchr/testify/require"
)

func TestReportGenerator_ProcessConfig(t *testing.T) {
	tests := []struct {
		name string
		envs map[string]string
		want *Config
	}{
		{
			name: "Single input",
			envs: map[string]string{
				"BITRISE_TEST_DEPLOY_DIR": "test-dir",
				"verbose":                 "false",
			},
			want: &Config{
				TestDeployDir: "test-dir",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			envRepository := new(mocks.Repository)
			for key, value := range tt.envs {
				envRepository.On("Get", key).Return(value)
			}

			inputParser := stepconf.NewInputParser(envRepository)
			commandFactory := command.NewFactory(envRepository)
			generator := ReportGenerator{
				inputParser:    inputParser,
				commandFactory: commandFactory,
				exporter:       export.Exporter{},
				logger:         log.NewLogger(),
			}

			got, err := generator.ProcessConfig()
			require.NoError(t, err)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProcessConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReportGenerator_InstallDependencies(t *testing.T) {
	cmd := new(mocks.Command)
	cmd.On("RunAndReturnTrimmedCombinedOutput").Return("", nil).Once()

	arguments := []string{"install", "bitrise-io/XCTestHTMLReport@speed-improvements", "--no-link"}
	commandFactory := new(mocks.Factory)
	commandFactory.On("Create", "mint", arguments, mock.Anything).Return(cmd).Once()

	envRepository := env.NewRepository()
	generator := ReportGenerator{
		inputParser:    stepconf.NewInputParser(envRepository),
		commandFactory: command.NewFactory(envRepository),
		exporter:       export.NewExporter(commandFactory),
		logger:         log.NewLogger(),
	}

	err := generator.InstallDependencies()
	require.NoError(t, err)
}

func TestReportGenerator_Run(t *testing.T) {
	testDeployDir := setupRunEnvironment(t)

	cmd := new(mocks.Command)
	cmd.On("RunAndReturnTrimmedCombinedOutput").Return("", nil).Once()

	commandFactory := new(mocks.Factory)
	runFunc := func(args mock.Arguments) {
		arguments, ok := args.Get(1).([]string)
		if !ok {
			t.FailNow()
		}

		require.Equal(t, "run", arguments[0])
		require.Equal(t, "bitrise-io/XCTestHTMLReport@speed-improvements", arguments[1])
		require.Equal(t, "--output", arguments[2])
	}
	commandFactory.On("Create", "mint", mock.Anything, mock.Anything).Run(runFunc).Return(cmd).Once()

	envRepository := env.NewRepository()
	generator := ReportGenerator{
		inputParser:    stepconf.NewInputParser(envRepository),
		commandFactory: commandFactory,
		exporter:       export.NewExporter(command.NewFactory(envRepository)),
		logger:         log.NewLogger(),
	}

	result, err := generator.Run(Config{TestDeployDir: testDeployDir})
	require.NoError(t, err)

	fmt.Println(result)
}

func setupRunEnvironment(t *testing.T) string {
	tempDir := t.TempDir()
	xcresultPath := filepath.Join(tempDir, "test-scheme.xcresult")
	require.NoError(t, os.Mkdir(xcresultPath, 0755))

	imagePaths := []string{
		filepath.Join(xcresultPath, "a.png"),
		filepath.Join(xcresultPath, "b.png"),
	}
	for _, path := range imagePaths {
		require.NoError(t, os.WriteFile(path, []byte("abcd"), 0755))
	}

	return tempDir
}

func TestReportGenerator_Export(t *testing.T) {
	cmd := new(mocks.Command)
	cmd.On("RunAndReturnTrimmedCombinedOutput").Return("", nil).Once()

	value := "test-dir"
	arguments := []string{"add", "--key", "BITRISE_TEST_REPORT_DIR", "--value", value}
	commandFactory := new(mocks.Factory)
	commandFactory.On("Create", "envman", arguments, mock.Anything).Return(cmd).Once()

	exporter := export.NewExporter(commandFactory)
	envRepository := env.NewRepository()
	generator := ReportGenerator{
		inputParser:    stepconf.NewInputParser(envRepository),
		commandFactory: command.NewFactory(envRepository),
		exporter:       exporter,
		logger:         log.NewLogger(),
	}

	err := generator.Export(Result{TestReportDir: value})
	require.NoError(t, err)
}