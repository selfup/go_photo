package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type state int

const (
	stateMenu state = iota
	stateInput
	stateCopying
	stateDone
)

type model struct {
	state       state
	presets     []Preset
	cursor      int
	wipe        bool
	progress    progress.Model
	copyProgress CopyProgress
	inputs      []textinput.Model
	inputFocus  int
	isNewPreset bool
	err         error
	src         string
	dst         string
	presetName  string
}

type progressMsg CopyProgress

var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205"))

	selectedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("86"))

	dimStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))

	errorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("196"))
)

func initialModel() model {
	presets, _ := LoadPresets()

	p := progress.New(progress.WithDefaultGradient())

	nameInput := textinput.New()
	nameInput.Placeholder = "Preset name (leave empty for one-off)"
	nameInput.CharLimit = 32

	srcInput := textinput.New()
	srcInput.Placeholder = "Source path"
	srcInput.CharLimit = 256

	dstInput := textinput.New()
	dstInput.Placeholder = "Destination path"
	dstInput.CharLimit = 256

	return model{
		state:    stateMenu,
		presets:  presets,
		progress: p,
		inputs:   []textinput.Model{nameInput, srcInput, dstInput},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)

	case progressMsg:
		m.copyProgress = CopyProgress(msg)

		if m.copyProgress.Done {
			m.state = stateDone
		}

		return m, nil

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - 10

		if m.progress.Width > 50 {
			m.progress.Width = 50
		}

		return m, nil
	}

	if m.state == stateInput {
		return m.updateInputs(msg)
	}

	return m, nil
}

func (m model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.state {
	case stateMenu:
		return m.handleMenuKey(msg)
	case stateInput:
		return m.handleInputKey(msg)
	case stateCopying:
		if msg.String() == "esc" {
			return m, tea.Quit
		}
	case stateDone:
		if msg.String() == "enter" || msg.String() == " " {
			m.state = stateMenu
			m.copyProgress = CopyProgress{}
			m.err = nil
		}

		if msg.String() == "q" {
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) handleMenuKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	menuLen := len(m.presets) + 2

	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}

	case "down", "j":
		if m.cursor < menuLen-1 {
			m.cursor++
		}

	case "w":
		m.wipe = !m.wipe

	case "d":
		if m.cursor < len(m.presets) {
			name := m.presets[m.cursor].Name
			DeletePreset(name)
			m.presets, _ = LoadPresets()

			if m.cursor >= len(m.presets) && m.cursor > 0 {
				m.cursor--
			}
		}

	case "enter":
		return m.handleMenuSelect()
	}

	return m, nil
}

func (m model) handleMenuSelect() (tea.Model, tea.Cmd) {
	if m.cursor < len(m.presets) {
		preset := m.presets[m.cursor]
		m.src = preset.Src
		m.dst = preset.Dst
		m.presetName = preset.Name

		return m.startCopy()
	}

	if m.cursor == len(m.presets) {
		m.isNewPreset = true
		m.state = stateInput
		m.inputs[0].Focus()
		m.inputFocus = 0

		return m, textinput.Blink
	}

	if m.cursor == len(m.presets)+1 {
		m.isNewPreset = false
		m.state = stateInput
		m.inputs[1].Focus()
		m.inputFocus = 1

		return m, textinput.Blink
	}

	return m, nil
}

func (m model) handleInputKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.state = stateMenu
		m = m.withResetInputs()

		return m, nil

	case "tab", "down":
		return m.nextInput()

	case "shift+tab", "up":
		return m.prevInput()

	case "enter":
		if m.inputFocus == len(m.inputs)-1 {
			return m.submitInputs()
		}

		return m.nextInput()
	}

	return m, nil
}

func (m model) nextInput() (tea.Model, tea.Cmd) {
	start := 0

	if !m.isNewPreset {
		start = 1
	}

	m.inputs[m.inputFocus].Blur()
	m.inputFocus++

	if m.inputFocus >= len(m.inputs) {
		m.inputFocus = start
	}

	m.inputs[m.inputFocus].Focus()

	return m, textinput.Blink
}

func (m model) prevInput() (tea.Model, tea.Cmd) {
	start := 0

	if !m.isNewPreset {
		start = 1
	}

	m.inputs[m.inputFocus].Blur()
	m.inputFocus--

	if m.inputFocus < start {
		m.inputFocus = len(m.inputs) - 1
	}

	m.inputs[m.inputFocus].Focus()

	return m, textinput.Blink
}

func (m model) submitInputs() (tea.Model, tea.Cmd) {
	m.src = m.inputs[1].Value()
	m.dst = m.inputs[2].Value()
	m.presetName = m.inputs[0].Value()

	if m.src == "" || m.dst == "" {
		m.err = fmt.Errorf("source and destination paths are required")

		return m, nil
	}

	if m.isNewPreset && m.presetName != "" {
		AddPreset(Preset{
			Name: m.presetName,
			Src:  m.src,
			Dst:  m.dst,
		})

		m.presets, _ = LoadPresets()
	}

	m = m.withResetInputs()

	return m.startCopy()
}

