//go:build !windows
// +build !windows

package shader

import (
	"fmt"
	"github.com/go-gl/gl/v3.1/gles2"
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
	return gles2.Init()
}

func stepGraphics() {
	glfw.PollEvents()
	gles2.ClearColor(0.2, 0.2, 0.2, 1.0)
	gles2.Clear(gles2.COLOR_BUFFER_BIT)
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
	vertexShader, err := newShaderFromFile(fileBase+".vert", gles2.VERTEX_SHADER)
	if err != nil {
		log.Println("Couldn't compile vertex shader")
		return nil, err
	}
	fragmentShader, err := newShaderFromFile(fileBase+".frag", gles2.FRAGMENT_SHADER)
	if err != nil {
		log.Println("Couldn't compile fragment shader")
		return nil, err
	}
	prog := &program{
		handle: gles2.CreateProgram(),
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
	gles2.DeleteProgram(p.handle)
}

func (p *program) attach(shaders ...shader) {
	for _, s := range shaders {
		gles2.AttachShader(p.handle, s.handle)
		p.shaders = append(p.shaders, s)
	}
}

func (p *program) use() {
	gles2.UseProgram(p.handle)
}

func (p *program) link() error {
	gles2.LinkProgram(p.handle)
	return getGlError(p.handle, gles2.LINK_STATUS, gles2.GetProgramiv, gles2.GetProgramInfoLog,
		"PROGRAM::LINKING_FAILURE")
}

func (p *program) runProgram(uniforms map[string]float32, mu *sync.RWMutex) error {
	p.use()
	p.setUniform2fv("resolution", []float32{p.width, p.height}, 1)
	mu.RLock()
	for u, v := range uniforms {
		p.setUniform1f(u, v)
	}
	mu.RUnlock()
	gles2.BindVertexArray(p.rectHandle)
	gles2.DrawElements(gles2.TRIANGLE_FAN, 4, gles2.UNSIGNED_INT, unsafe.Pointer(nil))
	gles2.BindVertexArray(0)
	// should probably check for an error here, not sure what tho
	return nil
}

func (p *program) setUniform1f(name string, value float32) {
	chars := []uint8(name)
	loc := gles2.GetUniformLocation(p.handle, &chars[0])
	if loc == -1 {
		// log.Println(fmt.Sprintf("Couldn't find uniform 1f %s", name))
		return
	}
	gles2.Uniform1f(loc, value)
}

func (p *program) setUniform2fv(name string, value []float32, count int32) {
	chars := []uint8(name)
	loc := gles2.GetUniformLocation(p.handle, &chars[0])
	if loc == -1 {
		// log.Println(fmt.Sprintf("Couldn't find uniform 2f %s", name))
		return
	}
	gles2.Uniform2fv(loc, count, &value[0])
}

func newShaderFromFile(file string, sType uint32) (*shader, error) {
	src, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	handle := gles2.CreateShader(sType)
	glSrc, freeFn := gles2.Strs(string(src) + "\x00")
	defer freeFn()
	gles2.ShaderSource(handle, 1, glSrc, nil)
	gles2.CompileShader(handle)
	err = getGlError(handle, gles2.COMPILE_STATUS, gles2.GetShaderiv, gles2.GetShaderInfoLog,
		"SHADER::COMPILE_FAILURE::"+file)
	if err != nil {
		return nil, err
	}
	return &shader{handle: handle}, nil
}

func (s *shader) delete() {
	gles2.DeleteShader(s.handle)
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
	gles2.GenVertexArrays(1, &VAO)

	var VBO uint32
	gles2.GenBuffers(1, &VBO)

	var EBO uint32
	gles2.GenBuffers(1, &EBO)

	// Bind the Vertex Array Object first, then bind and set vertex buffer(s) and attribute pointers()
	gles2.BindVertexArray(VAO)

	// copy vertices data into VBO (it needs to be bound first)
	gles2.BindBuffer(gles2.ARRAY_BUFFER, VBO)
	gles2.BufferData(gles2.ARRAY_BUFFER, len(vertices)*4, gles2.Ptr(vertices), gles2.STATIC_DRAW)

	// copy indices into element buffer
	gles2.BindBuffer(gles2.ELEMENT_ARRAY_BUFFER, EBO)
	gles2.BufferData(gles2.ELEMENT_ARRAY_BUFFER, len(indices)*4, gles2.Ptr(indices), gles2.STATIC_DRAW)

	// position
	gles2.VertexAttribPointer(0, 3, gles2.FLOAT, false, 6*4, gles2.PtrOffset(0))
	gles2.EnableVertexAttribArray(0)

	// unbind the VAO (safe practice so we don't accidentally (mis)configure it later)
	gles2.BindVertexArray(0)

	// should probably check for an error here, not sure what tho

	return VAO
}

func (gs *GraphicsShader) ReadToPixels(pb unsafe.Pointer) error {
	gles2.ReadPixels(0, 0, gs.width, gs.height, gles2.RGBA, gles2.UNSIGNED_BYTE, pb)
	// should probably check for error
	return nil
}

func getGlError(glHandle uint32, checkTrueParam uint32, getObjIvFn getObjIv,
	getObjInfoLogFn getObjInfoLog, failMsg string) error {

	var success int32
	getObjIvFn(glHandle, checkTrueParam, &success)

	if success == gles2.FALSE {
		var logLength int32
		getObjIvFn(glHandle, gles2.INFO_LOG_LENGTH, &logLength)

		outMsg := gles2.Str(strings.Repeat("\x00", int(logLength)))
		getObjInfoLogFn(glHandle, logLength, nil, outMsg)

		return fmt.Errorf("%s: %s", failMsg, gles2.GoStr(outMsg))
	}

	return nil
}
