package prompt

import (
	survey "github.com/AlecAivazis/survey/v2"
)

type MultipleSelect struct {
	message string
	choices []string
}

func NewMultipleSelectPrompt(message string, choices []string) *MultipleSelect {
	return &MultipleSelect{
		message: message,
		choices: choices,
	}
}

func (p *MultipleSelect) Prompt() (result []string) {
	result = []string{}

	prompt := &survey.MultiSelect{
		Message: p.message,
		Options: p.choices,
	}

	survey.AskOne(prompt, &result)

	return result
}
