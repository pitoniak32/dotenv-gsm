package version

import (
	"runtime"
)

type AppVersionInfo struct {
	Arch      string `json:"arch"`
	Os        string `json:"os"`
	GoVersion string `json:"go_version"`
	Name      string `json:"name"`
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildDate string `json:"build_date"`
}

var (
	version     = "v0.0.0"
	commit      = "dirty"
	buildDate   = "1970-01-01T00:00:00Z"
	VersionInfo AppVersionInfo
)

func init() {
	VersionInfo = AppVersionInfo{
		Arch:      runtime.GOARCH,
		Os:        runtime.GOOS,
		GoVersion: runtime.Version(),
		Name:      "dotenv_gsm",
		Version:   version,
		Commit:    commit,
		BuildDate: buildDate,
	}
}
