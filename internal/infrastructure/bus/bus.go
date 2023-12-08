package bus

import (
	"github.com/polis-interactive/slate-tv/internal/domain"
	"time"
)

type bus struct {
	lightingService domain.LightingService
	renderService   domain.RenderService
	controlService  domain.ControlService
	graphicsService domain.GraphicsService
}

func NewBus() *bus {
	return &bus{}
}

func (b *bus) BindRenderService(r domain.RenderService) {
	b.renderService = r
}

func (b *bus) BindControlService(c domain.ControlService) {
	b.controlService = c
}

func (b *bus) BindGraphicsService(g domain.GraphicsService) {
	b.graphicsService = g
}

func (b *bus) BindLightingService(l domain.LightingService) {
	b.lightingService = l
}

func (b *bus) Startup() {
	b.controlService.Startup()
	time.Sleep(20 * time.Millisecond)
	b.graphicsService.Startup()
	time.Sleep(20 * time.Millisecond)
	b.renderService.Startup()
}

func (b *bus) Shutdown() {
	b.controlService.Shutdown()
	b.renderService.Shutdown()
	b.graphicsService.Shutdown()
}
