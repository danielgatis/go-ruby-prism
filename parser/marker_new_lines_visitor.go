package parser

type MarkNewlinesVisitor struct {
	BaseAbstractNodeVisitor
	source        *Source
	newlineMarked []bool
}

func NewMarkNewlinesVisitor(source *Source, newlineMarked []bool) *MarkNewlinesVisitor {
	return &MarkNewlinesVisitor{
		BaseAbstractNodeVisitor: BaseAbstractNodeVisitor{},
		source:                  source,
		newlineMarked:           newlineMarked,
	}
}

func (v *MarkNewlinesVisitor) VisitBlockNode(node *BlockNode) {
	oldNewlineMarked := v.newlineMarked
	v.newlineMarked = make([]bool, len(oldNewlineMarked))
	defer func() { v.newlineMarked = oldNewlineMarked }()
	node.Accept(v)
}

func (v *MarkNewlinesVisitor) VisitLambdaNode(node *LambdaNode) {
	oldNewlineMarked := v.newlineMarked
	v.newlineMarked = make([]bool, len(oldNewlineMarked))
	defer func() { v.newlineMarked = oldNewlineMarked }()
	node.Accept(v)
}

func (v *MarkNewlinesVisitor) VisitIfNode(node *IfNode) {
	node.SetNewLineFlag(v.source, v.newlineMarked)
	node.Accept(v)
}

func (v *MarkNewlinesVisitor) VisitUnlessNode(node *UnlessNode) {
	node.SetNewLineFlag(v.source, v.newlineMarked)
	node.Accept(v)
}

func (v *MarkNewlinesVisitor) VisitStatementsNode(node *StatementsNode) {
	for _, child := range node.Body {
		child.SetNewLineFlag(v.source, v.newlineMarked)
		child.Accept(v)
	}
}

func (v *MarkNewlinesVisitor) DefaultVisit(node Node) {
	node.VisitChildNodes(v)
}
