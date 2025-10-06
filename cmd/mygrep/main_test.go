package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/mmarchesotti/build-your-own-grep/internal/nfasimulator"
)

func TestMatchLine(t *testing.T) {
	testCases := []struct {
		name             string
		line             []byte
		pattern          string
		expectedMatch    bool
		expectedCaptures [][]byte
	}{
		{
			name:          "Literal: Simple match",
			line:          []byte("abc"),
			pattern:       "a",
			expectedMatch: true,
		},
		{
			name:          "Start Anchor (^): Match at beginning",
			line:          []byte("abc"),
			pattern:       "^a",
			expectedMatch: true,
		},
		{
			name:          "Combination: Match a literal and a digit",
			line:          []byte("a1c"),
			pattern:       `a\d`,
			expectedMatch: true,
		},
		{
			name:          "Match one or more times: codecrafters #02",
			line:          []byte("caaats"),
			pattern:       `ca+at`,
			expectedMatch: true,
		},

		{
			name:             "Capture: Single simple group",
			line:             []byte("hello world"),
			pattern:          "w(o)rld",
			expectedMatch:    true,
			expectedCaptures: [][]byte{[]byte("world"), []byte("o")},
		},
		{
			name:             "Capture: Multiple groups",
			line:             []byte("abcde"),
			pattern:          "a(b)c(d)e",
			expectedMatch:    true,
			expectedCaptures: [][]byte{[]byte("abcde"), []byte("b"), []byte("d")},
		},
		{
			name:             "Capture: Nested groups",
			line:             []byte("axyzb"),
			pattern:          "a(x(y)z)b",
			expectedMatch:    true,
			expectedCaptures: [][]byte{[]byte("axyzb"), []byte("xyz"), []byte("y")},
		},
		{
			name:             "Capture: Group with quantifier",
			line:             []byte("ababab"),
			pattern:          "(ab)+",
			expectedMatch:    true,
			expectedCaptures: [][]byte{[]byte("ababab"), []byte("ab")},
		},
		{
			name:             "Capture: Full line match",
			line:             []byte("test"),
			pattern:          "(test)",
			expectedMatch:    true,
			expectedCaptures: [][]byte{[]byte("test"), []byte("test")},
		},
		{
			name:             "Capture: No match should return no captures",
			line:             []byte("abc"),
			pattern:          "(d)",
			expectedMatch:    false,
			expectedCaptures: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualMatch, actualCaptures, err := nfasimulator.Simulate(tc.line, tc.pattern)
			if err != nil {
				t.Errorf("error '%s':", err)
			}

			if actualMatch != tc.expectedMatch {
				t.Errorf("Pattern '%s' on line '%s': expected match %v, but got %v",
					tc.pattern, string(tc.line), tc.expectedMatch, actualMatch)
			}

			if tc.expectedCaptures != nil {
				if !reflect.DeepEqual(actualCaptures, tc.expectedCaptures) {
					t.Errorf("Pattern '%s' on line '%s': incorrect captures", tc.pattern, string(tc.line))
					t.Errorf("  got:  %v", byteSlicesToStrings(actualCaptures))
					t.Errorf("  want: %v", byteSlicesToStrings(tc.expectedCaptures))
				}
			}
		})
	}
}

func byteSlicesToStrings(bss [][]byte) []string {
	ss := make([]string, len(bss))
	for i, bs := range bss {
		ss[i] = string(bs)
	}
	return ss
}

func TestSimulateWithFile(t *testing.T) {
	testCases := []struct {
		name          string
		fileContent   string
		pattern       string
		expectedMatch bool
	}{
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
			fileBytes := []byte(tc.fileContent)

			actualMatch, _, err := nfasimulator.Simulate(fileBytes, tc.pattern)
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
