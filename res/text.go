package res

import (
	"github.com/tinne26/etxt"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

var BigFontName = "x16y32pxGridGazer.ttf"
var DefFontName = "x12y16pxLineLinker.ttf"
var SmallFontName = "x8y12pxTheStrongGamer.ttf"

var Text *etxt.Renderer
var BigFont *sfnt.Font
var Font *sfnt.Font
var SmallFont *sfnt.Font

func init() {
	Text = etxt.NewRenderer()

	b, err := FS.ReadFile(BigFontName)
	if err != nil {
		panic(err)
	}
	BigFont, err = opentype.Parse(b)
	if err != nil {
		panic(err)
	}

	b, err = FS.ReadFile(DefFontName)
	if err != nil {
		panic(err)
	}
	Font, err = opentype.Parse(b)
	if err != nil {
		panic(err)
	}

	b, err = FS.ReadFile(SmallFontName)
	if err != nil {
		panic(err)
	}
	SmallFont, err = opentype.Parse(b)
	if err != nil {
		panic(err)
	}

	Text.SetFont(BigFont)
	Text.SetSize(32)
	Text.Utils().SetCache8MiB()
}
