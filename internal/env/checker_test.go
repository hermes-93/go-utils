package env_test

import (
	"os"
	"testing"

	"github.com/hermes-93/go-utils/internal/env"
)

func TestCheckAllSet(t *testing.T) {
	os.Setenv("GO_UTILS_TEST_A", "val1")
	os.Setenv("GO_UTILS_TEST_B", "val2")
	defer os.Unsetenv("GO_UTILS_TEST_A")
	defer os.Unsetenv("GO_UTILS_TEST_B")

	missing := env.Check([]string{"GO_UTILS_TEST_A", "GO_UTILS_TEST_B"})
	if len(missing) != 0 {
		t.Fatalf("expected no missing vars, got: %v", missing)
	}
}

func TestCheckMissingVar(t *testing.T) {
	os.Unsetenv("GO_UTILS_DEFINITELY_NOT_SET_XYZ")
	missing := env.Check([]string{"GO_UTILS_DEFINITELY_NOT_SET_XYZ"})
	if len(missing) != 1 || missing[0] != "GO_UTILS_DEFINITELY_NOT_SET_XYZ" {
		t.Fatalf("expected one missing var, got: %v", missing)
	}
}

func TestCheckEmpty(t *testing.T) {
	os.Setenv("GO_UTILS_EMPTY_VAR", "")
	defer os.Unsetenv("GO_UTILS_EMPTY_VAR")

	missing := env.Check([]string{"GO_UTILS_EMPTY_VAR"})
	if len(missing) != 1 {
		t.Fatalf("empty string env var should be treated as missing, got: %v", missing)
	}
}

func TestParseFile(t *testing.T) {
	f, err := os.CreateTemp("", "env-test-*.env")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	f.WriteString("# comment\n\nDB_HOST=\nDB_PORT=5432\nSECRET_KEY\n  SPACED  \n")
	f.Close()

	vars, err := env.ParseFile(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	expected := []string{"DB_HOST", "DB_PORT", "SECRET_KEY", "SPACED"}
	if len(vars) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, vars)
	}
	for i, v := range vars {
		if v != expected[i] {
			t.Fatalf("index %d: expected %q, got %q", i, expected[i], v)
		}
	}
}

func TestParseFileMissing(t *testing.T) {
	_, err := env.ParseFile("/nonexistent/path/to/file.env")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
