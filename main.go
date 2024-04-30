package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

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
}

type printMessage string

type gasGolfResult struct {
	name      string
	selector  string
	timeTaken time.Duration
}

func initialModel() model {
	return model{
		title: "yoke CLI",
		choices: []string{
			"go gas golfing",
			"print message",
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
				return m, tea.Batch(gasGolf(m.input.Value()), m.spinner.Tick)
			}
			if m.cursor == 0 {
				m.title = "Gas Golfing\n"
				m.interactive = true
				m.input.Placeholder = "Function Selector"
				return m, tea.Batch(m.input.Focus(), textinput.Blink)
			} else {
				return m, tea.Batch(createMessage("Print Message", 2), m.spinner.Tick)
			}
		case "m":
			m.interactive = false
			m.title = "Yoke CLI"
			return m, nil
		}
	case printMessage:
		m.processing = false
		m.message = string(msg)
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case gasGolfResult:
		m.processing = false
		m.message = fmt.Sprintf("\nname:\t%s\nselector:\t%s\ntime taken:\t%v", msg.name, msg.selector, msg.timeTaken)
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) View() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%s\n", m.title))

	if m.interactive {
		writeInteractiveLayout(&sb, m)
	} else {
		writeMenuLayout(&sb, m)
	}

	sb.WriteString("\nPress q to quit or m for menu.\n")

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
		name, selector, timeTaken := golf.SearchFuncSelector(funcPattern, runtime.NumCPU())
		return gasGolfResult{
			name:      name,
			selector:  selector,
			timeTaken: timeTaken,
		}
	}
}

func createMessage(msg string, seconds int) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(time.Duration(seconds) * time.Second)
		return printMessage(msg)
	}
}

func main() {
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Printf("error occurred: %v", err)
		os.Exit(1)
	}
}
