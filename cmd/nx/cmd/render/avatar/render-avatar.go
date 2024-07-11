package avatar

import (
	"fmt"
	"os"
	"slices"

	"github.com/spf13/cobra"

	"xabbo.b7c.io/nx"
	gd "xabbo.b7c.io/nx/gamedata"
	"xabbo.b7c.io/nx/imager"
	"xabbo.b7c.io/nx/web"

	_root "xabbo.b7c.io/nx/cmd/nx/cmd"
	_parent "xabbo.b7c.io/nx/cmd/nx/cmd/render"
	"xabbo.b7c.io/nx/cmd/nx/spinner"
	"xabbo.b7c.io/nx/cmd/nx/util"
)

var Cmd = &cobra.Command{
	Use:  "avatar [figure]",
	Args: cobra.RangeArgs(0, 1),
	RunE: runRenderAvatar,
}

var opts struct {
	dir        int
	headDir    int
	action     string
	expression string
	userName   string
	handItem   int
	headOnly   bool
	outputName string
	noColor    bool
	verbose    bool
	outFormat  string
}

var validFormats = []string{"png", "svg"}

func init() {
	f := Cmd.Flags()
	f.IntVarP(&opts.dir, "dir", "d", 2, "The direction of the avatar (0-7)")
	f.IntVarP(&opts.headDir, "head-dir", "H", 2, "The direction of the avatar's head (0-7)")
	f.StringVarP(&opts.action, "action", "a", "std", "The action of the avatar")
	f.StringVarP(&opts.expression, "expression", "e", "", "The expression of the avatar")
	f.StringVarP(&opts.userName, "user", "u", "", "The name of the user to fetch a figure for")
	f.BoolVar(&opts.headOnly, "head-only", false, "Render head only")
	f.StringVarP(&opts.outputName, "output", "o", "", "The name of the output file")
	f.BoolVar(&opts.noColor, "no-color", false, "Do not color figure parts")
	f.BoolVarP(&opts.verbose, "verbose", "v", false, "Verbose output")
	f.StringVarP(&opts.outFormat, "format", "f", "png", "Output format")

	_parent.Cmd.AddCommand(Cmd)
}

func runRenderAvatar(cmd *cobra.Command, args []string) (err error) {
	// Match body direction if head direction not set
	if !cmd.Flags().Lookup("head-dir").Changed {
		opts.headDir = opts.dir
	}

	api := nx.NewApiClient(_root.Host)

	figureSpecified := len(args) > 0
	userSpecified := opts.userName != ""

	if !figureSpecified && !userSpecified {
		return fmt.Errorf("no figure or user specified")
	}

	if figureSpecified && userSpecified {
		return fmt.Errorf("only one of either figure or user may be specified")
	}

	if !slices.Contains(validFormats, opts.outFormat) {
		return fmt.Errorf("invalid output format %q, must be %s",
			opts.outFormat, util.CommaList(validFormats, "or"))
	}

	cmd.SilenceUsage = true

	if !slices.Contains(nx.AvatarActions, nx.AvatarState(opts.action)) {
		return fmt.Errorf("invalid action %q, must be one of %s",
			opts.action, util.CommaList(nx.AvatarActions, "or"))
	}

	if opts.expression != "" && !slices.Contains(nx.AvatarExpressions, nx.AvatarState(opts.expression)) {
		return fmt.Errorf("invalid expression %q, must be one of %s",
			opts.expression, util.CommaList(nx.AvatarExpressions, "or"))
	}

	vars := map[string]any{}
	vars["dir"] = opts.dir
	vars["hdir"] = opts.headDir
	vars["act"] = opts.action
	if opts.expression == "" {
		vars["expr"] = "ntr" // Neutral
	} else {
		vars["expr"] = opts.expression
	}

	var figureString string
	if len(args) > 0 {
		figureString = args[0]
		if opts.outputName == "" {
			opts.outputName = figureString
		}
	} else {
		var user web.User
		err = spinner.DoErr("Loading user...", func() (err error) {
			user, err = api.GetUserByName(opts.userName)
			if err != nil {
				return
			}
			figureString = user.FigureString
			return nil
		})
		if err != nil {
			return
		}
		fmt.Println(figureString)

		if opts.outputName == "" {
			opts.outputName = "$name-$act-$expr-$dir-$hdir"
		}
		vars["name"] = user.Name
	}

	vars["figure"] = figureString

	opts.outputName = os.Expand(opts.outputName, func(s string) (ret string) {
		if value, ok := vars[s]; ok {
			ret = fmt.Sprint(value)
		}
		return
	})
	fileName := opts.outputName
	switch opts.outFormat {
	case "png", "svg":
		fileName += "." + opts.outFormat
	}

	mgr := gd.NewManager(_root.Host)
	renderer := imager.NewAvatarImager(mgr)

	var figure nx.Figure
	err = figure.Parse(figureString)
	if err != nil {
		return
	}

	err = util.LoadGameData(mgr, "Loading game data...",
		gd.GameDataFigure, gd.GameDataFigureMap,
		gd.GameDataVariables, gd.GameDataAvatar)
	if err != nil {
		return
	}

	parts, err := renderer.Parts(figure)
	if err != nil {
		return
	}

	libraries := map[string]struct{}{}

	for _, part := range parts {
		libraries[part.LibraryName] = struct{}{}
	}

	err = spinner.DoErr("Loading figure part libraries...", func() error {
		for lib := range libraries {
			err = mgr.LoadFigureParts(lib)
			if err != nil {
				return err
			}
			if opts.verbose {
				spinner.Printf("Loaded %s\n", lib)
			}
		}
		return nil
	})
	if err != nil {
		return
	}

	avatar := imager.Avatar{
		Figure:        figure,
		Direction:     opts.dir,
		HeadDirection: opts.headDir,
		Action:        nx.AvatarState(opts.action),
		Expression:    nx.AvatarState(opts.expression),
		HeadOnly:      opts.headOnly,
	}

	anim, err := renderer.Compose(avatar)
	if err != nil {
		return
	}

	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	switch opts.outFormat {
	case "png":
		encoder := imager.NewEncoderPNG()
		encoder.EncodeFrame(f, anim, 0, 0)
	case "svg":
		encoder := imager.NewEncoderSVG()
		encoder.EncodeFrame(f, anim, 0, 0)
	default:
		return
	}

	fmt.Printf("output: %s\n", fileName)
	return
}
