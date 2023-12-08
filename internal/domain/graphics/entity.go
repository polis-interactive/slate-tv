package graphics

import (
	"errors"
	"fmt"
	"github.com/polis-interactive/go-lighting-utils/pkg/graphicsShader"
	"github.com/polis-interactive/slate-tv/internal/domain"
	"github.com/polis-interactive/slate-tv/internal/util"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type graphics struct {
	shaderPath        string
	reloadOnUpdate    bool
	pixelSize         int
	graphicsFrequency time.Duration
	gs                *graphicsShader.GraphicsShader
	pb                *util.PixelBuffer
	mu                *sync.RWMutex
	wg                *sync.WaitGroup
	bus               Bus
	lastTimeStep      time.Time
	speed             float64
	inputMap          graphicsShader.UniformDict
	shaderFiles       []string
	shutdowns         chan struct{}
}

func newGraphics(cfg Config, bus Bus) (*graphics, error) {
	log.Println("graphics, newGraphics: creating")
	pixelSize := cfg.GetGraphicsPixelSize()
	if !cfg.GetGraphicsDisplayOutput() {
		pixelSize = 1
	}

	inputs := cfg.GetInputTypes()
	inputMap := make(graphicsShader.UniformDict)
	for _, input := range inputs {
		inputMap[graphicsShader.UniformKey(input)] = 0.0
	}
	inputMap["time"] = 0

	return &graphics{
		reloadOnUpdate:    cfg.GetGraphicsReloadOnUpdate(),
		graphicsFrequency: cfg.GetGraphicsFrequency(),
		shaderFiles:       cfg.GetShaderFiles(),
		pixelSize:         pixelSize,
		gs:                nil,
		pb:                nil,
		bus:               bus,
		mu:                &sync.RWMutex{},
		wg:                &sync.WaitGroup{},
		inputMap:          inputMap,
		speed:             0.5,
	}, nil
}

func (g *graphics) startup() {

	log.Println("Graphics, startup; starting")

	if g.shutdowns == nil {
		g.shutdowns = make(chan struct{})
		g.wg.Add(1)
		go g.runMainLoop()
	}

	log.Println("Graphics, startup; started")
}

func (g *graphics) shutdown() {

	log.Println("Graphics, shutdown; shutting down")

	if g.shutdowns != nil {
		close(g.shutdowns)
		g.wg.Wait()
		g.shutdowns = nil
	}
	log.Println("Graphics, shutdown; finished")
}

func (g *graphics) runMainLoop() {
	for {
		err := g.runGraphicsLoop()
		if err != nil {
			log.Println(fmt.Sprintf("Graphics, Main Loop: received error; %s", err.Error()))
		}
		select {
		case _, ok := <-g.shutdowns:
			if !ok {
				goto CloseMainLoop
			}
		case <-time.After(5 * time.Second):
			log.Println("Graphics, Main Loop: retrying window")
		}
	}

CloseMainLoop:
	log.Println("Graphics runMainLoop, Main Loop: closed")
	g.wg.Done()
}

func (g *graphics) stepTime() {
	g.mu.Lock()
	defer g.mu.Unlock()
	nt := time.Now()
	timeMultiplier := 32.0*math.Pow(g.speed, 4.0) -
		(45+1/3)*math.Pow(g.speed, 3.0) +
		20*math.Pow(g.speed, 2.0) -
		(1+2/3)*g.speed
	elapsed := nt.Sub(g.lastTimeStep).Seconds() * timeMultiplier
	g.inputMap["time"] += float32(elapsed)
	g.lastTimeStep = nt
}

func (g *graphics) runGraphicsLoop() error {

	grid := g.bus.GetLightGrid()
	gridWidth := grid.MaxX - grid.MinX + 1
	gridHeight := grid.MaxY - grid.MinY + 1

	gridWidth = gridWidth * g.pixelSize
	gridHeight = gridHeight * g.pixelSize

	g.mu.Lock()
	g.pb = util.NewPixelBuffer(gridWidth, gridHeight, grid.MinX, grid.MinY, g.pixelSize)
	g.mu.Unlock()

	g.lastTimeStep = time.Now()
	g.inputMap["time"] = 0.0

	gs, err := g.setupGraphicsShader(int32(gridWidth), int32(gridHeight))
	if err != nil {
		return err
	}

	g.gs = gs

	ticker := time.NewTicker(g.graphicsFrequency)

	defer func(g *graphics, t *time.Ticker) {
		t.Stop()
		g.gs.Cleanup()
		g.gs = nil
		g.mu.Lock()
		g.pb.BlackOut()
		g.mu.Unlock()
	}(g, ticker)

	err = func() error {
		for {
			select {
			case _, ok := <-g.shutdowns:
				if !ok {
					return nil
				}
			case <-ticker.C:
				g.stepTime()
				if g.reloadOnUpdate {
					err = g.gs.ReloadShader()
					if err != nil {
						return err
					}
				}
				err = g.doRunShader()
				if err != nil {
					return err
				}
				err = g.gs.ReadToPixels(g.pb.GetUnsafePointer())
				if err != nil {
					return err
				}
			}
		}
	}()

	g.gs = nil
	return err

}

func (g *graphics) doRunShader() error {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.gs.RunShader()
}

func (g *graphics) getShaderPath() (string, error) {
	basePath, err := os.Getwd()
	if err != nil {
		return "", errors.New("COULDN'T GET CWD")
	}
	if !strings.Contains(basePath, domain.Program) {
		return "", errors.New(fmt.Sprintf("PATH DOES NOT INCLUDE PROGRAM %s", domain.Program))
	}
	dataPath := strings.Split(basePath, domain.Program)[0]
	dataPath = filepath.Join(dataPath, domain.Program, "data")
	if _, err := os.Stat(dataPath); errors.Is(err, os.ErrNotExist) {
		return "", errors.New(fmt.Sprintf("PATH DOES NOT EXIST: %s", dataPath))
	}
	return dataPath, nil
}

func (g *graphics) setupGraphicsShader(width int32, height int32) (*graphicsShader.GraphicsShader, error) {
	path, err := g.getShaderPath()
	if err != nil {
		return nil, err
	}
	gs, err := graphicsShader.NewGraphicsShader(path, width, height, g.inputMap, g.mu)
	if err != nil {
		return nil, err
	}
	for i, s := range g.shaderFiles {
		err = gs.AttachShader(graphicsShader.ShaderKey(rune(i)), s)
		if err != nil {
			return nil, err
		}
	}

	return gs, nil
}
