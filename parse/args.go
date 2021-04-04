package parse

import (
	"fmt"
	"regexp"
	"strings"
)

var IntegerPattern *regexp.Regexp = regexp.MustCompile(`[^\w|\d|\.|\,|'|"|_|-]?(-?\d+)(\W?[A-Z]{3})?`)
var DecimalPattern *regexp.Regexp = regexp.MustCompile(`[^\w|\d|\.|\,|'|"|_|-]?(-?\d+(\.\d{1,2})?)(\W?[A-Z]{3})?`)

var ItemNamePattern *regexp.Regexp = regexp.MustCompile(`[a-zA-Z][a-zA-Z_-]+`)

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

func MakeOptionalArgToken(id int, shortName, longName string) *TokenPattern {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(fmt.Sprintf("-%s=?", shortName)),
		regexp.MustCompile(fmt.Sprintf("--%s=?", longName)),
	}
	return &TokenPattern{id, fmt.Sprintf("--%s", longName), patterns}

}

func ExtractArgValue(token string) string {
	idx := strings.Index(token, "=")
	if idx < 0 {
		return ""
	}
	return token[idx+1:]
}

func parseOptionalArg(ctx *ParseContext, tok *TokenPattern, argPattern *regexp.Regexp, argPatternName string) (value string, suggestion AutoSuggestion) {
	nextToken, hasNext := ctx.nextToken()

	if !hasNext {
		// Missing flag token
		return "", makeAutoSuggestion(false, "", []*TokenPattern{tok})
	} else {

		if !tok.Matches(nextToken) {
			return "", makeAutoSuggestion(false, nextToken, []*TokenPattern{tok})
		}

		// Parse next token
		ctx.moveToNextToken()

		nextToken, hasNext = ctx.nextToken()

		if !hasNext {
			// Missing arg value
			return "", AutoSuggestion{false, "", []string{fmt.Sprintf("<%s>", argPatternName)}}
		}

		ctx.moveToNextToken()
		return nextToken, AutoSuggestion{true, "", []string{}}
	}
}
