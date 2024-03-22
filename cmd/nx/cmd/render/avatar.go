package render

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"os"
	"slices"

	"github.com/disintegration/imaging"
	"github.com/phrozen/blend"
	"github.com/spf13/cobra"

	"github.com/b7c/nx"
	"github.com/b7c/nx/render"
	"github.com/b7c/nx/web"

	root "cli/cmd"
	"cli/spinner"
	"cli/util"
)

var renderAvatarCmd = &cobra.Command{
	Use:  "avatar [figure]",
	Args: cobra.RangeArgs(0, 1),
	RunE: runRenderAvatar,
}

var (
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
)

var validFormats = []string{"png", "svg"}

func init() {
	renderAvatarCmd.Flags().IntVarP(&dir, "dir", "d", 2, "The direction of the avatar (0-7)")
	renderAvatarCmd.Flags().IntVarP(&headDir, "head-dir", "H", 2, "The direction of the avatar's head (0-7)")
	renderAvatarCmd.Flags().StringVarP(&action, "action", "a", "std", "The action of the avatar")
	renderAvatarCmd.Flags().StringVarP(&expression, "expression", "e", "", "The expression of the avatar")
	renderAvatarCmd.Flags().StringVarP(&userName, "user", "u", "", "The name of the user to fetch a figure for")
	renderAvatarCmd.Flags().BoolVar(&headOnly, "head-only", false, "Render head only")
	renderAvatarCmd.Flags().StringVarP(&outputName, "output", "o", "", "The name of the output file")
	renderAvatarCmd.Flags().BoolVar(&noColor, "no-color", false, "Do not color figure parts")
	renderAvatarCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	renderAvatarCmd.Flags().StringVarP(&outFormat, "format", "f", "png", "Output format")

	renderCmd.AddCommand(renderAvatarCmd)
}

func runRenderAvatar(cmd *cobra.Command, args []string) (err error) {
	// Match body direction if head direction not set
	if !cmd.Flags().Lookup("head-dir").Changed {
		headDir = dir
	}

	api := nx.NewApiClient(root.Host)

	figureSpecified := len(args) > 0
	userSpecified := userName != ""

	if !figureSpecified && !userSpecified {
		return fmt.Errorf("no figure or user specified")
	}

	if figureSpecified && userSpecified {
		return fmt.Errorf("only one of either figure or user may be specified")
	}

	if !slices.Contains(validFormats, outFormat) {
		return fmt.Errorf("invalid output format %q, must be %s",
			outFormat, util.CommaList(validFormats, "or"))
	}

	cmd.SilenceUsage = true

	if !slices.Contains(nx.AllActions, nx.AvatarState(action)) {
		return fmt.Errorf("invalid action %q, must be one of %s",
			action, util.CommaList(nx.AllActions, "or"))
	}

	if expression != "" && !slices.Contains(nx.AllExpressions, nx.AvatarState(expression)) {
		return fmt.Errorf("invalid expression %q, must be one of %s",
			expression, util.CommaList(nx.AllExpressions, "or"))
	}

	vars := map[string]any{}
	vars["dir"] = dir
	vars["hdir"] = headDir
	vars["act"] = action
	if expression == "" {
		vars["expr"] = "ntr" // Neutral
	} else {
		vars["expr"] = expression
	}

	var figureString string
	if len(args) > 0 {
		figureString = args[0]
		if outputName == "" {
			outputName = figureString
		}
	} else {
		var user web.User
		err = spinner.DoErr("Loading user...", func() (err error) {
			user, err = api.GetUserByName(userName)
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

		if outputName == "" {
			outputName = "$name-$act-$expr-$dir-$hdir"
		}
		vars["name"] = user.Name
	}

	vars["figure"] = figureString

	outputName = os.Expand(outputName, func(s string) (ret string) {
		if value, ok := vars[s]; ok {
			ret = fmt.Sprint(value)
		}
		return
	})
	fileName := outputName
	switch outFormat {
	case "png", "svg":
		fileName += "." + outFormat
	}

	mgr := nx.NewGamedataManager(root.Host)
	renderer := render.NewFigureRenderer(mgr)

	var figure nx.Figure
	err = figure.Parse(figureString)
	if err != nil {
		return
	}

	err = util.LoadGamedata(mgr, "Loading game data...",
		nx.GamedataFigure, nx.GamedataFigureMap,
		nx.GamedataVariables, nx.GamedataAvatar)
	if err != nil {
		return
	}

	// Load required libraries to render the figure
	// renderer.LoadRequiredLibraries(figure)

	// renderFigure := RenderFigure{
	//     Figure: figure,
	//     Direction: 2,
	//     ...,
	// }
	// img, origin := renderer.RenderImage(figure)

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
			if verbose {
				spinner.Printf("Loaded %s\n", lib)
			}
		}
		return nil
	})
	if err != nil {
		return
	}

	avatar := nx.Avatar{
		Figure:        figure,
		Direction:     dir,
		HeadDirection: headDir,
		Action:        nx.AvatarState(action),
		Expression:    nx.AvatarState(expression),
		HeadOnly:      headOnly,
	}

	sprites, err := renderer.Sprites(avatar)
	if err != nil {
		return
	}

	frameSize := image.Rectangle{}
	for _, sprite := range sprites {
		frameSize = frameSize.Union(sprite.Asset.Image.Bounds().Sub(sprite.Offset))
	}

	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return
	}
	defer f.Close()

	if verbose {
		fmt.Printf("Frame size: %v\n", frameSize.Size())
	}

	switch outFormat {
	case "png":
		img := image.NewRGBA(frameSize)
		for _, sprite := range sprites {
			src := sprite.Asset.Image
			if sprite.Color != color.White && !noColor {
				src = blend.BlendNewImage(src, image.NewUniform(sprite.Color), blend.Multiply)
			}
			if sprite.FlipH {
				src = imaging.FlipH(src)
			}
			if verbose {
				fmt.Printf("Drawing %s at %v\n", sprite.Asset.Name, sprite.Offset.Mul(-1))
			}
			draw.Draw(img,
				sprite.Asset.Image.Bounds().Sub(sprite.Offset),
				src,
				image.Point{},
				draw.Over)
		}

		err = png.Encode(f, img)
		if err != nil {
			return
		}

	case "svg":
		writeSvg(f, sprites)
	}

	fmt.Printf("output: %s\n", fileName)

	return
}

