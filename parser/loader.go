package parser

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"strings"
)

const prismHeader = "PRISM"
const majorVersion = 1
const minorVersion = 3
const patchVersion = 0

func load(serialized []byte, source []byte) (*ParseResult, error) {
	buff := NewSeekableBuffer(serialized)
	src := NewSource(source)

	// check header
	header := make([]byte, 5)
	_, err := buff.Read(header)
	if err != nil {
		return nil, fmt.Errorf("error reading header: %w", err)
	}

	if !bytes.Equal(header, []byte(prismHeader)) {
		return nil, fmt.Errorf("invalid prism header")
	}

	// check version
	version := make([]byte, 3)
	_, err = buff.Read(version)
	if err != nil {
		return nil, fmt.Errorf("error reading version: %w", err)
	}

	if !bytes.Equal(version, []byte{majorVersion, minorVersion, patchVersion}) {
		return nil, fmt.Errorf("invalid version number")
	}

	// check location
	var location = make([]byte, 1)
	_, err = buff.Read(location)
	if err != nil {
		return nil, fmt.Errorf("error reading location: %w", err)
	}

	if !bytes.Equal(location, []byte{0}) {
		return nil, fmt.Errorf("requires no location fields in the serialized output")
	}

	// load the encoding
	encodingLength, err := loadVarUInt(buff)
	if err != nil {
		return nil, fmt.Errorf("error reading encoding length: %w", err)
	}

	encodingNameBytes := make([]byte, encodingLength)
	if _, err := buff.Read(encodingNameBytes); err != nil {
		return nil, fmt.Errorf("error reading encoding name: %w", err)
	}

	encodingName := string(encodingNameBytes)
	encodingCharset := getEncodingCharset(encodingName)

	// load the source start line
	startLine, err := loadVarSInt(buff)
	if err != nil {
		return nil, fmt.Errorf("error reading start line: %w", err)
	}

	src.SetStartLine(startLine)

	// load the source line offsets
	lineOffsets, err := loadLineOffsets(buff)
	if err != nil {
		return nil, fmt.Errorf("error reading line offsets: %w", err)
	}

	src.SetLineOffsets(lineOffsets)

	// load the comments
	comments, err := loadComments(buff)
	if err != nil {
		return nil, fmt.Errorf("error reading comments: %w", err)
	}

	// load the magic comments
	magicComments, err := loadMagicComments(buff)
	if err != nil {
		return nil, fmt.Errorf("error reading magic comments: %w", err)
	}

	// load the data location
	dataLocation, err := loadOptionalLocation(buff)
	if err != nil {
		return nil, fmt.Errorf("error reading data location: %w", err)
	}

	// load errors
	errors, err := loadErrors(buff)
	if err != nil {
		return nil, fmt.Errorf("error reading errors: %w", err)
	}

	// load warnings
	warnings, err := loadWarnings(buff)
	if err != nil {
		return nil, fmt.Errorf("error reading warnings: %w", err)
	}

	// load constant pool
	constantPool, err := loadConstantPool(buff, source, encodingCharset)
	if err != nil {
		return nil, fmt.Errorf("error reading constant pool: %w", err)
	}

	// load nodes
	var node Node

	if len(errors) == 0 {
		n, err := loadNode(buff, source, constantPool)
		if err != nil {
			return nil, fmt.Errorf("error reading node: %w", err)
		}
		node = n

		left := constantPool.BufferOffset - buff.Position()
		if left != 0 {
			return nil, fmt.Errorf(
				"expected to consume all bytes while deserializing but there were %d bytes left",
				left,
			)
		}

		newlineMarked := make([]bool, 1+src.LineCount())
		visitor := NewMarkNewlinesVisitor(src, newlineMarked)
		node.Accept(visitor)
	}

	return NewParseResult(node, comments, magicComments, dataLocation, errors, warnings, src), nil
}

func getEncodingCharset(encodingName string) string {
	encodingName = strings.ToLower(encodingName)
	if encodingName == "ascii-8bit" {
		return "US-ASCII"
	}

	return encodingName
}

