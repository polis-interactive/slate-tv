package main

import (
	"fmt"
	"github.com/polis-interactive/slate-tv/internal/domain"
	"github.com/polis-interactive/slate-tv/internal/domain/control"
	"log"
	"os"
	"os/signal"
	"periph.io/x/periph/conn/physic"
	"syscall"
)

type testGuiConfig struct{}

func (t testGuiConfig) GetInputPins() []domain.InputPin {
	panic("unused")
}

func (t testGuiConfig) GetReadFrequency() physic.Frequency {
	panic("unused")
}

func (t testGuiConfig) GetReadVoltage() physic.ElectricPotential {
	panic("unused")
}

func (t testGuiConfig) GetInputTolerance() float64 {
	return 0.001
}

var _ control.Config = (*testGuiConfig)(nil)

func (t testGuiConfig) GetControlType() domain.ControlType {
	return domain.ControlTypes.GUI
}

func (t testGuiConfig) GetInputTypes() []domain.InputType {
	return []domain.InputType{
		domain.InputTypes.BRIGHTNESS,
		domain.InputTypes.PROGRAM,
		domain.InputTypes.SPEED,
		domain.InputTypes.VALUE,
	}
}

type testGuiBus struct{}

var _ control.Bus = (*testGuiBus)(nil)

func (t testGuiBus) HandleControlInputChange(state *domain.InputState) {
	log.Println(fmt.Sprintf("%s: %f", state.InputType, state.InputValue))
}

func main() {
	conf := &testGuiConfig{}
	bus := &testGuiBus{}
	c, err := control.NewService(conf, bus)
	if err != nil {
		log.Println("Unable to create control service")
	}
	c.Startup()

	s := make(chan os.Signal)
	signal.Notify(s, os.Interrupt, syscall.SIGTERM)
	<-s

	c.Shutdown()
}
