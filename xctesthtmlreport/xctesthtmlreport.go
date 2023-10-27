package xctesthtmlreport

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/google/go-github/github"

	"github.com/bitrise-io/go-utils/filedownloader"
	"github.com/bitrise-io/go-utils/retry"
)

const toolCmd = "xchtmlreport"

type versionProvider struct {
	localVersionProvider  func() string
	remoteVersionProvider func() string
}

func New() (string, error) {
	ctx := context.Background()
	client := github.NewClient(nil)

	release, _, err := client.Repositories.GetLatestRelease(ctx, "bitrise-io", "XCTestHTMLReport")
	if err != nil {
		return "", err
	}

	shouldDownload := shouldDownload(versionProvider{
		localVersionProvider: func() string {
			return getLocalVersion()
		},
		remoteVersionProvider: func() string {
			return *release.Name
		},
	})

	if !shouldDownload {
		return toolCmd, nil
	}

	return downloadRelease(release)
}

func shouldDownload(provider versionProvider) bool {
	localVersion, err := semver.NewVersion(strings.TrimSpace(provider.localVersionProvider()))
	if err != nil {
		return true
	}

	remoteVersion, err := semver.NewVersion(strings.TrimSpace(provider.remoteVersionProvider()))
	if err != nil {
		return false
	}

	return localVersion.Compare(remoteVersion) < 0
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

func downloadRelease(release *github.RepositoryRelease) (string, error) {
	temp, err := os.MkdirTemp("", toolCmd)
	if err != nil {
		return "", err
	}
	toolPath := fmt.Sprintf("%s/xchtmlreport-bitrise", temp)
	downloader := filedownloader.New(retry.NewHTTPClient().StandardClient())
	if err := downloader.Get(toolPath, *release.Assets[0].BrowserDownloadURL); err != nil {
		return "", err
	}
	err = os.Chmod(toolPath, 0755)
	if err != nil {
		return "", err
	}

	return toolPath, nil
}
