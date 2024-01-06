package res

import (
	"runtime/debug"
	"strings"
)

var EbiOS = ""

func init() {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		EbiOS = "(unknown)"
		return
	}

	for _, d := range bi.Deps {
		if strings.HasSuffix(d.Path, "ebiten/v2") {
			EbiOS = d.Version
		}
	}
}
