// Package generator reads prism's config.yml and renders the generated Go
// sources for the parser package. It replaces the Ruby/ERB pipeline that
// prism ships in templates/template.rb, so that regenerating this project
// only requires a Go toolchain.
package generator

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// COMMON_FLAGS_COUNT mirrors Prism::Template::COMMON_FLAGS_COUNT. The first
// two flag bits are reserved for flags common to every node, so per-node flag
// values start at this offset.
const CommonFlagsCount = 2

// removeOnErrorTypes mirrors Prism::Template::REMOVE_ON_ERROR_TYPES, which is
// derived from PRISM_SERIALIZE_ONLY_SEMANTICS_FIELDS. That environment
// variable is not set for this project's generation, so "on error" kinds are
// retained.
const removeOnErrorTypes = false

// Config is the root of config.yml.
type Config struct {
	Errors   []NamedEntry `yaml:"errors"`
	Warnings []NamedEntry `yaml:"warnings"`
	Tokens   []NamedEntry `yaml:"tokens"`
	Flags    []Flags      `yaml:"flags"`
	Nodes    []Node       `yaml:"nodes"`
}

// NamedEntry is an errors/warnings/tokens entry. These are written in
// config.yml either as a bare string or as a mapping with a name key.
type NamedEntry struct {
	Name    string
	Comment string
}

// UnmarshalYAML accepts both the scalar and mapping spellings.
func (n *NamedEntry) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.ScalarNode {
		return value.Decode(&n.Name)
	}

	var entry struct {
		Name    string `yaml:"name"`
		Comment string `yaml:"comment"`
	}

	if err := value.Decode(&entry); err != nil {
		return err
	}

	n.Name = entry.Name
	n.Comment = entry.Comment

	return nil
}

// Flags is a named set of flag bits.
type Flags struct {
	Name    string `yaml:"name"`
	Values  []Flag `yaml:"values"`
	Comment string `yaml:"comment"`
}

// Flag is one bit within a Flags set.
type Flag struct {
	Name    string `yaml:"name"`
	Comment string `yaml:"comment"`
}

// Node is a single node type in the syntax tree.
type Node struct {
	Name    string  `yaml:"name"`
	Fields  []Field `yaml:"fields"`
	Flags   string  `yaml:"flags"`
	Comment string  `yaml:"comment"`

	// FlagSet is resolved from Flags after the whole config is loaded.
	FlagSet *Flags `yaml:"-"`
}

// Field is one field on a node.
type Field struct {
	Name    string `yaml:"name"`
	Type    string `yaml:"type"`
	Comment string `yaml:"comment"`

	// Kind holds the resolved node kind for node-typed fields. It is empty
	// when the field refers to a generic node, which happens both when no
	// kind is given and when the kind is a union of several types.
	Kind string `yaml:"-"`

	// RawKind is the unprocessed kind entry from config.yml, which may be a
	// string or a list mixing strings and "on error" mappings.
	RawKind yaml.Node `yaml:"kind"`
}

// Load reads and resolves a config.yml.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	flagsByName := make(map[string]*Flags, len(config.Flags))
	for i := range config.Flags {
		flagsByName[config.Flags[i].Name] = &config.Flags[i]
	}

	for i := range config.Nodes {
		node := &config.Nodes[i]

		if node.Flags != "" {
			flagSet, ok := flagsByName[node.Flags]
			if !ok {
				return nil, fmt.Errorf("node %s references unknown flags %q", node.Name, node.Flags)
			}

			node.FlagSet = flagSet
		}

		for j := range node.Fields {
			field := &node.Fields[j]

			kind, err := resolveKind(&field.RawKind)
			if err != nil {
				return nil, fmt.Errorf("node %s field %s: %w", node.Name, field.Name, err)
			}

			field.Kind = kind
		}
	}

	return config, nil
}

// resolveKind reduces a field's kind entry to a single specific node type, or
// to the empty string when the field should be typed as a generic node.
//
// This mirrors the normalization in prism's template.rb: the kind is coerced
// into a list, the two placeholder kinds collapse to a generic node, "on
// error" mappings are unwrapped, and only a list that reduces to exactly one
// concrete type yields a specific kind.
func resolveKind(raw *yaml.Node) (string, error) {
	if raw == nil || raw.Kind == 0 {
		return "", nil
	}

	var entries []*yaml.Node

	switch raw.Kind {
	case yaml.SequenceNode:
		entries = raw.Content
	default:
		entries = []*yaml.Node{raw}
	}

	kinds := make([]string, 0, len(entries))

	for _, entry := range entries {
		switch entry.Kind {
		case yaml.ScalarNode:
			switch entry.Value {
			case "non-void expression", "pattern expression":
				// The full list of types is impractically long, so these
				// are represented as a generic node.
				kinds = append(kinds, "Node")
			default:
				kinds = append(kinds, entry.Value)
			}

		case yaml.MappingNode:
			var onError struct {
				Kind string `yaml:"on error"`
			}

			if err := entry.Decode(&onError); err != nil {
				return "", fmt.Errorf("failed to decode kind: %w", err)
			}

			if !removeOnErrorTypes {
				kinds = append(kinds, onError.Kind)
			}

		default:
			return "", fmt.Errorf("unexpected kind node kind %v", entry.Kind)
		}
	}

	// Only a single concrete type produces a specific kind. A union of
	// several types, or the generic placeholder, stays generic.
	if len(kinds) == 1 && kinds[0] != "Node" {
		return kinds[0], nil
	}

	return "", nil
}

