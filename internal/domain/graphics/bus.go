package graphics

import "github.com/polis-interactive/slate-tv/internal/util"

type Bus interface {
	GetLightGrid() util.Grid
}
