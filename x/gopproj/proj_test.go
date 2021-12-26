package gopproj

import "testing"

// -----------------------------------------------------------------------------

func TestParseOne(t *testing.T) {
	proj, next, err := ParseOne("a.go", "b.go", "abc")
	if err != nil || len(next) != 1 || next[0] != "abc" {
		t.Fatal("ParseOne failed:", proj, next, err)
	}
}

func TestParseAll_wildcard1(t *testing.T) {
	projs, err := ParseAll("*.go")
	if err != nil || len(projs) != 1 {
		t.Fatal("ParseAll failed:", projs, err)
	}
	if proj, ok := projs[0].(*FilesProj); !ok || len(proj.Files) != 1 || proj.Files[0] != "*.go" {
		t.Fatal("ParseAll failed:", projs)
	}
}

func TestParseAll_wildcard2(t *testing.T) {
	projs, err := ParseAll("t/*.go")
	if err != nil || len(projs) != 1 {
		t.Fatal("ParseAll failed:", projs, err)
	}
	if proj, ok := projs[0].(*FilesProj); !ok || len(proj.Files) != 1 || proj.Files[0] != "t/*.go" {
		t.Fatal("ParseAll failed:", projs)
	}
}

func TestParseAll_multiFiles(t *testing.T) {
	projs, err := ParseAll("a.gop", "b.go")
	if err != nil || len(projs) != 1 {
		t.Fatal("ParseAll failed:", projs, err)
	}
	if proj, ok := projs[0].(*FilesProj); !ok || len(proj.Files) != 2 || proj.Files[0] != "a.gop" {
		t.Fatal("ParseAll failed:", proj)
	}
	projs[0].projObj()
}

func TestParseAll_multiProjs(t *testing.T) {
	projs, err := ParseAll("a/...", "./a/...", "/a")
	if err != nil || len(projs) != 3 {
		t.Fatal("ParseAll failed:", projs, err)
	}
	if proj, ok := projs[0].(*PkgPathProj); !ok || proj.Path != "a/..." {
		t.Fatal("ParseAll failed:", proj)
	}
	if proj, ok := projs[1].(*DirProj); !ok || proj.Dir != "./a/..." {
		t.Fatal("ParseAll failed:", proj)
	}
	if proj, ok := projs[2].(*DirProj); !ok || proj.Dir != "/a" {
		t.Fatal("ParseAll failed:", proj)
	}
	for _, proj := range projs {
		proj.projObj()
	}
}

func TestParseAllErr(t *testing.T) {
	_, err := ParseAll("a/...", "./a/...", "/a", "*.go")
	if err != ErrMixedFilesProj {
		t.Fatal("ParseAll:", err)
	}
}

// -----------------------------------------------------------------------------
