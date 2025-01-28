package parser

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
)

type Comment struct {
	Typpe int
	Loc   *Location
}

func NewComment(typpe int, loc *Location) *Comment {
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

func NewMagicComment(keyLocation, valueLocation *Location) *MagicComment {
	return &MagicComment{KeyLocation: keyLocation, ValueLocation: valueLocation}
}

type ConstantPool struct {
	Source          []byte
	BufferOffset    int
	Cache           []string
	EncodingCharset string
}

func NewConstantPool(source []byte, bufferOffset, length int, encodingCharset string) *ConstantPool {
	return &ConstantPool{
		Source:          source,
		BufferOffset:    bufferOffset,
		Cache:           make([]string, length),
		EncodingCharset: encodingCharset,
	}
}

func (cp *ConstantPool) Get(buff *SeekableBuffer, oneBasedIndex int) (string, error) {
	index := oneBasedIndex - 1

	// Check if the constant is already cached
	if constant := cp.Cache[index]; constant != "" {
		return constant, nil
	}

	// Calculate the offset and read start/length
	offset := cp.BufferOffset + index*8

	// Read `start` and `length` values
	var start, length int
	err := binary.Read(bytes.NewReader(buff.Bytes()[offset:offset+4]), binary.LittleEndian, &start)
	if err != nil {
		return "", fmt.Errorf("failed to read start: %w", err)
	}
	err = binary.Read(bytes.NewReader(buff.Bytes()[offset+4:offset+8]), binary.LittleEndian, &length)
	if err != nil {
		return "", fmt.Errorf("failed to read length: %w", err)
	}

	// Extract the bytes for the constant
	bytesData := make([]byte, length)
	if uint32(start) <= 0x7FFFFFFF {
		copy(bytesData, cp.Source[start:start+length])
	} else {
		// Read from the buffer if the high bit of `start` is set
		position := buff.Position()
		buff.SetPosition(start & 0x7FFFFFFF)
		if _, err := buff.Read(bytesData); err != nil {
			return "", fmt.Errorf("failed to read bytes: %w", err)
		}
		buff.SetPosition(position)
	}

	// Convert the bytes to a string using the loader
	constant, err := bytesToName(bytesData, cp.EncodingCharset)
	if err != nil {
		return "", fmt.Errorf("failed to convert bytes to name: %w", err)
	}

	// Cache the result
	cp.Cache[index] = constant

	return constant, nil
}

type ParseResult struct {
	Value         Node
	Comments      []*Comment
	MagicComments []*MagicComment
	DataLocation  *Location
	Errors        []*Error
	Warnings      []*Warning
	Source        *Source
}

func NewParseResult(
	value Node,
	comments []*Comment,
	magicComments []*MagicComment,
	dataLocation *Location,
	errors []*Error,
	warnings []*Warning,
	source *Source,
) *ParseResult {
	return &ParseResult{
		Value:         value,
		Comments:      comments,
		MagicComments: magicComments,
		DataLocation:  dataLocation,
		Errors:        errors,
		Warnings:      warnings,
		Source:        source,
	}
}
