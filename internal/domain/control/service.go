package control

import (
	"github.com/polis-interactive/slate-tv/internal/domain"
	"log"
	"sync"
)

type service struct {
	controller *Controller
	mu         *sync.Mutex
}

var _ domain.ControlService = (*service)(nil)

func NewService(cfg Config, bus Bus) (*service, error) {

	log.Println("Control, NewService: creating")

	c, err := newController(cfg, bus)
	if err != nil {
		log.Println("Control, NewService: error creating render")
		return nil, err
	}

	log.Println("Control, NewService: created")
	return &service{
		controller: c,
		mu:         &sync.Mutex{},
	}, nil
}

func (s *service) Startup() {
	log.Println("ControlService Startup: starting")
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.controller != nil {
		s.controller.startup()
	}
}

func (s *service) Shutdown() {
	log.Println("ControlService Shutdown: shutting down")
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.controller != nil {
		s.controller.shutdown()
	}
}

func (s *service) GetControllerStates() []domain.InputState {
	s.controller.mu.RLock()
	defer s.controller.mu.RUnlock()
	inputStates := make([]domain.InputState, len(s.controller.inputStates))
	for i := range s.controller.inputStates {
		inputStates[i] = s.controller.inputStates[i]
	}
	return inputStates
}
