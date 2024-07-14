package engine

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Headers []string
	Subtext [5]string
	Spinner spinner.Model
	Done    bool
}

type Header string
type Subtext string
type Done bool

var (
	currentStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("65")).Bold(true)
	doneStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	subStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("65")).Faint(true).PaddingLeft(2)
)

func InitialModel() Model {
	spinner := spinner.New(
		spinner.WithSpinner(spinner.Dot),
		spinner.WithStyle(currentStyle),
	)
	return Model{
		Headers: []string{},
		Subtext: [5]string{"", "", "", "", ""},
		Spinner: spinner,
	}
}

func (model Model) Init() tea.Cmd {
	return model.Spinner.Tick
}

func (model Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return model, tea.Quit
		}
	case Header:
		model.Headers = append(model.Headers, string(msg))
		return model, nil
	case Subtext:
		model.Subtext = [5]string{
			model.Subtext[1],
			model.Subtext[2],
			model.Subtext[3],
			model.Subtext[4],
			string(msg),
		}
		return model, nil
	case Done:
		model.Done = bool(msg)
		return model, tea.Quit
	case spinner.TickMsg:
		var cmd tea.Cmd
		model.Spinner, cmd = model.Spinner.Update(msg)
		return model, cmd
	}
	return model, nil
}

func (model Model) View() string {
	view := ""

	if model.Done {
		for _, Header := range model.Headers {
			view += doneStyle.Render("✓ "+Header) + "\n"
		}
		return view
	}

	if len(model.Headers) > 1 {
		for _, Header := range model.Headers[:len(model.Headers)-1] {
			view += doneStyle.Render("✓ "+Header) + "\n"
		}
	}
	if len(model.Headers) > 0 {
		view += model.Spinner.View() + currentStyle.Render(model.Headers[len(model.Headers)-1]) + "\n"
	}
	for _, Subtext := range model.Subtext {
		view += subStyle.Render(Subtext) + "\n"
	}
	return view
}
