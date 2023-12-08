package domain

import (
	"github.com/polis-interactive/slate-tv/internal/util"
	"periph.io/x/periph/experimental/devices/ads1x15"
	"sync"
)

const Program = "slate-tv"

type RenderType string

const (
	ws2812Render   = "WS2812_RENDER"
	terminalRender = "TERMINAL_RENDER"
	windowRender   = "WINDOW_RENDER"
)

var RenderTypes = struct {
	WS2812   RenderType
	TERMINAL RenderType
	WINDOW   RenderType
}{
	WS2812:   ws2812Render,
	TERMINAL: terminalRender,
	WINDOW:   windowRender,
}

type RenderService interface {
	Startup()
	Shutdown()
}

type ControlType string

const (
	noneControl ControlType = "GUI_CONTROL"
	adcControl  ControlType = "ADC_CONTROL"
)

var ControlTypes = struct {
	NONE ControlType
	ADC  ControlType
}{
	NONE: noneControl,
	ADC:  adcControl,
}

type ControlService interface {
	Startup()
	Shutdown()
	GetControllerStates() []InputState
}

type InputType string

const (
	input1 InputType = "input1"
)

var InputTypes = struct {
	INPUT1 InputType
}{
	INPUT1: input1,
}

type InputPin struct {
	InputType InputType
	Pin       ads1x15.Channel
}

type InputPair struct {
	InputNumber int
	InputValue  float64
}

type InputState struct {
	InputType  InputType
	InputValue float64
}

type GraphicsService interface {
	Startup()
	Shutdown()
	HandleInputChange(*InputState)
	GetPb() (pb *util.PixelBuffer, preLockedMutex *sync.RWMutex)
}

type LightingService interface {
	GetLightCount() int
	GetGrid() util.Grid
	GetLights() (lights []util.Light, preLockedMutex *sync.RWMutex)
}
