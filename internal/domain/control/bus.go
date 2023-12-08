package control

import "github.com/polis-interactive/slate-tv/internal/domain"

type Bus interface {
	HandleControlInputChange(*domain.InputState)
}
