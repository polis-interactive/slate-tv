package graphics

import (
	"github.com/polis-interactive/slate-tv/internal/domain"
	"time"
)

type Config interface {
	GetGraphicsReloadOnUpdate() bool
	GetGraphicsDisplayOutput() bool
	GetGraphicsPixelSize() int
	GetShaderFiles() []string
	GetGraphicsFrequency() time.Duration
	GetInputTypes() []domain.InputType
}
