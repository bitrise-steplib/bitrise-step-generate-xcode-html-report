package xctesthtmlreport

import (
	"fmt"
	"os"
)

func (b *BitriseXchtmlGenerator) Install() error {
	b.logger.Printf("Checking release versions for Bitrise-XCTestHTMLReport")

	versionOverride := b.envRepository.Get(BitriseXcHTMLReportVersionEnvKey)
	localVersion := b.getLocalVersion()
	shouldDownload := localVersion == "" || versionOverride != ""

	if !shouldDownload {
		b.logger.Printf("Using pre-installed Bitrise-XCTestHTMLReport, version: %s", localVersion)
		b.toolPath = toolCmd
		return nil
	}

	if versionOverride == "" {
		versionOverride = defaultRemoteVersion
	}

	b.logger.Printf("Downloading %s version of Bitrise-XCTestHTMLReport", versionOverride)
	path, err := b.downloadRelease(versionOverride)
	if err != nil {
		return err
	}
	b.logger.Printf("Downloading %s version of Bitrise-XCTestHTMLReport is finished", versionOverride)
	b.toolPath = path
	return nil
}

func (b *BitriseXchtmlGenerator) getLocalVersion() string {
	cmd := b.commandFactory.Create(toolCmd, []string{"--version"}, nil)
	localVersion, err := cmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return ""
	}

	return localVersion
}

func (b *BitriseXchtmlGenerator) downloadRelease(version string) (string, error) {
	temp, err := os.MkdirTemp("", toolCmd)
	if err != nil {
		return "", err
	}
	toolPath := fmt.Sprintf("%s/xchtmlreport-bitrise", temp)
	if err := b.downloader.Get(toolPath, fmt.Sprintf("https://github.com/bitrise-io/XCTestHTMLReport/releases/download/%s/xchtmlreport-bitrise", version)); err != nil {
		return "", err
	}
	err = os.Chmod(toolPath, 0755)
	if err != nil {
		return "", err
	}

	return toolPath, nil
}
