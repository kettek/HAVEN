package res

import (
	"github.com/tinne26/etxt"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

type Font struct {
	*sfnt.Font
	Size int
}

var BigFontName = "x16y32pxGridGazer.ttf"
var DefFontName = "Mx437_IBM_VGA_9x16.ttf"
var SmallFontName = "Mx437_IBM_EGA_8x8.ttf"

var Text *etxt.Renderer
var BigFont Font
var DefFont Font
var SmallFont Font

func init() {
	Text = etxt.NewRenderer()

	b, err := FS.ReadFile(BigFontName)
	if err != nil {
		panic(err)
	}
	f, err := opentype.Parse(b)
	if err != nil {
		panic(err)
	}
	BigFont = Font{
		Font: f,
		Size: 32,
	}

	b, err = FS.ReadFile(DefFontName)
	if err != nil {
		panic(err)
	}
	f, err = opentype.Parse(b)
	if err != nil {
		panic(err)
	}

	DefFont = Font{
		Font: f,
		Size: 16,
	}

	b, err = FS.ReadFile(SmallFontName)
	if err != nil {
		panic(err)
	}
	f, err = opentype.Parse(b)
	if err != nil {
		panic(err)
	}
	SmallFont = Font{
		Font: f,
		Size: 8,
	}

	Text.SetFont(BigFont.Font)
	Text.SetSize(float64(BigFont.Size))
	Text.Utils().SetCache8MiB()
}