func bytesToName(bytes []byte, encodingCharset string) (string, error) {
	// Convert bytes to a string using the specified encoding
	str, err := decodeBytesToString(bytes, encodingCharset)
	if err != nil {
		return "", fmt.Errorf("failed to decode bytes: %w", err)
	}

	// Return the interned string (Go does not have direct intern support like Java)
	return str, nil
}

// decodeBytesToString decodes bytes into a string using the specified encoding
func decodeBytesToString(bytes []byte, encodingCharset string) (string, error) {
	// Example for UTF-8. For other encodings, you may need to use external libraries.
	if encodingCharset == "UTF-8" || encodingCharset == "utf-8" {
		return string(bytes), nil
	}
	return "", fmt.Errorf("unsupported encoding: %s", encodingCharset)
}

// variable-length encoding
func loadVarUInt(buff *SeekableBuffer) (int, error) {
	result := uint32(0)
	shift := 0

	for {
		byte, err := buff.ReadByte()
		if err != nil {
			return 0, fmt.Errorf("failed to read byte: %w", err)
		}

		result += uint32(byte&0x7F) << shift
		shift += 7

		if (byte & 0x80) == 0 {
			break
		}
	}

	return int(result), nil
}

func loadVarSInt(buff *SeekableBuffer) (int, error) {
	// Decode the unsigned varint first
	x, err := loadVarUInt(buff)
	if err != nil {
		return 0, err
	}

	// Perform the ZigZag decoding
	return (x >> 1) ^ -(x & 1), nil
}

func loadInt(buff *SeekableBuffer) (int, error) {
	data := make([]byte, 4)
	_, err := buff.Read(data)
	if err != nil {
		return 0, fmt.Errorf("failed to read int: %w", err)
	}

	value := int32(data[0])<<24 | int32(data[1])<<16 | int32(data[2])<<8 | int32(data[3])
	return int(value), nil
}

func loadLineOffsets(buff *SeekableBuffer) ([]int, error) {
	// Decode the count of line offsets
	count, err := loadVarUInt(buff)
	if err != nil {
		return nil, fmt.Errorf("failed to read count: %w", err)
	}

	// Initialize a slice to hold the line offsets
	lineOffsets := make([]int, count)

	// Read each line offset
	for i := 0; i < count; i++ {
		lineOffset, err := loadVarUInt(buff)
		if err != nil {
			return nil, fmt.Errorf("failed to read line offset at index %d: %w", i, err)
		}
		lineOffsets[i] = lineOffset
	}

	return lineOffsets, nil
}

func loadLocation(buff *SeekableBuffer) (*Location, error) {
	// Decode the start offset
	startOffset, err := loadVarUInt(buff)
	if err != nil {
		return nil, fmt.Errorf("failed to read the start offset: %w", err)
	}

	// Decode the length
	length, err := loadVarUInt(buff)
	if err != nil {
		return nil, fmt.Errorf("failed to read the length: %w", err)
	}

	return NewLocation(startOffset, length), nil
}

func loadOptionalLocation(buff *SeekableBuffer) (*Location, error) {
	isPresent, err := buff.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("failed to read byte: %w", err)
	}

	if isPresent != 0 {
		location, err := loadLocation(buff)
		if err != nil {
			return nil, fmt.Errorf("failed to read location: %w", err)
		}

		return location, nil
	}

	return nil, nil
}

func loadComments(buff *SeekableBuffer) ([]*Comment, error) {
	count, err := loadVarUInt(buff)
	if err != nil {
		return nil, fmt.Errorf("failed to read count: %w", err)
	}

	comments := make([]*Comment, count)

	for i := range count {
		typpe, err := loadVarUInt(buff)
		if err != nil {
			return nil, fmt.Errorf("failed to read comment type at index %d: %w", i, err)
		}

		loc, err := loadLocation(buff)
		if err != nil {
			return nil, fmt.Errorf("failed to read location at index %d: %w", i, err)
		}

		comment := NewComment(typpe, loc)
		comments[i] = comment
	}

	return comments, nil
}

