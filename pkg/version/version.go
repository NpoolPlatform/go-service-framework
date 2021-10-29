package version

import (
	"bytes"
	"runtime"
	"text/template"
	"time"

	"golang.org/x/xerrors"
)

var versionTemplate = `
	Version:      {{.Version}}
	Go version:   {{.GoVersion}}
	Built:        {{.BuildTime}}
	OS/Arch:      {{.Os}}/{{.Arch}}
	BranchCommit: {{.Branch}}-{{.Commit}}`

var (
	// Version holds the current version of traefik.
	gitVersion = "0.0.1"
	// BuildDate holds the build date of traefik.
	buildDate = "I don't remember exactly"
	// StartDate holds the start date of traefik.
	StartDate = time.Now()
	// Branch holds the compiled branch
	gitBranch = "master"
	// Commit hold the commit hash of compiled code
	gitCommit = "N/A"
)

func GetVersion() (string, error) {
	tmpl, err := template.New("").Parse(versionTemplate)
	if err != nil {
		return "", xerrors.Errorf("fail to parse version template: %v", err)
	}

	v := struct {
		Version   string
		GoVersion string
		BuildTime string
		Os        string
		Arch      string
		Branch    string
		Commit    string
	}{
		Version:   gitVersion,
		GoVersion: runtime.Version(),
		BuildTime: buildDate,
		Os:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		Branch:    gitBranch,
		Commit:    gitCommit,
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, v)
	if err != nil {
		return "", xerrors.Errorf("fail to parse version")
	}

	return buf.String(), nil
}

type Version struct {
	GitVersion string `json:"git_version"`
	BuildDate  string `json:"build_date"`
	GitBranch  string `json:"git_branch"`
	GitCommit  string `json:"git_commit"`
}

func (v Version) Equal(ver Version) bool {
	return v.GitVersion == ver.GitVersion &&
		v.GitBranch == ver.GitBranch &&
		v.GitCommit == ver.GitCommit
}

func MyVersion() Version {
	return Version{
		GitVersion: gitVersion,
		BuildDate:  buildDate,
		GitBranch:  gitBranch,
		GitCommit:  gitCommit,
	}
}
