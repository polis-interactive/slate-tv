package main

import (
	"fmt"
	"github.com/polis-interactive/slate-tv/internal/domain"
	"github.com/polis-interactive/slate-tv/internal/domain/control"
	"log"
	"os"
	"os/signal"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/experimental/devices/ads1x15"
	"syscall"
)

type testAdcConfig struct{}

func (t testAdcConfig) GetInputTolerance() float64 {
	return 0.001
}

var _ control.Config = (*testAdcConfig)(nil)

func (t testAdcConfig) GetControlType() domain.ControlType {
	return domain.ControlTypes.ADC
}

func (t testAdcConfig) GetInputTypes() []domain.InputType {
	panic("not in use")
}

func (t testAdcConfig) GetInputPins() []domain.InputPin {
	return []domain.InputPin{
		{
			InputType: domain.InputTypes.INPUT1,
			Pin:       ads1x15.Channel0,
		},
	}
}

func (t testAdcConfig) GetReadFrequency() physic.Frequency {
	return physic.Hertz * 20
}

func (t testAdcConfig) GetReadVoltage() physic.ElectricPotential {
	return physic.MilliVolt * 3300
}

type inputState struct {
	minValue float64
	maxValue float64
}

type testAdcBus struct {
	minMaxState map[domain.InputType]*inputState
}

var _ control.Bus = (*testAdcBus)(nil)

func (t *testAdcBus) HandleControlInputChange(state *domain.InputState) {
	ip := state.InputType
	iv := state.InputValue
	if _, ok := t.minMaxState[ip]; !ok {
		t.minMaxState[ip] = &inputState{
			minValue: iv,
			maxValue: iv,
		}
		log.Println(fmt.Sprintf("%s: %f", state.InputType, state.InputValue))
		return
	}
	if t.minMaxState[ip].minValue > iv {
		t.minMaxState[ip].minValue = iv
		log.Println(fmt.Sprintf("%s:  new min %f", state.InputType, state.InputValue))
	}
	if t.minMaxState[ip].maxValue < iv {
		t.minMaxState[ip].maxValue = iv
		log.Println(fmt.Sprintf("%s:  new max %f", state.InputType, state.InputValue))
	}
}

func main() {
	conf := &testAdcConfig{}
	bus := &testAdcBus{
		minMaxState: make(map[domain.InputType]*inputState),
	}
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
