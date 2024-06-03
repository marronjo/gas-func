package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/marronjo/yoke/address"
	"github.com/marronjo/yoke/golf"
)

type model struct {
	title       string
	choices     []string
	cursor      int
	message     string
	processing  bool
	spinner     spinner.Model
	input       textinput.Model
	interactive bool
	err         error
}

type addressResult struct {
	address string
}

type gasGolfResult struct {
	name      string
	selector  string
	timeTaken time.Duration
}

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

func initialModel() model {
	return model{
		title: "yoke CLI",
		choices: []string{
			"go gas golfing",
			"generate address",
		},
		spinner: spinner.New(),
		input:   textinput.New(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.processing = false
			return m, tea.Quit
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			if m.interactive {
				m.processing = true
				if m.cursor == 0 {
					return m, tea.Batch(gasGolf(m.input.Value()), m.spinner.Tick)
				} else {
					return m, tea.Batch(generateAddress(m.input.Value()), m.spinner.Tick)
				}
			}
			if m.cursor == 0 {
				m.title = "Gas Golfing\n"
				m.interactive = true
				m.input.Placeholder = "Function Selector"
				return m, tea.Batch(m.input.Focus(), textinput.Blink)
			} else {
				m.title = "Address Generator\n"
				m.interactive = true
				m.input.Placeholder = "Private Key"
				return m, tea.Batch(m.input.Focus(), textinput.Blink)
			}
		case "esc":
			m.input.Reset()
			m.message = ""
			m.interactive = false
			m.title = "Yoke CLI"
			return m, nil
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case gasGolfResult:
		m.processing = false
		m.message = fmt.Sprintf("\nname:\t%s\nselector:\t%s\ntime taken:\t%v", msg.name, msg.selector, msg.timeTaken)
	case addressResult:
		m.processing = false
		m.message = fmt.Sprintf("\naddress:\t%s", msg.address)
	case errMsg:
		m.err = msg
		return m, tea.Quit
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) View() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%s\n", m.title))

	if m.err != nil {
		return fmt.Sprintf("\nError: %v\n\n", m.err)
	}

	if m.interactive {
		writeInteractiveLayout(&sb, m)
	} else {
		writeMenuLayout(&sb, m)
	}

	sb.WriteString("\nPress q to quit or esc for menu.\n")

	return sb.String()
}

func writeInteractiveLayout(sb *strings.Builder, m model) {
	sb.WriteString(m.input.View() + "\n")

	if m.processing {
		sb.WriteString("\n" + m.spinner.View())
	}

	if m.message != "" {
		sb.WriteString(fmt.Sprintf("%s\n", m.message))
	}
}

func writeMenuLayout(sb *strings.Builder, m model) {
	sb.WriteString("What would you like to do ?\n\n")
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = "*"
		}
		sb.WriteString(fmt.Sprintf("%s %s\n", cursor, choice))
	}
}

func gasGolf(funcPattern string) tea.Cmd {
	return func() tea.Msg {
		result, err := golf.SearchFuncSelector(funcPattern)
		if err != nil {
			return errMsg{err}
		}
		return gasGolfResult{
			name:      result.Name,
			selector:  result.Selector,
			timeTaken: result.TimeTaken,
		}
	}
}

func generateAddress(privateKey string) tea.Cmd {
	return func() tea.Msg {
		result, err := address.GenerateAddressFromPrivateKey(privateKey)
		if err != nil {
			return errMsg{err}
		}
		return addressResult{
			address: result,
		}
	}
}

func main() {
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Printf("error occurred: %v", err)
		os.Exit(1)
	}
}
