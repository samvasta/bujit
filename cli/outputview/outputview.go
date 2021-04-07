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
	case output.HorizontalRule:
		return HorizontalRuleView(i)
	default:
		return termenv.
			String(fmt.Sprintf("!!Type %T not supported yet.!!", i)).
			Background(termenv.ANSIBrightRed).
			Foreground(termenv.ANSIBrightWhite).
			String()
	}
}

var wordsRegex = regexp.MustCompile(`\s`)

func WrappedString(text string, startCol, minCol, maxCol int, terminalWidth int) string {

	if maxCol > terminalWidth {
		maxCol = terminalWidth
	}

	if minCol >= maxCol {
		panic("MinCol must be < MaxCol")
	}
	if startCol >= maxCol {
		panic("StartCol must be < MaxCol")
	}
	if terminalWidth < 1 {
		panic("TerminalWidth must be at least 1")
	}

	splitPoints := []int{}
	for _, idx := range wordsRegex.FindAllStringIndex(text, -1) {
		splitPoints = append(splitPoints, idx[0])
	}
	splitPoints = append(splitPoints, len(text))

	sb := strings.Builder{}

	outputColStart := startCol
	start := 0
	currentIdx := 0
	for currentIdx < len(splitPoints) {
		maxWidth := maxCol - outputColStart
		end := splitPoints[currentIdx]

		// handle case where word is longer than the max width
		for end > start && end-start > maxWidth {
			// print next chunk
			sb.WriteString(text[start:start+maxWidth] + "\n")

			start += maxWidth

			// recompute
			outputColStart = minCol
			maxWidth = maxCol - outputColStart
			// pad out to line start
			sb.WriteString(strings.Repeat(" ", outputColStart))
		}

		// handle case where multiple words can fit on one line
		for currentIdx+1 < len(splitPoints) && splitPoints[currentIdx+1]-start < maxWidth {
			currentIdx++
			end = splitPoints[currentIdx]
		}

		// print next chunk
		sb.WriteString(text[start:end] + "\n")
		// reset start to the min col
		outputColStart = minCol
		if currentIdx+1 < len(splitPoints) {
			// pad out to line start for next line
			sb.WriteString(strings.Repeat(" ", outputColStart))
		}

		// move start to the next unprinted char so we don't repeat anything
		start = end + 1
		// move to next split point
		currentIdx++
	}

	return sb.String()
}

func TextView(t output.Text) string {
	size, _ := ts.GetSize()
	output := termenv.
		String(WrappedString(t.Text, 0, t.Indent*2, size.Col(), size.Col())).
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

func HorizontalRuleView(hr output.HorizontalRule) string {
	size, _ := ts.GetSize()
	output := termenv.String(strings.Repeat(hr.RuleChar, size.Col())).
		Foreground(TerminalColor(hr.Style.Color))

	if hr.Style.IsBold {
		output = output.Bold()
	}
	if hr.Style.IsItalic {
		output = output.Italic()
	}
	if hr.Style.IsUnderline {
		output = output.Underline()
	}

	return output.String()
}

func UnorderedListView(ul output.UnorderedList) string {
	size, _ := ts.GetSize()
	sb := strings.Builder{}
	for _, item := range ul.Items {
		sb.WriteString(strings.Repeat("  ", ul.Indent))
		sb.WriteString(ul.BulletChar + "  ")
		text := termenv.String(WrappedString(item.Text, ul.Indent*3+1, ul.Indent*5, size.Col(), size.Col())).Foreground(TerminalColor(item.Style.Color))
		if item.Style.IsBold {
			text = text.Bold()
		}
		if item.Style.IsItalic {
			text = text.Italic()
		}
		if item.Style.IsUnderline {
			text = text.Underline()
		}

		sb.WriteString(text.String())
	}
	return sb.String()
}
