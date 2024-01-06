package util

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"

	"github.com/dustin/go-humanize"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/b7c/nx"
	"github.com/b7c/nx/web"
)

var tableStyle = table.Style{
	Name: "Simple",
	Box: table.BoxStyle{
		TopSeparator:    " ",
		BottomSeparator: " ",
		LeftSeparator:   " ",
		RightSeparator:  " ",
		PaddingLeft:     "\u200B",
		PaddingRight:    " ",
		MiddleVertical:  "│",
	},
	Options: table.Options{
		DrawBorder:      true,
		SeparateColumns: true,
	},
}

var red = text.Colors{text.FgRed}
var green = text.Colors{text.FgGreen}
var faint = text.Colors{text.Faint}

var dateTimeFormat = "1 January 2006 3:04:05 am"

func makeTableWriter() table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(tableStyle)
	t.SuppressTrailingSpaces()
	return t
}

type Prop struct {
	Name  string
	Value any
}

func RenderProperties(props []Prop) {
	maxNameLength := 0

	rows := make([]table.Row, 0, len(props))
	for _, prop := range props {
		maxNameLength = max(maxNameLength, len(prop.Name))
		rows = append(rows, table.Row{prop.Name, prop.Value})
	}

	lastColumnWidth := 100
	fd := int(os.Stdout.Fd())
	if term.IsTerminal(fd) {
		w, _, err := term.GetSize(fd)
		if err == nil {
			lastColumnWidth = w - maxNameLength - 5
		}
	}

	for _, prop := range props {
		value := fmt.Sprint(prop.Value)
		var lines []string
		if lastColumnWidth <= 0 {
			lines = []string{value}
		} else {
			lines = strings.Split(text.WrapText(value, lastColumnWidth), "\n")
		}

		for i, line := range lines {
			name := prop.Name
			if i > 0 {
				name = ""
			}
			fmt.Printf("%*s │ %s\n", maxNameLength+2, name, line)
		}
	}
}

func RenderFurniInfo(f *nx.FurniInfo) {
	props := []Prop{
		{"Name", f.Name},
		{"Description", f.Description},
		{"Identifier", f.Identifier},
		{"Type", f.Type},
		{"Kind", f.Kind},
		{"Revision", f.Revision},
		{"Line", f.Line},
		{"Environment", f.Environment},
		{"Category", f.Category},
		{"Default direction", f.DefaultDir},
		{"Size", fmt.Sprintf("%d x %d", f.XDim, f.YDim)},
		{"Part colors", f.PartColors},
		{"Offer ID", f.OfferId},
		{"Buyout", f.Buyout},
		{"BC", f.BC},
		{"Excluded dynamic", f.ExcludedDynamic},
		{"Custom params", f.CustomParams},
		{"Special type", f.SpecialType},
		{"Can stand on", f.CanStandOn},
		{"Can sit on", f.CanSitOn},
		{"Can lay on", f.CanLayOn},
	}
	RenderProperties(props)
}

func RenderUserInfo(u web.User) {
	props := []Prop{}
	props = append(props, Prop{"Name", u.Name})
	if u.ProfileVisible {
		if u.LastAccessTime != nil {
			status := ""
			if u.Online {
				status = green.Sprint("Online")
			} else {
				status = red.Sprint("Offline")
			}
			sinceLastAccess := humanize.Time(*u.LastAccessTime)
			props = append(props, Prop{"Status", status})
			props = append(props, Prop{"Last access",
				fmt.Sprintf("%s (%s)",
					u.LastAccessTime.Local().Format(dateTimeFormat),
					sinceLastAccess),
			})
		} else {
			props = append(props, Prop{"Status", faint.Sprint("Invisible")})
		}
	}

	props = append(props,
		Prop{"Unique ID", u.UniqueId},
		Prop{"Created", u.MemberSince.Local().Format(dateTimeFormat)},
		Prop{"Figure", u.FigureString})

	motto := strings.TrimSpace(u.Motto)
	if motto == "" {
		motto = faint.Sprint("(no motto)")
	}
	props = append(props, Prop{"Motto", motto})
	badges := make([]string, len(u.SelectedBadges))
	for i := range u.SelectedBadges {
		badges[i] = u.SelectedBadges[i].Name
	}
	if len(badges) > 0 {
		props = append(props, Prop{"Selected badges", strings.Join(badges, faint.Sprint(", "))})
	} else {
		props = append(props, Prop{"Selected badges", faint.Sprint("(none)")})
	}

	RenderProperties(props)
}

func RenderUserInfo2(u web.User) {
	t := makeTableWriter()

	rows := []table.Row{}

	rows = append(rows, []any{"Name", u.Name})
	if u.ProfileVisible {
		if u.LastAccessTime != nil {
			status := ""
			if u.Online {
				status = green.Sprint("Online")
			} else {
				status = red.Sprint("Offline")
			}
			sinceLastAccess := humanize.Time(*u.LastAccessTime)
			rows = append(rows, []any{"Status", status})
			rows = append(rows, []any{"Last access",
				fmt.Sprintf("%s (%s)",
					u.LastAccessTime.Local().Format(dateTimeFormat),
					sinceLastAccess),
			})
		} else {
			t.AppendRow([]any{"Status", faint.Sprint("Invisible")})
		}
	}

	rows = append(rows,
		[]any{"Unique ID", u.UniqueId},
		[]any{"Created", u.MemberSince.Local().Format(dateTimeFormat)},
		[]any{"Figure", u.FigureString})

	motto := strings.TrimSpace(u.Motto)
	if motto == "" {
		motto = faint.Sprint("(no motto)")
	}
	rows = append(rows, []any{"Motto", motto})
	badges := make([]string, len(u.SelectedBadges))
	for i := range u.SelectedBadges {
		badges[i] = u.SelectedBadges[i].Name
	}
	if len(badges) > 0 {
		rows = append(rows, []any{"Selected badges", strings.Join(badges, faint.Sprint(", "))})
	} else {
		rows = append(rows, []any{"Selected badges", faint.Sprint("(none)")})
	}

	t.AppendRows(rows)
	t.Render()
}
