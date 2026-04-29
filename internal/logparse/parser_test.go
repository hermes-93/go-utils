package logparse_test

import (
	"bufio"
	"bytes"
	"strings"
	"testing"

	"github.com/hermes-93/go-utils/internal/logparse"
)

func parse(input, minLevel, format string, fields []string) string {
	scanner := bufio.NewScanner(strings.NewReader(input))
	var out bytes.Buffer
	p := logparse.Parser{MinLevel: minLevel, Format: format, Fields: fields}
	p.Process(scanner, &out)
	return out.String()
}

func TestProcessTextFormat(t *testing.T) {
	input := `{"time":"2024-01-01T10:00:00Z","level":"info","message":"hello world"}`
	result := parse(input, "", "text", nil)
	if !strings.Contains(result, "hello world") {
		t.Fatalf("expected message in output, got: %q", result)
	}
	if !strings.Contains(result, "[INFO]") {
		t.Fatalf("expected level in output, got: %q", result)
	}
}

func TestProcessJSONFormat(t *testing.T) {
	input := `{"level":"info","message":"hello"}`
	result := parse(input, "", "json", nil)
	if !strings.Contains(result, `"message"`) {
		t.Fatalf("expected JSON output, got: %q", result)
	}
}

func TestLevelFilterDebugExcluded(t *testing.T) {
	input := `{"level":"debug","message":"debug msg"}
{"level":"error","message":"error msg"}`
	result := parse(input, "error", "text", nil)
	if strings.Contains(result, "debug msg") {
		t.Fatal("debug message should be filtered when minLevel=error")
	}
	if !strings.Contains(result, "error msg") {
		t.Fatalf("error message should appear, got: %q", result)
	}
}

func TestLevelFilterInfoIncluded(t *testing.T) {
	input := `{"level":"info","message":"info msg"}
{"level":"warn","message":"warn msg"}`
	result := parse(input, "info", "text", nil)
	if !strings.Contains(result, "info msg") {
		t.Fatalf("info msg should appear, got: %q", result)
	}
	if !strings.Contains(result, "warn msg") {
		t.Fatalf("warn msg should appear, got: %q", result)
	}
}

func TestRawLinePassthrough(t *testing.T) {
	input := `not json at all`
	result := parse(input, "", "text", nil)
	if !strings.Contains(result, "not json at all") {
		t.Fatalf("raw line should pass through, got: %q", result)
	}
}

func TestExtraFields(t *testing.T) {
	input := `{"level":"info","message":"request","service":"api","status":"200"}`
	result := parse(input, "", "text", []string{"service", "status"})
	if !strings.Contains(result, "service=api") {
		t.Fatalf("expected service=api in output, got: %q", result)
	}
	if !strings.Contains(result, "status=200") {
		t.Fatalf("expected status=200 in output, got: %q", result)
	}
}

func TestEmptyLinesSkipped(t *testing.T) {
	input := "\n\n{\"level\":\"info\",\"message\":\"msg\"}\n\n"
	result := parse(input, "", "text", nil)
	lines := strings.Split(strings.TrimSpace(result), "\n")
	if len(lines) != 1 {
		t.Fatalf("expected 1 output line, got %d: %q", len(lines), result)
	}
}
