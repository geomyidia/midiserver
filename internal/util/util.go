package util

import "fmt"

const na string = "N/A"

// Version data
type Version struct {
	Semantic   string
	BuildDate  string
	GitCommit  string
	GitBranch  string
	GitSummary string
	GoVersion  string
	GoArch     string
}

// BuildString ...
func BuildString(version *Version) string {
	if version.GitCommit == "" {
		return na
	}
	return fmt.Sprintf("%s@%s, %s", version.GitBranch, version.GitCommit, version.BuildDate)
}

// GoVersionString ...
func GoVersionString(version *Version) string {
	if version.GoVersion == "" {
		return na
	}
	return version.GoVersion
}

// GoArchString ...
func GoArchString(version *Version) string {
	if version.GoArch == "" {
		return na
	}
	return version.GoArch
}

// VersionString ...
func VersionString(version *Version) string {
	if version.Semantic == "" {
		return na
	}
	return version.Semantic
}

// VersionedBuildString ...
func VersionedBuildString(version *Version) string {
	v := version.Semantic
	gc := version.GitCommit
	if v == "" {
		v = na
	}
	if gc == "" {
		gc = na
	}
	return fmt.Sprintf("%s/%s@%s, built: %s", v, version.GitBranch, gc, version.BuildDate)
}
