package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/net/html"
)

const (
	BASE_URI                       = "https://www.php.net/manual/en"
	STRING_FUNCTIONS FunctionGroup = "strings"
	ARRAY_FUNCTIONS  FunctionGroup = "arrays"
)

type FunctionGroup string
type model struct {
	choices []string
}

func getFunctionGroup(funcType FunctionGroup) []string {
	var functions []string
	resp, err := http.Get(fmt.Sprintf("%s/ref.%s.php", BASE_URI, funcType))
	if err != nil {
		panic("Error connecting to PHP.net")
	}
	defer resp.Body.Close()

	tokenizer := html.NewTokenizer(resp.Body)
	for {
		tokenType := tokenizer.Next()

		if tokenType == html.ErrorToken {
			if err := tokenizer.Err(); err == io.EOF {
				break
			}
			log.Fatalf("Error tokenizing html: %v", tokenizer.Err())
		}

		if tokenType == html.StartTagToken {
			token := tokenizer.Token()
			if token.Data == "ul" {
				for _, t := range token.Attr {
					if t.Val == "chunklist chunklist_reference" {
						// we have found the functions reference list
					}
				}
			}
		}
	}

	return functions
}

func getPhpFunctions() []string {
	stringFuncs := getFunctionGroup(STRING_FUNCTIONS)
	arrayFuncs := getFunctionGroup(ARRAY_FUNCTIONS)
	return append(stringFuncs, arrayFuncs...)
}

func initialModel() model {
	return model{
		choices: getPhpFunctions(),
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
		}
	}

	return m, nil
}

func (m model) View() string {
	s := ""
	for _, choice := range m.choices {
		s += fmt.Sprintf("%s\n", choice)
	}
	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
