package main

import (
	"github.com/TheOrnyx/psx-go/emulator"
	"runtime"

	"github.com/TheOrnyx/psx-go/cdrom"
	"github.com/TheOrnyx/psx-go/cpu"
	"github.com/TheOrnyx/psx-go/gpu"
	"github.com/TheOrnyx/psx-go/log"
	"github.com/TheOrnyx/psx-go/memory"
	"github.com/TheOrnyx/psx-go/renderer"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	runtime.LockOSThread()
	bios, err := memory.NewBios("./data/SCPH1001.BIN")
	if err != nil {
		log.Panicf("Failed to create Bios: %v", err)
	}

	err = sdl.Init(sdl.INIT_VIDEO)
	if err != nil {
		log.Panicf("Failed to initialize SDL: %v", err)
	}

	renderer, err := renderer.NewRenderer()
	if err != nil {
		sdl.Quit()
		log.Panicf("Failed to initialize renderer: %v", err)
	}

	gpu := gpu.NewGPU(renderer)
	defer gpu.Quit()

	cdrom, err := cdrom.NewCDROM("./data/Roms/tests/PeterLemon/HelloWorld/16BPP/HelloWorld16BPP.exe")
	if err != nil {
		gpu.Quit()
		log.Panicf("Failed to create CDROM: %v", err)
	}

	bus := memory.NewBus(bios, &gpu, &cdrom)
	cpu := cpu.NewCPU(bus)

	emu := emulator.Emulator{
		Cpu:      cpu,
		Gpu:      &gpu,
		Renderer: renderer,
		Bus:      bus,
		Cdrom:    &cdrom,
	}
	defer emu.Quit()

	for {
		for i := 0; i < 1000000; i++ {
			emu.Step()
		}

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				return

			case *sdl.KeyboardEvent:
				if t.Type == sdl.KEYDOWN {
					keyCode := t.Keysym.Sym

					if keyCode == sdl.K_ESCAPE {
						return
					}
				}
			}
		}
	}
}
