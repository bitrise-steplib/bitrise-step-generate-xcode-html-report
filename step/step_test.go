package step

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/bitrise-io/go-steputils/v2/export"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-steplib/bitrise-step-generate-xcode-html-report/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
				"test_result_dir":   "test-dir",
				"xcresult_patterns": "pattern.xcresult",
				"verbose":           "false",
			},
			want: &Config{
				TestDeployDir: "test-dir",
				XcresultPatterns: []string{
					"pattern.xcresult",
				},
			},
		},
		{
			name: "Empty input",
			envs: map[string]string{
				"test_result_dir":   "test-dir",
				"xcresult_patterns": "",
				"verbose":           "false",
			},
			want: &Config{
				TestDeployDir:    "test-dir",
				XcresultPatterns: nil,
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
				envRepository:  env.NewRepository(),
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
		envRepository:  env.NewRepository(),
		inputParser:    stepconf.NewInputParser(envRepository),
		commandFactory: commandFactory,
		exporter:       export.NewExporter(command.NewFactory(envRepository)),
		logger:         log.NewLogger(),
	}

	err := generator.InstallDependencies()
	require.NoError(t, err)

	cmd.AssertExpectations(t)
	commandFactory.AssertExpectations(t)
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
		envRepository:  envRepository,
		inputParser:    stepconf.NewInputParser(envRepository),
		commandFactory: commandFactory,
		exporter:       export.NewExporter(command.NewFactory(envRepository)),
		logger:         log.NewLogger(),
	}

	result, err := generator.Run(Config{TestDeployDir: testDeployDir})
	require.NoError(t, err)
	require.NotEqual(t, "", result.HtmlReportDir)

	cmd.AssertExpectations(t)
	commandFactory.AssertExpectations(t)
}

func TestReportGenerator_Run_ReportFolderHandling(t *testing.T) {
	tests := []struct {
		name    string
		path    func() string
		wantErr bool
	}{
		{
			name: "Existing folder",
			path: func() string {
				path := filepath.Join(t.TempDir(), "random-folder")
				require.NoError(t, os.MkdirAll(path, 0755))

				return path
			},
		},
		{
			name: "Missing folder",
			path: func() string {
				return "/path/that/does/not/exist"
			},
			wantErr: true,
		},
		{
			name: "File input",
			path: func() string {
				path := filepath.Join(t.TempDir(), "file.log")
				require.NoError(t, os.WriteFile(path, []byte("a"), 0755))

				return path
			},
			wantErr: true,
		},
		{
			name: "No input",
			path: func() string {
				return ""
			},
		},
	}
	for _, tt := range tests {
		dir := tt.path()

		envRepository := new(mocks.Repository)
		envRepository.On("Get", "BITRISE_HTML_REPORT_DIR").Return(dir)

		generator := ReportGenerator{
			envRepository:  envRepository,
			inputParser:    nil,
			commandFactory: nil,
			exporter:       export.Exporter{},
			logger:         nil,
		}
		rootDir, err := generator.htmlReportsRootDir()
		if tt.wantErr {
			assert.NotNil(t, err)
			continue
		} else {
			require.NoError(t, err)
		}

		if dir != "" {
			assert.Equal(t, dir, rootDir)
		} else {
			assert.NotEqual(t, "", rootDir)
		}
	}
}

func TestReportGenerator_Export(t *testing.T) {
	cmd := new(mocks.Command)
	cmd.On("RunAndReturnTrimmedCombinedOutput").Return("", nil).Once()

	value := "test-dir"
	arguments := []string{"add", "--key", "BITRISE_HTML_REPORT_DIR", "--value", value}
	commandFactory := new(mocks.Factory)
	commandFactory.On("Create", "envman", arguments, mock.Anything).Return(cmd).Once()

	exporter := export.NewExporter(commandFactory)
	envRepository := env.NewRepository()
	generator := ReportGenerator{
		envRepository:  envRepository,
		inputParser:    stepconf.NewInputParser(envRepository),
		commandFactory: command.NewFactory(envRepository),
		exporter:       exporter,
		logger:         log.NewLogger(),
	}

	err := generator.Export(Result{HtmlReportDir: value})
	require.NoError(t, err)
}

func setupRunEnvironment(t *testing.T) string {
	tempDir := t.TempDir()
	xcresultPath := filepath.Join(tempDir, "test-scheme", "test-scheme.xcresult")
	require.NoError(t, os.MkdirAll(xcresultPath, 0755))

	imagePaths := []string{
		filepath.Join(xcresultPath, "a.png"),
		filepath.Join(xcresultPath, "b.png"),
	}
	for _, path := range imagePaths {
		require.NoError(t, os.WriteFile(path, []byte("abcd"), 0755))
	}

	return tempDir
}
