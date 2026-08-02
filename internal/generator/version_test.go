package generator

import "testing"

// TestReadVersion checks that the prism version header is parsed, since the
// deserializer rejects any tree whose version does not match.
func TestReadVersion(t *testing.T) {
	version, err := ReadVersion("../../prism/include/prism/version.h")
	if err != nil {
		t.Fatalf("read version: %v", err)
	}

	if version.Major == 0 && version.Minor == 0 && version.Patch == 0 {
		t.Fatal("parsed an all-zero version")
	}
}
