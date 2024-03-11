package parser

import (
	"encoding/binary"

	"github.com/rotisserie/eris"
)

type constantPool struct {
	source       []byte
	serialized   []byte
	bufferOffset uint32
	cache        []*string
}

func newConstantPool(source []byte, serialized []byte, bufferOffset, length uint32) *constantPool {
	return &constantPool{
		source:       source,
		serialized:   serialized,
		bufferOffset: bufferOffset,
		cache:        make([]*string, length),
	}
}

func (cp *constantPool) Get(buff *buffer, oneBasedIndex uint32) (string, error) {
	index := oneBasedIndex - 1
	constant := cp.cache[index]

	if constant == nil {
		offset := cp.bufferOffset + index*8
		start := binary.LittleEndian.Uint32(cp.serialized[offset : offset+4])
		length := binary.LittleEndian.Uint32(cp.serialized[offset+4 : offset+8])
		bytes := make([]byte, length)

		if start <= 0x7FFFFFFF {
			copy(bytes, cp.source[start:start+length])
		} else {
			position := buff.position()
			buff.setPosition(int(start & 0x7FFFFFFF))
			_, err := buff.read(bytes)
			if err != nil {
				return "", eris.Wrapf(err, "error reading bytes")
			}
			buff.setPosition(position)
		}

		var str = string(bytes)
		cp.cache[index] = &str
	}

	return *cp.cache[index], nil
}
