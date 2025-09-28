// Create this file right next to your main.go
package main

import "testing"

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

	// The test runner iterates through each case in the table
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Run the function we're testing
			actualMatch := matchLine(tc.line, tc.pattern)

			// Compare the actual result with what we expected
			if actualMatch != tc.expectedMatch {
				t.Errorf("Pattern '%s' on line '%s': expected match %v, but got %v",
					tc.pattern, string(tc.line), tc.expectedMatch, actualMatch)
			}
		})
	}
}
