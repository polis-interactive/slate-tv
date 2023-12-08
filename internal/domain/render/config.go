package render

import (
	"github.com/polis-interactive/slate-tv/internal/domain"
	"github.com/polis-interactive/slate-tv/internal/util"
	"time"
)

type ws2812RenderConfig interface {
	GetGpioPin() util.GpioPinType
	GetStripType() util.StripType
	GetGamma() float32
}

type windowRenderConfig interface {
	GetTileSize() int
}

type baseRenderConfig interface {
	GetRenderType() domain.RenderType
	GetRenderFrequency() time.Duration
}

type Config interface {
	ws2812RenderConfig
	windowRenderConfig
	baseRenderConfig
}
