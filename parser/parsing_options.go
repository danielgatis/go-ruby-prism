package parser

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type SyntaxVersion byte

const (
	SyntaxVersionLatest SyntaxVersion = iota
	SyntaxVersionV3_3
	SyntaxVersionV3_4
)

type CommandLine int

const (
	CommandA CommandLine = iota
	CommandE
	CommandL
	CommandN
	CommandP
	CommandX
)

func serializeParserOptions(
	filepath []byte,
	line int,
	encoding []byte,
	frozenStringLiteral bool,
	commandLine []CommandLine,
	version SyntaxVersion,
	encodingLocked bool,
	mainScript bool,
	partialScript bool,
	scopes [][][]byte,
) ([]byte, error) {
	output := new(bytes.Buffer)

	// filepath
	if _, err := output.Write(serializeInt(len(filepath))); err != nil {
		return nil, fmt.Errorf("failed to serialize filepath length: %w", err)
	}
	if _, err := output.Write(filepath); err != nil {
		return nil, fmt.Errorf("failed to serialize filepath: %w", err)
	}

	// line
	l := 1
	if line > 0 {
		l = line
	}

	if _, err := output.Write(serializeInt(l)); err != nil {
		return nil, fmt.Errorf("failed to serialize line: %w", err)
	}

	// encoding
	if _, err := output.Write(serializeInt(len(encoding))); err != nil {
		return nil, fmt.Errorf("failed to serialize encoding length: %w", err)
	}
	if _, err := output.Write(encoding); err != nil {
		return nil, fmt.Errorf("failed to serialize encoding: %w", err)
	}

	// frozenStringLiteral
	if frozenStringLiteral {
		output.WriteByte(1)
	} else {
		output.WriteByte(0)
	}

	// command line
	commandLineByte, err := serializeEnumSet(commandLine)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize command line: %w", err)
	}
	output.WriteByte(commandLineByte)

	// version
	output.WriteByte(byte(version))

	// encodingLocked
	if encodingLocked {
		output.WriteByte(1)
	} else {
		output.WriteByte(0)
	}

	// mainScript
	if mainScript {
		output.WriteByte(1)
	} else {
		output.WriteByte(0)
	}

	// partialScript
	if partialScript {
		output.WriteByte(1)
	} else {
		output.WriteByte(0)
	}

	// scopes
	if _, err := output.Write(serializeInt(len(scopes))); err != nil {
		return nil, fmt.Errorf("failed to serialize scopes length: %w", err)
	}
	for _, scope := range scopes {
		if _, err := output.Write(serializeInt(len(scope))); err != nil {
			return nil, fmt.Errorf("failed to serialize scope length: %w", err)
		}
		for _, local := range scope {
			if _, err := output.Write(serializeInt(len(local))); err != nil {
				return nil, fmt.Errorf("failed to serialize local length: %w", err)
			}
			if _, err := output.Write(local); err != nil {
				return nil, fmt.Errorf("failed to serialize local: %w", err)
			}
		}
	}

	return output.Bytes(), nil
}

func serializeEnumSet(set []CommandLine) (byte, error) {
	var result byte
	for _, value := range set {
		if (1 << value) > 127 {
			return 0, fmt.Errorf("value is too large to serialize")
		}
		result |= (1 << value)
	}
	return result, nil
}

func serializeInt(n int) []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, uint32(n))
	return bytes
}
