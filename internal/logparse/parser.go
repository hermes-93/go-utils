package logparse

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

var levelOrder = map[string]int{
	"debug":   0,
	"info":    1,
	"warn":    2,
	"warning": 2,
	"error":   3,
	"fatal":   4,
}

type Parser struct {
	MinLevel string
	Fields   []string
	Format   string
}

func (p *Parser) Process(scanner *bufio.Scanner, out io.Writer) {
	minOrder := 0
	if p.MinLevel != "" {
		if ord, ok := levelOrder[strings.ToLower(p.MinLevel)]; ok {
			minOrder = ord
		}
	}

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		raw := tryParseJSON(line)
		if raw == nil {
			raw = tryParseLogfmt(line)
		}
		if raw == nil {
			fmt.Fprintf(out, "[raw] %s\n", line)
			continue
		}

		lvl := strings.ToLower(strField(raw, "level", "severity", "msg_level"))
		if lvl != "" {
			if ord, ok := levelOrder[lvl]; ok && ord < minOrder {
				continue
			}
		}

		if p.Format == "json" {
			data, _ := json.Marshal(raw)
			fmt.Fprintf(out, "%s\n", data)
		} else {
			p.printText(out, raw)
		}
	}
}

func tryParseJSON(line string) map[string]interface{} {
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(line), &raw); err != nil {
		return nil
	}
	return raw
}

// tryParseLogfmt parses lines in key=value or key="value with spaces" format.
func tryParseLogfmt(line string) map[string]interface{} {
	if !strings.Contains(line, "=") {
		return nil
	}
	raw := make(map[string]interface{})
	i := 0
	for i < len(line) {
		for i < len(line) && line[i] == ' ' {
			i++
		}
		start := i
		for i < len(line) && line[i] != '=' && line[i] != ' ' {
			i++
		}
		key := line[start:i]
		if key == "" || i >= len(line) || line[i] != '=' {
			return nil
		}
		i++

		var value string
		if i < len(line) && line[i] == '"' {
			i++
			start = i
			for i < len(line) && line[i] != '"' {
				if line[i] == '\\' {
					i++
				}
				i++
			}
			value = line[start:i]
			if i < len(line) {
				i++
			}
		} else {
			start = i
			for i < len(line) && line[i] != ' ' {
				i++
			}
			value = line[start:i]
		}
		raw[key] = value
	}
	if len(raw) == 0 {
		return nil
	}
	return raw
}

func (p *Parser) printText(out io.Writer, raw map[string]interface{}) {
	ts := strField(raw, "time", "timestamp", "ts", "@timestamp")
	if ts != "" {
		if t, err := time.Parse(time.RFC3339, ts); err == nil {
			ts = t.Format("15:04:05")
		}
	}
	lvl := strings.ToUpper(strField(raw, "level", "severity", "msg_level"))
	msg := strField(raw, "message", "msg", "log")

	var parts []string
	if ts != "" {
		parts = append(parts, ts)
	}
	if lvl != "" {
		parts = append(parts, "["+lvl+"]")
	}
	if msg != "" {
		parts = append(parts, msg)
	}
	for _, f := range p.Fields {
		if v := strField(raw, f); v != "" {
			parts = append(parts, f+"="+v)
		}
	}
	fmt.Fprintf(out, "%s\n", strings.Join(parts, " "))
}

func strField(raw map[string]interface{}, keys ...string) string {
	for _, k := range keys {
		if v, ok := raw[k]; ok {
			return fmt.Sprintf("%v", v)
		}
	}
	return ""
}
