package convert

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
	"xabbo.b7c.io/nx/cmd/nx/spinner"
	gd "xabbo.b7c.io/nx/gamedata"
	"xabbo.b7c.io/nx/gamedata/origins"

	_parent "xabbo.b7c.io/nx/cmd/nx/cmd/figure"
)

var Cmd = &cobra.Command{
	Use:   "convert [figure]",
	Short: "Origins figure converter",
	Long:  "Converts Origins figure strings to its modern representation",
	Args:  cobra.ExactArgs(1),
	RunE:  run,
}

func init() {
	_parent.Cmd.AddCommand(Cmd)
}

func run(cmd *cobra.Command, args []string) (err error) {
	cmd.SilenceUsage = true

	originsFigure := strings.TrimSpace(args[0])
	if len(originsFigure) % 5 != 0 {
		return origins.ErrInvalidFigureStringLength
	}

	for _, c := range originsFigure {
		if c < '0' || c > '9' {
			return origins.ErrNonNumericFigureString
		}
	}

	spinner.Start()
	defer spinner.Stop()

	gdm := gd.NewManager("www.habbo.com")

	spinner.Message("Loading modern figure data...")
	err = gdm.Load(gd.GameDataFigure)
	if err != nil {
		return fmt.Errorf("failed to load modern figure data: %w", err)
	}

	spinner.Message("Loading origins figure data...")
	ofd, err := loadOriginsFigureData()
	if err != nil {
		return fmt.Errorf("failed to load origins figure data: %w", err)
	}

	colorMap := origins.MakeColorMap(gdm.Figure())
	converter := origins.NewFigureConverter(ofd, colorMap)

	figure, err := converter.Convert(originsFigure)
	if err != nil {
		return
	}

	spinner.Stop()
	cmd.Printf("%s\n", figure.String())
	return
}

func loadOriginsFigureData() (fd *origins.FigureData, err error) {
	res, err := http.Get("http://origins-gamedata.habbo.com/figuredata/1")
	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, errors.New(res.Status)
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	fd, err = origins.ParseFigureData(b)
	return
}
