package main

import (
	"context"
	"fmt"
	"os"

	"github.com/danielgatis/go-ruby-prism/parser"
)

var _ parser.Visitor = (*visitor)(nil)

type visitor struct {
	parser.DefaultVisitor
}

func NewVisitor() *visitor {
	return &visitor{}
}

func (v *visitor) Visit(node parser.Node) {
	fmt.Printf("🔍 Visiting node: %T\n", node)
	v.DefaultVisitor.Visit(node)
}

func (v *visitor) Traverse(node parser.Node) {
	v.Visit(node)
	for _, child := range node.ChildNodes() {
		v.Traverse(child)
	}
}

func main() {
	ctx := context.Background()

	fmt.Println("🚀 Starting Ruby AST traversal with Visitor pattern...")

	p, _ := parser.NewParser(ctx)
	defer p.Close(ctx)

	source := "puts 'Hello, World!'"
	fmt.Printf("💡 Parsing source: %s\n", source)
	result, err := p.Parse(ctx, []byte(source))
	if err != nil {
		fmt.Println("❌ Error parsing Ruby code:", err)
		os.Exit(1)
	}

	fmt.Printf("🌲 Traversing the AST...\n")
	visitor := NewVisitor()
	visitor.Traverse(result.Value)
	fmt.Println("✅ AST traversal complete!")
}
