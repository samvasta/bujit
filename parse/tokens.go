package parse

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"samvasta.com/bujit/models"
)

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

func makeFlagToken(id int, shortName, longName string) *TokenPattern {
	displayName := fmt.Sprintf("--%s", longName)
	patterns := []*regexp.Regexp{
		regexp.MustCompile(regexp.QuoteMeta(fmt.Sprintf("-%s", shortName))),
		regexp.MustCompile(regexp.QuoteMeta(displayName)),
	}

	return &TokenPattern{id, displayName, patterns}
}

var TokenizePattern = regexp.MustCompile(`[^\s"'=]+=?|"([^"]*)"|'([^']*)'`)

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

func (tok *TokenPattern) BestMatch(test string) string {
	var bestSoFar string
	bestSoFarMatchLen := 0
	for _, pattern := range tok.Patterns {
		prefix, _ := pattern.LiteralPrefix()

		if len(test) > len(prefix) {
			continue
		}

		if len(prefix) > 0 {
			lenOfMatch := LengthOfMatch(test, prefix)

			if lenOfMatch == len(test) && len(test) == len(prefix) {
				// The token must be 100% literal and be an exact match of the test str
				return prefix
			} else {
				submatches := pattern.FindStringSubmatch(test)

				if len(submatches) > 0 && len(submatches[0]) == len(test) {
					// The token starts with literals and has an expandable section, but still matches the test str completely
					return prefix
				}
			}

			if lenOfMatch > bestSoFarMatchLen {
				bestSoFarMatchLen = lenOfMatch
				bestSoFar = prefix
			}
		}
	}
	if bestSoFarMatchLen == 0 {
		return tok.DisplayName
	}

	return bestSoFar
}

func PossibleMatches(test string, possibleTokens []*TokenPattern) (exactMatch *TokenPattern, possibleMatches []*TokenPattern) {

	// Copy possible tokens
	var matchLens map[int]int = make(map[int]int, len(possibleTokens))

	//Remove bad tokens & store match lengths for sorting
	for i, tok := range possibleTokens {
		isTokenAlreadyAdded := false
		for _, pattern := range tok.Patterns {
			prefix, isComplete := pattern.LiteralPrefix()

			if isComplete && len(prefix) < len(test) {
				continue
			}

			if len(prefix) > 0 {
				lenOfMatch := LengthOfMatch(test, prefix)
				matchLens[i] = lenOfMatch
				if !isTokenAlreadyAdded && lenOfMatch > 0 {
					possibleMatches = append(possibleMatches, tok)
					isTokenAlreadyAdded = true
				}
				if lenOfMatch == len(test) && len(test) == len(prefix) {
					// The token must be 100% literal and be an exact match of the test str
					exactMatch = tok
				} else {
					submatches := pattern.FindStringSubmatch(test)

					if len(submatches) > 0 && len(submatches[0]) == len(test) {
						// The token starts with literals and has an expandable section, but still matches the test str completely
						exactMatch = tok
					}
				}
			} else if pattern.MatchString(test) {
				matchLens[i] = 0
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

func itemNameValue(tokenStr string) string {
	return strings.Trim(tokenStr, `'"`)
}

func moneyValue(tokenStr string) models.Money {
	submatches := DecimalPattern.FindStringSubmatch(tokenStr)

	if len(submatches) > 1 {
		//Guaranteed to never get a parsing error because if submatch exists then the value must fit the float pattern
		value, _ := strconv.ParseFloat(submatches[1], 64)
		return models.MakeMoney(value)
	}

	return models.MakeMoney(0)
}
