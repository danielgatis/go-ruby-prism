package parser

import (
	"math/big"

	"github.com/rotisserie/eris"
)

func loadVarUInt(buff *buffer) (uint32, error) {
	result := uint32(0)
	shift := 0

	for {
		byte, err := buff.readByte()
		if err != nil {
			return 0, eris.Wrap(err, "error reading byte")
		}

		result += uint32(byte&0x7F) << shift
		shift += 7

		if (byte & 0x80) == 0 {
			break
		}
	}

	return result, nil
}

func loadVarSInt(buff *buffer) (int32, error) {
	x, err := loadVarUInt(buff)
	if err != nil {
		return 0, eris.Wrap(err, "error reading VarUInt")
	}

	return int32((x >> 1) ^ (-(x & 1))), nil
}

func loadLineOffsets(buff *buffer) ([]uint32, error) {
	count, err := loadVarUInt(buff)
	if err != nil {
		return nil, eris.Wrap(err, "error reading line offsets count")
	}

	lineOffsets := make([]uint32, count)
	for i := range count {
		lineOffsets[i], err = loadVarUInt(buff)
		if err != nil {
			return nil, eris.Wrap(err, "error reading line offset")
		}
	}

	return lineOffsets, nil
}

func loadLocation(buff *buffer) (*Location, error) {
	startOffset, err := loadVarUInt(buff)
	if err != nil {
		return nil, eris.Wrap(err, "error reading location startOffset")
	}

	length, err := loadVarUInt(buff)
	if err != nil {
		return nil, eris.Wrap(err, "error reading location length")
	}

	return NewLocation(startOffset, length), nil
}

func loadComments(buff *buffer) ([]*Comment, error) {
	count, err := loadVarUInt(buff)
	if err != nil {
		return nil, eris.Wrap(err, "error reading comments count")
	}

	comments := make([]*Comment, count)

	for i := range count {
		typpe, err := loadVarUInt(buff)
		if err != nil {
			return nil, eris.Wrap(err, "error reading comment type")
		}

		loc, err := loadLocation(buff)
		if err != nil {
			return nil, eris.Wrap(err, "error reading comment value location")
		}

		comment := NewComment(typpe, loc)
		comments[i] = comment
	}

	return comments, nil
}

func loadMagicComments(buff *buffer) ([]*MagicComment, error) {
	count, err := loadVarUInt(buff)
	if err != nil {
		return nil, eris.Wrap(err, "error reading magic comments count")
	}

	comments := make([]*MagicComment, count)

	for i := range count {
		keyLocation, err := loadLocation(buff)
		if err != nil {
			return nil, eris.Wrap(err, "error reading magic comment key location")
		}

		valueLocation, err := loadLocation(buff)
		if err != nil {
			return nil, eris.Wrap(err, "error reading magic comment value location")
		}

		comment := NewMagicComment(keyLocation, valueLocation)
		comments[i] = comment
	}

	return comments, nil
}

func loadOptionalLocation(buff *buffer) (*Location, error) {
	nextByte, err := buff.readByte()
	if err != nil {
		return nil, eris.Wrap(err, "error reading optional location")
	}

	if nextByte == 0 {
		return nil, nil
	}

	location, err := loadLocation(buff)
	if err != nil {
		return nil, eris.Wrap(err, "error reading location")
	}

	return location, nil
}

func loadSynErrors(buff *buffer) ([]*SyntaxError, error) {
	count, err := loadVarUInt(buff)
	if err != nil {
		return nil, eris.Wrap(err, "error reading errors count")
	}

	synErrs := make([]*SyntaxError, count)
	for i := range count {
		errorType, err := buff.readByte()
		if err != nil {
			return nil, eris.Wrap(err, "error reading error type")
		}

		messageBytes, err := loadEmbeddedStr(buff)
		if err != nil {
			return nil, eris.Wrap(err, "error reading error message")
		}

		location, err := loadLocation(buff)
		if err != nil {
			return nil, eris.Wrap(err, "error reading error location")
		}

		errorLevel, err := buff.readByte()
		if err != nil {
			return nil, eris.Wrap(err, "error reading error level")
		}

		synErrs[i] = NewSyntaxError(
			string(messageBytes),
			location,
			SyntaxErrorLevel(errorLevel),
			SyntaxErrorTypes[(errorType&0xFF)],
		)
	}

	return synErrs, nil
}

