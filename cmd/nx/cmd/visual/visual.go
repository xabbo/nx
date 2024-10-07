package info

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/jedib0t/go-pretty/v6/list"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"

	gd "xabbo.io/nx/gamedata"
	"xabbo.io/nx/res"

	_root "xabbo.io/nx/cmd/nx/cmd"
	"xabbo.io/nx/cmd/nx/spinner"
	"xabbo.io/nx/cmd/nx/util"
)

var Cmd = &cobra.Command{
	Use:     "visual [flags] <identifier>",
	Aliases: []string{"vis", "v"},
	Short:   "Displays furni visualization information.",
	Args:    cobra.ExactArgs(1),
	RunE:    run,
}

var opts struct {
	size   int
	frames bool
}

func init() {
	f := Cmd.Flags()
	f.IntVarP(&opts.size, "size", "s", 0, "The visualization size to print. Prints all sizes if not specified.")
	f.BoolVar(&opts.frames, "frames", false, "Print all frame indices in frame sequences.")

	_root.Cmd.AddCommand(Cmd)
}

func run(cmd *cobra.Command, args []string) (err error) {
	cmd.SilenceUsage = true

	mgr := gd.NewManager(_root.Host)
	identifier := args[0]

	err = spinner.DoErr("Loading game data...", func() (err error) {
		return mgr.Load(gd.GameDataVariables, gd.GameDataFurni)
	})
	if err != nil {
		return
	}

	fi, ok := mgr.Furni()[identifier]
	if !ok {
		return errors.New("furni info not found")
	}

	err = spinner.DoErr("Loading furni library...", func() (err error) {
		return mgr.LoadFurni(identifier)
	})
	if err != nil {
		return
	}

	split := strings.Split(identifier, "*")
	libName := split[0]

	lib, ok := mgr.Library(libName).(res.FurniLibrary)
	if !ok {
		return fmt.Errorf("failed to load library: %q", identifier)
	}

	index := lib.Index()
	if index == nil {
		return errors.New("failed to load index")
	}
	cmd.Printf("Furni name: %s\n", fi.Name)
	cmd.Printf("Visualization type: %s\n", index.Visualization)

	l := list.NewWriter()
	l.SetStyle(list.StyleConnectedLight)
	l.SetOutputMirror(cmd.OutOrStdout())
	l.UnIndentAll()

	visualizations := lib.Visualizations()
	if opts.size == 0 {
		sizes := maps.Keys(visualizations)
		slices.Sort(sizes)
		for _, size := range sizes {
			vis := visualizations[size]
			l.AppendItem(fmt.Sprintf("Size: %d", vis.Size))
			l.Indent()
			printVisualization(l, vis)
			l.UnIndent()
		}
	} else {
		vis, ok := visualizations[opts.size]
		if !ok {
			return fmt.Errorf("no visualization for size: %d", opts.size)
		}
		printVisualization(l, vis)
	}

	l.Render()
	return
}

func printVisualization(l list.Writer, vis *res.Visualization) {
	// angle
	l.AppendItem(fmt.Sprintf("Angle: %d", vis.Angle))

	// directions
	dirs := maps.Keys(vis.Directions)
	slices.Sort(dirs)
	l.AppendItem(fmt.Sprintf("Directions: %s", util.CommaList(dirs, "")))

	// layers
	l.AppendItem(fmt.Sprintf("Layers: %d", vis.LayerCount))
	l.Indent()
	layerIds := maps.Keys(vis.Layers)
	slices.Sort(layerIds)
	for _, layerId := range layerIds {
		l.AppendItem(fmt.Sprintf("Layer %d", layerId))
	}
	l.UnIndent()

	// colors
	l.AppendItem(fmt.Sprintf("Colors: %d", len(vis.Colors)))
	l.Indent()
	colorIds := maps.Keys(vis.Colors)
	slices.Sort(colorIds)
	for _, colorId := range colorIds {
		color := vis.Colors[colorId]
		l.AppendItem(fmt.Sprint(color.Id))
		l.Indent()
		layerIds := maps.Keys(color.Layers)
		slices.Sort(layerIds)
		for _, layerId := range layerIds {
			layer := color.Layers[layerId]
			l.AppendItem(fmt.Sprintf("%d: %s", layer.Id, layer.Color))
		}
		l.UnIndent()
	}
	l.UnIndent()

	// animations
	l.AppendItem(fmt.Sprintf("Animations: %d", len(vis.Animations)))
	l.Indent()
	animationIds := maps.Keys(vis.Animations)
	slices.Sort(animationIds)
	for _, animationId := range animationIds {
		printAnimation(l, vis.Animations[animationId])
	}
	l.UnIndent()
}

func printAnimation(l list.Writer, anim *res.Animation) {
	if anim.TransitionTo != nil {
		l.AppendItem(fmt.Sprintf("Animation %d (transition -> %d)", anim.Id,
			anim.TransitionTo.Id))
	} else {
		l.AppendItem(fmt.Sprintf("Animation %d", anim.Id))
	}
	l.Indent()
	l.AppendItem(fmt.Sprintf("Layers: %d", len(anim.Layers)))
	l.Indent()
	for _, layer := range anim.Layers {
		printAnimationLayer(l, layer)
	}
	l.UnIndent()
	l.UnIndent()
}

func printAnimationLayer(l list.Writer, layer *res.AnimationLayer) {
	l.AppendItem(fmt.Sprintf("Layer %d", layer.Id))
	l.Indent()
	l.AppendItem(fmt.Sprintf("Loop count: %d", layer.LoopCount))
	l.AppendItem(fmt.Sprintf("Frame repeat: %d", layer.FrameRepeat))
	l.AppendItem(fmt.Sprintf("Random: %d", layer.Random))
	l.AppendItem(fmt.Sprintf("Sequences: %d", len(layer.FrameSequences)))
	l.Indent()
	for i, seq := range layer.FrameSequences {
		printFrameSequence(l, i, seq)
	}
	l.UnIndent()
	l.UnIndent()
}

func printFrameSequence(l list.Writer, i int, seq res.FrameSequence) {
	count := util.Pluralize(len(seq), "frame", "s")
	if opts.frames {
		l.AppendItem(fmt.Sprintf("%d: %s [%s]", i, count, util.CommaList(seq, "")))
	} else {
		l.AppendItem(fmt.Sprintf("%d: %s", i, count))
	}
}
