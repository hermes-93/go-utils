package env

import (
	"bufio"
	"os"
	"strings"
)

func Check(vars []string) []string {
	var missing []string
	for _, v := range vars {
		if os.Getenv(v) == "" {
			missing = append(missing, v)
		}
	}
	return missing
}

func ParseFile(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var vars []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key := line
		if idx := strings.IndexByte(line, '='); idx >= 0 {
			key = strings.TrimSpace(line[:idx])
		}
		if key != "" {
			vars = append(vars, key)
		}
	}
	return vars, scanner.Err()
}
