package render

import (
	"errors"
	"fmt"
	"github.com/polis-interactive/slate-tv/internal/domain"
	"log"
	"reflect"
	"sync"
	"time"
)

type render interface {
	startup()
	shutdown()
	runMainLoop()
	runRenderLoop() error
	runRender() error
}

type baseRender struct {
	bus             Bus
	renderFrequency time.Duration
	render          render
	ledCount        int
	mu              *sync.RWMutex
	wg              *sync.WaitGroup
	shutdowns       chan struct{}
}

func newRender(cfg Config, bus Bus) (render, error) {

	log.Println("render, newRender: creating")

	base := &baseRender{
		bus:             bus,
		renderFrequency: cfg.GetRenderFrequency(),
		mu:              &sync.RWMutex{},
		wg:              &sync.WaitGroup{},
	}

	var r render

	switch cfg.GetRenderType() {
	case domain.RenderTypes.WINDOW:
	case domain.RenderTypes.WS2812:
		r = newWs2812Render(base, cfg)
	case domain.RenderTypes.TERMINAL:
		r = newTerminalRender(base)
	default:
		return nil, errors.New("sign, newRender: Invalid render type")
	}

	base.render = r
	return r, nil
}

func (r *baseRender) startup() {

	log.Println(fmt.Sprintf("%s, startup: starting", reflect.TypeOf(r.render)))

	if r.shutdowns == nil {
		r.ledCount = r.bus.GetLightCount()
		r.shutdowns = make(chan struct{})
		r.wg.Add(1)
		go r.render.runMainLoop()
		log.Println(fmt.Sprintf("%s, startup: running", reflect.TypeOf(r.render)))
	}

}

func (r *baseRender) shutdown() {
	log.Println(fmt.Sprintf("%s, shutdown: shutting down", reflect.TypeOf(r.render)))
	if r.shutdowns != nil {
		close(r.shutdowns)
		r.wg.Wait()
		r.shutdowns = nil
	}
	log.Println(fmt.Sprintf("%s, shutdown: done", reflect.TypeOf(r.render)))
}

func (r *baseRender) runRenderLoop() error {
	for {
		select {
		case _, ok := <-r.shutdowns:
			if !ok {
				return nil
			}
		case <-time.After(r.renderFrequency):
			err := r.render.runRender()
			if err != nil {
				return err
			}
		}
	}
}
