package parser

import "encoding/json"

type ParseResult struct {
	Value         Node
	Comments      []*Comment
	MagicComments []*MagicComment
	DataLocation  *Location
	SynError      []*SyntaxError
	SynWarnings   []*SyntaxWarning
}

func NewParseResult(
	value Node,
	comments []*Comment,
	magicComments []*MagicComment,
	dataLocation *Location,
	synError []*SyntaxError,
	synWarnings []*SyntaxWarning,
) *ParseResult {
	return &ParseResult{
		Value:         value,
		Comments:      comments,
		MagicComments: magicComments,
		DataLocation:  dataLocation,
		SynError:      synError,
		SynWarnings:   synWarnings,
	}
}

func (p *ParseResult) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"value":         p.Value,
		"comments":      p.Comments,
		"magicComments": p.MagicComments,
		"dataLocation":  p.DataLocation,
		"synError":      p.SynError,
		"synWarnings":   p.SynWarnings,
	})
}
