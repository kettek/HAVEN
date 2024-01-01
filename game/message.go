package game

import (
	"image/color"
	"time"
)

type Message struct {
	Text     string
	Duration time.Duration
	Color    color.NRGBA
	X        int
	Y        int
	start    time.Time
	id       int
}

var messageID int
