package cli

// A simple program demonstrating the text input component from the Bubbles
// component library.

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
	"github.com/olekukonko/ts"
	"samvasta.com/bujit/actions"
	"samvasta.com/bujit/cli/customtext"
	"samvasta.com/bujit/models"
	"samvasta.com/bujit/parse"
	"samvasta.com/bujit/session"
)

// color is a helper for returning colors.
var color func(s string) termenv.Color = termenv.ColorProfile().Color

func StartInteractive() {
	p := tea.NewProgram(initialModel())

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}

type tickMsg struct{}
type errMsg error

type model struct {
	textInput    customtext.Model
	suggestion   parse.AutoSuggestion
	result       actions.ActionResult
	consequences []*actions.Consequence
	session      *session.Session
	history      []string        // previous commands
	historyLog   strings.Builder // output from previous commands
	err          error
}

func initialModel() model {
	size, _ := ts.GetSize()
	ti := customtext.NewModel()
	ti.Placeholder = "help"
	ti.SuggestionColor = "10"
	ti.Focus()
	ti.CharLimit = 0
	ti.Width = size.Col() - len(ti.Prompt) - 1

	session := session.InMemorySession(models.MigrateSchema)
	return model{
		textInput: ti,
		session:   &session,
		err:       nil,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *model) pushHistory() {
	m.textInput.SetCursorMode(customtext.CursorHide)
	m.history = append(m.history, m.textInput.Value())
	m.historyLog.WriteString(m.textInput.View())
	m.historyLog.WriteString("\n")
	m.textInput.SetCursorMode(customtext.CursorBlink)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			fallthrough
		case tea.KeyEsc:
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)

	currentText := m.textInput.Value()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.pushHistory()
			if m.session != nil {
				action, _ := parse.ParseExpression(currentText, m.session)
				if action != nil {
					result, consequences := action.Execute()
					m.result = result
					m.consequences = consequences
					m.historyLog.WriteString(result.Output)
					m.historyLog.WriteString("\n")
					if result.IsSuccessful {
						m.historyLog.WriteString("success!")
					} else {
						m.historyLog.WriteString("Failed")
					}
					m.historyLog.WriteString("\n")
				}
			}
			m.textInput.Reset()
			m.suggestion = parse.EmptySuggestions
			// m.history.WriteString("\n")
			// m.history.WriteString(result.Output)
			// m.history.WriteString("\n\n")
			return m, cmd
		default:
			m.textInput.SetSuggestValue("")
			if len(currentText) > 0 {
				_, suggestion := parse.ParseExpression(currentText, m.session)
				m.suggestion = suggestion
				for _, next := range suggestion.NextArgs {
					if strings.HasPrefix(next, currentText) {
						m.textInput.SetSuggestValue(strings.TrimPrefix(next, currentText))
						break
					}
				}
			} else {
				m.suggestion = parse.EmptySuggestions
			}
		}
	}

	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s\n\n%s\n\n%s\n\n%s",
		m.historyLog.String(),
		m.textInput.View(),
		termenv.String(strings.Join(m.suggestion.NextArgs, "  ")).
			Foreground(color("10")).
			Background(color("")).
			String(),
		"(esc to quit)",
	) + "\n"
}
