package xctesthtmlreport

import (
	"fmt"
	"github.com/bitrise-io/go-utils/filedownloader"
	"github.com/bitrise-io/go-utils/retry"
	"os"
	"os/exec"
)

func (b *BitriseXchtmlGenerator) Install() error {
	b.Logger.Printf("Checking release versions for Bitrise-XCTestHTMLReport")

	versionOverride := b.EnvRepository.Get(BitriseXcHTMLReportVersionEnvKey)
	localVersion := getLocalVersion()
	shouldDownload := localVersion == "" || versionOverride != ""

	if !shouldDownload {
		b.Logger.Printf("Using pre-installed Bitrise-XCTestHTMLReport, version: %s", localVersion)
		b.toolPath = toolCmd
		return nil
	}

	if versionOverride == "" {
		versionOverride = defaultRemoteVersion
	}

	b.Logger.Printf("Downloading %s version of Bitrise-XCTestHTMLReport", versionOverride)
	path, err := downloadRelease(versionOverride)
	if err != nil {
		return err
	}
	b.Logger.Printf("Downloading %s version of Bitrise-XCTestHTMLReport is finished", versionOverride)
	b.toolPath = path
	return nil
}

func getLocalVersion() string {
	if !commandExists(toolCmd) {
		return ""
	}

	localVersion, err := exec.Command(toolCmd, "--version").Output()
	if err != nil {
		return ""
	}

	return string(localVersion)
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func downloadRelease(version string) (string, error) {
	temp, err := os.MkdirTemp("", toolCmd)
	if err != nil {
		return "", err
	}
	toolPath := fmt.Sprintf("%s/xchtmlreport-bitrise", temp)
	downloader := filedownloader.New(retry.NewHTTPClient().StandardClient())
	if err := downloader.Get(toolPath, fmt.Sprintf("https://github.com/bitrise-io/XCTestHTMLReport/releases/download/%s/xchtmlreport-bitrise", version)); err != nil {
		return "", err
	}
	err = os.Chmod(toolPath, 0755)
	if err != nil {
		return "", err
	}

	return toolPath, nil
}