func loadMagicComments(buff *SeekableBuffer) ([]*MagicComment, error) {
	// Decode the count of magic comments
	count, err := loadVarUInt(buff)
	if err != nil {
		return nil, fmt.Errorf("failed to read count: %w", err)
	}

	magicComments := make([]*MagicComment, count)

	// Read each magic comment
	for i := 0; i < count; i++ {
		keyLocation, err := loadLocation(buff)
		if err != nil {
			return nil, fmt.Errorf("failed to read key location at index %d: %w", i, err)
		}

		valueLocation, err := loadLocation(buff)
		if err != nil {
			return nil, fmt.Errorf("failed to read value location at index %d: %w", i, err)
		}

		// Create a MagicComment and add it to the slice
		magicComments[i] = NewMagicComment(keyLocation, valueLocation)
	}

	return magicComments, nil
}

func loadErrors(buff *SeekableBuffer) ([]*Error, error) {
	count, err := loadVarUInt(buff)
	if err != nil {
		return nil, fmt.Errorf("failed to read count: %w", err)
	}

	errors := make([]*Error, count)

	for i := 0; i < count; i++ {
		errType, err := loadVarUInt(buff)
		if err != nil {
			return nil, fmt.Errorf("failed to read error type at index %d: %w", i, err)
		}

		message, err := loadEmbeddedString(buff)
		if err != nil {
			return nil, fmt.Errorf("failed to read error message at index %d: %w", i, err)
		}

		location, err := loadLocation(buff)
		if err != nil {
			return nil, fmt.Errorf("failed to read location at index %d: %w", i, err)
		}

		level, err := loadInt(buff)
		if err != nil {
			return nil, fmt.Errorf("failed to read error level at index %d: %w", i, err)
		}

		errors[i] = NewError(string(message), location, ErrorLevel(level), ErrorType(errType))
	}

	return errors, nil
}

func loadWarnings(buff *SeekableBuffer) ([]*Warning, error) {
	count, err := loadVarUInt(buff)
	if err != nil {
		return nil, fmt.Errorf("failed to read count: %w", err)
	}

	warnings := make([]*Warning, count)

	for i := 0; i < count; i++ {
		warningType, err := loadVarUInt(buff)
		if err != nil {
			return nil, fmt.Errorf("failed to read warning type at index %d: %w", i, err)
		}

		message, err := loadEmbeddedString(buff)
		if err != nil {
			return nil, fmt.Errorf("failed to read warning message at index %d: %w", i, err)
		}

		location, err := loadLocation(buff)
		if err != nil {
			return nil, fmt.Errorf("failed to read location at index %d: %w", i, err)
		}

		level, err := loadInt(buff)
		if err != nil {
			return nil, fmt.Errorf("failed to read warning level at index %d: %w", i, err)
		}

		warnings[i] = NewWarning(string(message), location, WarningLevel(level), WarningType(warningType-289))
	}

	return warnings, nil
}

func loadConstantPool(buff *SeekableBuffer, source []byte, encodingCharset string) (*ConstantPool, error) {
	constantPoolBufferOffset, err := loadInt(buff)
	if err != nil {
		return nil, fmt.Errorf("failed to read constant pool buff offset: %w", err)
	}

	constantPoolLength, err := loadVarUInt(buff)
	if err != nil {
		return nil, fmt.Errorf("failed to read constant pool length: %w", err)
	}

	return NewConstantPool(source, constantPoolBufferOffset, constantPoolLength, encodingCharset), nil
}

func loadEmbeddedString(buff *SeekableBuffer) ([]byte, error) {
	length, err := loadVarUInt(buff)
	if err != nil {
		return nil, fmt.Errorf("error reading embedded string length: %w", err)
	}

	strBytes := make([]byte, length)
	_, err = buff.Read(strBytes)
	if err != nil {
		return nil, fmt.Errorf("error reading embedded string bytes: %w", err)
	}

	return strBytes, nil
}

func loadString(buff *SeekableBuffer, src []byte) (string, error) {
	b, err := buff.ReadByte()
	if err != nil {
		return "", fmt.Errorf("error reading string type: %w", err)
	}

	switch b {
	case 1:
		start, err := loadVarUInt(buff)
		if err != nil {
			return "", fmt.Errorf("error reading string start: %w", err)
		}

		length, err := loadVarUInt(buff)
		if err != nil {
			return "", fmt.Errorf("error reading string length: %w", err)
		}

		bytes := make([]byte, length)
		copy(bytes, src[start:start+length])
		return string(bytes), nil
	case 2:
		strBytes, err := loadEmbeddedString(buff)
		if err != nil {
			return "", fmt.Errorf("error reading embedded string: %w", err)
		}

		return string(strBytes), nil
	default:
		return "", fmt.Errorf("invalid string type: %d", b)
	}
}

