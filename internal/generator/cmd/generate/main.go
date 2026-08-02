// Command generate renders the parser package's generated sources from
// prism's config.yml. It replaces the Ruby/ERB pipeline, so regenerating this
// project requires only a Go toolchain.
package main

import (
	"flag"
	"log"

	"github.com/danielgatis/go-ruby-prism/internal/generator"
)

func main() {
	config := flag.String("config", "config.yml", "path to prism's config.yml")
	version := flag.String("version", "prism/include/prism/version.h", "path to prism's version.h")
	out := flag.String("out", "parser", "directory to write the generated files into")
	flag.Parse()

	if err := generator.Generate(*config, *version, *out); err != nil {
		log.Fatal(err)
	}
}
