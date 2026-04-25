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
		opt := huh.NewOption(item.Label, item.ID)
		if item.Enabled {
			opt = opt.Selected(true)
		}
		options[i] = opt
	}

	var selected []string
	field := huh.NewMultiSelect[string]().
		Title(title).
		Options(options...).
		Value(&selected)

	err := huh.NewForm(huh.NewGroup(field)).
		WithTheme(flackTheme()).
		Run()
	if err != nil {
		return nil, fmt.Errorf("selection cancelled: %w", err)
	}
	return selected, nil
}

func SingleSelect(title string, items []Selectable) (string, error) {
	options := make([]huh.Option[string], len(items))
	for i, item := range items {
		options[i] = huh.NewOption(item.Label, item.ID)
	}

	var selected string
	field := huh.NewSelect[string]().
		Title(title).
		Options(options...).
		Value(&selected)

	err := huh.NewForm(huh.NewGroup(field)).
		WithTheme(flackTheme()).
		Run()
	if err != nil {
		return "", fmt.Errorf("selection cancelled: %w", err)
	}
	return selected, nil
}

func flackTheme() *huh.Theme {
	t := huh.ThemeBase16()

	t.Focused.UnselectedOption = t.Focused.UnselectedOption.UnsetForeground()
	t.Focused.UnselectedPrefix = t.Focused.UnselectedPrefix.UnsetForeground()
	t.Focused.Option = t.Focused.Option.UnsetForeground()

	t.Blurred.UnselectedOption = t.Blurred.UnselectedOption.UnsetForeground()
	t.Blurred.UnselectedPrefix = t.Blurred.UnselectedPrefix.UnsetForeground()

	return t
}