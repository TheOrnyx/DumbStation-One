package main

import (
	"log"

	"github.com/TheOrnyx/psx-go/cpu"
	"github.com/TheOrnyx/psx-go/memory"
)

func main() {
	bios, err := memory.NewBios("./data/SCPH1001.BIN")
	if err != nil {
		log.Panic(err)
	}

	bus := memory.NewBus(bios)
	cpu := cpu.NewCPU(bus)
	
	for  {
		cpu.RunNextInstruction()
	}
}
