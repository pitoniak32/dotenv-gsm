package version

import (
	"runtime"
)

type AppVersionInfo struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	Name      string `json:"name"`
	Arch      string `json:"arch"`
	Os        string `json:"os"`
	GoVersion string `json:"go_version"`
	BuildDate string `json:"build_date"`
}

var (
	version     = "UNKNOWN"
	commit      = "UNKNOWN"
	buildDate   = "UNKNOWN"
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
