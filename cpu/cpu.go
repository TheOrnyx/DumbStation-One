package cpu

import (
	"github.com/TheOrnyx/psx-go/log"
	"github.com/TheOrnyx/psx-go/memory"
)

// The CPU struct
// NOTES -
// All registers are 32-bits wide (uint32)
type CPU struct {
	pc              uint32    // The program counter
	regs            Registers // the rest of the registers
	bus             *memory.Bus // the memory bus
	nextInstruction Instruction // the next instruction, used to simulate branch delay shot
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
	cpu.pc = 0xbfc00000           // reset to beginning of the BIOS
	cpu.regs = Registers{zero: 0} // TODO - probably reset to garbage but idc
	cpu.nextInstruction = Instruction(0x0) // NOP
	log.Info("Reset CPU state")
}

// GetReg get the register at index - just calls the reg's method
func (cpu *CPU) GetReg(index uint32) uint32 {
	return cpu.regs.GetReg(index)
}

// SetReg sets the register at index to val - just calls the reg's method
func (cpu *CPU) SetReg(index, val uint32) {
	cpu.regs.SetReg(index, val)
}

// RunNextInstruction run the next instruction
func (cpu *CPU) RunNextInstruction() {
	instruction := cpu.nextInstruction
	cpu.nextInstruction = Instruction(cpu.load32(cpu.pc))
	
	cpu.pc += 4 // increment pc by 4 bytes
	cpu.decodeAndExecuteInstr(instruction)
}

// decodeAndExecuteInstr decode and execute an instruction
// TODO - switch from binary to hex cuz nicer
func (cpu *CPU) decodeAndExecuteInstr(instruction Instruction) {
	switch instruction.function() {
	case 0x00: // SPECIAL
		cpu.executeSubInstr(instruction)
	case 0x02: // J
		cpu.jump(instruction)
	case 0x0f: // LUID
		cpu.loadUpperImmediate(instruction)
	case 0x0d: // ORI
		cpu.orImmediate(instruction)
	case 0x2b: // SW
		cpu.storeWord(instruction)
	case 0x09: // ADDIU
		cpu.addImmediateUnsigned(instruction)
	default:
		log.Panicf("Unknown instruction - 0x%08x", instruction)
	}
}

// executeSubInstr decode and execute sub instruction (special)
//
// NOTE - This is a seperate function cuz I didn't wanna have a like nested switch
func (cpu *CPU) executeSubInstr(instruction Instruction) {
	switch instruction.subFunction() {
	case 0x00: // shift left logical
		cpu.shiftLeftLogical(instruction)
	default:
		log.Panicf("Unknown sub instruction - 0x%08x", instruction)
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
func (cpu *CPU) Store32(addr, val uint32) {
	err := cpu.bus.Store32(addr, val)
	if err != nil {
		// log.Infof("%+v", cpu.regs)
		log.Panicf("Store32 Failed - %v", err)
	}
}
