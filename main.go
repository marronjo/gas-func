package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
	message  string
}

type printMessage string

func initialModel() model {
	return model{
		choices: []string{
			"go gas golfing",
			"print message",
		},
		selected: make(map[int]struct{}),
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
				if m.cursor == 0 {
					return m, createMessage("Gas Golfing", 4)
				} else {
					return m, createMessage("Print Message", 2)
				}
			}
		}
	case printMessage:
		m.message = string(msg)
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

	if m.message != "" {
		s += fmt.Sprintf("Message : [%s]\n", m.message)
	}

	s += "\nPress q or ctrl+c to quit.\n"

	return s
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
