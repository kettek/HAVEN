package settings

import "github.com/hajimehoshi/ebiten/v2"

var MayoMode ebiten.Filter = ebiten.FilterLinear
var ClarityMode ebiten.Filter = ebiten.FilterNearest

var FilterMode = ClarityMode

var StackShading bool = true