func loadSynWarnings(buff *buffer) ([]*SyntaxWarning, error) {
	count, err := loadVarUInt(buff)
	if err != nil {
		return nil, eris.Wrap(err, "error reading warnings count")
	}

	synWarnings := make([]*SyntaxWarning, count)
	for i := range count {
		warningType, err := buff.readByte()
		if err != nil {
			return nil, eris.Wrap(err, "error reading warning type")
		}

		messageBytes, err := loadEmbeddedStr(buff)
		if err != nil {
			return nil, eris.Wrap(err, "error reading warning message")
		}

		location, err := loadLocation(buff)
		if err != nil {
			return nil, eris.Wrap(err, "error reading warning location")
		}

		warningLevel, err := buff.readByte()
		if err != nil {
			return nil, eris.Wrap(err, "error reading warning level")
		}

		synWarnings[i] = NewSyntaxWarning(
			string(messageBytes),
			location,
			SyntaxWarningLevel(warningLevel),
			SyntaxWarningTypes[(warningType&0xFF)-224],
		)
	}

	return synWarnings, nil
}

func loadEmbeddedStr(buff *buffer) ([]byte, error) {
	length, err := loadVarUInt(buff)
	if err != nil {
		return nil, eris.Wrap(err, "error reading embedded string length")
	}

	strBytes := make([]byte, length)
	_, err = buff.read(strBytes)
	if err != nil {
		return nil, eris.Wrap(err, "error reading embedded string")
	}

	return strBytes, nil
}

func loadStr(buff *buffer, src []byte) ([]byte, error) {
	b, err := buff.readByte()
	if err != nil {
		return nil, eris.Wrap(err, "error reading string type")
	}

	switch b {
	case 1:
		start, err := loadVarUInt(buff)
		if err != nil {
			return nil, eris.Wrap(err, "error reading string start")
		}

		length, err := loadVarUInt(buff)
		if err != nil {
			return nil, eris.Wrap(err, "error reading string length")
		}

		bytes := make([]byte, length)
		copy(bytes, src[start:start+length])
		return bytes, nil
	case 2:
		strBytes, err := loadEmbeddedStr(buff)
		if err != nil {
			return nil, eris.Wrap(err, "error reading embedded string")
		}

		return strBytes, nil
	default:
		return nil, eris.Wrapf(err, "invalid string type: %d", b)
	}
}

func loadOptionalNode(buff *buffer, src []byte, pool *constantPool) (Node, error) {
	nextByte, err := buff.readByte()
	if err != nil {
		return nil, eris.Wrap(err, "error reading optional node")
	}

	if nextByte == 0 {
		return nil, nil
	}

	buff.setPosition(buff.position() - 1)
	node, err := loadNode(buff, src, pool)
	if err != nil {
		return nil, eris.Wrap(err, "error reading node")
	}

	return node, nil
}

func loadInteger(buff *buffer) (*big.Int, error) {
	negative, err := buff.readByte()
	if err != nil {
		return nil, eris.Wrap(err, "error reading integer sign")
	}

	length, err := loadVarUInt(buff)
	if err != nil {
		return nil, eris.Wrap(err, "error reading integer words length")
	}

	if length == 0 {
		return nil, eris.Wrapf(err, "invalid integer words length: %d", length)
	}

	firstWord, err := loadVarUInt(buff)
	if err != nil {
		return nil, eris.Wrap(err, "error reading integer first word")
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
			return nil, eris.Wrap(err, "error reading word")
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

func loadConstant(buff *buffer, pool *constantPool) (string, error) {
	idx, err := loadVarUInt(buff)
	if err != nil {
		return "", eris.Wrap(err, "error reading constant index")
	}

	constant, err := pool.Get(buff, idx)
	if err != nil {
		return "", eris.Wrap(err, "error getting constant")
	}

	return constant, nil
}

func loadOptionalConstant(buff *buffer, pool *constantPool) (*string, error) {
	nextByte, err := buff.readByte()
	if err != nil {
		return nil, eris.Wrap(err, "error reading optional node")
	}

	if nextByte == 0 {
		return nil, nil
	}

	buff.setPosition(buff.position() - 1)
	constant, err := loadConstant(buff, pool)
	if err != nil {
		return nil, eris.Wrap(err, "error reading node")
	}

	return &constant, nil
}

func loadConstants(buff *buffer, pool *constantPool) ([]string, error) {
	count, err := loadVarUInt(buff)
	if err != nil {
		return nil, eris.Wrap(err, "error reading constants count")
	}

	constants := make([]string, count)
	for i := range count {
		constants[i], err = loadConstant(buff, pool)
		if err != nil {
			return nil, eris.Wrap(err, "error reading constant")
		}
	}

	return constants, nil
}

func loadFlags(buff *buffer) (int16, error) {
	flags, err := loadVarUInt(buff)
	if err != nil {
		return 0, eris.Wrap(err, "error reading flags")
	}

	if flags > 0x7FFF {
		return 0, eris.Wrapf(err, "invalid flags: %d", flags)
	}

	return int16(flags), nil
}
