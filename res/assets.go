package res

import (
	"bytes"
	"embed"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//go:embed *.png
//go:embed *.ttf
var FS embed.FS

var loadedSpriteStacks = make(map[string][]*ebiten.Image)

var TileWidth = 13.0
var TileHeight = 13.0
var TileHalfWidth = 6.0
var TileHalfHeight = 6.0
var TileYStep = 9.0
var TileXStep = 4.5
var FontName = "x16y32pxGridGazer.ttf"

func LoadSpriteStack(sprite string) ([]*ebiten.Image, error) {
	if layers, ok := loadedSpriteStacks[sprite]; ok {
		return layers, nil
	}
	b, err := FS.ReadFile(sprite + ".png")
	if err != nil {
		b, err = FS.ReadFile("missing.png")
		if err != nil {
			return nil, err
		}
	}
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	eimg := ebiten.NewImageFromImage(img)

	var layers []*ebiten.Image
	for x := 0; x < img.Bounds().Max.X; x += int(TileWidth) {
		layers = append(layers, eimg.SubImage(image.Rect(x, 0, x+int(TileWidth), int(TileHeight))).(*ebiten.Image))
	}
	loadedSpriteStacks[sprite] = layers
	return layers, nil
}

var Font font.Face

func init() {
	b, err := FS.ReadFile(FontName)
	if err != nil {
		panic(err)
	}
	tt, err := opentype.Parse(b)
	if err != nil {
		panic(err)
	}
	Font, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    32,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		panic(err)
	}
}
