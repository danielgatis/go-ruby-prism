package parser

import (
	"bytes"
	"encoding/binary"

	"github.com/rotisserie/eris"
)

type parseOptions struct {
	filepath            string
	line                int
	encoding            string
	frozenStringLiteral bool
	commandLine         []int
	version             int
	scopes              [][][]byte
}

func newParseOptions() *parseOptions {
	return &parseOptions{}
}

func (o *parseOptions) bytes() ([]byte, error) {
	buffer := new(bytes.Buffer)

	// filepath
	if _, err := buffer.Write(strToBytes(o.filepath)); err != nil {
		return nil, eris.Wrap(err, "failed to write the length of the filepath")
	}

	// line
	if _, err := buffer.Write(intToBytes(o.line)); err != nil {
		return nil, eris.Wrap(err, "failed to write the line")
	}

	// encoding
	if _, err := buffer.Write(strToBytes(o.encoding)); err != nil {
		return nil, eris.Wrap(err, "failed to write the length of the encoding")
	}

	// frozenStringLiteral
	if _, err := buffer.Write(boolToBytes(o.frozenStringLiteral)); err != nil {
		return nil, eris.Wrap(err, "failed to write the frozenStringLiteral")
	}

	// command line
	if _, err := buffer.Write(bitsetToBytes(o.commandLine)); err != nil {
		return nil, eris.Wrap(err, "failed to write the command line")
	}

	// version
	if err := buffer.WriteByte(byte(o.version)); err != nil {
		return nil, eris.Wrap(err, "failed to write the version")
	}

	// scopes
	if _, err := buffer.Write(b3ToBytes(o.scopes)); err != nil {
		return nil, eris.Wrap(err, "failed to write the length of the scopes")
	}

	return buffer.Bytes(), nil
}

func intToBytes(n int) []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, uint32(n))
	return bytes
}

func boolToBytes(b bool) []byte {
	if b {
		return []byte{1}
	} else {
		return []byte{0}
	}
}

func strToBytes(s string) []byte {
	return b1ToBytes([]byte(s))
}

func b1ToBytes(s []byte) []byte {
	buffer := new(bytes.Buffer)
	buffer.Write(intToBytes(len(s)))
	buffer.Write(s)
	return buffer.Bytes()
}

func b2ToBytes(b2 [][]byte) []byte {
	buffer := new(bytes.Buffer)
	buffer.Write(intToBytes(len(b2)))

	for _, b1 := range b2 {
		buffer.Write(b1ToBytes(b1))
	}

	return buffer.Bytes()
}

func b3ToBytes(b3 [][][]byte) []byte {
	buffer := new(bytes.Buffer)
	buffer.Write(intToBytes(len(b3)))

	for _, b2 := range b3 {
		buffer.Write(b2ToBytes(b2))
	}

	return buffer.Bytes()
}

func bitsetToBytes(bitset []int) []byte {
	var bt byte = 0
	for _, b := range bitset {
		bt |= byte(b)
	}

	return []byte{bt}
}
