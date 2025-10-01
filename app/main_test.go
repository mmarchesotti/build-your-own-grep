// Create this file right next to your main.go
package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mmarchesotti/build-your-own-grep/app/nfasimulator"
)

func TestMatchLine(t *testing.T) {
	testCases := []struct {
		name          string
		line          []byte
		pattern       string
		expectedMatch bool
	}{
		// Basic Literal Matches
		{
			name:          "Literal: Simple match",
			line:          []byte("abc"),
			pattern:       "a",
			expectedMatch: true,
		},
		{
			name:          "Literal: No match",
			line:          []byte("abc"),
			pattern:       "d",
			expectedMatch: false,
		},
		{
			name:          "Literal: Match anywhere in string",
			line:          []byte("xbyc"),
			pattern:       "b",
			expectedMatch: true,
		},

		// Digit '\d'
		{
			name:          `Digit (\d): Match`,
			line:          []byte("a1c"),
			pattern:       `\d`,
			expectedMatch: true,
		},
		{
			name:          `Digit (\d): No match`,
			line:          []byte("abc"),
			pattern:       `\d`,
			expectedMatch: false,
		},

		// Alphanumeric '\w'
		{
			name:          `Alphanumeric (\w): Match letter`,
			line:          []byte("1a2"),
			pattern:       `\w`,
			expectedMatch: true, // Matches '1'
		},
		{
			name:          `Alphanumeric (\w): No match`,
			line:          []byte("$#%"),
			pattern:       `\w`,
			expectedMatch: false,
		},

		// Start Anchor '^'
		{
			name:          "Start Anchor (^): Match at beginning",
			line:          []byte("abc"),
			pattern:       "^a",
			expectedMatch: true,
		},
		{
			name:          "Start Anchor (^): Fails when not at beginning",
			line:          []byte("bac"),
			pattern:       "^a",
			expectedMatch: false,
		},

		// End Anchor '$'
		{
			name:          "End Anchor ($): Match at end",
			line:          []byte("abc"),
			pattern:       "c$",
			expectedMatch: true,
		},
		{
			name:          "End Anchor ($): Fails when not at end",
			line:          []byte("acb"),
			pattern:       "c$",
			expectedMatch: false,
		},

		// Positive Character Group '[...]'
		{
			name:          "Positive Group: Match found",
			line:          []byte("axbyc"),
			pattern:       "[xyz]",
			expectedMatch: true, // Matches 'x'
		},
		{
			name:          "Positive Group: No match",
			line:          []byte("abc"),
			pattern:       "[xyz]",
			expectedMatch: false,
		},

		// Negative Character Group '[^...]'
		{
			name:          "Negative Group: Match found",
			line:          []byte("xay"),
			pattern:       "[^xyz]",
			expectedMatch: true, // Matches 'a'
		},
		{
			name:          "Negative Group: No match",
			line:          []byte("xyz"),
			pattern:       "[^xyz]",
			expectedMatch: false,
		},

		// Combination of patterns
		{
			name:          "Combination: Match a literal and a digit",
			line:          []byte("a1c"),
			pattern:       `a\d`,
			expectedMatch: true,
		},
		{
			name:          "Combination: Fails on wrong order",
			line:          []byte("1ac"),
			pattern:       `a\d`,
			expectedMatch: false,
		},

		// Match one or more times
		{
			name:          "Match one or more times: Match triple letter a in the middle of word",
			line:          []byte("caaats"),
			pattern:       `ca+ts`,
			expectedMatch: true,
		},
		{
			name:          "Match one or more times: Match triple letter a in the end of word",
			line:          []byte("caaa"),
			pattern:       `ca+`,
			expectedMatch: true,
		},
		{
			name:          "Match one or more times: Fails on no character matching",
			line:          []byte("caat"),
			pattern:       `cat+`,
			expectedMatch: false,
		},
		{
			name:          "Match one or more times: codecrafters #02",
			line:          []byte("caaats"),
			pattern:       `ca+at`,
			expectedMatch: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualMatch, err := nfasimulator.Simulate(tc.line, tc.pattern)
			if err != nil {
				t.Errorf("error '%s':", err)
			}

			if actualMatch != tc.expectedMatch {
				t.Errorf("Pattern '%s' on line '%s': expected match %v, but got %v",
					tc.pattern, string(tc.line), tc.expectedMatch, actualMatch)
			}
		})
	}
}

func TestSimulateWithFile(t *testing.T) {
	testCases := []struct {
		name          string
		fileContent   string
		pattern       string
		expectedMatch bool
	}{
		// Basic Literal Matches from a file
		{
			name:          "File - Literal: Simple match",
			fileContent:   "apple\nbanana\ncherry",
			pattern:       "banana",
			expectedMatch: true,
		},
		{
			name:          "File - Literal: No match",
			fileContent:   "apple\nbanana\ncherry",
			pattern:       "durian",
			expectedMatch: false,
		},

		// Regex pattern matching
		{
			name:          "File - Regex: Match word starting with 'app'",
			fileContent:   "application\napplepie\napplesauce",
			pattern:       `^appl.*`,
			expectedMatch: true,
		},
		{
			name:          "File - Regex: Match lines with digits",
			fileContent:   "line one\nline 2\nline three",
			pattern:       `\d`,
			expectedMatch: true,
		},
		{
			name:          "File - Regex: Fails when not at the start",
			fileContent:   "An application\nAn applepie",
			pattern:       `^appl.*`,
			expectedMatch: false,
		},

		// Edge Cases
		{
			name:          "File - Edge Case: Empty file",
			fileContent:   "",
			pattern:       "a",
			expectedMatch: false,
		},
		{
			name:          "File - Edge Case: Pattern matches entire file content",
			fileContent:   "supercalifragilisticexpialidocious",
			pattern:       ".*",
			expectedMatch: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filePath := createTestFile(t, tc.fileContent)

			fileBytes, err := os.ReadFile(filePath)
			if err != nil {
				t.Fatalf("Failed to read temp file %s: %v", filePath, err)
			}

			actualMatch, err := nfasimulator.Simulate(fileBytes, tc.pattern)
			if err != nil {
				t.Errorf("Simulate returned an unexpected error: %v", err)
			}

			if actualMatch != tc.expectedMatch {
				t.Errorf("Pattern '%s' on file with content '%s': expected match %v, but got %v",
					tc.pattern, tc.fileContent, tc.expectedMatch, actualMatch)
			}
		})
	}
}

func createTestFile(t *testing.T, content string) string {
	t.Helper()

	tempDir := t.TempDir()

	filePath := filepath.Join(tempDir, "testfile.txt")

	err := os.WriteFile(filePath, []byte(content), 0666)
	if err != nil {
		t.Fatalf("Failed to write to temporary file %s: %v", filePath, err)
	}

	return filePath
}
