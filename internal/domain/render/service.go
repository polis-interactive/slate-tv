package render

import (
	"github.com/polis-interactive/slate-tv/internal/domain"
	"log"
	"sync"
)

type service struct {
	render render
	mu     *sync.Mutex
}

var _ domain.RenderService = (*service)(nil)

func NewService(cfg Config, bus Bus) (*service, error) {

	log.Println("Render, NewService: creating")

	r, err := newRender(cfg, bus)
	if err != nil {
		log.Println("Render, NewService: error creating render")
		return nil, err
	}

	log.Println("Render, NewService: created")
	return &service{
		render: r,
		mu:     &sync.Mutex{},
	}, nil
}

func (s *service) Startup() {
	log.Println("RenderService Startup: starting")
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.render != nil {
		s.render.startup()
	}
}

func (s *service) Shutdown() {
	log.Println("RenderService Shutdown: shutting down")
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.render != nil {
		s.render.shutdown()
	}
}
