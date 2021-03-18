package parse

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

const CurrencySymbols string = "$€£¥₿ɱŁ"

var IntegerPattern *regexp.Regexp = regexp.MustCompile(fmt.Sprintf(`[%s]?-?\d+`, CurrencySymbols))
var DecimalPattern *regexp.Regexp = regexp.MustCompile(fmt.Sprintf(`[%s]?-?\d+(\.\d{1,2})?`, CurrencySymbols))

var ItemNamePattern *regexp.Regexp = regexp.MustCompile(`[a-zA-Z][a-zA-Z_-]+`)

type TokenPattern struct {
	Id          int
	DisplayName string
	Patterns    []*regexp.Regexp
}

func (tok *TokenPattern) Matches(input string) bool {
	for _, pattern := range tok.Patterns {
		idx := pattern.FindStringIndex(input)
		if idx != nil && idx[0] == 0 && idx[1] == len(input) {
			return true
		}
	}
	return false
}

func (token *TokenPattern) LiteralPatternStrings() (patterns []string) {
	for _, pattern := range token.Patterns {
		prefix, _ := pattern.LiteralPrefix()

		if len(prefix) > 0 {
			patterns = append(patterns, prefix)
		}
	}
	return patterns
}

func DisplayNames(tokens []*TokenPattern) (names []string) {
	for _, token := range tokens {
		names = append(names, token.DisplayName)
	}
	return names
}

func MakeLiteralToken(id int, literalOptions ...string) *TokenPattern {
	var patterns []*regexp.Regexp

	for _, option := range literalOptions {
		patterns = append(patterns, regexp.MustCompile(strings.ToUpper(option)))
		patterns = append(patterns, regexp.MustCompile(strings.ToLower(option)))
	}

	return &TokenPattern{id, literalOptions[0], patterns}
}

var ArgumentTokenPattern *regexp.Regexp = regexp.MustCompile("[a-zA-Z0-9]+")

func MakeArgToken(id int, displayName string, pattern *regexp.Regexp) *TokenPattern {
	var sb strings.Builder

	sb.WriteRune('<')
	for i, match := range ArgumentTokenPattern.FindAllString(displayName, -1) {
		if i > 0 {
			sb.WriteRune('-')
		}
		sb.WriteString(match)
	}
	sb.WriteRune('>')

	return &TokenPattern{id, sb.String(), []*regexp.Regexp{pattern}}
}

var TokenizePattern = regexp.MustCompile(`[^\s"']+|"([^"]*)"|'([^']*)'`)

func Tokenize(input string) (tokens []string) {
	matches := TokenizePattern.FindAllStringSubmatch(input, -1)
	for _, match := range matches {
		if match[2] != "" {
			tokens = append(tokens, match[2])
		} else if match[1] != "" {
			tokens = append(tokens, match[1])
		} else {
			tokens = append(tokens, match[0])
		}
	}
	return tokens
}

func LengthOfMatch(a, b string) int {
	var minLen int
	if len(a) < len(b) {
		minLen = len(a)
	} else {
		minLen = len(b)
	}

	for i := 0; i < minLen; i++ {
		if a[i] != b[i] {
			return i
		}
	}
	return minLen
}

func PossibleMatches(test string, possibleTokens []*TokenPattern) (exactMatch *TokenPattern, possibleMatches []*TokenPattern) {
	exactMatch = nil

	// Copy possible tokens
	var matchLens map[int]int = make(map[int]int, len(possibleTokens))

	//Remove bad tokens & store match lengths for sorting
	for i, tok := range possibleTokens {
		for _, pattern := range tok.Patterns {
			prefix, _ := pattern.LiteralPrefix()

			fmt.Printf(prefix)

			if len(prefix) > 0 {
				lenOfMatch := LengthOfMatch(test, prefix)
				matchLens[i] = lenOfMatch
				if lenOfMatch > 0 {
					possibleMatches = append(possibleMatches, tok)
				}
				if lenOfMatch == len(test) {
					exactMatch = tok
				}
			}
		}
	}

	// Sort by match length
	sort.Slice(possibleMatches, func(i, j int) bool {
		return matchLens[i] < matchLens[j]
	})

	return exactMatch, possibleMatches
}
