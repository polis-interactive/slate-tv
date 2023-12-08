package main

import (
	"github.com/polis-interactive/slate-tv/data"
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
		LightingConfig: &application.LightingConfig{
			BoardConfiguration:  data.BoardConfiguration,
			DisallowedPositions: data.BoardDisallowedPositions,
		},
		ControlConfig: &application.ControlConfig{
			ControlType:    domain.ControlTypes.NONE,
			InputTolerance: 0.001,
		},
		AdcConfig: &application.AdcConfig{
			InputPins: []domain.InputPin{
				{
					InputType: domain.InputTypes.INPUT1,
					Pin:       ads1x15.Channel0,
				},
			},
			ReadFrequency: physic.Hertz * 33,
			ReadVoltage:   physic.MilliVolt * 3300,
		},
		GraphicsConfig: &application.GraphicsConfig{
			ShaderFiles: []string{
				"basic",
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
		InputTypes: []domain.InputType{
			domain.InputTypes.INPUT1,
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
