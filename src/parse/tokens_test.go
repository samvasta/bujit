package parse

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeToken(t *testing.T) {
	token := MakeLiteralToken(1, "test", "alt")

	literalPatterns := token.LiteralPatternStrings()
	assert.Contains(t, literalPatterns, "test")
	assert.Contains(t, literalPatterns, "TEST")
	assert.Contains(t, literalPatterns, "alt")
	assert.Contains(t, literalPatterns, "ALT")
}

func TestMakeArgToken(t *testing.T) {
	testCase := func(displayName, expectedDisplayName string) func(t *testing.T) {
		return func(t *testing.T) {
			token := MakeArgToken(123, displayName, IntegerPattern)

			assert.Equal(t, token.DisplayName, expectedDisplayName)
			assert.Equal(t, token.Id, 123)
			assert.Len(t, token.Patterns, 1)
			assert.Same(t, token.Patterns[0], IntegerPattern)
		}
	}

	t.Run("one word name", testCase("name", "<name>"))
	t.Run("multiple word name", testCase("multi word name", "<multi-word-name>"))
	t.Run("preformatted name", testCase("<name>", "<name>"))
	t.Run("preformatted multi word name", testCase("<multi-word-name>", "<multi-word-name>"))
	t.Run("Removes symbols", testCase("n&%@#!$a(\"{][me", "<n-a-me>"))
}

func TestDisplayNames(t *testing.T) {
	tokens := []*TokenPattern{
		MakeLiteralToken(1, "one", "uno"),
		MakeArgToken(2, "two", IntegerPattern),
		MakeLiteralToken(3, "three"),
	}

	displayNames := DisplayNames(tokens)

	assert.Len(t, displayNames, 3)
	assert.Contains(t, displayNames, "one")
	assert.Contains(t, displayNames, "<two>")
	assert.Contains(t, displayNames, "three")
}

func TestCurrencySymbolsWorkInRegex(t *testing.T) {
	testCase := func(symbol string) func(t *testing.T) {
		return func(t *testing.T) {
			result := regexp.MustCompile(fmt.Sprintf("[%s]", CurrencySymbols)).MatchString(symbol)
			assert.Truef(t, result, "%s was not recognized as a currency symbol, but should have been", symbol)
		}
	}

	for _, s := range CurrencySymbols {
		currencySymbol := string(s)
		t.Run(currencySymbol, testCase(currencySymbol))
	}
}

func TestIntegerPattern(t *testing.T) {
	testCase := func(input string, isValid bool) func(t *testing.T) {
		return func(t *testing.T) {
			result := IntegerPattern.FindString(input)
			if isValid {
				assert.Equal(t, result, input)
			} else {
				assert.NotEqual(t, result, input)
			}
		}
	}

	t.Run("positive int", testCase("1234", true))
	t.Run("negative int", testCase("-1234", true))
	t.Run("positive decimal", testCase("12.34", false))
	t.Run("negative decimal", testCase("-12.34", false))

	t.Run("positive int with currency", testCase("$1234", true))
	t.Run("negative int with currency", testCase("$-1234", true))
	t.Run("positive decimal with currency", testCase("$12.34", false))
	t.Run("negative decimal with currency", testCase("$-12.34", false))

	t.Run("words", testCase("not an int", false))
}

func TestTokenMatches(t *testing.T) {
	testCase := func(token TokenPattern, test string, expected bool) func(t *testing.T) {
		return func(t *testing.T) {
			result := token.Matches(test)
			if expected {
				assert.Truef(t, result, "Token with pattern %s does not match '%s'", token.Patterns[0].String(), test)
			} else {
				assert.Falsef(t, result, "Token with pattern %s matched '%s' but should not have", token.Patterns[0].String(), test)
			}
		}
	}

	literal := TokenPattern{1, "Test Literal", []*regexp.Regexp{regexp.MustCompile("literal1!")}}
	expandable := TokenPattern{1, "Test Expandable", []*regexp.Regexp{regexp.MustCompile("ab+!\\d+")}}

	t.Run("Matches literal", testCase(literal, "literal1!", true))
	t.Run("Matches expandable", testCase(expandable, "abbb!1234", true))

	t.Run("Does not match literal", testCase(expandable, "!1laretil", false))
	t.Run("Does not match expandable", testCase(expandable, "literal1!", false))

}

