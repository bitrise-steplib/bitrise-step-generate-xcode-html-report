package xctesthtmlreport

import (
	"context"
	"os/exec"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/google/go-github/github"
)

// latest release:
// https://docs.github.com/en/free-pro-team@latest/rest/releases/releases?apiVersion=2022-11-28#get-the-latest-release

type Path string

const toolCmd = "xchtmlreport"

type versionProvider struct {
	localVersionProvider  func() string
	remoteVersionProvider func() string
}

func New() (Path, error) {
	shouldDownload := shouldDownload(versionProvider{
		localVersionProvider: func() string {
			return getLocalVersion()
		},
		remoteVersionProvider: func() string {
			return getRemoteVersion()
		},
	})

	if !shouldDownload {
		return toolCmd, nil
	}

	return "", nil
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

// Check if the command is installed:
func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func getRemoteVersion() string {
	ctx := context.Background()
	client := github.NewClient(nil)

	release, _, err := client.Repositories.GetLatestRelease(ctx, "bitrise-io", "XCTestHTMLReport")
	if err != nil {
		return ""
	}

	return *release.Name
}
