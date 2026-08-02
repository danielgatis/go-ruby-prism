package parser

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestParseRubySyntax checks that representative Ruby constructs parse without
// diagnostics and produce the expected root node. Unlike the generated tests,
// which only check the shape of the types, these assert on parsing behaviour.
func TestParseRubySyntax(t *testing.T) {
	cases := []struct {
		name   string
		source string
	}{
		{"empty", ""},
		{"literal", "1"},
		{"string interpolation", `"a #{1 + 2} b"`},
		{"nested interpolation", `"a #{"b #{c}"} d"`},
		{"class with methods", "class A\n  def b; 1; end\nend"},
		{"module nesting", "module A\n  module B\n    class C; end\n  end\nend"},
		{"blocks", "[1,2].each { |x| puts x }"},
		{"do block", "[1,2].each do |x|\n  puts x\nend"},
		{"keyword args", "def f(a, b = 1, *c, d:, e: 2, **f, &g); end"},
		{"pattern matching", "case x\nin [1, *rest] then rest\nin {k: Integer => n} then n\nend"},
		{"rescue", "begin\n  raise 'x'\nrescue => e\n  retry\nensure\n  puts 1\nend"},
		{"heredoc", "x = <<~HEREDOC\n  text\nHEREDOC"},
		{"lambda", "->(x) { x * 2 }"},
		{"numbered params", "[1,2].map { _1 * 2 }"},
		{"safe navigation", "a&.b&.c"},
		{"multiple assignment", "a, b, *c = 1, 2, 3"},
		{"conditional modifiers", "puts 1 if true\nputs 2 unless false"},
		{"ternary", "x = a ? b : c"},
		{"ranges", "(1..10).to_a; (1...10).to_a"},
		{"symbols and hashes", "{a: 1, **rest}"},
		{"shorthand hash", "a = 1; {a:}"},
		{"global and class vars", "$g = 1; @@c = 2; @i = 3"},
		{"defined and yield", "def f; yield if defined?(yield); end"},
		{"operators", "a <=> b; a === b; a =~ /re/; a ** b"},
		{"frozen string comment", "# frozen_string_literal: true\nx = 'a'"},
		{"endless method", "def square(x) = x * x"},
		{"rightward assignment", "42 => x"},
	}

	ctx := context.Background()

	p, err := NewParser(ctx)
	if err != nil {
		t.Fatalf("failed to create parser: %v", err)
	}
	defer p.Close(ctx)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := p.Parse(ctx, []byte(tc.source))
			if err != nil {
				t.Fatalf("parse: %v", err)
			}

			if len(result.Errors) != 0 {
				t.Errorf("unexpected errors: %+v", result.Errors)
			}

			if result.Value == nil {
				t.Fatal("nil root node")
			}

			// Every parse produces a program at the root.
			if _, ok := interface{}(result.Value).(*ProgramNode); !ok {
				t.Errorf("root node is %T, want *ProgramNode", result.Value)
			}
		})
	}
}

// TestParseReportsInvalidSyntax checks that broken input is reported as
// diagnostics rather than as a hard failure.
func TestParseReportsInvalidSyntax(t *testing.T) {
	cases := []string{
		"def broken(",
		"class",
		"if true",
		"end",
		"[1, 2",
	}

	ctx := context.Background()

	p, err := NewParser(ctx)
	if err != nil {
		t.Fatalf("failed to create parser: %v", err)
	}
	defer p.Close(ctx)

	for _, source := range cases {
		result, err := p.Parse(ctx, []byte(source))
		if err != nil {
			t.Errorf("%q: unexpected hard error: %v", source, err)
			continue
		}

		if len(result.Errors) == 0 {
			t.Errorf("%q: expected syntax errors, got none", source)
		}
	}
}

// TestParseRailsFixtures parses the vendored Rails source. It is the broadest
// correctness signal available here: thousands of files of real Ruby, which
// between them exercise far more of the grammar than a hand-written list.
func TestParseRailsFixtures(t *testing.T) {
	root := filepath.Join("..", "test", "fixtures", "rails-main")

	if _, err := os.Stat(root); os.IsNotExist(err) {
		t.Skip("rails fixtures are not present")
	}

	var paths []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".rb") {
			paths = append(paths, path)
		}

		return nil
	})
	if err != nil {
		t.Fatalf("walk fixtures: %v", err)
	}

	if len(paths) == 0 {
		t.Skip("no ruby files in fixtures")
	}

	ctx := context.Background()

	p, err := NewParser(ctx)
	if err != nil {
		t.Fatalf("failed to create parser: %v", err)
	}
	defer p.Close(ctx)

	failures := 0

	for _, path := range paths {
		source, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}

		result, err := p.Parse(ctx, source)
		if err != nil {
			t.Errorf("%s: %v", path, err)

			failures++
			if failures > 10 {
				t.Fatal("too many failures, stopping")
			}

			continue
		}

		if len(result.Errors) != 0 {
			t.Errorf("%s: unexpected diagnostics: %+v", path, result.Errors)

			failures++
			if failures > 10 {
				t.Fatal("too many failures, stopping")
			}
		}
	}

	t.Logf("parsed %d files", len(paths))
}

// TestParseResultIsJSONSerializable checks that a tree round-trips through
// JSON, which is how the library is commonly consumed.
func TestParseResultIsJSONSerializable(t *testing.T) {
	ctx := context.Background()

	p, err := NewParser(ctx)
	if err != nil {
		t.Fatalf("failed to create parser: %v", err)
	}
	defer p.Close(ctx)

	result, err := p.Parse(ctx, []byte("class A\n  def b(c) = c + 1\nend"))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	encoded, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if _, ok := decoded["value"]; !ok {
		t.Error("serialized result has no value field")
	}
}
