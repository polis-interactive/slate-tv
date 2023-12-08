package bus

import (
	"github.com/polis-interactive/slate-tv/internal/domain"
)

func (b *bus) HandleControlInputChange(state *domain.InputState) {
	b.graphicsService.HandleInputChange(state)
}
