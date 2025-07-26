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

	fmt.Println("🚀 Starting Ruby code parsing to JSON...")

	p, _ := parser.NewParser(ctx)
	defer p.Close(ctx)

	source := "puts 'Hello, World!'"
	fmt.Println("💡 Parsing source:", source)
	result, err := p.Parse(ctx, []byte(source))
	if err != nil {
		fmt.Println("❌ Error parsing Ruby code:", err)
		os.Exit(1)
	}

	fmt.Println("📝 Converting parse result to JSON...")
	jsonResult, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println("✅ JSON output:")
	fmt.Println(string(jsonResult))
}
