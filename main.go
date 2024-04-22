package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/marronjo/yoke/golf"
)

type model struct {
	choices    []string
	cursor     int
	selected   map[int]struct{}
	message    string
	processing bool
	spinner    spinner.Model
}

type printMessage string

type gasGolfResult struct {
	name      string
	selector  string
	timeTaken time.Duration
}

func initialModel() model {
	return model{
		choices: []string{
			"go gas golfing",
			"print message",
		},
		selected: make(map[int]struct{}),
		spinner:  spinner.New(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
				m.processing = true
				if m.cursor == 0 {
					return m, tea.Batch(gasGolf("mint%d(uint256,address)"), m.spinner.Tick)
				} else {
					return m, tea.Batch(createMessage("Print Message", 2), m.spinner.Tick)
				}
			}
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
		m.message = fmt.Sprintf("name: %s\nselector: %s\ntime taken: %v\n", msg.name, msg.selector, msg.timeTaken)
	}
	return m, nil
}

func (m model) View() string {
	s := "What would you like to do ?\n\n"
	for i, choice := range m.choices {

		cursor := "  "
		if m.cursor == i {
			cursor = "->"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	if m.processing {
		s += m.spinner.View()
	}

	if m.message != "" {
		s += fmt.Sprintf("\n%s\n", m.message)
	}

	s += "\nPress q or ctrl+c to quit.\n"

	return s
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