func TestTokenize(t *testing.T) {
	test := `This is a string that "will be" highlighted when your 'regular expression' matches something.`

	tokens := Tokenize(test)

	expected := []string{
		"This",
		"is",
		"a",
		"string",
		"that",
		"will be",
		"highlighted",
		"when",
		"your",
		"regular expression",
		"matches",
		"something.",
	}

	for i, token := range tokens {
		assert.Equalf(t, token, expected[i], "Token %d did not match expected. Found \"%s\" but expected \"%s\"", i, token, expected[i])
	}
}
func TestTokenizeWithQuotedStrings(t *testing.T) {
	test := "\"This should be one token\""

	tokens := Tokenize(test)

	assert.Len(t, tokens, 1)

	test = "'This should be one token'"

	tokens = Tokenize(test)

	assert.Len(t, tokens, 1)
}

func TestLengthOfMatch(t *testing.T) {

	testCase := func(str1, str2 string, output int) func(t *testing.T) {
		return func(t *testing.T) {
			result := LengthOfMatch(str1, str2)
			assert.Equalf(t, result, output, "'%s' vs '%s' match length should have been 3 but got %d", str1, str2, result)
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
	var tokens []*TokenPattern = []*TokenPattern{
		{0, "abc", []*regexp.Regexp{regexp.MustCompile(`abc`)}},                     // only alpha literals
		{1, "123", []*regexp.Regexp{regexp.MustCompile(`123`)}},                     // only number literals
		{2, "$#@!", []*regexp.Regexp{regexp.MustCompile(regexp.QuoteMeta("$#@!"))}}, // only symbol literals
		{3, "alpha nums", []*regexp.Regexp{regexp.MustCompile(`alpha\d+`)}},         // alpha literals plus expandable
		{4, "192aaaa...", []*regexp.Regexp{regexp.MustCompile(`193a*`)}},            // number literals plus expandable
		{5, "digits", []*regexp.Regexp{regexp.MustCompile(`\d+`)}},                  // only expandable numbers
		{6, "letters", []*regexp.Regexp{regexp.MustCompile(`[a-zA-Z]`)}},            // only expandable alphas
	}

	testCase := func(input string, expectedTokens []*TokenPattern) func(t *testing.T) {
		var tokenMap map[int]*TokenPattern = make(map[int]*TokenPattern)

		for _, token := range expectedTokens {
			tokenMap[token.Id] = token
		}

		return func(t *testing.T) {
			_, matches := PossibleMatches(input, tokens)

			var extraTokens strings.Builder
			for _, token := range matches {
				_, ok := tokenMap[token.Id]
				if !ok {
					extraTokens.WriteString("(")
					for _, pattern := range token.Patterns {
						extraTokens.WriteString(pattern.String())
						extraTokens.WriteString(" or ")
					}
					extraTokens.WriteString("), ")
				}
			}

			var missingTokens strings.Builder
			for _, token := range expectedTokens {
				_, ok := tokenMap[token.Id]
				if !ok {
					extraTokens.WriteString("(")
					for _, pattern := range token.Patterns {
						extraTokens.WriteString(pattern.String())
						extraTokens.WriteString(" or ")
					}
					extraTokens.WriteString("), ")
				}
			}

			assert.Greater(t, len(missingTokens.String()), 0, "Missing tokens: %s", missingTokens.String())
			assert.Greater(t, len(extraTokens.String()), 0, "Extra tokens: %s", extraTokens.String())
		}
	}

	t.Run("Alpha Literals", func(t *testing.T) { testCase("abc", []*TokenPattern{tokens[0]}) })
	t.Run("Number Literals", func(t *testing.T) { testCase("123", []*TokenPattern{tokens[1]}) })
	t.Run("Symbol Literals", func(t *testing.T) { testCase("$#", []*TokenPattern{tokens[2]}) })

	t.Run("Partial alpha literals", func(t *testing.T) { testCase("a", []*TokenPattern{tokens[0], tokens[3]}) })

	t.Run("Partial Number Literals", func(t *testing.T) { testCase("1", []*TokenPattern{tokens[1], tokens[4]}) })

	t.Run("No matches", func(t *testing.T) { testCase("nothingshouldmatch", []*TokenPattern{}) })

}

func TestPossibleMatches2(t *testing.T) {
	tokens := []*TokenPattern{
		MakeLiteralToken(1, "fun"),
		MakeLiteralToken(2, "fan"),
		MakeLiteralToken(3, "fact"),
		MakeArgToken(4, "Arg", ItemNamePattern),
	}

	exactMatch, possibleMatches := PossibleMatches("fact", tokens)

	assert.Equal(t, tokens[2], exactMatch)
	assert.Contains(t, possibleMatches, tokens[0])
	assert.Contains(t, possibleMatches, tokens[1])
	assert.NotContains(t, possibleMatches, tokens[3])
}
