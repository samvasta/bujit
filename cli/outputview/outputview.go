package outputview

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/muesli/termenv"
	"github.com/olekukonko/ts"
	"samvasta.com/bujit/models/output"
)

func TerminalColor(col output.ColorHint) termenv.Color {
	switch col {
	case output.Body:
		return termenv.Ascii.Color("")
	case output.Subtle:
		return termenv.ANSIBrightBlack
	case output.Primary:
		return termenv.ANSIBrightCyan
	case output.Success:
		return termenv.ANSIBrightGreen
	case output.Info:
		return termenv.ANSIBrightBlue
	case output.Warning:
		return termenv.ANSIBrightYellow
	case output.Error:
		return termenv.ANSIBrightRed
	default:
		return termenv.ANSIBlack
	}
}

func View(item interface{}) string {
	switch i := item.(type) {
	case []output.Helper:
		sb := strings.Builder{}
		for _, helper := range i {
			sb.WriteString(View(helper))
			sb.WriteString("\n")
		}
		return sb.String()
	case output.Text:
		return TextView(i)
	case output.UnorderedList:
		return UnorderedListView(i)
	default:
		return termenv.
			String(fmt.Sprintf("!!Type %T not supported yet.!!", i)).
			Background(termenv.ANSIBrightRed).
			Foreground(termenv.ANSIBrightWhite).
			String()
	}
}

func WrappedString(text string, initialIndent, indentLevel int) string {
	size, _ := ts.GetSize()
	sb := strings.Builder{}

	sb.WriteString(strings.Repeat("  ", initialIndent-indentLevel))
	chunkSize := size.Col() - (indentLevel)*2

	prevIdx := 0
	for idx := chunkSize - initialIndent*2; idx < len(text); idx += chunkSize {
		sb.WriteString(strings.Repeat("  ", indentLevel))
		sb.WriteString(text[prevIdx:idx])
		sb.WriteString("\n")
		prevIdx = idx
	}
	return sb.String()
}

var wordsRegex = regexp.MustCompile(`\s`)

func BetterWrappedStringBlock(text string, initialCol, minCol, maxCol int) string {
	splitPoints := []int{}
	for _, idx := range wordsRegex.FindAllStringIndex(text, -1) {
		splitPoints = append(splitPoints, idx[0])
	}
	splitPoints = append(splitPoints, len(text))

	size, _ := ts.GetSize()
	sb := strings.Builder{}

	if maxCol > size.Col() {
		maxCol = size.Col()
	}

	maxWidth := maxCol - minCol

	// Pad out to initial col
	// sb.WriteString(strings.Repeat(" ", initialCol))

	start := 0
	currentIdx := 0
	for currentIdx < len(splitPoints) {
		end := splitPoints[currentIdx]

		// handle case where word is longer than the max width
		for end > start && end-start > maxWidth {
			sb.WriteString(strings.Repeat(" ", minCol))
			sb.WriteString(text[start:start+maxWidth] + "\n")
			start += maxWidth
		}

		// handle case where multiple words can fit on one line
		for end < start && end-start < maxWidth {
			currentIdx++
			end = splitPoints[currentIdx]
		}

		sb.WriteString(strings.Repeat(" ", minCol))
		sb.WriteString(text[start:end] + "\n")
		start += maxWidth
		currentIdx++
	}

	return sb.String()
}

func WrappedStringBlock(text string, initialCol, minCol, maxCol int) string {
	size, _ := ts.GetSize()
	sb := strings.Builder{}

	if maxCol > size.Col() {
		maxCol = size.Col()
	}

	chunkSize := maxCol - minCol

	sb.WriteString(strings.Repeat(" ", initialCol))
	prevIdx := 0
	idx := chunkSize - (initialCol - minCol)
	sb.WriteString(text[prevIdx:idx])
	sb.WriteRune('\n')

	for prevIdx, idx = idx, idx+chunkSize; idx < len(text); idx += chunkSize {
		sb.WriteString(strings.Repeat(" ", minCol))
		sb.WriteString(text[prevIdx:idx])
		sb.WriteRune('\n')
		prevIdx = idx
	}
	return sb.String()
}

func TextView(t output.Text) string {
	output := termenv.
		String(WrappedString(t.Text, 0, t.Indent)).
		Foreground(TerminalColor(t.Style.Color))

	if t.Style.IsBold {
		output = output.Bold()
	}
	if t.Style.IsItalic {
		output = output.Italic()
	}
	if t.Style.IsUnderline {
		output = output.Underline()
	}

	return output.String()
}

func UnorderedListView(ul output.UnorderedList) string {
	return BetterWrappedStringBlock("hello goodbye. captain! ok abcdefghijklmnopqrstuvwxyz", 5, 5, 20)
}
