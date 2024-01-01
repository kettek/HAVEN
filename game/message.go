package game

import (
	"image/color"
	"time"

	"golang.org/x/image/font"
)

type Message struct {
	Text     string
	Duration time.Duration
	Color    color.NRGBA
	X        int
	Y        int
	Font     font.Face
	start    time.Time
	id       int
}

var messageID int
