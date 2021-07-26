package app

import "github.com/geomyidia/erl-midi-server/internal/util"

// Versioning data
var (
	version    string
	buildDate  string
	gitCommit  string
	gitBranch  string
	gitSummary string
)

// VersionData stuff for things
func VersionData() *util.Version {
	return &util.Version{
		Semantic:   version,
		BuildDate:  buildDate,
		GitCommit:  gitCommit,
		GitBranch:  gitBranch,
		GitSummary: gitSummary,
	}
}

// BuildString ...
func BuildString() string {
	return util.BuildString(VersionData())
}

// VersionString ...
func VersionString() string {
	return util.VersionString(VersionData())
}

// VersionedBuildString ...
func VersionedBuildString() string {
	return util.VersionedBuildString(VersionData())
}
