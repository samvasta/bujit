package parse

import (
	"regexp"
	"strings"
	"testing"
)

func TestLengthOfMatch(t *testing.T) {

	testCase := func(str1, str2 string, output int) func(t *testing.T) {
		return func(t *testing.T) {
			result := LengthOfMatch(str1, str2)
			if result != output {
				t.Errorf("'%s' vs '%s' match length should have been 3 but got %d", str1, str2, result)
			}
		}
	}

	t.Run("Partial match", testCase("ABCDEF", "ABCdef", 3))
	t.Run("No match", testCase("ABCDEF", "BCDEF", 0))
	t.Run("Full match on str2", testCase("ABCDEF", "ABC", 3))
	t.Run("Full match on str1", testCase("ABC", "ABCdef", 3))
	t.Run("Empty strings", testCase("", "", 0))
	t.Run("Full Match", testCase("ABCDEF", "ABCDEF", 6))
}

func TestPossibleMatches(t *testing.T) {
	var tokens []*Token = []*Token{
		{0, regexp.MustCompile(`abc`)},                    // only alpha literals
		{0, regexp.MustCompile(`123`)},                    // only number literals
		{0, regexp.MustCompile(regexp.QuoteMeta("$#@!"))}, // only symbol literals
		{0, regexp.MustCompile(`alpha\d+`)},               // alpha literals plus expandable
		{0, regexp.MustCompile(`193a*`)},                  // number literals plus expandable
		{0, regexp.MustCompile(`\d+`)},                    // only expandable numbers
		{0, regexp.MustCompile(`[a-zA-Z]`)},               // only expandable alphas
	}

	testCase := func(input string, expectedTokens []*Token) func(t *testing.T) {
		var tokenMap map[int]*Token = make(map[int]*Token)

		for _, token := range expectedTokens {
			tokenMap[token.Id] = token
		}

		return func(t *testing.T) {
			result := PossibleMatches(input, tokens)

			var extraTokens strings.Builder
			for _, token := range result {
				_, ok := tokenMap[token.Id]
				if !ok {
					extraTokens.WriteString(token.Pattern.String())
					extraTokens.WriteString(", ")
				}
			}

			var missingTokens strings.Builder
			for _, token := range expectedTokens {
				_, ok := tokenMap[token.Id]
				if !ok {
					missingTokens.WriteString(token.Pattern.String())
					missingTokens.WriteString(", ")
				}
			}

			if len(missingTokens.String()) > 0 {
				t.Errorf("Missing tokens: %s", missingTokens.String())
			} else if len(extraTokens.String()) > 0 {
				t.Errorf("Extra tokens: %s", extraTokens.String())
			}
		}
	}

	t.Run("Alpha Literals", func(t *testing.T) { testCase("abc", []*Token{tokens[0]}) })
	t.Run("Number Literals", func(t *testing.T) { testCase("123", []*Token{tokens[1]}) })
	t.Run("Symbol Literals", func(t *testing.T) { testCase("$#", []*Token{tokens[2]}) })

	t.Run("Partial alpha literals", func(t *testing.T) { testCase("a", []*Token{tokens[0], tokens[3]}) })

	t.Run("Partial Number Literals", func(t *testing.T) { testCase("1", []*Token{tokens[1], tokens[4]}) })

	t.Run("No matches", func(t *testing.T) { testCase("nothingshouldmatch", []*Token{}) })

}
