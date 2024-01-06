package figure

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/v6/list"
	"github.com/spf13/cobra"

	"github.com/b7c/nx"

	root "cli/cmd"
	"cli/spinner"
	"cli/util"
)

var infoCmd = &cobra.Command{
	Use:  "info",
	RunE: runInfo,
}

var (
	userName        string
	showParts       bool
	showIdentifiers bool
	showColors      bool
	showAll         bool
)

func init() {
	figureCmd.AddCommand(infoCmd)

	infoCmd.Flags().StringVarP(&userName, "user", "u", "", "User to load figure for")
	infoCmd.Flags().BoolVarP(&showIdentifiers, "identifiers", "i", false, "Show clothing furni identifiers")
	infoCmd.Flags().BoolVarP(&showParts, "parts", "p", false, "Show individual figure parts")
	infoCmd.Flags().BoolVarP(&showColors, "colors", "c", false, "Show figure part colors")
	infoCmd.Flags().BoolVar(&showAll, "all", false, "Show all information")
}

func runInfo(cmd *cobra.Command, args []string) (err error) {
	if len(args) < 1 && userName == "" {
		return fmt.Errorf("no figure or user specified")
	}

	cmd.SilenceUsage = true

	if showAll {
		showIdentifiers = true
		showParts = true
		showColors = true
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
			api := nx.NewApiClient(root.Host)
			user, err := api.GetUserByName(userName)
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

	mgr := nx.NewGamedataManager(root.Host)
	err = util.LoadGamedata(mgr, "Loading game data...",
		nx.GamedataFigure, nx.GamedataFigureMap, nx.GamedataFurni, nx.GamedataTexts, nx.GamedataVariables)
	if err != nil {
		return err
	}

	partCountMap := make(map[int]int)
	clothingMap := make(map[int]nx.FurniInfo)
	for _, f := range mgr.Furni {
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

	for _, part := range figure.Parts {
		setGroup := mgr.Figure.Sets[part.Type]
		set := setGroup[part.Id]

		if typeName, ok := mgr.Texts["avatareditor.category."+string(part.Type)]; ok {
			l.AppendItem(fmt.Sprintf("%s (%s)", typeName, part.Type))
		} else {
			l.AppendItem(fmt.Sprintf("%s", part.Type))
		}

		l.Indent()

		if fi, ok := clothingMap[part.Id]; ok {
			if showIdentifiers {
				l.AppendItem(fmt.Sprintf("%4d: %s [%s]", part.Id, fi.Name, fi.Identifier))
			} else {
				l.AppendItem(fmt.Sprintf("%4d: %s", part.Id, fi.Name))
			}
		} else {
			l.AppendItem(fmt.Sprintf("%4d", part.Id))
		}

		if showParts {
			l.Indent()
			for _, piece := range set.Parts {
				mapPart := nx.FigureMapPart{Type: piece.Type, Id: piece.Id}
				if lib, ok := mgr.FigureMap.Parts[mapPart]; ok {
					l.AppendItem(fmt.Sprintf("%s-%d [%s]", piece.Type, piece.Id, lib.Name))
				} else {
					l.AppendItem(fmt.Sprintf("%s-%d", piece.Type, piece.Id))
				}
			}
			l.UnIndent()
		}

		if showColors {
			palette := mgr.Figure.PaletteFor(part.Type)
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
