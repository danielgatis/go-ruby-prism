package parser

type Location struct {
	StartOffset int
	Length      int
}

func NewLocation(startOffset, length int) *Location {
	return &Location{StartOffset: startOffset, Length: length}
}

func (l *Location) EndOffset() int {
	return l.StartOffset + l.Length
}
