package generator

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestGeneratedFilesAreUpToDate regenerates into a temporary directory from
// the repository's own config.yml and compares the result against the
// committed files. It fails when someone edits a generated file by hand or
// changes config.yml without running go generate.
func TestGeneratedFilesAreUpToDate(t *testing.T) {
	dir := t.TempDir()

	if err := Generate("../../config.yml", "../../prism/include/prism/version.h", dir); err != nil {
		t.Fatalf("generate: %v", err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read generated dir: %v", err)
	}

	if len(entries) == 0 {
		t.Fatal("generator produced no files")
	}

	for _, entry := range entries {
		name := entry.Name()

		want, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			t.Fatalf("read regenerated %s: %v", name, err)
		}

		got, err := os.ReadFile(filepath.Join("../../parser", name))
		if err != nil {
			t.Fatalf("read committed %s: %v", name, err)
		}

		if string(got) != string(want) {
			t.Errorf("%s is out of date; run: go generate ./internal/generator/", name)
		}
	}
}

// TestGeneratedCodeCompiles builds the parser package so a template change
// that produces invalid Go fails here rather than in a downstream build.
func TestGeneratedCodeCompiles(t *testing.T) {
	cmd := exec.Command("go", "build", "../../parser")

	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("parser package does not build: %v\n%s", err, out)
	}
}
