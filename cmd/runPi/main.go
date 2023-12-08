package main

import (
	"github.com/polis-interactive/slate-tv/internal/application"
	"github.com/polis-interactive/slate-tv/internal/domain"
	"github.com/polis-interactive/slate-tv/internal/util"
	"log"
	"os"
	"os/signal"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/experimental/devices/ads1x15"
	"syscall"
	"time"
)

func main() {
	conf := &application.Config{
		GridDefinition: util.GridDefinition{
			Rows:         3,
			Columns:      11,
			LedPerCell:   8,
			LedPerScoot:  2,
			RowExtension: 0,
		},
		ControlConfig: &application.ControlConfig{
			ControlType:    domain.ControlTypes.ADC,
			InputTolerance: 0.001,
		},
		AdcConfig: &application.AdcConfig{
			InputPins: []domain.InputPin{
				{
					InputType: domain.InputTypes.BRIGHTNESS,
					Pin:       ads1x15.Channel0,
				},
				{
					InputType: domain.InputTypes.SPEED,
					Pin:       ads1x15.Channel1,
				},
				{
					InputType: domain.InputTypes.PROGRAM,
					Pin:       ads1x15.Channel2,
				},
				{
					InputType: domain.InputTypes.VALUE,
					Pin:       ads1x15.Channel3,
				},
			},
			ReadFrequency: physic.Hertz * 33,
			ReadVoltage:   physic.MilliVolt * 3300,
		},
		GraphicsConfig: &application.GraphicsConfig{
			ShaderFiles: []string{
				"checkerboard", "stripe-wheel-spread",
				"bar-hoppin",
			},
			ReloadOnUpdate: false,
			DisplayOutput:  false,
			PixelSize:      1,
			Frequency:      33 * time.Millisecond,
		},
		RenderConfig: &application.RenderConfig{
			RenderType:      domain.RenderTypes.WS2812,
			RenderFrequency: 33 * time.Millisecond,
		},
		Ws2812Config: &application.Ws2812Config{
			GpioPin:   util.GpioPinTypes.GPIO18,
			StripType: util.StripTypes.WS2811RGB,
			Gamma:     1.2,
		},
	}

	app, err := application.NewApplication(conf)
	if err != nil {
		panic(err)
	}

	err = app.Startup()
	if err != nil {
		log.Println("Main: failed to startup, shutting down")
		err2 := app.Shutdown()
		if err2 != nil {
			log.Println("Main: issue shutting down; ", err2)
		}
		panic(err)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	err = app.Shutdown()
	if err != nil {
		log.Println("Main: issue shutting down; ", err)
	}

}
