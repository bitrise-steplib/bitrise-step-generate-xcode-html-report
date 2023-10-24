package step

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/go-steputils/v2/export"
	"github.com/bitrise-io/go-steputils/v2/stepconf"
	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/log"
)

const (
	testReportDirKey = "BITRISE_TEST_REPORT_DIR"
)

type Input struct {
	TestDeployDir    string `env:"test_result_dir,required"`
	XcresultPatterns string `env:"xcresult_patterns"`
	Verbose          bool   `env:"verbose,opt[true,false]"`
}

type Config struct {
	TestDeployDir    string
	XcresultPatterns []string
}

type Result struct {
	TestReportDir string
}

type ReportGenerator struct {
	inputParser    stepconf.InputParser
	commandFactory command.Factory
	exporter       export.Exporter
	logger         log.Logger
}

func NewReportGenerator(inputParser stepconf.InputParser, commandFactory command.Factory, exporter export.Exporter, logger log.Logger) ReportGenerator {
	return ReportGenerator{
		inputParser:    inputParser,
		commandFactory: commandFactory,
		exporter:       exporter,
		logger:         logger,
	}
}

func (r *ReportGenerator) ProcessConfig() (*Config, error) {
	var input Input
	err := r.inputParser.Parse(&input)
	if err != nil {
		return &Config{}, err
	}

	stepconf.Print(input)
	r.logger.Println()
	r.logger.EnableDebugLog(input.Verbose)

	patterns := strings.Split(input.XcresultPatterns, "\n")
	var filteredPatterns []string

	for _, pattern := range patterns {
		if pattern == "" {
			continue
		}

		if !strings.HasSuffix(pattern, ".xcresult") {
			return nil, fmt.Errorf("pattern (%s) must filter for xcresult files", pattern)
		}

		filteredPatterns = append(filteredPatterns, pattern)
	}

	return &Config{
		TestDeployDir:    input.TestDeployDir,
		XcresultPatterns: filteredPatterns,
	}, nil
}

func (r *ReportGenerator) InstallDependencies() error {
	params := []string{"install", "bitrise-io/XCTestHTMLReport@speed-improvements", "--no-link"}
	cmd := r.commandFactory.Create("mint", params, nil)

	_, err := cmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install XCTestHTMLReport: %w", err)
	}

	return nil
}

func (r *ReportGenerator) Run(config Config) (Result, error) {
	r.logger.Infof("Collecting xcresult files")

	patterns := []string{
		fmt.Sprintf("%s/**/*.xcresult", config.TestDeployDir),
	}
	if 0 < len(config.XcresultPatterns) {
		patterns = config.XcresultPatterns
	}

	paths, err := collectFilesWithPatterns(patterns)
	if err != nil {
		return Result{}, fmt.Errorf("failed to find all xcresult files: %w", err)
	}

	if len(paths) == 0 {
		r.logger.Printf("No files found.")

		return Result{
			TestReportDir: "",
		}, nil
	}

	r.logger.Printf("List of files:")
	for _, path := range paths {
		r.logger.Printf("- %s", path)
	}

	rootDir, err := testReportsRootDir()
	if err != nil {
		return Result{}, fmt.Errorf("failed to create test report directory: %w", err)
	}

	r.logger.Println()
	r.logger.Infof("Generating reports")

	for _, path := range paths {
		if err := r.generateTestReport(rootDir, path); err != nil {
			r.logger.Errorf("failed to generate test report (%s): %w", path, err)
		}
	}

	r.logger.Println()
	r.logger.Donef("Finished")

	return Result{
		TestReportDir: rootDir,
	}, nil
}

func (r *ReportGenerator) Export(result Result) error {
	return r.exporter.ExportOutput(testReportDirKey, result.TestReportDir)
}

func (r *ReportGenerator) generateTestReport(rootDir string, xcresultPath string) error {
	baseName := strings.TrimSuffix(filepath.Base(xcresultPath), filepath.Ext(xcresultPath))
	dirPath := filepath.Join(rootDir, baseName)
	err := os.Mkdir(dirPath, 0755)
	if err != nil {
		if os.IsExist(err) {
			r.logger.Warnf("Html report already exists for %s at %s", baseName, dirPath)
			return nil
		}
		return err
	}

	params := []string{"run", "bitrise-io/XCTestHTMLReport@speed-improvements", "--output", dirPath, xcresultPath}
	cmd := r.commandFactory.Create("mint", params, nil)
	_, err = cmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to export html: %w", err)
	}

	if err := moveAssets(xcresultPath, dirPath); err != nil {
		return fmt.Errorf("failed to move assets: %w", err)
	}

	return nil
}

func collectFilesWithPatterns(patterns []string) ([]string, error) {
	// Go does not have a set, so a map will help filter out duplicate results.
	allMatches := map[string]struct{}{}

	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, err
		}

		for _, match := range matches {
			allMatches[match] = struct{}{}
		}
	}

	var paths []string
	for key := range allMatches {
		paths = append(paths, key)
	}

	return paths, nil
}

func testReportsRootDir() (string, error) {
	return os.MkdirTemp("", "test-reports")
}

func moveAssets(xcresultPath string, testReportDir string) error {
	entries, err := os.ReadDir(xcresultPath)
	if err != nil {
		return err
	}

	assetFolder := filepath.Join(testReportDir, filepath.Base(xcresultPath))
	if err := os.Mkdir(assetFolder, 0755); err != nil {
		return err
	}

	for _, entry := range entries {
		// The assets are dumped into the root, so we do not need the folders.
		if entry.IsDir() {
			continue
		}

		extension := filepath.Ext(entry.Name())
		// We want to only move the useful assets which are images and videos.
		if extension == ".plist" || extension == ".log" {
			continue
		}

		oldPath := filepath.Join(xcresultPath, entry.Name())
		newPath := filepath.Join(assetFolder, filepath.Base(entry.Name()))
		if err := os.Rename(oldPath, newPath); err != nil {
			return err
		}
	}

	return nil
}
