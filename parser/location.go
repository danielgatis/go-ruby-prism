package parser

import "encoding/json"

type Location struct {
	StartOffset uint32
	Length      uint32
}

func NewLocation(startOffset uint32, length uint32) *Location {
	return &Location{
		StartOffset: startOffset,
		Length:      length,
	}
}

func (l *Location) EndOffset() uint32 {
	return l.StartOffset + l.Length
}

func (l *Location) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"startOffset": l.StartOffset,
		"length":      l.Length,
	})
}
