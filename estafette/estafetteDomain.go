package estafette

import manifest "github.com/estafette/estafette-ci-manifest"

// CiBuilderEvent represents a finished estafette build
type CiBuilderEvent struct {
	JobName string `json:"job_name"`
}

// CiBuilderLogLine represents a line logged by the ci builder
type CiBuilderLogLine struct {
	Time     string `json:"time"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

// CiBuilderParams contains the parameters required to create a ci builder job
type CiBuilderParams struct {
	RepoSource           string
	RepoFullName         string
	RepoURL              string
	RepoBranch           string
	RepoRevision         string
	EnvironmentVariables map[string]string
	Track                string
	AutoIncrement        int
	VersionNumber        string
	HasValidManifest     bool
	Manifest             manifest.EstafetteManifest
}
