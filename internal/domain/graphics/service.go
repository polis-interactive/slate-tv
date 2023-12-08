package graphics

import (
	"fmt"
	"github.com/polis-interactive/go-lighting-utils/pkg/graphicsShader"
	"github.com/polis-interactive/slate-tv/internal/domain"
	"github.com/polis-interactive/slate-tv/internal/util"
	"log"
	"math"
	"sync"
)

type service struct {
	graphics *graphics
	mu       *sync.Mutex
}

var _ domain.GraphicsService = (*service)(nil)

func NewService(cfg Config, bus Bus) (*service, error) {
	log.Println("Graphics, NewService: creating")

	g, err := newGraphics(cfg, bus)
	if err != nil {
		log.Println("Graphics, NewService: error creating graphics")
		return nil, err
	}
	return &service{
		graphics: g,
		mu:       &sync.Mutex{},
	}, nil
}

func (s *service) Startup() {
	log.Println("RenderService Startup: starting")
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.graphics != nil {
		s.graphics.startup()
	}
}

func (s *service) Reset() {
	log.Println("RenderService Startup: resetting")
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.graphics != nil {
		s.graphics.shutdown()
		s.graphics.startup()
	}
}

func (s *service) Shutdown() {
	log.Println("GraphicsService Shutdown: shutting down")
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.graphics != nil {
		s.graphics.shutdown()
	}
}

func (s *service) HandleInputChange(state *domain.InputState) {
	s.graphics.mu.Lock()
	defer s.graphics.mu.Unlock()
	if s.graphics.gs == nil {
		return
	}
	programCount := float64(len(s.graphics.shaderFiles))
	log.Println(state.InputValue)
	log.Println(programCount)
	if programCount == 1 {
		return
	}
	var programKey graphicsShader.ShaderKey
	if state.InputValue == 1 {
		log.Println("1.0?")
		programKey = graphicsShader.ShaderKey(rune(programCount - 1))
	} else {
		selectProgram := int(math.Floor(programCount * state.InputValue))
		log.Println(programCount * state.InputValue)
		log.Println(math.Floor(programCount * state.InputValue))
		log.Println(selectProgram)
		programKey = graphicsShader.ShaderKey(rune(selectProgram))
	}
	err := s.graphics.gs.SetShader(programKey)
	if err != nil {
		log.Println(fmt.Sprintf(
			"GraphicsService, HandleInputChange - Program: couldn't set to %s with error %s",
			programKey, err.Error(),
		))
	}
}

func (s *service) GetPb() (pb *util.PixelBuffer, preLockedMutex *sync.RWMutex) {
	s.graphics.mu.RLock()
	return s.graphics.pb, s.graphics.mu
}
