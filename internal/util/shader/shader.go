package shader

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

func GetShaderQualifiedPath(shaderName string, programName string) (string, error) {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)
	if !strings.Contains(basePath, programName) {
		return "", errors.New(fmt.Sprintf("PATH DOES NOT INCLUDE PROGRAM %s", programName))
	}
	dataPath := strings.Split(basePath, programName)[0]
	dataPath = filepath.Join(dataPath, programName, "data", shaderName)
	fragPath := dataPath + ".frag"
	if _, err := os.Stat(fragPath); errors.Is(err, os.ErrNotExist) {
		return "", errors.New("COULDN'T FIND FRAGMENT SHADER")
	}
	vertPath := dataPath + ".vert"
	if _, err := os.Stat(vertPath); errors.Is(err, os.ErrNotExist) {
		return "", errors.New("COULDN'T FIND VERTEX SHADER")
	}
	return dataPath, nil
}

type shader struct {
	handle uint32
}

type program struct {
	handle     uint32
	shaders    []shader
	rectHandle uint32
	width      float32
	height     float32
}

type GraphicsShader struct {
	shaderPath  string
	width       int32
	height      int32
	window      *windowProxy
	program     *program
	uniformDict map[string]float32
	mu          *sync.RWMutex
}

func NewGraphicsShader(
	shaderPath string, width int, height int, uniforms map[string]float32, mu *sync.RWMutex,
) (*GraphicsShader, error) {

	runtime.LockOSThread()

	gs := &GraphicsShader{
		shaderPath: shaderPath,
		width:      int32(width),
		height:     int32(height),
		mu:         mu,
	}

	gs.uniformDict = uniforms

	err := glfwInit()
	if err != nil {
		gs.Cleanup()
		log.Fatalln("failed to inifitialize glfw:", err)
		return nil, err
	}

	window, err := newWindow(shaderPath, width, height)
	if err != nil {
		gs.Cleanup()
		log.Fatalln("failed to create glfw window:", err)
		return nil, err
	}
	window.MakeContextCurrent()

	gs.window = window

	err = glInit()
	if err != nil {
		gs.Cleanup()
		log.Fatalln("failed to create gl context:", err)
		return nil, err
	}

	p, err := newProgram(shaderPath, float32(width), float32(height))
	if err != nil {
		gs.Cleanup()
		return nil, err
	}

	gs.program = p

	window.SetKeyCallback(windowKeyCallback)

	return gs, nil
}

func (gs *GraphicsShader) ReloadShader() error {
	if gs.program != nil {
		gs.program.delete()
	}
	p, err := newProgram(gs.shaderPath, float32(gs.width), float32(gs.height))
	if err != nil {
		return err
	}
	gs.program = p
	return nil
}

func (gs *GraphicsShader) RunShader() error {

	if gs.window.ShouldClose() {
		return errors.New("force close window")
	}

	stepGraphics()

	err := gs.program.runProgram(gs.uniformDict, gs.mu)
	if err != nil {
		return errors.New("shader failed")
	}

	gs.window.SwapBuffers()

	return nil
}

func (gs *GraphicsShader) DisplayShader() error {
	return nil
}

func (gs *GraphicsShader) Cleanup() {
	if gs.program != nil {
		gs.program.delete()
	}
	glfwTerminate()
	runtime.UnlockOSThread()
}

type getObjIv func(uint32, uint32, *int32)
type getObjInfoLog func(uint32, int32, *int32, *uint8)
