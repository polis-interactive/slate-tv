package control

import (
	"errors"
	"fmt"
	"github.com/polis-interactive/slate-tv/internal/domain"
	"log"
	"math"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/experimental/conn/analog"
	"periph.io/x/periph/experimental/devices/ads1x15"
	"periph.io/x/periph/host"
	"sync"
	"time"
)

type adcController struct {
	*Controller
	inputPins []domain.InputPin
	freq      physic.Frequency
	volt      physic.ElectricPotential
}

var _ controllerImpl = (*adcController)(nil)

func newAdcController(c *Controller, cfg AdcConfig) (*adcController, error) {

	inputPins := cfg.GetInputPins()
	inputStates := make([]domain.InputState, len(inputPins))
	for i, inputPin := range inputPins {
		inputStates[i].InputType = inputPin.InputType
		inputStates[i].InputValue = -1.0
	}
	c.inputStates = inputStates

	if _, err := host.Init(); err != nil {
		return nil, err
	}

	bus, err := i2creg.Open("")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to open I²C: %v", err))
	}
	err = bus.Close()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to close I2C: %v", err))
	}

	a := &adcController{
		Controller: c,
		inputPins:  inputPins,
		freq:       cfg.GetReadFrequency(),
		volt:       cfg.GetReadVoltage(),
	}
	return a, nil
}

func (a *adcController) RunMainLoop() {

	log.Println("adcControl, Main Loop: running")

	for {
		err := a.tryRunMainLoop()
		if err != nil {
			log.Println(fmt.Sprintf("adcControl, main loop: error %v", err))
		}
		select {
		case _, ok := <-a.shutdowns:
			if !ok {
				goto CloseAdcLoop
			}
		case <-time.After(5 * time.Second):
			log.Println("adcControl, Main Loop: retrying")
		}
	}
CloseAdcLoop:
	log.Println("adcControl, Main Loop: closed")
	a.wg.Done()
}

func (a *adcController) tryRunMainLoop() error {

	wg := &sync.WaitGroup{}
	chIn := make([]chan domain.InputPair, 0)
	var chAll chan domain.InputPair

	bus, err := i2creg.Open("")
	if err != nil {
		log.Fatalf("failed to open I²C: %v", err)
	}
	defer bus.Close()

	adc, err := ads1x15.NewADS1115(bus, &ads1x15.DefaultOpts)
	if err != nil {
		log.Fatalln(err)
	}
	defer adc.Halt()

	for i, input := range a.inputPins {
		var ch chan domain.InputPair
		ch, err = a.createAdcChannel(adc, input, i, wg)
		if err != nil {
			goto CleanUp
		}
		chIn = append(chIn, ch)
	}

	chAll = mergeInputEvents(chIn...)

	for {
		select {
		case _, ok := <-a.shutdowns:
			if !ok {
				goto CleanUp
			}
		case e, ok := <-chAll:
			if !ok {
				err = errors.New("event channel suddenly closed")
				goto CleanUp
			}
			a.SetInputValue(e.InputNumber, e.InputValue)
		}
	}

CleanUp:
	for i := range chIn {
		if chIn[i] != nil {
			close(chIn[i])
		}
	}
	wg.Wait()
	return err
}

func (a *adcController) createAdcChannel(
	adc *ads1x15.Dev, input domain.InputPin, inputNumber int, wg *sync.WaitGroup,
) (chan domain.InputPair, error) {
	pin, err := adc.PinForChannel(input.Pin, a.volt, a.freq, ads1x15.BestQuality)
	if err != nil {
		return nil, errors.New(
			fmt.Sprintf(
				"unable to start adc for input %s with err %v",
				input.InputType, err,
			),
		)
	}
	wg.Add(1)
	chIn := pin.ReadContinuous()
	chOut := make(chan domain.InputPair, 5)
	go func() {
		defer func() {
			wg.Done()
			err2 := pin.Halt()
			if err2 != nil {
				log.Println(fmt.Sprintf("failed to halt pin for input %s with err %v", input.InputType, err2))
			}
		}()
		for {
			select {
			case _, ok := <-a.shutdowns:
				if !ok {
					return
				}
			case rd, ok := <-chIn:
				if !ok {
					close(chOut)
					return
				}
				raw := a.getValueFromAnalog(rd)
				chOut <- domain.InputPair{
					InputNumber: inputNumber,
					InputValue:  math.Max(math.Min(raw, 1.0), 0.0),
				}
			}
		}
	}()
	return chOut, nil
}

func (a *adcController) getValueFromAnalog(sample analog.Sample) float64 {
	n := float64(sample.V) / float64(a.volt)
	return n
}
