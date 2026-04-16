package profile

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/olekukonko/tablewriter"
	tw "github.com/olekukonko/tablewriter/tw"

	"github.com/YPares/flack/internal/nix"
)

type Entry struct {
	Name        string
	Active      bool
	AttrPath    string
	OriginalURL string
	Priority    int
	URL         string
}

func List(profile string) ([]Entry, error) {
	pl, err := nix.ProfileListJSON(profile)
	if err != nil {
		return nil, err
	}
	var entries []Entry
	for name, el := range pl.Elements {
		entries = append(entries, Entry{
			Name:        name,
			Active:      el.Active,
			AttrPath:    el.AttrPath,
			OriginalURL: el.OriginalURL,
			Priority:    el.Priority,
			URL:         el.URL,
		})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Priority < entries[j].Priority
	})
	return entries, nil
}

func nextPriority(entries []Entry, step int) int {
	if len(entries) == 0 {
		return 0
	}
	return entries[0].Priority - step
}

func Push(profile, ref string, step int) error {
	entries, err := List(profile)
	if err != nil {
		return err
	}
	prio := nextPriority(entries, step)
	return nix.ProfileAdd(profile, ref, prio)
}

func Pop(profile string) error {
	entries, err := List(profile)
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		return fmt.Errorf("profile is empty")
	}
	return nix.ProfileRemove(profile, entries[0].Name)
}

func Upgrade(profile string, names []string, refresh bool) error {
	return nix.ProfileUpgrade(profile, names, refresh)
}

func FormatEntries(entries []Entry) string {
	if len(entries) == 0 {
		return "(empty profile)"
	}
	var b strings.Builder
	t := tablewriter.NewTable(&b,
		tablewriter.WithHeader([]string{"NAME", "PRIORITY", "LOCKED"}),
		tablewriter.WithRendition(tw.Rendition{
			Borders:  tw.BorderNone,
			Symbols:  tw.NewSymbols(tw.StyleNone),
			Settings: tw.Settings{Separators: tw.SeparatorsNone, Lines: tw.LinesNone},
		}),
	)
	for _, e := range entries {
		t.Append(e.Name, fmt.Sprintf("%d", e.Priority), e.URL)
	}
	t.Render()
	return b.String()
}

type EntryJSON struct {
	Name        string `json:"name"`
	Active      bool   `json:"active"`
	AttrPath    string `json:"attrPath"`
	OriginalURL string `json:"originalUrl"`
	Priority    int    `json:"priority"`
	URL         string `json:"url"`
}

func EntriesToJSON(entries []Entry) (string, error) {
	out := make([]EntryJSON, len(entries))
	for i, e := range entries {
		out[i] = EntryJSON{
			Name:        e.Name,
			Active:      e.Active,
			AttrPath:    e.AttrPath,
			OriginalURL: e.OriginalURL,
			Priority:    e.Priority,
			URL:         e.URL,
		}
	}
	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}