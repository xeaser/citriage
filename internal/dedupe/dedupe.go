package dedupe

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

var timestampRegex = regexp.MustCompile(`^\[.*?\]\s*`)

// Process removes timestamps and deduplicates consecutive identical log lines.
// Parameters:
//
//	logData []byte: the log data to process.
//
// Returns:
//
//	string: the processed log with deduplicated lines and error messages if any.
func Process(logData []byte) string {
	var builder strings.Builder
	scanner := bufio.NewScanner(bytes.NewReader(logData))

	var previousLine string
	var duplicateCount int = 1

	firstLine := true

	for scanner.Scan() {
		currentLine := timestampRegex.ReplaceAllString(scanner.Text(), "")

		if firstLine {
			previousLine = currentLine
			firstLine = false
			continue
		}

		if currentLine == previousLine {
			duplicateCount++
		} else {
			flush(duplicateCount, previousLine, &builder)
			previousLine = currentLine
			duplicateCount = 1
		}
	}

	flush(duplicateCount, previousLine, &builder)

	if err := scanner.Err(); err != nil {
		builder.WriteString(fmt.Sprintf("\n--- ERROR READING LOG: %v ---", err))
	}

	return builder.String()
}

func flush(duplicateCount int, previousLine string, builder *strings.Builder) {
	if previousLine == "" {
		return
	}
	if duplicateCount > 1 {
		builder.WriteString(fmt.Sprintf("%s (repeated %d times)\n", previousLine, duplicateCount))
	} else {
		builder.WriteString(previousLine + "\n")
	}
}
