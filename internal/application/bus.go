package application

import (
	"github.com/polis-interactive/slate-tv/internal/domain"
	"github.com/polis-interactive/slate-tv/internal/domain/control"
	"github.com/polis-interactive/slate-tv/internal/domain/graphics"
	"github.com/polis-interactive/slate-tv/internal/domain/render"
)

type applicationBus interface {
	Startup()
	Shutdown()
	BindRenderService(r domain.RenderService)
	BindControlService(b domain.ControlService)
	BindGraphicsService(g domain.GraphicsService)
	BindLightingService(stateService domain.LightingService)
	render.Bus
	control.Bus
	graphics.Bus
}
