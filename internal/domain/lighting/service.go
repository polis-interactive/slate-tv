package lighting

import (
	"fmt"
	"github.com/polis-interactive/slate-tv/internal/util"
	"log"
	"sync"
)

type Service struct {
	lights    []util.Light
	lastLight int
	grid      util.Grid
	mu        *sync.RWMutex
}

func NewService(cfg Config) *Service {

	log.Println("Lighting, NewService: creating")

	lights, lastLight := generateLights(cfg.GetBoardConfiguration(), cfg.GetDisallowedPositions())

	grid := util.Grid{}
	for _, l := range lights {
		grid.TrySetMinMax(l.Position)
	}

	log.Println(fmt.Sprintf(
		"MinP (%d, %d); MaxP (%d, %d); last led %d",
		grid.MinX, grid.MinY, grid.MaxX, grid.MaxY, lastLight,
	))

	log.Println("Lighting, NewService: created")

	return &Service{
		lights:    lights,
		lastLight: lastLight,
		grid:      grid,
		mu:        &sync.RWMutex{},
	}
}

func (s *Service) GetLightCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.lastLight
}

func (s *Service) GetGrid() util.Grid {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.grid
}

func (s *Service) GetLights() (lights []util.Light, preLockedMutex *sync.RWMutex) {
	s.mu.RLock()
	return s.lights, s.mu
}