// GoCamelCase converts a snake_case config name to CamelCase.
//
// This mirrors the ERB helper exactly:
//
//	string.gsub(/_([a-z])/) { $1.upcase }.gsub(/^([a-z])/) { $1.upcase }
//
// Only an underscore followed by a *lowercase* letter is removed. Names that
// are already uppercase, such as the CONTAINS_FORWARDING flag values, keep
// their underscores.
func GoCamelCase(s string) string {
	var b strings.Builder
	b.Grow(len(s))

	runes := []rune(s)

	for i := 0; i < len(runes); i++ {
		r := runes[i]

		if r == '_' && i+1 < len(runes) && runes[i+1] >= 'a' && runes[i+1] <= 'z' {
			b.WriteRune(runes[i+1] - 32)
			i++

			continue
		}

		b.WriteRune(r)
	}

	out := b.String()
	if out != "" && out[0] >= 'a' && out[0] <= 'z' {
		out = string(out[0]-32) + out[1:]
	}

	return out
}

// GoProp is the exported Go field name for a config field.
func (f Field) GoProp() string {
	if f.Name == "arguments" {
		return "Arguments"
	}

	return GoCamelCase(f.Name)
}

// GoName is the identifier used for the field as a constructor parameter and
// as a JSON key.
func (f Field) GoName() string {
	if f.Name == "arguments" {
		return "arguments"
	}

	return f.Name
}

// GoType is the Go type for a config field.
func (f Field) GoType() string {
	switch f.Type {
	case "node", "node?":
		if f.Kind == "" {
			return "Node"
		}

		return "*" + f.Kind
	case "node[]":
		return "[]Node"
	case "string":
		return "RubyString"
	case "constant":
		return "string"
	case "constant?":
		return "*string"
	case "constant[]":
		return "[]string"
	case "location":
		return "Location"
	case "location?":
		return "*Location"
	case "uint8":
		return "uint8"
	case "uint32":
		return "uint32"
	case "integer":
		return "int64"
	case "double":
		return "float64"
	default:
		panic(fmt.Sprintf("unknown field type: %q", f.Type))
	}
}

// NeedsSerializedLength mirrors Node#needs_serialized_length?. Only DefNode
// emits its serialized length so implementations can skip it for lazy parsing.
func (n Node) NeedsSerializedLength() bool {
	return n.Name == "DefNode"
}

// ReadExpr is the Go expression that deserializes this field inside the
// generated switch in gen_deserialize.go.
func (f Field) ReadExpr() string {
	switch f.Type {
	case "node":
		if f.Kind == "" {
			return "readRequiredNode()"
		}

		return "readRequiredNode().(*" + f.Kind + ")"

	case "node?":
		if f.Kind == "" {
			return "func() Node { if n := readOptionalNode(); n != nil { return n }; return nil }()"
		}

		return "func() *" + f.Kind + " { if n := readOptionalNode(); n != nil { return n.(*" + f.Kind + ") }; return nil }()"

	case "string":
		return "buffer.ReadStringField(flags)"
	case "node[]":
		return "func() []Node { count := buffer.ReadVarInt(); nodes := make([]Node, count); for i := 0; i < count; i++ { nodes[i] = readRequiredNode() }; return nodes }()"
	case "constant":
		return "readRequiredConstant()"
	case "constant?":
		return "readOptionalConstant()"
	case "constant[]":
		return "func() []string { count := buffer.ReadVarInt(); result := make([]string, count); for i := 0; i < count; i++ { result[i] = readRequiredConstant() }; return result }()"
	case "location":
		return "buffer.ReadLocation()"
	case "location?":
		return "buffer.ReadOptionalLocation()"
	case "uint8":
		return "buffer.ReadRawByte()"
	case "uint32":
		return "uint32(buffer.ReadVarInt())"
	case "integer":
		return "readInteger()"
	case "double":
		return "buffer.ReadDouble()"
	default:
		panic(fmt.Sprintf("unknown field type: %q", f.Type))
	}
}

// CommentLines renders a comment the way prism's ConfigComment does: each
// line is prefixed with a space and has trailing whitespace removed.
func CommentLines(comment string) []string {
	if comment == "" {
		return nil
	}

	// Ruby's each_line keeps the trailing newline on every line but the
	// last, and a trailing newline therefore does not produce a final empty
	// line. Splitting on "\n" would, so drop one trailing newline first.
	comment = strings.TrimSuffix(comment, "\n")

	lines := strings.Split(comment, "\n")
	out := make([]string, 0, len(lines))

	for _, line := range lines {
		out = append(out, strings.TrimRight(" "+line, " \t\r\n"))
	}

	return out
}

// CamelCase is the CamelCase form of a flag name, matching Flag#camelcase.
func (f Flag) CamelCase() string {
	parts := strings.Split(f.Name, "_")
	for i, part := range parts {
		if part == "" {
			continue
		}

		parts[i] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
	}

	return strings.Join(parts, "")
}
