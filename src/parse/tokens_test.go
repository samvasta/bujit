package parse

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func TestMakeToken(t *testing.T) {
	token := MakeLiteralToken(1, "test", "alt")
	if token.Pattern.String() != "test|TEST|alt|ALT" {
		t.Errorf("Expected test|TEST|alt|ALT, got %s", token.Pattern.String())
	}
}

func TestCurrencySymbolsWorkInRegex(t *testing.T) {
	testCase := func(symbol string) func(t *testing.T) {
		return func(t *testing.T) {
			result := regexp.MustCompile(fmt.Sprintf("[%s]", CurrencySymbols)).MatchString(symbol)
			if !result {
				t.Errorf("%s was not recognized in regex", symbol)
			}
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
			if result != input && isValid {
				t.Errorf("%s did not match the integer pattern", input)
			} else if result == input && !isValid {
				t.Errorf("%s matched the integer pattern but shouldn't", input)
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
			if result != expected {
				if expected {
					t.Errorf("Token with pattern %s does not match '%s'", token.Pattern.String(), test)
				} else {
					t.Errorf("Token with pattern %s matched '%s' but should not have", token.Pattern.String(), test)
				}
			}
		}
	}

	literal := TokenPattern{1, "Test Literal", regexp.MustCompile("literal1!")}
	expandable := TokenPattern{1, "Test Expandable", regexp.MustCompile("ab+!\\d+")}

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
		if expected[i] != token {
			t.Errorf("Token %d did not match expected. Found \"%s\" but expected \"%s\"", i, token, expected[i])
		}
	}
}
func TestTokenizeWithQuotedStrings(t *testing.T) {
	test := "\"This should be one token\""

	tokens := Tokenize(test)

	if len(tokens) != 1 {
		t.Error()
	}

	test = "'This should be one token'"

	tokens = Tokenize(test)

	if len(tokens) != 1 {
		t.Error()
	}
}

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
	var tokens []*TokenPattern = []*TokenPattern{
		{0, "abc", regexp.MustCompile(`abc`)},                     // only alpha literals
		{1, "123", regexp.MustCompile(`123`)},                     // only number literals
		{2, "$#@!", regexp.MustCompile(regexp.QuoteMeta("$#@!"))}, // only symbol literals
		{3, "alpha nums", regexp.MustCompile(`alpha\d+`)},         // alpha literals plus expandable
		{4, "192aaaa...", regexp.MustCompile(`193a*`)},            // number literals plus expandable
		{5, "digits", regexp.MustCompile(`\d+`)},                  // only expandable numbers
		{6, "letters", regexp.MustCompile(`[a-zA-Z]`)},            // only expandable alphas
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

	t.Run("Alpha Literals", func(t *testing.T) { testCase("abc", []*TokenPattern{tokens[0]}) })
	t.Run("Number Literals", func(t *testing.T) { testCase("123", []*TokenPattern{tokens[1]}) })
	t.Run("Symbol Literals", func(t *testing.T) { testCase("$#", []*TokenPattern{tokens[2]}) })

	t.Run("Partial alpha literals", func(t *testing.T) { testCase("a", []*TokenPattern{tokens[0], tokens[3]}) })

	t.Run("Partial Number Literals", func(t *testing.T) { testCase("1", []*TokenPattern{tokens[1], tokens[4]}) })

	t.Run("No matches", func(t *testing.T) { testCase("nothingshouldmatch", []*TokenPattern{}) })

}