//
// func init() {
// 	renderfigurecmd.flags().intvarp(&dir, "dir", "d", 2, "the direction of the avatar (0-7)")
// 	renderfigurecmd.flags().stringvarp(&gesture, "gesture", "g", "std", "the gesture of the avatar ("+strings.join(validgestures, ", ")+")")
// 	renderfigurecmd.flags().stringvarp(&username, "user", "u", "", "a user to fetch")
// 	renderfigurecmd.flags().intvarp(&handitem, "handitem", "h", 0, "a hand item id")
// 	renderfigurecmd.flags().stringvarp(&outputname, "output", "o", "", "the output filename\ndefaults to the input figure string or user name")
// 	rendercmd.addcommand(renderfigurecmd)
// }
//
// type figure []figurepart
//
// func parsefigure(figurestring string) (figure, error) {
// 	figure := figure{}
// 	parts := make(map[figureparttype]struct{})
// 	for _, partstr := range strings.split(figurestring, ".") {
// 		partsplit := strings.split(partstr, "-")
// 		if len(partsplit) < 2 || len(partsplit) > 4 {
// 			return nil, fmt.errorf("invalid figure part: %q", partstr)
// 		}
// 		parttype := figureparttype(partsplit[0])
// 		if !isvalidfigureparttype(parttype) {
// 			return nil, fmt.errorf("invalid figure part type: %q", parttype)
// 		}
// 		if _, exist := parts[parttype]; exist {
// 			return nil, fmt.errorf("duplicate part type: %q", parttype)
// 		}
// 		parts[parttype] = struct{}{}
// 		partid, err := strconv.atoi(partsplit[1])
// 		if err != nil {
// 			return nil, fmt.errorf("invalid figure part id: %q", partsplit[1])
// 		}
// 		colors := []int{}
// 		for _, coloridstr := range partsplit[2:] {
// 			colorid, err := strconv.atoi(coloridstr)
// 			if err != nil {
// 				return nil, fmt.errorf("invalid color: %q", coloridstr)
// 			}
// 			colors = append(colors, colorid)
// 		}
// 		figure = append(figure, figurepart{
// 			type:   parttype,
// 			id:     partid,
// 			colors: colors,
// 		})
// 	}
// 	return figure, nil
// }
//
// type figurepart struct {
// 	type   figureparttype
// 	id     int
// 	colors []int
// }
//
// type renderpart struct {
// 	lib           gd.partlibrary
// 	parttype      string
// 	figureset     gd.figurepartset
// 	figurepart    gd.figurepart
// 	color         color.color
// 	partpalette   *gd.figurepartpalette
// 	partcolor     *gd.figurepartcolor
// 	assetname     string
// 	originalimage image.image
// 	image         image.image
// 	offset        image.point
// 	hidden        bool
// 	handitem      int
// 	fliph         bool
// }
//
// func runrenderfigure(cmd *cobra.command, args []string) (err error) {
// 	defer spinner.stop()
//
// 	if dir < 0 || dir > 7 {
// 		return fmt.errorf("invalid direction: %d", dir)
// 	}
//
// 	if !slices.contains(validgestures, gesture) {
// 		return fmt.errorf("invalid gesture: %q", gesture)
// 	}
//
// 	var fallbackgesture string
// 	switch gesture {
// 	case "lsp":
// 		fallbackgesture = "lay"
// 	default:
// 		fallbackgesture = "std"
// 	}
//
// 	var figurestr string
// 	if username != "" {
// 		spinner.message("fetching profile...")
// 		spinner.start()
// 		res, err := http.get(fmt.sprintf("https://%s/api/public/users?name=%s", root.host, url.queryescape(username)))
// 		if err != nil {
// 			return err
// 		}
// 		if res.statuscode != http.statusok {
// 			return fmt.errorf("server responded %s", res.status)
// 		}
// 		body, err := io.readall(res.body)
// 		res.body.close()
// 		if err != nil {
// 			return err
// 		}
// 		u := web.user{}
// 		err = json.unmarshal(body, &u)
// 		if err != nil {
// 			return err
// 		}
// 		figurestr = u.figurestring
// 		if outputname == "" {
// 			outputname = u.name
// 		}
// 	} else {
// 		if len(args) == 0 {
// 			return fmt.errorf("no figure string specified")
// 		}
// 		figurestr = args[0]
// 		if outputname == "" {
// 			outputname = figurestr
// 		}
// 	}
//
// 	cmd.silenceusage = true
//
// 	figure, err := parsefigure(figurestr)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.printf("input: %s\n", figurestr)
//
// 	gdm := gd.newmanager(root.host)
//
// 	spinner.start()
// 	spinner.message("loading game data...")
// 	gdm.load(gd.variables, gd.figure)
//
// 	flashclienturl := gdm.variables[varflashclienturl]
// 	figuremapurl := gdm.variables[vardynamicavatardownloadurl]
// 	figuremapurl = strings.replaceall(figuremapurl, "${flash.client.url}", flashclienturl)
//
// 	spinner.message("loading figure map...")
// 	res, err := http.get(figuremapurl)
// 	if err != nil {
// 		return
// 	}
// 	data, err := io.readall(res.body)
// 	if err != nil {
// 		return
// 	}
// 	res.body.close()
//
// 	figuremap := gd.figuremap{}
// 	err = xml.unmarshal(data, &figuremap)
// 	if err != nil {
// 		return
// 	}
//
// 	figurepartsets := make(map[figureparttype]map[int]gd.figurepartset)
// 	for _, sets := range gdm.figure.sets {
// 		t := figureparttype(sets.type)
// 		if _, ok := figurepartsets[t]; !ok {
// 			figurepartsets[t] = make(map[int]gd.figurepartset)
// 		}
// 		for _, set := range sets.sets {
// 			figurepartsets[t][set.id] = set
// 		}
// 	}
//
// 	figuremap.populatemaps()
//
// 	// load libraries
// 	libs := make(map[string]*swfx.swf)
// 	requiredlibs := map[string]gd.partlibrary{}
// 	if handitem > 0 {
// 		requiredlibs["hh_human_item"] = gd.partlibrary{
// 			id: "hh_human_item",
// 		}
// 	}
//
// 	hiddenlayers := map[string]struct{}{}
//
// 	renderparts := []*renderpart{}
//
// 	for _, figurepart := range figure {
// 		set := figurepartsets[figurepart.type][figurepart.id]
// 		fmt.printf("%s-%d\n", figurepart.type, figurepart.id)
// 		for _, hiddenlayer := range set.hiddenlayers {
// 			hiddenlayers[hiddenlayer.parttype] = struct{}{}
// 		}
// 		var currentlib *gd.partlibrary
// 		for _, part := range set.parts {
// 			lib, ok := figuremap.typeidmap[part.type][strconv.itoa(part.id)]
// 			if !ok {
// 				if currentlib != nil {
// 					fmt.printf("  %s-%d: %s (missing; assumed)\n", part.type, part.id, currentlib.id)
// 					lib = *currentlib
// 				} else {
// 					fmt.printf("  %s-%d: missing\n", part.type, part.id)
// 					continue
// 				}
// 			} else {
// 				currentlib = &lib
// 				fmt.printf("  %s-%d: %s\n", part.type, part.id, lib.id)
// 			}
// 			c := color.rgba{r: 255, g: 255, b: 255}
// 			colorindex := part.colorindex - 1
// 			var palette *gd.figurepartpalette
// 			var partcolor *gd.figurepartcolor
// 			if part.type != "ey" && colorindex >= 0 && colorindex < len(figurepart.colors) {
// 				sets := gdm.figure.findsets(string(figurepart.type))
// 				if sets == nil {
// 					continue
// 				}
// 				palette = gdm.figure.findpalette(sets.paletteid)
// 				if palette == nil {
// 					continue
// 				}
// 				partcolor = palette.findcolor(figurepart.colors[colorindex])
// 				if partcolor == nil {
// 					continue
// 				}
// 				colorvalue, err := strconv.parseint(partcolor.value, 16, 32)
// 				if err != nil {
// 					continue
// 				}
// 				c = color.rgba{
// 					r: uint8(colorvalue >> 16),
// 					g: uint8(colorvalue >> 8),
// 					b: uint8(colorvalue),
// 				}
// 			}
// 			renderparts = append(renderparts, &renderpart{
// 				lib:         lib,
// 				parttype:    string(figurepart.type),
// 				figureset:   set,
// 				figurepart:  part,
// 				partpalette: palette,
// 				partcolor:   partcolor,
// 				color:       c,
// 			})
// 			requiredlibs[lib.id] = lib
// 		}
// 	}
//
// 	for _, lib := range requiredlibs {
// 		swf, err := loadpartlib(gdm, lib.id)
// 		if err != nil {
// 			return err
// 		}
// 		spinner.stop()
// 		cmd.printf("loaded %s.swf\n", lib.id)
// 		spinner.start()
// 		libs[lib.id] = swf
// 	}
//
// 	var actions gd.avataractions
// 	if handitem > 0 {
// 		spinner.message("loading avatar actions...")
// 		err := loadactions(gdm, &actions)
// 		if err != nil {
// 			return err
// 		}
// 		handitemid := -1
// 		for _, action := range actions.actions {
// 			if action.id == "carryitem" {
// 				search := strconv.itoa(handitem)
// 				for _, param := range action.params {
// 					if param.id == search {
// 						handitemid, err = strconv.atoi(param.value)
// 						if err != nil {
// 							handitemid = -1
// 						}
// 						break
// 					}
// 				}
// 				break
// 			}
// 		}
// 		if handitemid != -1 {
// 			renderparts = append(renderparts, &renderpart{
// 				lib:      requiredlibs["hh_human_item"],
// 				handitem: handitemid,
// 				color:    color.white,
// 				figurepart: gd.figurepart{
// 					type: "ri",
// 				},
// 			})
// 		}
// 	}
// 	spinner.stop()
//
// 	// resolve assets and calculate required canvas size
// 	bounds := image.rectangle{}
// 	for _, renderpart := range renderparts {
// 		var lib gd.partlibrary
// 		var libswf *swfx.swf
// 		var chrid swfx.characterid
// 		if renderpart.handitem > 0 {
// 			var ok bool
// 			// special case
// 			lib = renderpart.lib
// 			libswf, ok = libs[lib.id]
// 			if !ok {
// 				cmd.printf("lib %q missing!\n", lib.id)
// 				continue
// 			}
// 			renderpart.assetname = fmt.sprintf("h_crr_ri_%d_%d_0", renderpart.handitem, dir)
// 			libassetname := "hh_human_item_" + renderpart.assetname
// 			chrid, ok = libswf.symbols[libassetname]
// 			if !ok {
// 				cmd.printf("symbol %q not found\n", libassetname)
// 				continue
// 			}
// 		} else {
// 			var ok bool
// 			if _, hidden := hiddenlayers[renderpart.figurepart.type]; hidden {
// 				renderpart.hidden = true
// 				continue
// 			}
// 			figurepart := renderpart.figurepart
// 			lib = renderpart.lib
// 			libswf, ok = libs[lib.id]
// 			if !ok {
// 				cmd.printf("lib %q missing!\n", lib.id)
// 				continue
// 			}
// 			gst := gesture
// 			if gesture == "crr" && figurepart.type[0] == 'l' {
// 				gst = "std"
// 			}
//
// 			tryassets := []string{
// 				fmt.sprintf("h_%s_%s_%d_%d_0", gst, figurepart.type, figurepart.id, dir),
// 				fmt.sprintf("h_%s_%s_%d_%d_0", fallbackgesture, figurepart.type, figurepart.id, dir),
// 			}
// 			if dir > 3 && dir < 7 {
// 				flipdir := 0
// 				switch dir {
// 				case 4:
// 					flipdir = 2
// 				case 5:
// 					flipdir = 1
// 				case 6:
// 					flipdir = 0
// 				}
// 				tryassets = slices.insert(tryassets, 1,
// 					fmt.sprintf("h_%s_%s_%d_%d_0", gst, figurepart.type, figurepart.id, flipdir))
// 				tryassets = append(tryassets,
// 					fmt.sprintf("h_%s_%s_%d_%d_0", fallbackgesture, figurepart.type, figurepart.id, flipdir))
// 			}
// 			var resolved bool
// 			for i, name := range tryassets {
// 				if chrid, ok = libswf.symbols[lib.id+"_"+name]; ok {
// 					if len(tryassets) > 2 && (i == 1 || i == 3) {
// 						renderpart.fliph = true
// 					}
// 					renderpart.assetname = name
// 					resolved = true
// 					break
// 				}
// 			}
// 			if !resolved {
// 				cmd.printf("asset for %s/%s not found\n", lib.id, figurepart.type)
// 				continue
// 			}
// 		}
//
// 		chrtag, ok := libswf.characters[chrid]
// 		if !ok {
// 			cmd.printf("character %d not found\n", chrid)
// 			continue
// 		}
// 		imgtag := chrtag.(*swfx.definebitslossless2)
// 		chrid, ok = libswf.symbols[lib.id+"_manifest"]
// 		if !ok {
// 			cmd.printf("manifest for %q not found!\n", lib.id)
// 			continue
// 		}
// 		manifesttag, ok := libswf.characters[chrid].(*swfx.definebinarydata)
// 		if !ok {
// 			cmd.printf("manifest tag for %q not found!\n", lib.id)
// 			continue
// 		}
// 		manifest := gd.manifest{}
// 		decoder := xml.newdecoder(bytes.newreader(manifesttag.data))
// 		decoder.charsetreader = charset.newreaderlabel
// 		err := decoder.decode(&manifest)
// 		if err != nil {
// 			return err
// 		}
// 		assetimg, err := imgtag.decode()
// 		if err != nil {
// 			return err
// 		}
// 		var asset *gd.manifestasset
// 		if len(manifest.libraries) == 0 {
// 			return fmt.errorf("no libraries in manifest for %s", lib.id)
// 		}
// 		for _, a := range manifest.libraries[0].assets {
// 			if a.name == renderpart.assetname {
// 				asset = &a
// 				break
// 			}
// 		}
// 		if asset == nil {
// 			cmd.printf("asset %q not found\n", renderpart.assetname)
// 			continue
// 		}
// 		xys := strings.split(asset.params[0].value, ",")
// 		offsetx, err := strconv.atoi(xys[0])
// 		if err != nil {
// 			return err
// 		}
// 		offsety, err := strconv.atoi(xys[1])
// 		if err != nil {
// 			return err
// 		}
// 		renderpart.offset = image.pt(offsetx, offsety)
// 		renderpart.originalimage = assetimg
// 		renderpart.image = blend.blendnewimage(assetimg, image.newuniform(renderpart.color), blend.multiply)
// 		if renderpart.fliph || (dir > 3 && dir < 7) {
// 			renderpart.offset = image.pt(offsetx*-1+assetimg.bounds().dx()+64, offsety)
// 			if renderpart.fliph {
// 				renderpart.image = imaging.fliph(renderpart.image)
// 			}
// 		}
// 		bounds = bounds.union(renderpart.image.bounds().sub(renderpart.offset))
// 	}
//
// 	layers := []string{
// 		"lh",  // left hand
// 		"ls",  // left shoulder
// 		"lc",  // left coat
// 		"bd",  // body
// 		"sh",  // shoes
// 		"lg",  // legs
// 		"ch",  // chest
// 		"wa",  // waist
// 		"cc",  // coat
// 		"ca",  // chest accessory
// 		"hd",  // head
// 		"fc",  // face
// 		"ey",  // eye
// 		"hr",  // hair
// 		"hrb", // hair below
// 		"fa",  // face accessory
// 		"ea",  // eye accessory
// 		"ha",  // hat
// 		"he",  // head accessory
// 		"li",  // left item
// 		"ri",  // right item
// 		"rh",  // right hand
// 		"rs",  // right shoulder
// 		"rc",  // right coat
// 	}
//
// 	layerorder := map[string]int{}
// 	for i, l := range layers {
// 		layerorder[l] = i
// 	}
//
// 	// render figure
// 	spinner.stop()
// 	fmt.printf("origin is %+v\n", bounds.min)
// 	img := image.newrgba(bounds.sub(bounds.min))
// 	slices.sortfunc(renderparts, func(a, b *renderpart) int {
// 		n := layerorder[a.figurepart.type] - layerorder[b.figurepart.type]
// 		if n == 0 {
// 			n = a.figurepart.id - b.figurepart.id
// 		}
// 		return n
// 	})
//
// 	// fsvg, err := os.openfile(outputname+".svg", os.o_rdwr|os.o_create|os.o_trunc, 0755)
// 	// if err != nil {
// 	// 	return
// 	// }
// 	// defer fsvg.close()

