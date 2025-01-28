package parser

import (
	"io"
)

// SeekableBuffer is a buffer that can be read from and written to at an arbitrary position.
type SeekableBuffer struct {
	data []byte
	pos  int
}

// NewSeekableBuffer initializes a SeekableBuffer with the given data.
func NewSeekableBuffer(data []byte) *SeekableBuffer {
	return &SeekableBuffer{
		data: data,
		pos:  0,
	}
}

// Position returns the current position in the buffer.
func (sb *SeekableBuffer) Position() int {
	return sb.pos
}

// Read reads data from the buffer starting at the current position.
func (sb *SeekableBuffer) Read(p []byte) (int, error) {
	if sb.pos >= len(sb.data) {
		return 0, io.EOF
	}

	n := copy(p, sb.data[sb.pos:])
	sb.pos += n
	return n, nil
}

// ReadByte reads a single byte from the buffer at the current position.
func (sb *SeekableBuffer) ReadByte() (byte, error) {
	if sb.pos >= len(sb.data) {
		return 0, io.EOF
	}

	b := sb.data[sb.pos]
	sb.pos++
	return b, nil
}

// Seek adjusts the current position.
func (sb *SeekableBuffer) Seek(offset int) {
	newPos := sb.pos + int(offset)
	sb.pos = newPos
}

// Write writes data to the buffer at the current position.
func (sb *SeekableBuffer) Write(p []byte) (int, error) {
	if sb.pos >= len(sb.data) {
		sb.data = append(sb.data, p...)
	} else {
		sb.data = append(sb.data[:sb.pos], append(p, sb.data[sb.pos:]...)...)
	}

	sb.pos += len(p)
	return len(p), nil
}

// Bytes returns the contents of the buffer.
func (sb *SeekableBuffer) Bytes() []byte {
	return sb.data
}

// Len returns the length of the buffer.
func (sb *SeekableBuffer) Len() int {
	return len(sb.data)
}

// Get returns the byte at the given position in the buffer.
func (sb *SeekableBuffer) Get(pos int) byte {
	return sb.data[pos]
}

// SetPosition sets the current position in the buffer.
func (sb *SeekableBuffer) SetPosition(pos int) {
	sb.pos = pos
}
