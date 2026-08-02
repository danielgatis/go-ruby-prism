package generator

import (
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestRenderMatchesCommitted renders each template and compares it against the
// committed generated file. The committed files were produced by the previous
// Ruby/ERB pipeline, so a byte-identical match proves the Go port is faithful.
func TestRenderMatchesCommitted(t *testing.T) {
	config, err := Load(oracleConfig())
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	cases := map[string]string{
		"gen_visitor.go.tmpl":     "../../parser/gen_visitor.go",
		"gen_nodes.go.tmpl":       "../../parser/gen_nodes.go",
		"gen_deserialize.go.tmpl": "../../parser/gen_deserialize.go",
	}

	for tmplName, goldenPath := range cases {
		t.Run(tmplName, func(t *testing.T) {
			if _, err := os.Stat(filepath.Join("templates", tmplName)); err != nil {
				t.Skipf("template not written yet: %v", err)
			}

			rendered, err := Render(tmplName, config)
			if err != nil {
				t.Fatalf("render: %v", err)
			}

			formatted, err := format.Source(rendered)
			if err != nil {
				t.Fatalf("format: %v", err)
			}

			golden, err := os.ReadFile(goldenPath)
			if err != nil {
				t.Fatalf("read golden: %v", err)
			}

			if string(formatted) == string(golden) {
				return
			}

			gotLines := strings.Split(string(formatted), "\n")
			wantLines := strings.Split(string(golden), "\n")

			for i := 0; i < len(gotLines) && i < len(wantLines); i++ {
				if gotLines[i] != wantLines[i] {
					t.Fatalf("first difference at line %d:\n  got:  %q\n  want: %q\n(total lines got=%d want=%d)",
						i+1, gotLines[i], wantLines[i], len(gotLines), len(wantLines))
				}
			}

			t.Fatalf("output length differs: got %d lines, want %d lines", len(gotLines), len(wantLines))
		})
	}
}

// oracleConfig returns the config.yml that the committed generated files were
// produced from. The repository's own config.yml has drifted ahead of those
// files, so it cannot serve as the fidelity oracle for the port.
func oracleConfig() string {
	if path := os.Getenv("PRISM_ORACLE_CONFIG"); path != "" {
		return path
	}

	return "testdata/config-1.4.0.yml"
}
