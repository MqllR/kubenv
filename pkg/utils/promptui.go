package utils

import (
	"sort"
	"strings"

	"github.com/manifoldco/promptui"
)

// Prompt use manifoldco/promptui to select an item of an array
func Prompt(label string, items []string) (string, error) {
	sort.Strings(items)

	searcher := func(input string, index int) bool {
		item := items[index]
		name := strings.Replace(strings.ToLower(item), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:    label,
		Items:    items,
		Searcher: searcher,
	}

	_, selectedItem, err := prompt.Run()

	if err != nil {
		return "", err
	}

	return selectedItem, nil
}
