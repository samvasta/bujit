package parse

import (
	"regexp"
	"sort"
)

type Token struct {
	Id      int
	Pattern *regexp.Regexp
}

func (tok *Token) Matches(input string) bool {
	return tok.Pattern.MatchString(input)
}

type LangNode struct {
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

func PossibleMatches(test string, possibleTokens []*Token) (ret []*Token) {
	// Copy possible tokens
	var matchLens map[int]int = make(map[int]int, len(possibleTokens))

	//Remove bad tokens & store match lengths for sorting
	for i, tok := range possibleTokens {
		prefix, _ := tok.Pattern.LiteralPrefix()

		if len(prefix) > 0 {
			lenOfMatch := LengthOfMatch(test, prefix)
			matchLens[i] = lenOfMatch
			if lenOfMatch > 0 {
				ret = append(ret, tok)
			}
		}
	}

	sort.Slice(ret, func(i, j int) bool {
		return matchLens[i] < matchLens[j]
	})

	return ret
}
