package tui

import (
	"fmt"

	"github.com/charmbracelet/huh"
)

type Selectable struct {
	ID      string
	Label   string
	Enabled bool
}

func MultiSelect(title string, items []Selectable) ([]string, error) {
	options := make([]huh.Option[string], len(items))
	for i, item := range items {
		options[i] = huh.NewOption(item.Label, item.ID)
	}

	var selected []string
	field := huh.NewMultiSelect[string]().
		Title(title).
		Options(options...).
		Value(&selected)

	err := huh.NewForm(huh.NewGroup(field)).Run()
	if err != nil {
		return nil, fmt.Errorf("selection cancelled: %w", err)
	}
	return selected, nil
}