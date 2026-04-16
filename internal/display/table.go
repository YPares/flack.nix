package display

import (
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	tw "github.com/olekukonko/tablewriter/tw"
	"golang.org/x/term"
)

func termWidth() int {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || w <= 0 {
		return 80
	}
	return w
}

type Table struct {
	headers []string
	rows    [][]string
}

func NewTable(headers []string) *Table {
	return &Table{headers: headers}
}

func (t *Table) Row(cells ...string) {
	t.rows = append(t.rows, cells)
}

func (t *Table) Render() string {
	if len(t.rows) == 0 {
		return ""
	}
	var b strings.Builder
	opts := []tablewriter.Option{
		tablewriter.WithHeader(t.headers),
		tablewriter.WithRendition(tw.Rendition{
			Borders:  tw.BorderNone,
			Symbols:  tw.NewSymbols(tw.StyleNone),
			Settings: tw.Settings{Separators: tw.SeparatorsNone, Lines: tw.LinesNone},
		}),
		tablewriter.WithRowAutoWrap(tw.WrapBreak),
		tablewriter.WithHeaderAutoWrap(tw.WrapNormal),
		tablewriter.WithMaxWidth(termWidth()),
	}
	tbl := tablewriter.NewTable(&b, opts...)
	for _, row := range t.rows {
		vals := make([]any, len(row))
		for i, v := range row {
			vals[i] = v
		}
		tbl.Append(vals...)
	}
	tbl.Render()
	return b.String()
}