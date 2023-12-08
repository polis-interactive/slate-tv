package bus

import (
	"github.com/polis-interactive/slate-tv/internal/util"
)

func (b *bus) GetLightGrid() util.Grid {
	return b.lightingService.GetGrid()
}