func (m model) startCopy() (tea.Model, tea.Cmd) {
	m.state = stateCopying

	progressChan := make(chan CopyProgress)

	go runWithProgress(m.src, m.dst, m.wipe, progressChan)

	return m, func() tea.Msg {
		for p := range progressChan {
			return progressMsg(p)
		}

		return nil
	}
}

func (m model) withResetInputs() model {
	for i := range m.inputs {
		m.inputs[i].SetValue("")
		m.inputs[i].Blur()
	}

	m.inputFocus = 0

	return m
}

func (m model) updateInputs(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.inputs[m.inputFocus], cmd = m.inputs[m.inputFocus].Update(msg)

	return m, cmd
}

func (m model) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("go_photo"))
	b.WriteString("\n\n")

	switch m.state {
	case stateMenu:
		b.WriteString(m.viewMenu())
	case stateInput:
		b.WriteString(m.viewInput())
	case stateCopying:
		b.WriteString(m.viewCopying())
	case stateDone:
		b.WriteString(m.viewDone())
	}

	return b.String()
}

func (m model) viewMenu() string {
	var b strings.Builder

	b.WriteString("Presets:\n")

	for i, preset := range m.presets {
		cursor := "  "

		if m.cursor == i {
			cursor = "> "
			b.WriteString(selectedStyle.Render(cursor + preset.Name))
		} else {
			b.WriteString(cursor + preset.Name)
		}

		b.WriteString("\n")
	}

	newPresetIdx := len(m.presets)
	customIdx := len(m.presets) + 1

	if m.cursor == newPresetIdx {
		b.WriteString(selectedStyle.Render("> [+ New preset]"))
	} else {
		b.WriteString("  [+ New preset]")
	}

	b.WriteString("\n")

	if m.cursor == customIdx {
		b.WriteString(selectedStyle.Render("> [Custom paths...]"))
	} else {
		b.WriteString("  [Custom paths...]")
	}

	b.WriteString("\n\n")

	wipeStatus := "[ ]"

	if m.wipe {
		wipeStatus = "[x]"
	}

	b.WriteString(fmt.Sprintf("%s Delete source after copy (w to toggle)\n", wipeStatus))
	b.WriteString("\n")
	b.WriteString(dimStyle.Render("Enter: Select  d: Delete preset  q: Quit"))

	return b.String()
}

func (m model) viewInput() string {
	var b strings.Builder

	if m.isNewPreset {
		b.WriteString("New Preset:\n\n")
		b.WriteString("Name: ")
		b.WriteString(m.inputs[0].View())
		b.WriteString("\n\n")
	} else {
		b.WriteString("Custom Import:\n\n")
	}

	b.WriteString("Source:      ")
	b.WriteString(m.inputs[1].View())
	b.WriteString("\n\n")
	b.WriteString("Destination: ")
	b.WriteString(m.inputs[2].View())
	b.WriteString("\n\n")

	if m.err != nil {
		b.WriteString(errorStyle.Render(m.err.Error()))
		b.WriteString("\n\n")
	}

	b.WriteString(dimStyle.Render("Tab: Next field  Enter: Submit  Esc: Cancel"))

	return b.String()
}

func (m model) viewCopying() string {
	var b strings.Builder

	title := "Importing"

	if m.presetName != "" {
		title = "Importing " + m.presetName
	}

	b.WriteString(title + "\n\n")

	if m.copyProgress.File != "" {
		action := "Copying"

		if m.wipe {
			action = "Moving"
		}

		b.WriteString(fmt.Sprintf("%s: %s\n", action, m.copyProgress.File))
	}

	percent := 0.0

	if m.copyProgress.Total > 0 {
		percent = float64(m.copyProgress.Current) / float64(m.copyProgress.Total)
	}

	b.WriteString(m.progress.ViewAs(percent))
	b.WriteString(fmt.Sprintf(" %d/%d\n\n", m.copyProgress.Current, m.copyProgress.Total))
	b.WriteString(dimStyle.Render("Esc: Cancel"))

	return b.String()
}

func (m model) viewDone() string {
	var b strings.Builder

	if m.copyProgress.Error != nil {
		b.WriteString(errorStyle.Render("Error: " + m.copyProgress.Error.Error()))
		b.WriteString("\n\n")
	} else {
		action := "Copied"

		if m.wipe {
			action = "Moved"
		}

		b.WriteString(fmt.Sprintf("Done! %s %d files.\n\n", action, m.copyProgress.Current))
	}

	b.WriteString(dimStyle.Render("Enter: Back to menu  q: Quit"))

	return b.String()
}

func runTUI() error {
	p := tea.NewProgram(initialModel())

	_, err := p.Run()

	return err
}
