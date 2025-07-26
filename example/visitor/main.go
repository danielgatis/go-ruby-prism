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
	fmt.Printf("%T\n", node)
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

	p, _ := parser.NewParser(ctx)
	defer p.Close(ctx)

	source := "puts 'Hello, World!'"
	result, err := p.Parse(ctx, []byte(source))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	visitor := NewVisitor()
	visitor.Traverse(result.Value)
}
