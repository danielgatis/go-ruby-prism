package parser

import (
	"bytes"

	"github.com/rotisserie/eris"
)

const prismHeader = "PRISM"
const majorVersion = 0
const minorVersion = 24
const patchVersion = 0

func deserialize(serialized []byte, source []byte) (*ParseResult, error) {
	buff := newBuffer(serialized)

	// check header
	header := make([]byte, 5)
	_, err := buff.read(header)
	if err != nil {
		return nil, eris.Wrap(err, "error reading header")
	}

	if !bytes.Equal(header, []byte(prismHeader)) {
		return nil, eris.New("invalid prism header")
	}

	// check version
	version := make([]byte, 3)
	_, err = buff.read(version)
	if err != nil {
		return nil, eris.Wrap(err, "error reading version")
	}

	if !bytes.Equal(version, []byte{majorVersion, minorVersion, patchVersion}) {
		return nil, eris.New("invalid version number")
	}

	// check location
	var location = make([]byte, 1)
	_, err = buff.read(location)
	if err != nil {
		return nil, eris.Wrap(err, "error reading location")
	}

	if !bytes.Equal(location, []byte{0}) {
		return nil, eris.New("requires no location fields in the serialized output")
	}

	// reading encoding and discard it is always UTF-8
	encodingLen, err := loadVarUInt(buff)
	if err != nil {
		return nil, eris.Wrap(err, "error reading encoding length")
	}

	var encoding = make([]byte, encodingLen)
	_, err = buff.read(encoding)
	if err != nil {
		return nil, eris.Wrap(err, "error reading encoding")
	}

	// skip start line and line offsets
	_, err = loadVarSInt(buff)
	if err != nil {
		return nil, eris.Wrap(err, "error reading start line")
	}

	_, err = loadLineOffsets(buff)
	if err != nil {
		return nil, eris.Wrap(err, "error reading line offsets")
	}

	// load magic comments
	comments, err := loadComments(buff)
	if err != nil {
		return nil, eris.Wrap(err, "error reading comments")
	}

	// load magic comments
	magicComments, err := loadMagicComments(buff)
	if err != nil {
		return nil, eris.Wrap(err, "error reading magic comments")
	}

	// load optional locations
	dataLocation, err := loadOptionalLocation(buff)
	if err != nil {
		return nil, eris.Wrap(err, "error reading data location")
	}

	// load syntax errors
	synErrors, err := loadSynErrors(buff)
	if err != nil {
		return nil, eris.Wrap(err, "error reading syntax errors")
	}

	// load syntax warnings
	synWarnings, err := loadSynWarnings(buff)
	if err != nil {
		return nil, eris.Wrap(err, "error reading syntax warnings")
	}

	// build constant pool
	constantPoolBufferOffset := buff.readUInt32()
	if err != nil {
		return nil, eris.Wrap(err, "error reading constant pool buff offset")
	}

	constantPoolLength, err := loadVarUInt(buff)
	if err != nil {
		return nil, eris.Wrap(err, "error reading constant pool length")
	}

	constantPool := newConstantPool(
		source,
		serialized,
		constantPoolBufferOffset,
		constantPoolLength,
	)

	// load first node
	node, err := loadNode(buff, source, constantPool)
	if err != nil {
		return nil, eris.Wrap(err, "error reading first node")
	}

	// build parse result
	return NewParseResult(
		node,
		comments,
		magicComments,
		dataLocation,
		synErrors,
		synWarnings,
	), nil
}
