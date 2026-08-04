// Copyright (C) 2025 T-Force I/O
// SPDX-License-Identifier: MIT

package dotenv

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type envEntry struct {
	key   string
	value string
}

// Parse .env file and set these environment variables for current process.
func LoadDotEnv(path string, override bool) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open %s: %w", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0
	var entries []envEntry

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, err := parseEnvLine(line)
		if err != nil {
			return fmt.Errorf("error parsing %s at line %d: %w", path, lineNum, err)
		}

		if key != "" {
			entries = append(entries, envEntry{key, value})
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading %s: %w", path, err)
	}

	for _, e := range entries {
		if override || os.Getenv(e.key) == "" {
			os.Setenv(e.key, e.value)
		}
	}

	return nil
}

// Parse a single KEY=VALUE line from a .env file.
func parseEnvLine(line string) (string, string, error) {
	idx := strings.Index(line, "=")
	if idx == -1 {
		return "", "", nil
	}

	key := strings.TrimSpace(line[:idx])
	value := strings.TrimSpace(line[idx+1:])

	if key == "" {
		return "", "", fmt.Errorf("empty key")
	}

	value = stripQuotes(value)

	return key, value, nil
}

// Remove surrounding quotes from a value if present.
func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
