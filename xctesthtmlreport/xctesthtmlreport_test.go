package xctesthtmlreport

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

// latest release:
// https://docs.github.com/en/free-pro-team@latest/rest/releases/releases?apiVersion=2022-11-28#get-the-latest-release

func Test_shouldDownload(t *testing.T) {
	tests := []struct {
		name                  string
		localVersionProvider  func() string
		remoteVersionProvider func() string
		expected              bool
	}{
		{
			name: "no local version",
			localVersionProvider: func() string {
				return ""
			},
			expected: true,
		},
		{
			name: "local and remote has the same version",
			localVersionProvider: func() string {
				return "3.4.5\n"
			},
			remoteVersionProvider: func() string {
				return "3.4.5\n"
			},
			expected: false,
		},
		{
			name: "remote has newer patch version",
			localVersionProvider: func() string {
				return "3.4.5"
			},
			remoteVersionProvider: func() string {
				return "3.4.6"
			},
			expected: true,
		},
		{
			name: "remote has newer minor version",
			localVersionProvider: func() string {
				return "3.4.5"
			},
			remoteVersionProvider: func() string {
				return "3.5.0"
			},
			expected: true,
		},
		{
			name: "remote has newer major version",
			localVersionProvider: func() string {
				return "3.4.5"
			},
			remoteVersionProvider: func() string {
				return "4.0.0"
			},
			expected: true,
		},
		{
			name: "local has newer patch version",
			localVersionProvider: func() string {
				return "3.4.6"
			},
			remoteVersionProvider: func() string {
				return "3.4.5"
			},
			expected: false,
		},
		{
			name: "local has newer minor version",
			localVersionProvider: func() string {
				return "3.5.0"
			},
			remoteVersionProvider: func() string {
				return "3.4.5"
			},
			expected: false,
		},
		{
			name: "local has newer major version",
			localVersionProvider: func() string {
				return "6.0.0"
			},
			remoteVersionProvider: func() string {
				return "3.4.5"
			},
			expected: false,
		},
		{
			name: "local version is not a semver version",
			localVersionProvider: func() string {
				return "foo"
			},
			expected: true,
		},
		{
			name: "remote version is not a semver version",
			localVersionProvider: func() string {
				return "1.2.3"
			},
			remoteVersionProvider: func() string {
				return "foo"
			},
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := shouldDownload(versionProvider{
				localVersionProvider:  test.localVersionProvider,
				remoteVersionProvider: test.remoteVersionProvider,
			})
			require.Equal(t, test.expected, r)
		})
	}

	str := getLocalVersion()
	fmt.Print(str)

	str = getRemoteVersion()
	fmt.Print(str)
}
