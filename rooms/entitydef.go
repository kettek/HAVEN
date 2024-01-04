package rooms

import (
	"github.com/kettek/ebihack23/actors"
)

type EntityDef struct {
	Actor      string
	OnCreate   actors.CreateFunc
	OnInteract actors.InteractFunc
}

type EntityDefs map[string]EntityDef
