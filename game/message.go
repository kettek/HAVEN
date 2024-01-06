package game

import (
	"image/color"
	"time"

	"golang.org/x/image/font/sfnt"
)

type Message struct {
	Text     string
	Duration time.Duration
	Color    color.NRGBA
	X        int
	Y        int
	Font     *sfnt.Font
	start    time.Time
	id       int
}

var messageID int
