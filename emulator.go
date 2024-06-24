package main

import (
	"github.com/TheOrnyx/psx-go/cdrom"
	"github.com/TheOrnyx/psx-go/cpu"
	"github.com/TheOrnyx/psx-go/gpu"
	"github.com/TheOrnyx/psx-go/memory"
	"github.com/TheOrnyx/psx-go/renderer"
)

// Emulator - Basic struct for holding all the components of the emulator
type Emulator struct {
	Cpu      *cpu.CPU
	Gpu      *gpu.Gpu
	Renderer *renderer.Renderer
	Bus      *memory.Bus
	Cdrom    *cdrom.CDROM
}

// Step - step emulator once
func (e *Emulator) Step() {
	e.Cpu.RunNextInstruction()
}

// Quit - Quit the emulator and cleanup it's stuff
func (e *Emulator) Quit() {
	e.Gpu.Quit()
}
