package game

import (
	"image/color"
	"time"

	"github.com/kettek/ebihack23/res"
)

type Message struct {
	Text     string
	Duration time.Duration
	Color    color.NRGBA
	X        int
	Y        int
	Font     *res.Font
	start    time.Time
	id       int
}

var messageID int