func writeSvg(f io.StringWriter, sprites []render.Sprite) {
	f.WriteString(`<svg xmlns="http://www.w3.org/2000/svg"
		xmlns:svg="http://www.w3.org/2000/svg"
		xmlns:xlink="http://www.w3.org/1999/xlink"
		xmlns:inkscape="http://www.inkscape.org/namespaces/inkscape">`)

	f.WriteString(`<g inkscape:label="layer" id="layer1">`)
	for _, sprite := range sprites {
		buffer := &bytes.Buffer{}
		png.Encode(buffer, sprite.Asset.Image)
		b64 := base64.StdEncoding.EncodeToString(buffer.Bytes())
		offset := sprite.Offset.Mul(-1)
		attrs := ""
		if sprite.FlipH {
			attrs = `transform="scale(-1,1)"`
			// scale transforms around origin x=0,
			// so we need to translate back
			offset.X = offset.X*-1 - sprite.Asset.Image.Bounds().Dx()
		}
		f.WriteString(fmt.Sprintf(`<image id="%s"
			xlink:href="data:image/png;base64,%s"
			style="image-rendering:optimizeSpeed%s"
			x="%d" y="%d" width="%d" height="%d" %s />`,
			sprite.Asset.Name,
			b64,
			"",
			offset.X,
			offset.Y,
			sprite.Asset.Image.Bounds().Dx(),
			sprite.Asset.Image.Bounds().Dy(),
			attrs,
		))
	}
}

