package parser

import (
	"encoding/binary"
	"math"
)

type buffer struct {
	data  []byte
	index int
}

func newBuffer(data []byte) *buffer {
	return &buffer{
		data: data,
	}
}

func (b *buffer) readByte() (byte, error) {
	if len(b.data) == 0 {
		return 0, nil
	}

	b.index++
	return b.data[b.index-1], nil
}

func (b *buffer) read(p []byte) (int, error) {
	if len(b.data) == 0 {
		return 0, nil
	}

	n := copy(p, b.data[b.index:])
	b.index += n
	return n, nil
}

func (b *buffer) readUInt32() uint32 {
	if len(b.data) == 0 {
		return 0
	}

	bytes := b.data[b.index : b.index+4]
	b.index += 4
	return binary.LittleEndian.Uint32(bytes)
}

func (b *buffer) readFloat64() float64 {
	if len(b.data) == 0 {
		return 0
	}

	bytes := b.data[b.index : b.index+8]
	b.index += 8

	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}

func (b *buffer) position() int {
	return b.index
}

func (b *buffer) setPosition(index int) {
	b.index = index
}
