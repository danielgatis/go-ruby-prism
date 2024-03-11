package main

import (
	"context"
	"fmt"
	"os"

	parser "github.com/danielgatis/go-ruby-prism/parser"
)

type visitor struct{}

func newVisitor() *visitor {
	return &visitor{}
}

func (v *visitor) Visit(node parser.Node) {
	fmt.Printf("%T\n", node)
}

func (v *visitor) traverse(node parser.Node) {
	node.Accept(v)
	for _, child := range node.Children() {
		v.traverse(child)
	}
}

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

	visitor := newVisitor()
	visitor.traverse(result.Value)
}
