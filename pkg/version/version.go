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
	BranchCommit  {{.Branch}}-{{.Commit}}`

var (
	// Version holds the current version of traefik.
	Version = "0.0.1"
	// BuildDate holds the build date of traefik.
	BuildDate = "I don't remember exactly"
	// StartDate holds the start date of traefik.
	StartDate = time.Now()
	// Branch holds the compiled branch
	Branch = "master"
	// Commit hold the commit hash of compiled code
	Commit = "N/A"
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
		Version:   Version,
		GoVersion: runtime.Version(),
		BuildTime: BuildDate,
		Os:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		Branch:    Branch,
		Commit:    Commit,
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, v)
	if err != nil {
		return "", xerrors.Errorf("fail to parse version")
	}

	return buf.String(), nil
}
