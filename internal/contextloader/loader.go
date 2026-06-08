package contextloader

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const defaultRadius = 3

func EnrichSnippet(file string, line int, existing string) string {
	if strings.TrimSpace(existing) != "" {
		return existing
	}
	return LoadAround(file, line, defaultRadius)
}

func LoadAround(file string, line int, radius int) string {
	if line <= 0 {
		return ""
	}

	f, err := os.Open(file)
	if err != nil {
		return ""
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil || len(lines) == 0 {
		return ""
	}

	start := line - radius - 1
	if start < 0 {
		start = 0
	}
	end := line + radius
	if end > len(lines) {
		end = len(lines)
	}

	var b strings.Builder
	for i := start; i < end; i++ {
		b.WriteString(fmt.Sprintf("%4d | %s\n", i+1, lines[i]))
	}
	return b.String()
}
