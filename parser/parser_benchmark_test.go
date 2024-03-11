package parser_test

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/danielgatis/go-ruby-prism/parser"
)

func clone(repo, output string) {
	_, err := os.Stat(output)
	if err != nil {
		cmd := exec.Command("git", "clone", repo, output)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		_ = cmd.Run()
	}
}

func parseAllFiles(b *testing.B, output string) {
	ctx := context.Background()
	p, _ := parser.NewParser(ctx)
	defer p.Close(ctx)

	var matches []string
	_ = filepath.Walk(output, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(info.Name(), ".rb") {
			matches = append(matches, path)
		}

		return nil
	})

	for _, file := range matches {
		source, err := os.ReadFile(file)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		_, err = p.Parse(ctx, string(source))
		if err != nil {
			b.Fatalf(fmt.Sprintf("failed to parse file %s: %s", file, err))
		}
	}
}

func Benchmark_Parser(b *testing.B) {
	table := []struct {
		repo   string
		output string
	}{
		{repo: "https://github.com/rails/rails.git", output: "../tmp/rails"},
		{repo: "https://github.com/faker-ruby/faker", output: "../tmp/faker"},
		{repo: "https://github.com/rubocop/rubocop", output: "../tmp/rubocop"},
	}

	for _, v := range table {
		b.Run(v.repo, func(b *testing.B) {
			clone(v.repo, v.output)

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				parseAllFiles(b, v.output)
			}
		})
	}
}
