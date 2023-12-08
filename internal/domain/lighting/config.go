package lighting

import (
	"github.com/polis-interactive/slate-tv/internal/util"
)

type Config interface {
	GetBoardConfiguration() []util.BoardConfiguration
	GetDisallowedPositions() []util.Point
}
