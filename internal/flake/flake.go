package flake

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	tw "github.com/olekukonko/tablewriter/tw"

	"github.com/YPares/flack/internal/nix"
)

type InputInfo struct {
	Name        string
	Original    string
	Locked      string
	Rev         string
	LastModTime time.Time
	Type        string
}

func Inputs(flakePath string) ([]InputInfo, error) {
	meta, err := nix.FlakeMetadataJSON(flakePath)
	if err != nil {
		return nil, err
	}
	root, ok := meta.Locks.Nodes[meta.Locks.Root]
	if !ok {
		return nil, fmt.Errorf("root node not found in lockfile")
	}
	rootInputs, ok := root.Inputs.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("root node has no inputs")
	}
	var inputs []InputInfo
	for name, ref := range rootInputs {
		inputName, ok := ref.(string)
		if !ok {
			continue
		}
		node, ok := meta.Locks.Nodes[inputName]
		if !ok {
			continue
		}
		info := InputInfo{Name: name}
		if node.Original != nil {
			info.Original = formatOriginal(node.Original)
			info.Type = node.Original.Type
		}
		if node.Locked != nil {
			info.Locked = formatLocked(node.Locked)
			info.Rev = node.Locked.Rev
			if node.Locked.LastModified != 0 {
				info.LastModTime = time.Unix(node.Locked.LastModified, 0)
			}
			if info.Type == "" {
				info.Type = node.Locked.Type
			}
		}
		inputs = append(inputs, info)
	}
	slices.SortFunc(inputs, func(a, b InputInfo) int {
		return strings.Compare(a.Name, b.Name)
	})
	return inputs, nil
}

func formatOriginal(o *struct {
	Owner string `json:"owner,omitempty"`
	Repo  string `json:"repo,omitempty"`
	Type  string `json:"type,omitempty"`
	URL   string `json:"url,omitempty"`
	Ref   string `json:"ref,omitempty"`
}) string {
	if o.Owner != "" && o.Repo != "" {
		s := fmt.Sprintf("github:%s/%s", o.Owner, o.Repo)
		if o.Ref != "" {
			s += "/" + o.Ref
		}
		return s
	}
	if o.URL != "" {
		return o.URL
	}
	return fmt.Sprintf("%s:…", o.Type)
}

func formatLocked(l *struct {
	LastModified int64  `json:"lastModified"`
	NarHash     string `json:"narHash"`
	Owner       string `json:"owner,omitempty"`
	Repo        string `json:"repo,omitempty"`
	Rev         string `json:"rev"`
	Type        string `json:"type"`
	URL         string `json:"url,omitempty"`
}) string {
	if l.Owner != "" && l.Repo != "" {
		return fmt.Sprintf("github:%s/%s", l.Owner, l.Repo)
	}
	if l.URL != "" {
		return l.URL
	}
	return l.Type
}

func Update(flakePath string, inputs []string, refresh bool) error {
	return nix.FlakeUpdate(flakePath, inputs, refresh)
}

type InputJSON struct {
	Name        string `json:"name"`
	Original    string `json:"original"`
	Locked      string `json:"locked"`
	Rev         string `json:"rev"`
	Type        string `json:"type"`
	LastUpdated string `json:"lastUpdated"`
}

func InputsToJSON(inputs []InputInfo) (string, error) {
	out := make([]InputJSON, len(inputs))
	for i, inp := range inputs {
		lastUpdated := ""
		if !inp.LastModTime.IsZero() {
			lastUpdated = inp.LastModTime.Format(time.DateOnly)
		}
		out[i] = InputJSON{
			Name:        inp.Name,
			Original:    inp.Original,
			Locked:      inp.Locked,
			Rev:         inp.Rev,
			Type:        inp.Type,
			LastUpdated: lastUpdated,
		}
	}
	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func FormatInputs(inputs []InputInfo) string {
	if len(inputs) == 0 {
		return "(no inputs)"
	}
	var b strings.Builder
	t := tablewriter.NewTable(&b,
		tablewriter.WithHeader([]string{"INPUT", "SOURCE", "REV", "UPDATED"}),
		tablewriter.WithRendition(tw.Rendition{
			Borders:  tw.BorderNone,
			Symbols:  tw.NewSymbols(tw.StyleNone),
			Settings: tw.Settings{Separators: tw.SeparatorsNone, Lines: tw.LinesNone},
		}),
	)
	for _, inp := range inputs {
		lastUpdated := ""
		if !inp.LastModTime.IsZero() {
			lastUpdated = inp.LastModTime.Format(time.DateOnly)
		}
		rev := inp.Rev
		if len(rev) > 12 {
			rev = rev[:12]
		}
		t.Append(inp.Name, inp.Original, rev, lastUpdated)
	}
	t.Render()
	return b.String()
}