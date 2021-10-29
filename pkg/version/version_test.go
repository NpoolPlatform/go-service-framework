package version

import (
	"testing"
)

func TestVersion(t *testing.T) {
	_, err := GetVersion()
	if err != nil {
		t.Errorf("fail to get version: %v", err)
	}
}

func TestMyVersion(t *testing.T) {
	ver := MyVersion()
	t.Logf("my version: %v", ver)
}
