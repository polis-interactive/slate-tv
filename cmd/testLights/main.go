package main

import (
	"github.com/polis-interactive/slate-tv/internal/domain"
	"github.com/polis-interactive/slate-tv/internal/domain/render"
	"github.com/polis-interactive/slate-tv/internal/util"
	"log"
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type testRenderConfig struct {
}

func (t testRenderConfig) GetGpioPin() util.GpioPinType {
	return util.GpioPinTypes.GPIO18
}

func (t testRenderConfig) GetStripType() util.StripType {
	return util.StripTypes.WS2811RGB
}

func (t testRenderConfig) GetGamma() float32 {
	return 1.2
}

func (t testRenderConfig) GetTileSize() int {
	panic("unused")
}

func (t testRenderConfig) GetRenderType() domain.RenderType {
	return domain.RenderTypes.WS2812
}

func (t testRenderConfig) GetRenderFrequency() time.Duration {
	return time.Millisecond * 33
}

func (t testRenderConfig) GetGridDefinition() util.GridDefinition {
	return util.GridDefinition{
		Rows:         3,
		Columns:      11,
		LedPerCell:   8,
		LedPerScoot:  2,
		RowExtension: 0,
	}
}

var _ render.Config = (*testRenderConfig)(nil)

type testRenderBus struct {
	colors    []uint32
	increment int
}

func (t *testRenderBus) GetLightCount() int {
	return 196
}

func (t *testRenderBus) CopyLightsToColorBuffer(buff []util.Color) error {
	//TODO implement me
	panic("not used")
}

func (t *testRenderBus) CopyLightsToUint32Buffer(buff []uint32) error {
	for i := 0; i < 14; i++ {
		for j := 0; j < 14; j++ {
			position := (t.increment + i*2 + j) % 256
			buff[i*14+j] = wheelUint32(position)
		}
	}
	return nil
}

func min3(a, b, c float64) float64 {
	return math.Min(math.Min(a, b), c)
}

func wheelUint32(pos int) uint32 {
	//h := math.Floor(float64(pos*360) / 256)
	//kr := math.Mod(5+h*6, 6)
	//kg := math.Mod(3+h*6, 6)
	//kb := math.Mod(1+h*6, 6)

	//r := 1 - math.Max(min3(kr, 4-kr, 1), 0)
	//g := 1 - math.Max(min3(kg, 4-kg, 1), 0)
	//b := 1 - math.Max(min3(kb, 4-kb, 1), 0)

	c := util.Color{
		R: 0,
		G: 255,
		B: 0,
		W: 0,
	}
	return c.ToBits()
}

func min(i int, j int) int {
	if i <= j {
		return i
	} else {
		return j
	}
}

func max(i int, j int) int {
	if i <= j {
		return j
	} else {
		return i
	}
}

var _ render.Bus = (*testRenderBus)(nil)

func main() {
	conf := &testRenderConfig{}
	bus := &testRenderBus{
		colors:    make([]uint32, 196),
		increment: 0,
	}
	r, err := render.NewService(conf, bus)
	if err != nil {
		log.Println("Unable to create control service")
	}
	r.Startup()

	s := make(chan os.Signal)
	signal.Notify(s, os.Interrupt, syscall.SIGTERM)
	<-s

	r.Shutdown()
}