func loadFlags(buff *SeekableBuffer) (int16, error) {
	flags, err := loadVarUInt(buff)
	if err != nil {
		return 0, fmt.Errorf("failed to read flags: %w", err)
	}
	if flags < 0 || flags > 32767 {
		return 0, fmt.Errorf("flags out of range: %d", flags)
	}
	return int16(flags), nil
}

func loadConstant(buff *SeekableBuffer, constantPool *ConstantPool) (string, error) {
	idx, err := loadVarUInt(buff)
	if err != nil {
		return "", fmt.Errorf("failed to read constant index: %w", err)
	}

	constant, err := constantPool.Get(buff, idx)
	if err != nil {
		return "", fmt.Errorf("failed to get constant: %w", err)
	}

	return constant, nil
}

func loadOptionalConstant(buff *SeekableBuffer, constantPool *ConstantPool) (*string, error) {
	if buff.Get(buff.Position()) != 0 {
		str, err := loadConstant(buff, constantPool)
		if err != nil {
			return nil, fmt.Errorf("failed to read constant: %w", err)
		}

		return &str, nil
	} else {
		buff.Seek(1)
		return nil, nil
	}
}

func loadConstants(buff *SeekableBuffer, constantPool *ConstantPool) ([]string, error) {
	length, err := loadVarUInt(buff)

	if err != nil {
		return nil, fmt.Errorf("failed to read constants length: %w", err)
	}

	if length == 0 {
		return []string{}, nil
	}

	constants := make([]string, length)

	for i := 0; i < length; i++ {
		constant, err := loadConstant(buff, constantPool)
		if err != nil {
			return nil, fmt.Errorf("failed to read constant at index %d: %w", i, err)
		}
		constants[i] = constant
	}

	return constants, nil
}

func loadInteger(buff *SeekableBuffer) (*big.Int, error) {
	negative, err := buff.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("error reading integer sign: %w", err)
	}

	length, err := loadVarUInt(buff)
	if err != nil {
		return nil, fmt.Errorf("error reading integer words length: %w", err)
	}

	if length == 0 {
		return nil, fmt.Errorf("invalid integer length: %d", length)
	}

	firstWord, err := loadVarUInt(buff)
	if err != nil {
		return nil, fmt.Errorf("error reading first word: %w", err)
	}

	if length == 1 {
		if negative != 0 {
			return big.NewInt(-int64(firstWord)), nil
		} else {
			return big.NewInt(int64(firstWord)), nil
		}
	}

	result := big.NewInt(int64(firstWord))
	for index := 1; index < int(length); index++ {
		word, err := loadVarUInt(buff)
		if err != nil {
			return nil, fmt.Errorf("error reading word: %w", err)
		}

		temp := big.NewInt(int64(word))
		temp = temp.Lsh(temp, uint(index*32))
		result = result.Or(result, temp)
	}

	if negative != 0 {
		result = result.Neg(result)
	} else {
		result = result.Abs(result)
	}

	return result, nil
}

func loadDouble(buff *SeekableBuffer) (float64, error) {
	if len(buff.Bytes()) == 0 {
		return 0, fmt.Errorf("buffer is empty")
	}

	data := buff.Bytes()
	bytes := data[buff.Position() : buff.Position()+8]
	buff.Seek(8)

	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)

	return float, nil
}

func loadOptionalNode(buff *SeekableBuffer, source []byte, constantPool *ConstantPool) (Node, error) {
	if buff.Get(buff.Position()) != 0 {
		node, err := loadNode(buff, source, constantPool)
		if err != nil {
			return nil, fmt.Errorf("failed to load node: %w", err)

		}

		return node, nil
	} else {
		buff.Seek(1)
		return nil, nil
	}
}
