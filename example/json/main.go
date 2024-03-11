package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	parser "github.com/danielgatis/go-ruby-prism/parser"
)

func main() {
	ctx := context.Background()

	p, _ := parser.NewParser(ctx)
	defer p.Close(ctx)

	source := "puts 'Hello, World!'"
	result, err := p.Parse(ctx, source)
	if err != nil {
		fmt.Println(parser.ErrToStr(err))
		os.Exit(1)
	}

	jsonResult, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(jsonResult))
}
