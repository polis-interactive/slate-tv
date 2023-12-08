package application

import (
	"github.com/polis-interactive/slate-tv/internal/domain"
	"github.com/polis-interactive/slate-tv/internal/util"
	"periph.io/x/periph/conn/physic"
	"time"
)

type LightingConfig struct {
	BoardConfiguration  []util.BoardConfiguration
	DisallowedPositions []util.Point
}

func (l *LightingConfig) GetBoardConfiguration() []util.BoardConfiguration {
	return l.BoardConfiguration
}

func (l *LightingConfig) GetDisallowedPositions() []util.Point {
	return l.DisallowedPositions
}

type Ws2812Config struct {
	GpioPin   util.GpioPinType
	StripType util.StripType
	Gamma     float32
}

func (w *Ws2812Config) GetGpioPin() util.GpioPinType {
	return w.GpioPin
}

func (w *Ws2812Config) GetStripType() util.StripType {
	return w.StripType
}

func (w *Ws2812Config) GetGamma() float32 {
	return w.Gamma
}

type RenderConfig struct {
	RenderType      domain.RenderType
	RenderFrequency time.Duration
}

func (r *RenderConfig) GetRenderType() domain.RenderType {
	return r.RenderType
}

func (r *RenderConfig) GetRenderFrequency() time.Duration {
	return r.RenderFrequency
}

type WindowConfig struct {
	TileSize int
}

func (w *WindowConfig) GetTileSize() int {
	return w.TileSize
}

type AdcConfig struct {
	InputPins     []domain.InputPin
	ReadFrequency physic.Frequency
	ReadVoltage   physic.ElectricPotential
}

func (a *AdcConfig) GetInputPins() []domain.InputPin {
	return a.InputPins
}

func (a *AdcConfig) GetReadFrequency() physic.Frequency {
	return a.ReadFrequency
}
func (a *AdcConfig) GetReadVoltage() physic.ElectricPotential {
	return a.ReadVoltage
}

type ControlConfig struct {
	ControlType    domain.ControlType
	InputTolerance float64
	GrpcPort       int
}

func (c *ControlConfig) GetGrpcPort() int {
	return c.GrpcPort
}

func (c *ControlConfig) GetControlType() domain.ControlType {
	return c.ControlType
}

func (c *ControlConfig) GetInputTolerance() float64 {
	return c.InputTolerance
}

type GraphicsConfig struct {
	ShaderFiles    []string
	DisplayOutput  bool
	ReloadOnUpdate bool
	PixelSize      int
	Frequency      time.Duration
}

func (g *GraphicsConfig) GetShaderFiles() []string {
	return g.ShaderFiles
}

func (g *GraphicsConfig) GetGraphicsReloadOnUpdate() bool {
	return g.ReloadOnUpdate
}

func (g *GraphicsConfig) GetGraphicsDisplayOutput() bool {
	return g.DisplayOutput
}

func (g *GraphicsConfig) GetGraphicsPixelSize() int {
	return g.PixelSize
}

func (g *GraphicsConfig) GetGraphicsFrequency() time.Duration {
	return g.Frequency
}

type Config struct {
	*LightingConfig
	*RenderConfig
	*Ws2812Config
	*WindowConfig
	*AdcConfig
	*ControlConfig
	*GraphicsConfig
	InputTypes []domain.InputType
}

func (c *Config) GetInputTypes() []domain.InputType {
	return c.InputTypes
}
