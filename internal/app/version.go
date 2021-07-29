package app

import "github.com/geomyidia/erl-midi-server/internal/util"

// Versioning data
var (
	version    string
	buildDate  string
	gitCommit  string
	gitBranch  string
	gitSummary string
	goVersion  string
	goArch     string
)

// VersionData stuff for things
func VersionData() *util.Version {
	return &util.Version{
		Semantic:   version,
		BuildDate:  buildDate,
		GitCommit:  gitCommit,
		GitBranch:  gitBranch,
		GitSummary: gitSummary,
		GoVersion:  goVersion,
		GoArch:  goArch,
	}
}

// BuildString ...
func BuildString() string {
	return util.BuildString(VersionData())
}

// GoVersionString ...
func GoVersionString() string {
	return util.GoVersionString(VersionData())
}

// GoArchString ...
func GoArchString() string {
	return util.GoArchString(VersionData())
}

// VersionString ...
func VersionString() string {
	return util.VersionString(VersionData())
}

// VersionedBuildString ...
func VersionedBuildString() string {
	return util.VersionedBuildString(VersionData())
}
