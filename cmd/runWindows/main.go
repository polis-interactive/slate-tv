package main

import (
	"github.com/polis-interactive/slate-tv/data"
	"github.com/polis-interactive/slate-tv/internal/application"
	"github.com/polis-interactive/slate-tv/internal/domain"
	"log"
	"os"
	"os/signal"
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
			GrpcPort:       5000,
		},
		WindowConfig: &application.WindowConfig{
			TileSize: 0,
		},
		GraphicsConfig: &application.GraphicsConfig{
			ShaderFiles: []string{
				"checkerboard", "stripe-wheel-spread",
				"bar-hoppin",
			},
			ReloadOnUpdate: true,
			DisplayOutput:  true,
			PixelSize:      30,
			Frequency:      33 * time.Millisecond,
		},
		RenderConfig: &application.RenderConfig{
			RenderType:      domain.RenderTypes.TERMINAL,
			RenderFrequency: 10 * time.Second,
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
