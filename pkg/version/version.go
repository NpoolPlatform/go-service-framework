package version

import (
	"bytes"
	"runtime"
	"text/template"
	"time"

	"golang.org/x/xerrors"
)

var versionTemplate = `Version:      {{.Version}}
Go version:   {{.GoVersion}}
Built:        {{.BuildTime}}
OS/Arch:      {{.Os}}/{{.Arch}}`

var (
	// Version holds the current version of traefik.
	Version = "0.0.1"
	// BuildDate holds the build date of traefik.
	BuildDate = "I don't remember exactly"
	// StartDate holds the start date of traefik.
	StartDate = time.Now()
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
	}{
		Version:   Version,
		GoVersion: runtime.Version(),
		BuildTime: BuildDate,
		Os:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, v)
	if err != nil {
		return "", xerrors.Errorf("fail to parse version")
	}

	return buf.String(), nil
}
