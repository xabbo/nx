package extract

import (
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/spf13/cobra"

	"b7c.io/swfx"

	root "xabbo.io/nx/cmd/nx/cmd"
	"xabbo.io/nx/raw/nitro"
)

var opts struct {
	data   bool
	images bool
}

var Cmd = &cobra.Command{
	Use:   "extract [files...]",
	Short: "Extracts files. Currently supports Nitro archives.",
	RunE:  run,
}

func init() {
	root.Cmd.AddCommand(Cmd)

	f := Cmd.Flags()
	f.BoolVarP(&opts.data, "data", "d", false, "Extract binary data")
	f.BoolVarP(&opts.images, "images", "i", false, "Extract images")
}

func run(cmd *cobra.Command, args []string) (err error) {
	if !opts.data && !opts.images {
		opts.data = true
		opts.images = true
	}

	for _, file := range args {
		err := extractFile(file)
		if err != nil {
			cmd.PrintErrf("%s: %s\n", file, err)
		}
	}
	return
}

func extractFile(name string) (err error) {
	switch {
	case strings.HasSuffix(name, ".nitro"):
		err = extractNitro(name)
	case strings.HasSuffix(name, ".swf"):
		err = extractSwf(name)
	default:
		err = fmt.Errorf("unknown file format")
	}
	return
}

func extractSwf(name string) (err error) {
	f, err := os.Open(name)
	if err != nil {
		return
	}
	defer f.Close()

	outDir := ""
	if idx := strings.LastIndex(name, "."); idx > 0 {
		outDir = name[:idx]
	} else {
		outDir = name + "_extracted"
	}
	err = os.MkdirAll(outDir, 0755)
	if err != nil {
		return
	}

	swf, err := swfx.ReadSwf(f)
	if err != nil {
		return
	}

	for _, tag := range swf.Tags {
		switch tag := tag.(type) {
		case *swfx.DefineBitsLossless2:
			if !opts.images {
				continue
			}
			err = extractBitsLossless2(outDir, swf, tag)
			if err != nil {
				return
			}
		case *swfx.DefineBitsJpeg2:
			if !opts.images {
				continue
			}
			err = extractBitsJpeg2(outDir, swf, tag)
			if err != nil {
				return
			}
		case *swfx.DefineBinaryData:
			if !opts.data {
				continue
			}
			err = extractBinaryData(outDir, swf, tag)
			if err != nil {
				return
			}
		}
	}
	return
}

func extractBitsLossless2(outDir string, swf *swfx.Swf, tag *swfx.DefineBitsLossless2) error {
	var names []string
	var ok bool
	if names, ok = swf.ReverseSymbols[tag.CharacterId()]; !ok {
		names = []string{strconv.Itoa(int(tag.CharacterId()))}
	}

	for _, name := range names {
		outputFile := filepath.Join(outDir, name+".png")
		img, err := tag.Decode()
		if err != nil {
			return err
		}
		f, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			return err
		}
		defer f.Close()
		err = png.Encode(f, img)
		if err != nil {
			return err
		}
		fmt.Println(outputFile)
	}
	return nil
}

func extractBitsJpeg2(outDir string, swf *swfx.Swf, tag *swfx.DefineBitsJpeg2) error {
	var names []string
	var ok bool
	if names, ok = swf.ReverseSymbols[tag.CharacterId()]; !ok {
		names = []string{fmt.Sprintf("%d", tag.CharacterId())}
	}

	var ext string
	switch tag.ImageType() {
	case swfx.Jpeg:
		ext = ".jpg"
	case swfx.Png:
		ext = ".png"
	case swfx.Gif:
		ext = ".gif"
	default:
		return fmt.Errorf("unknown image type")
	}

	for _, name := range names {
		outputFile := filepath.Join(outDir, name+ext)
		f, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.Write(tag.ImageData)
		if err != nil {
			return err
		}
		fmt.Println(outputFile)
	}
	return nil
}

func extractBinaryData(outDir string, swf *swfx.Swf, tag *swfx.DefineBinaryData) (err error) {
	var names []string
	var ok bool
	if names, ok = swf.ReverseSymbols[tag.CharacterId()]; !ok {
		names = []string{strconv.Itoa(int(tag.CharacterId()))}
	}

	mtype := mimetype.Detect(tag.Data)
	ext := mtype.Extension()

	for _, name := range names {
		outputFile := filepath.Join(outDir, name+ext)

		f, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			return fmt.Errorf("failed to open file: %q", outputFile)
		}
		defer f.Close()

		_, err = f.Write(tag.Data)
		if err != nil {
			return err
		}
		fmt.Println(outputFile)
	}
	return nil
}

func extractNitro(name string) (err error) {
	dir := strings.TrimSuffix(name, filepath.Ext(name))
	if dir == name {
		dir += "_extracted"
	}

	err = os.Mkdir(dir, 0755)
	if err != nil {
		return
	}

	archive, err := readNitroArchive(name)

	for name, file := range archive.Files {
		fmt.Printf("%s: ", name)
		err = os.WriteFile(filepath.Join(dir, name), file.Data, 0644)
		if err != nil {
			return
		}
		fmt.Println("ok")
	}

	return
}

func readNitroArchive(name string) (archive nitro.Archive, err error) {
	f, err := os.Open(name)
	if err != nil {
		return
	}
	defer f.Close()

	return nitro.NewReader(f).ReadArchive()
}
