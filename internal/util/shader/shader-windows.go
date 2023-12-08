//go:build windows
// +build windows

package shader

import (
	"fmt"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"sync"
	"unsafe"
)

func glfwInit() error {
	return glfw.Init()
}

func glfwTerminate() {
	glfw.Terminate()
}

func glInit() error {
	return gl.Init()
}

func stepGraphics() {
	glfw.PollEvents()
	gl.ClearColor(0.2, 0.2, 0.2, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

type windowProxy struct {
	*glfw.Window
}

func newWindow(shaderPath string, width int, height int) (*windowProxy, error) {
	glfw.WindowHint(glfw.Resizable, glfw.False)
	_, name := filepath.Split(shaderPath)
	window, err := glfw.CreateWindow(width, height, name, nil, nil)
	if err != nil {
		return nil, err
	}
	return &windowProxy{window}, nil
}

func windowKeyCallback(
	window *glfw.Window, key glfw.Key, _ int, action glfw.Action, _ glfw.ModifierKey,
) {
	if key == glfw.KeyEscape && action == glfw.Press {
		window.SetShouldClose(true)
	}
}

func newProgram(fileBase string, width float32, height float32) (*program, error) {
	vertexShader, err := newShaderFromFile(fileBase+".vert", gl.VERTEX_SHADER)
	if err != nil {
		log.Println("Couldn't compile vertex shader")
		return nil, err
	}
	fragmentShader, err := newShaderFromFile(fileBase+".frag", gl.FRAGMENT_SHADER)
	if err != nil {
		log.Println("Couldn't compile fragment shader")
		return nil, err
	}
	prog := &program{
		handle: gl.CreateProgram(),
	}
	prog.attach(*vertexShader, *fragmentShader)

	err = prog.link()
	if err != nil {
		log.Println("Couldn't link shaders")
		return nil, err
	}

	prog.rectHandle = createFillRect()
	prog.width = width
	prog.height = height

	return prog, nil
}

func (p *program) delete() {
	for _, s := range p.shaders {
		s.delete()
	}
	gl.DeleteProgram(p.handle)
}

func (p *program) attach(shaders ...shader) {
	for _, s := range shaders {
		gl.AttachShader(p.handle, s.handle)
		p.shaders = append(p.shaders, s)
	}
}

func (p *program) use() {
	gl.UseProgram(p.handle)
}

func (p *program) link() error {
	gl.LinkProgram(p.handle)
	return getGlError(p.handle, gl.LINK_STATUS, gl.GetProgramiv, gl.GetProgramInfoLog,
		"PROGRAM::LINKING_FAILURE")
}

func (p *program) runProgram(uniformDict map[string]float32, mu *sync.RWMutex) error {
	p.use()
	p.setUniform2fv("resolution", []float32{p.width, p.height}, 1)
	mu.RLock()
	for u, v := range uniformDict {
		p.setUniform1f(u, v)
	}
	mu.RUnlock()
	gl.BindVertexArray(p.rectHandle)
	gl.DrawElements(gl.TRIANGLE_FAN, 4, gl.UNSIGNED_INT, unsafe.Pointer(nil))
	gl.BindVertexArray(0)
	// should probably check for an error here, not sure what tho
	return nil
}

func (p *program) setUniform1f(name string, value float32) {
	chars := []uint8(name)
	loc := gl.GetUniformLocation(p.handle, &chars[0])
	if loc == -1 {
		// log.Println(fmt.Sprintf("Couldn't find uniform 1f %s", name))
		return
	}
	gl.Uniform1f(loc, value)
}

func (p *program) setUniform2fv(name string, value []float32, count int32) {
	chars := []uint8(name)
	loc := gl.GetUniformLocation(p.handle, &chars[0])
	if loc == -1 {
		// log.Println(fmt.Sprintf("Couldn't find uniform 2f %s", name))
		return
	}
	gl.Uniform2fv(loc, count, &value[0])
}

func newShaderFromFile(file string, sType uint32) (*shader, error) {
	src, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	handle := gl.CreateShader(sType)
	glSrc, freeFn := gl.Strs(string(src) + "\x00")
	defer freeFn()
	gl.ShaderSource(handle, 1, glSrc, nil)
	gl.CompileShader(handle)
	err = getGlError(handle, gl.COMPILE_STATUS, gl.GetShaderiv, gl.GetShaderInfoLog,
		"SHADER::COMPILE_FAILURE::"+file)
	if err != nil {
		return nil, err
	}
	return &shader{handle: handle}, nil
}

func (s *shader) delete() {
	gl.DeleteShader(s.handle)
}

func createFillRect() uint32 {

	vertices := []float32{
		// bottom left
		-1.0, -1.0, 0.0, // position
		0.0, 0.0, 0.0, // Color

		// bottom right
		1.0, -1.0, 0.0,
		0.0, 0.0, 0.0,

		// top right
		1.0, 1.0, 0.0,
		0.0, 0.0, 0.0,

		// top left
		-1.0, 1.0, 0.0,
		0.0, 0.0, 0.0,
	}

	indices := []uint32{
		0, 1, 2, 3,
	}

	var VAO uint32
	gl.GenVertexArrays(1, &VAO)

	var VBO uint32
	gl.GenBuffers(1, &VBO)

	var EBO uint32
	gl.GenBuffers(1, &EBO)

	// Bind the Vertex Array Object first, then bind and set vertex buffer(s) and attribute pointers()
	gl.BindVertexArray(VAO)

	// copy vertices data into VBO (it needs to be bound first)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	// copy indices into element buffer
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	// position
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// unbind the VAO (safe practice so we don't accidentally (mis)configure it later)
	gl.BindVertexArray(0)

	// should probably check for an error here, not sure what tho

	return VAO
}

func (gs *GraphicsShader) ReadToPixels(pb unsafe.Pointer) error {
	gl.ReadPixels(0, 0, gs.width, gs.height, gl.RGBA, gl.UNSIGNED_BYTE, pb)
	// should probably check for error
	return nil
}

func getGlError(glHandle uint32, checkTrueParam uint32, getObjIvFn getObjIv,
	getObjInfoLogFn getObjInfoLog, failMsg string) error {

	var success int32
	getObjIvFn(glHandle, checkTrueParam, &success)

	if success == gl.FALSE {
		var logLength int32
		getObjIvFn(glHandle, gl.INFO_LOG_LENGTH, &logLength)

		outMsg := gl.Str(strings.Repeat("\x00", int(logLength)))
		getObjInfoLogFn(glHandle, logLength, nil, outMsg)

		return fmt.Errorf("%s: %s", failMsg, gl.GoStr(outMsg))
	}

	return nil
}
