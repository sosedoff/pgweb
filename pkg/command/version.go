package command

import (
	"fmt"
	"runtime"
	"strings"
)

const (
	// Version is the current Pgweb application version
	Version = "0.14.0"
)

var (
	// GitCommit contains the Git commit SHA for the binary
	GitCommit string

	// BuildTime contains the binary build time
	BuildTime string

	// BuildArch contains the OS architecture of the binary
	BuildArch string = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)

	// GoVersion contains the build time Go version
	GoVersion string

	// Info contains all version information
	Info VersionInfo
)

type VersionInfo struct {
	Version   string `json:"version"`
	GitCommit string `json:"git_sha"`
	BuildTime string `json:"build_time"`
	BuildArch string `json:"build_arch"`
	GoVersion string `json:"go_version"`
}

func init() {
	Info.Version = Version
	Info.GitCommit = GitCommit
	Info.BuildTime = BuildTime
	Info.BuildArch = BuildArch
	Info.GoVersion = GoVersion
}

func VersionString() string {
	chunks := []string{fmt.Sprintf("Pgweb v%s", Version)}

	if GitCommit != "" {
		chunks = append(chunks, fmt.Sprintf("(git: %s)", GitCommit))
	}
	if GoVersion != "" {
		chunks = append(chunks, fmt.Sprintf("(go: %s)", GoVersion))
	}
	if BuildTime != "" {
		chunks = append(chunks, fmt.Sprintf("(build time: %s)", BuildTime))
	}
	if BuildArch != "" {
		chunks = append(chunks, fmt.Sprintf("(arch: %s)", BuildArch))
	}

	return strings.Join(chunks, " ")
}
