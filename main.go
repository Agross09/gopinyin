package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Word represents a vocabulary entry
type Word struct {
	Pinyin      string
	Chinese     string
	Definition  string
	Example     string
}

// Model represents the application state
type model struct {
	words           []Word
	currentIndex    int
	showDetails     bool
	addingNewCard   bool
	inputs          []textinput.Model
	focusIndex      int
}

// Initial word list (placeholder for now)
var wordList = []Word{
	{
		Pinyin:     "ni hao",
		Chinese:    "你好",
		Definition: "Hello",
		Example:    "Ni hao, how are you?",
	},
	{
		Pinyin:     "xie xie",
		Chinese:    "谢谢",
		Definition: "Thank you",
		Example:    "Xie xie for your help.",
	},
	{
		Pinyin:     "zao",
		Chinese:    "早",
		Definition: "Morning",
		Example:    "Zao, good morning!",
	},
	{
		Pinyin:     "pengyou",
		Chinese:    "朋友",
		Definition: "Friend",
		Example:    "Wo de pengyou hen hao.",
	},
	{
		Pinyin:     "chi fan",
		Chinese:    "吃饭",
		Definition: "Eat meal",
		Example:    "Women qu chi fan.",
	},
	{
		Pinyin:     "hao",
		Chinese:    "好",
		Definition: "Good",
		Example:    "Hen hao, that's good!",
	},
	{
		Pinyin:     "shui",
		Chinese:    "水",
		Definition: "Water",
		Example:    "Wo yao yi bei shui.",
	},
	{
		Pinyin:     "ai",
		Chinese:    "爱",
		Definition: "Love",
		Example:    "Wo ai ni means I love you.",
	},
	{
		Pinyin:     "ren",
		Chinese:    "人",
		Definition: "Person",
		Example:    "Mei ge ren dou bu tong.",
	},
	{
		Pinyin:     "jia",
		Chinese:    "家",
		Definition: "Home/Family",
		Example:    "Wo de jia zai Beijing.",
	},
}

// Initialize the model
func initialModel() model {
	m := model{
		words:         wordList,
		currentIndex:  0,
		showDetails:   false,
		addingNewCard: false,
		inputs:        make([]textinput.Model, 4),
	}

	// Create text inputs for adding a new card
	for i := range m.inputs {
		t := textinput.New()
		t.Placeholder = []string{"Chinese Characters", "Pinyin", "Definition", "Example"}[i]
		t.Focus()

		switch i {
		case 0:
			t.Placeholder = "Chinese Characters (e.g. 你好)"
		case 1:
			t.Placeholder = "Pinyin (e.g. ni hao)"
		case 2:
			t.Placeholder = "Definition (e.g. Hello)"
		case 3:
			t.Placeholder = "Example Sentence (optional)"
		}

		m.inputs[i] = t
	}

	return m
}

// Update method handles user input and state changes
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// Handle adding new card state
	if m.addingNewCard {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyCtrlC, tea.KeyEsc:
				m.addingNewCard = false
				return m, nil

			case tea.KeyTab, tea.KeyShiftTab:
				// Change focus
				if msg.Type == tea.KeyTab {
					m.focusIndex++
				} else {
					m.focusIndex--
				}

				if m.focusIndex > len(m.inputs)-1 {
					m.focusIndex = 0
				} else if m.focusIndex < 0 {
					m.focusIndex = len(m.inputs) - 1
				}

				for i := 0; i < len(m.inputs); i++ {
					if i == m.focusIndex {
						m.inputs[i].Focus()
					} else {
						m.inputs[i].Blur()
					}
				}

			case tea.KeyEnter:
				// Save new card if at least Chinese, Pinyin, and Definition are filled
				if m.inputs[0].Value() != "" && m.inputs[1].Value() != "" && m.inputs[2].Value() != "" {
					newWord := Word{
						Chinese:     m.inputs[0].Value(),
						Pinyin:      m.inputs[1].Value(),
						Definition:  m.inputs[2].Value(),
						Example:     m.inputs[3].Value(),
					}
					m.words = append(m.words, newWord)
					m.addingNewCard = false
					m.currentIndex = len(m.words) - 1
					
					// Reset inputs
					for i := range m.inputs {
						m.inputs[i].Reset()
					}
					m.focusIndex = 0
				}
			}

			// Handle text input updates
			for i := range m.inputs {
				m.inputs[i], cmd = m.inputs[i].Update(msg)
				cmds = append(cmds, cmd)
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Normal navigation state
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "right", "l":
			m.currentIndex = (m.currentIndex + 1) % len(m.words)
			m.showDetails = false

		case "left", "h":
			m.currentIndex = (m.currentIndex - 1 + len(m.words)) % len(m.words)
			m.showDetails = false

		case " ", "enter":
			m.showDetails = !m.showDetails

		case "a":
			// Enter add card mode
			m.addingNewCard = true
			m.focusIndex = 0
			for i := range m.inputs {
				if i == 0 {
					m.inputs[i].Focus()
				} else {
					m.inputs[i].Blur()
				}
			}
		}
	}

	return m, nil
}

// View method renders the UI
func (m model) View() string {
	// Adding new card view
	if m.addingNewCard {
		s := "Add New Vocabulary Card\n\n"
		
		for _, input := range m.inputs {
			s += fmt.Sprintf("%s\n", input.View())
		}

		s += "\nTAB: Next field\n"
		s += "ENTER: Save card\n"
		s += "ESC: Cancel\n"
		
		return s
	}

	// Normal vocabulary view (similar to previous version)
	if len(m.words) == 0 {
		return "No words available.\n"
	}

	currentWord := m.words[m.currentIndex]
	
	s := fmt.Sprintf("Pinyin Vocab Flashcards (%d/%d)\n\n", m.currentIndex+1, len(m.words))
	s += fmt.Sprintf("Pinyin: %s\n", currentWord.Pinyin)
	s += fmt.Sprintf("Chinese: %s\n", currentWord.Chinese)
	
	if m.showDetails {
		s += fmt.Sprintf("\nDefinition: %s\n", currentWord.Definition)
		s += fmt.Sprintf("Example: %s\n", currentWord.Example)
	} else {
		s += "\n(Press SPACE or ENTER to show details)\n"
	}

	s += "\n\nControls:\n"
	s += "← / h : Previous word\n"
	s += "→ / l : Next word\n"
	s += "SPACE : Toggle details\n"
	s += "a     : Add new card\n"
	s += "q     : Quit\n"

	return s
}

// Initialize method
func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
