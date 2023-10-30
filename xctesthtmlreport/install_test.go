package xctesthtmlreport

import (
	"fmt"
	"github.com/bitrise-io/go-utils/v2/log"
	"github.com/bitrise-steplib/bitrise-step-generate-xcode-html-report/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
)

func TestBitriseXchtmlGenerator_Install(t *testing.T) {

	testCases := []struct {
		title           string
		localVersion    string
		versionOverride string
		downloadPath    string

		shouldDownload bool
	}{
		{
			title: "when override and local version is present, the tool is downloaded",

			localVersion:    "3.4.0",
			versionOverride: "2.1.3",

			downloadPath:   "https://github.com/bitrise-io/XCTestHTMLReport/releases/download/2.1.3/xchtmlreport-bitrise",
			shouldDownload: true,
		},
		{
			title: "when no local neither override version is present, the default tool is downloaded",

			localVersion:    "",
			versionOverride: "",

			downloadPath:   "https://github.com/bitrise-io/XCTestHTMLReport/releases/download/1.0.0/xchtmlreport-bitrise",
			shouldDownload: true,
		},
		{
			title: "when local version is present and no override, the tool is not downloaded",

			localVersion:    "1.1.3",
			versionOverride: "",

			shouldDownload: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			cmd := new(mocks.Command)
			cmd.On("RunAndReturnTrimmedCombinedOutput").Return(tc.localVersion, nil).Once()

			commandFactory := new(mocks.Factory)
			versionFunc := func(args mock.Arguments) {
				arguments, ok := args.Get(1).([]string)
				if !ok {
					t.FailNow()
				}

				require.Equal(t, "--version", arguments[0])
			}
			commandFactory.On("Create", toolCmd, mock.Anything, mock.Anything).Run(versionFunc).Return(cmd).Once()

			envRepository := new(mocks.Repository)
			envRepository.On("Get", "BITRISE_XCHTML_REPORT_VERSION").Return(tc.versionOverride)

			downloadFunc := func(args mock.Arguments) {
				destination, ok := args.Get(0).(string)
				if !ok {
					t.FailNow()
				}
				source, ok := args.Get(1).(string)
				if !ok {
					t.FailNow()
				}

				require.True(t, strings.HasSuffix(destination, "/xchtmlreport-bitrise"), fmt.Sprintf("Downloading to the wrong path: %s", destination))
				require.Equal(t, source, tc.downloadPath)
				_, err := os.Create(destination)
				require.NoError(t, err)
			}
			downloader := new(mocks.Downloader)
			if tc.shouldDownload {
				downloader.
					On("Get", mock.Anything, mock.Anything).
					Run(downloadFunc).
					Return(nil).
					Once()
			}

			bitriseXchtmlGenerator := BitriseXchtmlGenerator{
				envRepository:  envRepository,
				commandFactory: commandFactory,
				logger:         log.NewLogger(),
				downloader:     downloader,
			}
			err := bitriseXchtmlGenerator.Install()
			require.NoError(t, err)

			cmd.AssertExpectations(t)
			commandFactory.AssertExpectations(t)
			envRepository.AssertExpectations(t)
			if tc.shouldDownload {
				downloader.AssertExpectations(t)
			} else {
				downloader.AssertNotCalled(t, "Get", mock.Anything, mock.Anything)
			}
		})
	}
}
