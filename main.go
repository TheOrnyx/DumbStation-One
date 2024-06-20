package main

import (
	"github.com/TheOrnyx/psx-go/cpu"
	"github.com/TheOrnyx/psx-go/gpu"
	"github.com/TheOrnyx/psx-go/log"
	"github.com/TheOrnyx/psx-go/memory"
	"github.com/TheOrnyx/psx-go/renderer"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
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
	
	bus := memory.NewBus(bios, &gpu)
	cpu := cpu.NewCPU(bus)
	
	for  {
		for i := 0; i < 1000000; i++ {
			cpu.RunNextInstruction()
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