//
// 	// fsvg.writestring(`<svg xmlns="http://www.w3.org/2000/svg"
// 	// xmlns:svg="http://www.w3.org/2000/svg"
// 	// xmlns:xlink="http://www.w3.org/1999/xlink"
// 	// xmlns:inkscape="http://www.inkscape.org/namespaces/inkscape">`)
// 	// colors
// 	// fsvg.writestring(`<defs id="defcolors">`)
// 	// for _, renderpart := range renderparts {
// 	// 	if renderpart.partcolor != nil {
// 	// 		fsvg.writestring(fmt.sprintf(`
// 	// 		<filter
// 	//    style="color-interpolation-filters:srgb;"
// 	//    inkscape:label="simple blend"
// 	//    id="filtercolor%s"
// 	//    x="0"
// 	//    y="0"
// 	//    width="1"
// 	//    height="1">
// 	//   <feflood
// 	//      result="flood1"
// 	//      flood-color="#%s"
// 	//      flood-opacity="1"
// 	//      id="feflood322" />
// 	//   <feblend
// 	//      result="blend1"
// 	//      in="flood1"
// 	//      in2="sourcegraphic"
// 	//      mode="multiply"
// 	//      id="feblend324" />
// 	//   <fecomposite
// 	//      operator="in"
// 	//      in="blend1"
// 	//      in2="sourcegraphic"
// 	//      id="fecomposite326" />
// 	// </filter>
// 	// 		`, renderpart.partcolor.value, renderpart.partcolor.value))
// 	// 	}
// 	// }
// 	// fsvg.writestring(`</defs>`)
// 	// fsvg.writestring(`<g inkscape:label="layer" id="layer1">`)
//
// 	for _, renderpart := range renderparts {
// 		if renderpart.hidden || renderpart.image == nil {
// 			continue
// 		}
// 		draw.draw(img, renderpart.image.bounds().sub(renderpart.offset).sub(bounds.min), renderpart.image, image.point{}, draw.over)
//
// 		// buffer := &bytes.buffer{}
// 		// png.encode(buffer, renderpart.originalimage)
//
// 		// b64 := base64.stdencoding.encodetostring(buffer.bytes())
// 		// style := ""
// 		// if renderpart.partcolor != nil {
// 		// 	style = fmt.sprintf(";filter:url(#filtercolor%s)", renderpart.partcolor.value)
// 		// }
// 		// fsvg.writestring(fmt.sprintf(`<image id="%s" xlink:href="data:image/png;base64,%s" style="image-rendering:optimizespeed%s" x="%d" y="%d" width="%d" height="%d" />`,
// 		// 	renderpart.lib.id+"_"+renderpart.assetname,
// 		// 	b64,
// 		// 	style,
// 		// 	-renderpart.offset.x,
// 		// 	-renderpart.offset.y,
// 		// 	renderpart.image.bounds().dx(),
// 		// 	renderpart.image.bounds().dy(),
// 		// ))
// 	}
// 	// fsvg.writestring("</g></svg>")
//
// 	f, err := os.openfile(outputname+".png", os.o_rdwr|os.o_create|os.o_trunc, 0755)
// 	if err != nil {
// 		return
// 	}
// 	defer f.close()
// 	png.encode(f, img)
//
// 	spinner.stop()
// 	fmt.printf("output: %s\n", outputname+".png")
//
// 	return
// }
//
// func loadpartlib(mgr *gd.manager, name string) (swf *swfx.swf, err error) {
// 	spinner.message(fmt.sprintf("loading library %s...", name))
// 	cachedir, err := os.usercachedir()
// 	if err == nil {
// 		cachedir = filepath.join(cachedir, "nx")
// 	} else {
// 		cachedir = ".nx"
// 	}
// 	cachedir = filepath.join(cachedir, "swf", "figure")
// 	err = os.mkdirall(cachedir, 0755)
// 	if err != nil {
// 		return
// 	}
//
// 	filepath := filepath.join(cachedir, name+".swf")
// 	f, err := os.openfile(filepath, os.o_rdwr|os.o_create, 0755)
// 	if err != nil {
// 		return
// 	}
// 	defer f.close()
// 	stats, err := f.stat()
// 	if err != nil {
// 		return
// 	}
//
// 	if stats.size() == 0 {
// 		// download swf
// 		u := mgr.variables[varflashclienturl]
// 		u, err = url.joinpath(u, name+".swf")
// 		if err != nil {
// 			return
// 		}
// 		var res *http.response
// 		res, err = http.get(u)
// 		if err != nil {
// 			return
// 		}
// 		defer res.body.close()
// 		_, err = io.copy(f, res.body)
// 		if err != nil {
// 			return
// 		}
// 		f.seek(0, io.seekstart)
// 	}
//
// 	swf, err = swfx.readswf(f)
// 	return
// }
//
// func loadactions(mgr *gd.manager, actions *gd.avataractions) (err error) {
// 	u := mgr.variables[varflashclienturl]
// 	u, err = url.joinpath(u, "habboavataractions.xml")
// 	if err != nil {
// 		return
// 	}
// 	res, err := http.get(u)
// 	if err != nil {
// 		return
// 	}
// 	defer res.body.close()
// 	data, err := io.readall(res.body)
// 	if err != nil {
// 		return
// 	}
// 	return actions.unmarshalbytes(data)
// }
