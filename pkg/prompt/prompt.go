package prompt

import (
	"fmt"

	survey "github.com/AlecAivazis/survey/v2"
)

type Prompt struct {
	message string
	choices []string
}

func NewPrompt(message string, choices []string) *Prompt {
	return &Prompt{
		message: message,
		choices: choices,
	}
}

func (p *Prompt) PromptMultipleSelect() ([]string, error) {
	var result = []string{}

	prompt := &survey.MultiSelect{
		Message: p.message,
		Options: p.choices,
	}

	err := survey.AskOne(prompt, &result)
	if err != nil {
		return nil, fmt.Errorf("Prompt error: %s", err)
	}

	return result, nil
}

func (p *Prompt) PromptSelect() (string, error) {
	var result = ""

	prompt := &survey.Select{
		Message: p.message,
		Options: p.choices,
	}

	err := survey.AskOne(prompt, &result)
	if err != nil {
		return "", fmt.Errorf("Prompt error: %s", err)
	}

	return result, nil
}
