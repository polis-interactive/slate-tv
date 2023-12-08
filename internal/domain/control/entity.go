package control

import (
	"fmt"
	"github.com/polis-interactive/slate-tv/internal/domain"
	"log"
	"reflect"
	"sync"
)

type controllerImpl interface {
	RunMainLoop()
}

type Controller struct {
	impl           controllerImpl
	bus            Bus
	inputStates    []domain.InputState
	inputTolerance float64
	mu             *sync.RWMutex
	wg             *sync.WaitGroup
	shutdowns      chan struct{}
}

func newController(cfg Config, bus Bus) (*Controller, error) {

	base := &Controller{
		bus:            bus,
		inputTolerance: cfg.GetInputTolerance(),
		mu:             &sync.RWMutex{},
		wg:             &sync.WaitGroup{},
	}

	var err error = nil
	switch cfg.GetControlType() {
	case domain.ControlTypes.ADC:
		base.impl, err = newAdcController(base, cfg)
	case domain.ControlTypes.NONE:
		inputTypes := cfg.GetInputTypes()
		inputStates := make([]domain.InputState, len(inputTypes))
		for i, inputType := range inputTypes {
			inputStates[i].InputType = inputType
		}
		base.inputStates = inputStates
		base.impl, err = newNoneController(base, cfg)
	}

	return base, err
}

func (c *Controller) startup() {

	log.Println(fmt.Sprintf("%s, startup: starting", reflect.TypeOf(c.impl)))

	if c.shutdowns == nil {
		c.shutdowns = make(chan struct{})
		c.wg.Add(1)
		go c.impl.RunMainLoop()
		log.Println(fmt.Sprintf("%s, startup: running", reflect.TypeOf(c.impl)))
	}
}

func (c *Controller) shutdown() {
	log.Println(fmt.Sprintf("%s, shutdown: shutting down", reflect.TypeOf(c.impl)))
	if c.shutdowns != nil {
		close(c.shutdowns)
		c.wg.Wait()
		c.shutdowns = nil
	}
	log.Println(fmt.Sprintf("%s, shutdown: done", reflect.TypeOf(c.impl)))
}

func (c *Controller) SetInputValue(inputNumber int, inputValue float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	oldVal := c.inputStates[inputNumber].InputValue
	positiveDiff := oldVal + c.inputTolerance
	negativeBound := oldVal - c.inputTolerance
	if positiveDiff > inputValue && inputValue > negativeBound {
		return
	}
	c.inputStates[inputNumber].InputValue = inputValue
	c.bus.HandleControlInputChange(&domain.InputState{
		InputType:  c.inputStates[inputNumber].InputType,
		InputValue: inputValue,
	})
}

func (c *Controller) GetShutdowns() chan struct{} {
	return c.shutdowns
}

func (c *Controller) GetWg() *sync.WaitGroup {
	return c.wg
}
