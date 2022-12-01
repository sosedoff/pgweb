package command

const (
	// Version is the current Pgweb application version
	Version = "0.11.12"
)

var (
	// GitCommit contains the Git commit SHA for the binary
	GitCommit string

	// BuildTime contains the binary build time
	BuildTime string

	// GoVersion contains the build time Go version
	GoVersion string

	// Info contains all version information
	Info VersionInfo
)

type VersionInfo struct {
	Version   string `json:"version"`
	GitCommit string `json:"git_sha"`
	BuildTime string `json:"build_time"`
	GoVersion string `json:"go_version"`
}

func init() {
	Info.Version = Version
	Info.GitCommit = GitCommit
	Info.BuildTime = BuildTime
	Info.GoVersion = GoVersion
}
