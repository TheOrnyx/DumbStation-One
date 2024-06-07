package cpu

import (
	"github.com/TheOrnyx/psx-go/memory"
	"github.com/TheOrnyx/psx-go/log"
)

// The CPU struct
// NOTES -
// All registers are 32-bits wide (uint32)
type CPU struct {
	pc   uint32    // The program counter
	regs Registers // the rest of the registers
	bus  *memory.Bus
}

// NewCPU Create and return a new CPU that's been reset
func NewCPU(bus *memory.Bus) *CPU {
	cpu := new(CPU)
	cpu.bus = bus
	cpu.Reset()

	return cpu
}

// Reset reset the cpu to its initial state
func (cpu *CPU) Reset() {
	cpu.pc = 0xbfc00000 // reset to beginning of the BIOS
	cpu.regs = Registers{zero: 0} // TODO - probably reset to garbage but idc
	log.Info("Reset CPU state")
}

// RunNextInstruction run the next instruction
func (cpu *CPU) RunNextInstruction() {
	instruction := cpu.load32(cpu.pc)
	cpu.pc += 4 // increment pc by 4 bytes
	cpu.decodeAndExecuteInstr(Instruction(instruction))
}

// decodeAndExecuteInstr decode and execute an instruction
func (cpu *CPU) decodeAndExecuteInstr(instruction Instruction) {
	switch instruction.function() {
	case 0b001111:
		cpu.loadUpperImmediate(instruction)
	case 0b001101:
		cpu.orImmediate(instruction)
	case 0b101011:
		cpu.storeWord(instruction)
	default:
		
		log.Panicf("Unknown instruction - %X", instruction)
	}
}

// load32 Load and return the value at given address addr
func (cpu *CPU) load32(addr uint32) uint32 {
	data, err := cpu.bus.Load32(addr)
	if err != nil {
		log.Fatalf("%v", err)
	}

	return data
}

// Store32 store given value val into address addr
func (cpu *CPU) Store32(addr, val uint32)  {
	err := cpu.bus.Store32(addr, val)
	if err != nil {
		log.Panicf("Store32 Failed - %v", err)
	}
}
