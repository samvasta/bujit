package parse

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

const CurrencySymbols string = "$€£¥₿ɱŁ"

var IntegerPattern *regexp.Regexp = regexp.MustCompile(fmt.Sprintf(`[%s]?-?\d+`, CurrencySymbols))
var DecimalPattern *regexp.Regexp = regexp.MustCompile(fmt.Sprintf(`[%s]?-?\d+(\.\d{1,2})?`, CurrencySymbols))

var ItemNamePattern *regexp.Regexp = regexp.MustCompile(`[a-zA-Z][a-zA-Z_-]+`)

type TokenPattern struct {
	Id          int
	DisplayName string
	Pattern     *regexp.Regexp
}

func (tok *TokenPattern) Matches(input string) bool {
	idx := tok.Pattern.FindStringIndex(input)
	fmt.Print(idx)
	return idx != nil && idx[0] == 0 && idx[1] == len(input)
}

func DisplayNames(tokens []*TokenPattern) (names []string) {
	for _, token := range tokens {
		names = append(names, token.DisplayName)
	}
	return names
}

func MakeLiteralToken(id int, literalOptions ...string) *TokenPattern {
	var sb strings.Builder

	for i, option := range literalOptions {
		if i > 0 {
			sb.WriteRune('|')
		}
		// Lowercase
		for _, char := range option {
			sb.WriteRune(unicode.ToLower(char))
		}

		// Uppercase
		sb.WriteRune('|')
		for _, char := range option {
			sb.WriteRune(unicode.ToUpper(char))
		}
	}

	return &TokenPattern{id, literalOptions[0], regexp.MustCompile(sb.String())}
}

var ArgumentTokenPattern *regexp.Regexp = regexp.MustCompile("[ a-zA-Z0-9]+")

func MakeArgToken(id int, displayName string, pattern *regexp.Regexp) *TokenPattern {
	var finalName string
	var sb strings.Builder

	sb.WriteRune('<')
	for i, match := range ArgumentTokenPattern.FindAllString(displayName, -1) {
		if i > 0 {
			sb.WriteRune('-')
		}
		sb.WriteString(match)
	}
	sb.WriteRune('>')

	return &TokenPattern{id, finalName, pattern}
}

var TokenizePattern = regexp.MustCompile(regexp.QuoteMeta(`[^\s"']+|"([^"]*)"|'([^']*)'`))

func Tokenize(input string) (tokens []string) {
	matches := TokenizePattern.FindAllStringSubmatch(input, -1)
	fmt.Println(matches)
	for _, match := range matches {
		if match[1] != "" {
			tokens = append(tokens, match[1])
		} else {
			tokens = append(tokens, match[0])
		}
	}
	fmt.Println(tokens)
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
		prefix, _ := tok.Pattern.LiteralPrefix()

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

	// Sort by match length
	sort.Slice(possibleMatches, func(i, j int) bool {
		return matchLens[i] < matchLens[j]
	})

	return exactMatch, possibleMatches
}
