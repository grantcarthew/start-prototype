package cli

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// PromptHelper provides utilities for interactive CLI prompts
type PromptHelper struct {
	reader *bufio.Reader
}

// NewPromptHelper creates a new prompt helper
func NewPromptHelper() *PromptHelper {
	return &PromptHelper{
		reader: bufio.NewReader(os.Stdin),
	}
}

// Ask prompts the user with a question and returns their response
func (p *PromptHelper) Ask(question string) (string, error) {
	fmt.Print(question)
	response, err := p.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(response), nil
}

// AskWithDefault prompts with a default value shown in brackets
func (p *PromptHelper) AskWithDefault(question, defaultValue string) (string, error) {
	prompt := fmt.Sprintf("%s [%s]: ", question, defaultValue)
	response, err := p.Ask(prompt)
	if err != nil {
		return "", err
	}
	if response == "" {
		return defaultValue, nil
	}
	return response, nil
}

// AskYesNo prompts for yes/no confirmation
func (p *PromptHelper) AskYesNo(question string, defaultYes bool) (bool, error) {
	var prompt string
	if defaultYes {
		prompt = fmt.Sprintf("%s [Y/n]: ", question)
	} else {
		prompt = fmt.Sprintf("%s [y/N]: ", question)
	}

	response, err := p.Ask(prompt)
	if err != nil {
		return false, err
	}

	response = strings.ToLower(response)
	if response == "" {
		return defaultYes, nil
	}

	return response == "y" || response == "yes", nil
}

// AskChoice prompts for a choice from a list of options
func (p *PromptHelper) AskChoice(question string, options []string) (string, error) {
	fmt.Println(question)
	for i, opt := range options {
		fmt.Printf("  %d) %s\n", i+1, opt)
	}
	fmt.Println()

	prompt := fmt.Sprintf("Select [1-%d]: ", len(options))
	for {
		response, err := p.Ask(prompt)
		if err != nil {
			return "", err
		}

		var choice int
		_, err = fmt.Sscanf(response, "%d", &choice)
		if err == nil && choice >= 1 && choice <= len(options) {
			return options[choice-1], nil
		}

		fmt.Printf("Invalid choice. Please enter 1-%d: ", len(options))
	}
}

// ValidateName validates a name according to the naming rules
// Pattern: lowercase alphanumeric with hyphens, e.g., "claude", "gpt-4", "my-agent"
func (p *PromptHelper) ValidateName(name string) error {
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	pattern := regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
	if !pattern.MatchString(name) {
		return fmt.Errorf("invalid name. Use lowercase alphanumeric with hyphens.\n  Examples: claude, gpt-4, my-agent")
	}

	return nil
}

// AskValidatedName prompts for a name and validates it
func (p *PromptHelper) AskValidatedName(question string) (string, error) {
	for {
		name, err := p.Ask(question)
		if err != nil {
			return "", err
		}

		if err := p.ValidateName(name); err != nil {
			fmt.Printf("✗ %v\n\n", err)
			continue
		}

		return name, nil
	}
}

// AskOptional prompts for an optional value (can be empty)
func (p *PromptHelper) AskOptional(question string) (string, error) {
	response, err := p.Ask(fmt.Sprintf("%s (optional): ", question))
	if err != nil {
		return "", err
	}
	return response, nil
}

// PrintHeader prints a section header
func (p *PromptHelper) PrintHeader(title string) {
	fmt.Println()
	fmt.Println(title)
	fmt.Println(strings.Repeat("─", 60))
	fmt.Println()
}

// PrintSuccess prints a success message
func (p *PromptHelper) PrintSuccess(message string) {
	fmt.Printf("✓ %s\n", message)
}

// PrintError prints an error message
func (p *PromptHelper) PrintError(message string) {
	fmt.Printf("✗ %s\n", message)
}

// PrintWarning prints a warning message
func (p *PromptHelper) PrintWarning(message string) {
	fmt.Printf("⚠ %s\n", message)
}
