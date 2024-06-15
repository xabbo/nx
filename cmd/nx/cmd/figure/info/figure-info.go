package info

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/v6/list"
	"github.com/spf13/cobra"

	"xabbo.b7c.io/nx"
	gd "xabbo.b7c.io/nx/gamedata"

	_root "xabbo.b7c.io/nx/cmd/nx/cmd"
	_parent "xabbo.b7c.io/nx/cmd/nx/cmd/figure"
	"xabbo.b7c.io/nx/cmd/nx/spinner"
	"xabbo.b7c.io/nx/cmd/nx/util"
)

var Cmd = &cobra.Command{
	Use:  "info",
	RunE: runInfo,
}

var opts struct {
	userName        string
	showParts       bool
	showIdentifiers bool
	showColors      bool
	showAll         bool
}

func init() {
	f := Cmd.Flags()
	f.StringVarP(&opts.userName, "user", "u", "", "User to load figure for")
	f.BoolVarP(&opts.showIdentifiers, "identifiers", "i", false, "Show clothing furni identifiers")
	f.BoolVarP(&opts.showParts, "parts", "p", false, "Show individual figure parts")
	f.BoolVarP(&opts.showColors, "colors", "c", false, "Show figure part colors")
	f.BoolVar(&opts.showAll, "all", false, "Show all information")

	_parent.Cmd.AddCommand(Cmd)
}

func runInfo(cmd *cobra.Command, args []string) (err error) {
	if len(args) < 1 && opts.userName == "" {
		return fmt.Errorf("no figure or user specified")
	}

	cmd.SilenceUsage = true

	if opts.showAll {
		opts.showIdentifiers = true
		opts.showParts = true
		opts.showColors = true
	}

	l := list.NewWriter()
	l.SetStyle(list.StyleConnectedLight)
	l.SetOutputMirror(os.Stdout)
	l.UnIndent()

	figureString := ""
	if len(args) > 0 {
		figureString = args[0]
	} else {
		err = spinner.DoErr("Loading user...", func() error {
			api := nx.NewApiClient(_root.Host)
			user, err := api.GetUserByName(opts.userName)
			if err != nil {
				return err
			}
			figureString = user.FigureString
			return nil
		})
		if err != nil {
			return
		}
		fmt.Println(figureString)
	}

	figure := nx.Figure{}
	err = figure.Parse(figureString)
	if err != nil {
		return err
	}

	mgr := gd.NewManager(_root.Host)
	err = util.LoadGameData(mgr, "Loading game data...",
		gd.GameDataFigure, gd.GameDataFigureMap, gd.GameDataFurni, gd.GameDataTexts, gd.GameDataVariables)
	if err != nil {
		return err
	}

	partCountMap := make(map[int]int)
	clothingMap := make(map[int]gd.FurniInfo)
	for _, f := range mgr.Furni() {
		if f.SpecialType == nx.FurniTypeClothing {
			parts := strings.Split(f.CustomParams, ",")
			for _, s := range parts {
				s = strings.TrimSpace(s)
				if id, err := strconv.Atoi(s); err == nil {
					c := partCountMap[id]
					if c == 0 || len(parts) < c {
						partCountMap[id] = len(parts)
						clothingMap[id] = f
					}
				}
			}
		}
	}

	for _, part := range figure.Items {
		setGroup := mgr.Figure().Sets[part.Type]
		set := setGroup[part.Id]

		if typeName, ok := mgr.Texts()["avatareditor.category."+string(part.Type)]; ok {
			l.AppendItem(fmt.Sprintf("%s (%s)", typeName, part.Type))
		} else {
			l.AppendItem(fmt.Sprintf("%s", part.Type))
		}

		l.Indent()

		if fi, ok := clothingMap[part.Id]; ok {
			if opts.showIdentifiers {
				l.AppendItem(fmt.Sprintf("%4d: %s [%s]", part.Id, fi.Name, fi.Identifier))
			} else {
				l.AppendItem(fmt.Sprintf("%4d: %s", part.Id, fi.Name))
			}
		} else {
			l.AppendItem(fmt.Sprintf("%4d", part.Id))
		}

		if opts.showParts {
			l.Indent()
			for _, piece := range set.Parts {
				mapPart := nx.FigurePart{Type: piece.Type, Id: piece.Id}
				if lib, ok := mgr.FigureMap().Parts[mapPart]; ok {
					l.AppendItem(fmt.Sprintf("%s-%d [%s]", piece.Type, piece.Id, lib.Name))
				} else {
					l.AppendItem(fmt.Sprintf("%s-%d", piece.Type, piece.Id))
				}
			}
			l.UnIndent()
		}

		if opts.showColors {
			palette := mgr.Figure().PaletteFor(part.Type)
			for _, colorId := range part.Colors {
				if color, ok := palette[colorId]; ok {
					colorValue, err := strconv.ParseInt(color.Value, 16, 64)
					if err == nil {
						r := (colorValue >> 16) & 0xff
						g := (colorValue >> 8) & 0xff
						b := colorValue & 0xff
						l.AppendItem(fmt.Sprintf("%4d: #%06x \x1b[48;2;%d;%d;%dm  \x1b[0m",
							colorId, colorValue, r, g, b))
					}
				} else {
					l.AppendItem(fmt.Sprintf("%4d", colorId))
				}
			}
		}

		l.UnIndent()
	}

	l.Render()

	return
}
