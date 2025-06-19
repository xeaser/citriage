package dedupe

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

var timestampRegex = regexp.MustCompile(`^\[.*?\]\s*`)

func Process(logData []byte) string {
	var builder strings.Builder
	scanner := bufio.NewScanner(bytes.NewReader(logData))

	var previousLine string
	var duplicateCount int = 1

	firstLine := true

	flush := func() {
		if previousLine == "" {
			return
		}
		if duplicateCount > 1 {
			builder.WriteString(fmt.Sprintf("%s (repeated %d times)\n", previousLine, duplicateCount))
		} else {
			builder.WriteString(previousLine + "\n")
		}
	}

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
			flush()
			previousLine = currentLine
			duplicateCount = 1
		}
	}

	flush()

	if err := scanner.Err(); err != nil {
		builder.WriteString(fmt.Sprintf("\n--- ERROR READING LOG: %v ---", err))
	}

	return builder.String()
}
