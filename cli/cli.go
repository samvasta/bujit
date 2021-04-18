package cli

// A simple program demonstrating the text input component from the Bubbles
// component library.

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
	"github.com/olekukonko/ts"
	"samvasta.com/bujit/actions"
	"samvasta.com/bujit/cli/customtext"
	"samvasta.com/bujit/cli/outputview"
	"samvasta.com/bujit/models"
	"samvasta.com/bujit/parse"
	"samvasta.com/bujit/session"
)

// color is a helper for returning colors.
var color func(s string) termenv.Color = termenv.ColorProfile().Color

type history struct {
	result       actions.ActionResult   // result of last command action
	consequences []*actions.Consequence // consequences of last command action
	prevCommands []string
	err          error
	exit         bool
}

func StartInteractive() {

	session := session.InMemorySession(models.MigrateSchema)

	history := history{}

	for !history.exit {
		model := initialModel(&session, &history)
		// Get user command
		p := tea.NewProgram(model)

		if err := p.Start(); err != nil {
			log.Fatal(err)
		}

		if history.err == nil {
			// print output
			sb := strings.Builder{}

			sb.WriteString(outputview.View(history.result.Output, history.consequences))
			for _, c := range history.consequences {
				json, err := json.Marshal(c.Object)
				if err == nil {
					sb.WriteString(fmt.Sprintf("\n%s: %s", c.ConsequenceType, string(json)))
				}
			}
			sb.WriteString("\n")
			if history.result.IsSuccessful {
				sb.WriteString("success!")
			} else {
				sb.WriteString("Failed")
			}
			sb.WriteString("\n")

			fmt.Println(sb.String())
		}
	}
}

type model struct {
	session    *session.Session
	textInput  customtext.Model
	suggestion parse.AutoSuggestion
	history    *history
}

func initialModel(session *session.Session, history *history) model {
	size, _ := ts.GetSize()
	ti := customtext.NewModel()
	ti.Placeholder = "help"
	ti.SuggestionColor = "10"
	ti.Focus()
	ti.CharLimit = 0
	ti.Width = size.Col() - len(ti.Prompt) - 1

	return model{
		textInput: ti,
		session:   session,
		history:   history,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			fallthrough
		case tea.KeyEsc:
			m.history.exit = true
			return m, tea.Quit
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)

	currentText := m.textInput.Value()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:

			if m.session != nil {
				action, _ := parse.ParseExpression(currentText, m.session)
				if action != nil {
					result, consequences := action.Execute()
					m.history.result = result
					m.history.consequences = consequences
				} else {
					m.history.result = actions.ActionResult{IsSuccessful: false}
					m.history.consequences = []*actions.Consequence{}
				}
			}

			// have to turn off cursor so the view doesn't include it when we save to history
			m.textInput.SetCursorMode(customtext.CursorHide)
			m.history.prevCommands = append(m.history.prevCommands, m.textInput.Value())
			m.suggestion = parse.EmptySuggestions
			return m, tea.Quit
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
	sb := strings.Builder{}

	sb.WriteString(m.textInput.View())

	if !m.history.exit {
		sb.WriteString("\n")
		sb.WriteString(termenv.String(strings.Join(m.suggestion.NextArgs, "  ")).
			Foreground(color("10")).
			Background(color("")).
			String())
		sb.WriteString("\n")
		sb.WriteString("(esc to quit)")
	}
	return sb.String()
}
