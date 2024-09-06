package version

import "fmt"

const na string = "N/A"

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

// BuildString ...
func BuildString() string {
	if gitCommit == "" {
		return na
	}
	return fmt.Sprintf("%s@%s, built: %s", gitBranch, gitCommit, buildDate)
}

// GoVersionString ...
func GoVersionString() string {
	if goVersion == "" {
		return na
	}
	return goVersion
}

// GoArchString ...
func GoArchString() string {
	if goArch == "" {
		return na
	}
	return goArch
}

// VersionString ...
func VersionString() string {
	if version == "" {
		return na
	}
	return version
}

// VersionedBuildString ...
func VersionedBuildString() string {
	v := version
	gc := gitCommit
	if v == "" {
		v = na
	}
	if gc == "" {
		gc = na
	}
	return fmt.Sprintf("%s/%s", v, BuildString())
}
