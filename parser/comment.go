package parser

import "encoding/json"

type Comment struct {
	Typpe uint32
	Loc   *Location
}

func NewComment(typpe uint32, loc *Location) *Comment {
	return &Comment{
		Typpe: typpe,
		Loc:   loc,
	}
}

func (c *Comment) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type": c.Typpe,
		"loc":  c.Loc,
	})
}

type MagicComment struct {
	KeyLocation   *Location
	ValueLocation *Location
}

func NewMagicComment(keyLocation *Location, valueLocation *Location) *MagicComment {
	return &MagicComment{
		KeyLocation:   keyLocation,
		ValueLocation: valueLocation,
	}
}

func (c *MagicComment) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"keyLocation":   c.KeyLocation,
		"valueLocation": c.ValueLocation,
	})
}
