package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	Words "chinese_vocab/words"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joho/godotenv"
)

// Styling variables
var (
	// Color palette
	colorPrimary          = lipgloss.Color("#2C7BB6") // Soft blue
	colorSecondary        = lipgloss.Color("#D7191C") // Warm red
	colorAccent           = lipgloss.Color("#1A9641") // Dark green
	colorBackground       = lipgloss.Color("#F7F7F7") // Light gray background
	colorText             = lipgloss.Color("#333333") // Dark gray text
	transparentBackground = lipgloss.Color("transparent")

	// Styles
	titleStyle = lipgloss.NewStyle().
			Foreground(colorPrimary).
			Bold(true).
			Padding(0, 1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(colorAccent).
			Italic(true)

	subtitleRedStyle = lipgloss.NewStyle().
				Foreground(colorSecondary).
				Italic(true)

	subtitleDarkStyle = lipgloss.NewStyle().
				Foreground(colorPrimary).
				Italic(true)

	exampleTextStyle = lipgloss.NewStyle()

	cardStyle = lipgloss.NewStyle().
			Background(transparentBackground).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(colorPrimary).
			Padding(1, 2)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666")).
			Italic(true)

	inputStyle = lipgloss.NewStyle().
			Background(transparentBackground).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(colorPrimary).
			Padding(0, 1)

	focusedInputStyle = inputStyle.Copy().
				BorderForeground(colorSecondary)

	successStyle = lipgloss.NewStyle().
			Foreground(colorAccent).
			Bold(true)
)

// Model represents the application state
type model struct {
	words          []Words.Word
	currentIndex   int
	showDetails    bool
	addingNewCard  bool
	inputs         []textinput.Model
	focusIndex     int
	loadingExample bool // Indicates if we are currently loading an example
}

// ExampleMsg carries the example sentence fetched from OpenAI
type ExampleMsg struct {
	Index   int
	Example string
	Error   error
}

// Initialize the model
func initialModel() model {
	m := model{
		words:         Words.ExampleWords,
		currentIndex:  0,
		showDetails:   false,
		addingNewCard: false,
		inputs:        make([]textinput.Model, 4),
	}

	// Create text inputs for adding a new card
	for i := range m.inputs {
		t := textinput.New()
		t.Placeholder = []string{"Chinese Characters", "Pinyin", "Definition", "Example"}[i]
		t.Prompt = "» "
		t.CharLimit = 50
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
					newWord := Words.Word{
						Chinese:    m.inputs[0].Value(),
						Pinyin:     m.inputs[1].Value(),
						Definition: m.inputs[2].Value(),
						Example:    "",
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

	// Handle ExampleMsg (from OpenAI API)
	switch msg := msg.(type) {
	case ExampleMsg:
		m.loadingExample = false
		if msg.Error != nil {
			// Handle the error
			m.words[msg.Index].Example = fmt.Sprintf("Error: %v", msg.Error)
		} else {
			// Update the word's Example field
			m.words[msg.Index].Example = msg.Example
		}
		return m, nil
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
			m.loadingExample = false

		case "left", "h":
			m.currentIndex = (m.currentIndex - 1 + len(m.words)) % len(m.words)
			m.showDetails = false
			m.loadingExample = false

		case " ", "enter":
			m.showDetails = !m.showDetails
			if m.showDetails {
				m.loadingExample = true
				return m, fetchExample(m.words[m.currentIndex], m.currentIndex)
			} else {
				m.loadingExample = false
			}

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
		s := titleStyle.Render("Add New Vocabulary Card") + "\n\n"

		for i, input := range m.inputs {
			label := []string{
				"Chinese Characters:",
				"Pinyin:",
				"Definition:",
				"Example (optional):",
			}[i]

			// Render input with custom styling
			var renderedInput string
			if m.focusIndex == i {
				renderedInput = focusedInputStyle.Render(input.View())
			} else {
				renderedInput = inputStyle.Render(input.View())
			}

			s += fmt.Sprintf("%s\n%s\n",
				subtitleStyle.Render(label),
				renderedInput,
			)
		}

		s += "\n" + helpStyle.Render("TAB: Next field") + "\n"
		s += helpStyle.Render("ENTER: Save card") + "\n"
		s += helpStyle.Render("ESC: Cancel") + "\n"

		return s
	}

	// Normal vocabulary view
	if len(m.words) == 0 {
		return "No words available.\n"
	}

	currentWord := m.words[m.currentIndex]

	// Card content
	cardContent := fmt.Sprintf(
		"%s: %s\n%s: %s",
		subtitleStyle.Render("Pinyin"),
		exampleTextStyle.Render(currentWord.Pinyin),
		subtitleRedStyle.Render("Chinese"),
		exampleTextStyle.Render(currentWord.Chinese),
	)

	// Additional details
	var detailsContent string
	if m.showDetails {
		if m.loadingExample {
			detailsContent = "\n" + subtitleDarkStyle.Render("Loading example...")
		} else {
			detailsContent = fmt.Sprintf(
				"\n%s: %s\n\n%s\n%s",
				subtitleDarkStyle.Render("Definition"),
				exampleTextStyle.Render(currentWord.Definition),
				subtitleDarkStyle.Render("Example"),
				exampleTextStyle.Render(currentWord.Example),
			)
		}
	}

	// Combine everything
	view := lipgloss.NewStyle().
		Width(50).
		Align(lipgloss.Center).
		Render(
			titleStyle.Render("Pinyin Vocab Flashcards") +
				fmt.Sprintf(" (%d/%d)\n\n", m.currentIndex+1, len(m.words)) +
				cardStyle.Render(cardContent+detailsContent) +
				"\n\n" +
				helpStyle.Render("Controls:") + "\n" +
				fmt.Sprintf("← / h : Previous word \n") +
				fmt.Sprintf("→ / l : Next word     \n") +
				fmt.Sprintf("SPACE : Toggle details\n") +
				fmt.Sprintf("a     : Add new card  \n") +
				fmt.Sprintf("q     : Quit          \n"),
		)

	return view
}

// Initialize method
func (m model) Init() tea.Cmd {
	return textinput.Blink
}

// fetchExample makes an API call to OpenAI to get an example sentence
func fetchExample(word Words.Word, index int) tea.Cmd {
	return func() tea.Msg {
		apiKey := getAPIKey()

		// Create the request
		url := "https://api.openai.com/v1/chat/completions"
		model := "gpt-3.5-turbo"
		prompt := fmt.Sprintf("Give me an example phrase in Chinese, Pinyin, and English with the following word: %s", word.Chinese)

		requestBody := map[string]interface{}{
			"model": model,
			"messages": []map[string]string{
				{"role": "user", "content": prompt},
			},
		}

		requestData, err := json.Marshal(requestBody)
		if err != nil {
			return ExampleMsg{Index: index, Error: err}
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestData))
		if err != nil {
			return ExampleMsg{Index: index, Error: err}
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer  %s", apiKey))

		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			return ExampleMsg{Index: index, Error: err}
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return ExampleMsg{Index: index, Error: err}
		}

		// Parse the response
		var responseData struct {
			Choices []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			} `json:"choices"`
			Error struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		err = json.Unmarshal(body, &responseData)
		if err != nil {
			return ExampleMsg{Index: index, Error: err}
		}

		if responseData.Error.Message != "" {
			return ExampleMsg{Index: index, Error: fmt.Errorf(responseData.Error.Message)}
		}

		if len(responseData.Choices) == 0 {
			return ExampleMsg{Index: index, Error: fmt.Errorf("No choices returned")}
		}

		example := responseData.Choices[0].Message.Content

		return ExampleMsg{Index: index, Example: example}
	}
}

func getAPIKey() string {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY not set in environment variables")
	}
	return apiKey
}

func initEnv() {
	// Load .env file into environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	initEnv()
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Error running program: %v", err)
		os.Exit(1)
	}
}
