package env

import "testing"

func TestVersion(t *testing.T) {
	if Version() != "v"+MainVersion+".x" {
		t.Fatal("TestVersion failed:", Version())
	}
	buildVersion = "v1.0.0-beta1"
	if Version() != buildVersion {
		t.Fatal("TestVersion failed:", Version())
	}
	buildVersion = ""
}

func TestBuild(t *testing.T) {
	if BuildCommit() != "" {
		t.Fatal("BuildCommit failed:", BuildCommit())
	}
	if BuildInfo() != "() " {
		t.Fatal("BuildInfo failed:", BuildInfo())
	}
}
