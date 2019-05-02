package netup

import "time"

var (
	// Version stores version to log.
	Version = "none"

	// DownloadTimeout is timeout for download file.
	DownloadTimeout = 5 * time.Minute

	// GithubUser is username which be used for github's basic auth.
	GithubUser string

	// GithubVerbose enables log for github related operation.
	GithubVerbose bool

	// LogRotateCount is used for log rotation.
	LogRotateCount = 5

	// ExeRotateCount is used for executable files rotation.
	ExeRotateCount = 5
)
