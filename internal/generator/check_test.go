package generator

import (
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestRenderMatchesPrism14 renders the templates against prism 1.4.0 and
// compares the result with the files the Ruby/ERB pipeline produced from that
// same version.
//
// This is the fidelity check for the port away from ERB. Both inputs are
// pinned in testdata, so it keeps proving the templates behave like
// template.rb no matter which prism version the project itself tracks.
func TestRenderMatchesPrism14(t *testing.T) {
	config, err := Load("testdata/config-1.4.0.yml")
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	// The 1.4.0 files were generated with these version constants baked in.
	config.Version = Version{Major: 1, Minor: 4, Patch: 0}

	cases := map[string]string{
		"gen_visitor.go.tmpl":     "gen_visitor.go",
		"gen_nodes.go.tmpl":       "gen_nodes.go",
		"gen_deserialize.go.tmpl": "gen_deserialize.go",
	}

	for tmplName, goldenName := range cases {
		t.Run(tmplName, func(t *testing.T) {
			rendered, err := Render(tmplName, config)
			if err != nil {
				t.Fatalf("render: %v", err)
			}

			formatted, err := format.Source(rendered)
			if err != nil {
				t.Fatalf("format: %v", err)
			}

			golden, err := os.ReadFile(filepath.Join("testdata", "golden-1.4.0", goldenName))
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

// TestReadVersion checks that the prism version header is parsed, since the
// deserializer rejects any tree whose version does not match.
func TestReadVersion(t *testing.T) {
	version, err := ReadVersion("../../prism/include/prism/version.h")
	if err != nil {
		t.Fatalf("read version: %v", err)
	}

	if version.Major == 0 && version.Minor == 0 && version.Patch == 0 {
		t.Fatal("parsed an all-zero version")
	}
}
