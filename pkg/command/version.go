package command

const (
	// Version is the current Pgweb application version
	Version = "0.11.8"
)

var (
	// GitCommit contains the Git commit SHA for the binary
	GitCommit string

	// BuildTime contains the binary build time
	BuildTime string

	// GoVersion contains the build time Go version
	GoVersion string
)
