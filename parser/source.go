package parser

import (
	"fmt"
	"sort"
)

type Source struct {
	Bytes       []byte
	StartLine   int
	LineOffsets []int
}

func NewSource(bytes []byte) *Source {
	return &Source{
		Bytes:     bytes,
		StartLine: 1,
	}
}

func (s *Source) SetStartLine(startLine int) {
	s.StartLine = startLine
}

func (s *Source) SetLineOffsets(lineOffsets []int) {
	s.LineOffsets = lineOffsets
}

func (s *Source) Line(byteOffset int) (int, error) {
	line, err := s.FindLine(byteOffset)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate line number: %w", err)
	}
	return s.StartLine + line, nil
}

func (s *Source) FindLine(byteOffset int) (int, error) {
	if byteOffset >= len(s.Bytes) {
		byteOffset = len(s.Bytes) - 1
	}
	if byteOffset < 0 {
		return 0, fmt.Errorf("byteOffset must be non-negative")
	}
	index := sort.Search(len(s.LineOffsets), func(i int) bool {
		return s.LineOffsets[i] > byteOffset
	})
	if index > 0 && s.LineOffsets[index-1] == byteOffset {
		index--
	}
	line := index - 1
	if index == 0 || s.LineOffsets[line] > byteOffset {
		line++
	}
	if line < 0 || line >= len(s.LineOffsets) {
		return 0, fmt.Errorf("line index out of bounds")
	}
	return line, nil
}

func (s *Source) LineCount() int {
	return len(s.LineOffsets)
}
