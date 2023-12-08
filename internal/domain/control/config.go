package control

import (
	"github.com/polis-interactive/slate-tv/internal/domain"
	"periph.io/x/periph/conn/physic"
)

type WindowConfig interface {
	GetInputTypes() []domain.InputType
}

type AdcConfig interface {
	GetInputPins() []domain.InputPin
	GetReadFrequency() physic.Frequency
	GetReadVoltage() physic.ElectricPotential
}

type Config interface {
	WindowConfig
	AdcConfig
	GetControlType() domain.ControlType
	GetInputTolerance() float64
}
