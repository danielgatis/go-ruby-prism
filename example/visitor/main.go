package main

import (
	"context"
	"fmt"
	"os"

	"github.com/danielgatis/go-ruby-prism/parser"
)

var _ parser.AbstractNodeVisitor = (*visitor)(nil)

type visitor struct {
	parser.BaseAbstractNodeVisitor
}

func NewVisitor() *visitor {
	return &visitor{}
}

func (v *visitor) DefaultVisit(node parser.Node) {
	fmt.Printf("%T\n", node)
}

func (v *visitor) traverse(node parser.Node) {
	node.Accept(v)
	for _, child := range node.ChildNodes() {
		v.traverse(child)
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
	visitor.traverse(result.Value)
}
