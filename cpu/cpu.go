package cpu

import (
	"github.com/TheOrnyx/psx-go/log"
	"github.com/TheOrnyx/psx-go/memory"
)

// The CPU struct
// NOTES -
// All registers are 32-bits wide (uint32)
type CPU struct {
	pc               uint32           // The program counter
	nextPC           uint32           // the next value for the PC - used for simulating branch delay slot
	regs             Registers        // the rest of the registers
	outRegs          Registers        // second set of registers used to emulate load delay slot - this sucks
	loadReg          LoadRegPair      // the pair to use for loading
	copZeroRegs      CopZeroRegisters // Coprocessor zero's registers
	bus              *memory.Bus      // the memory bus
	nextInstruction  Instruction      // the next instruction, used to simulate branch delay shot
	hi               uint32           // HI register for division remainder and multiplication high result
	lo               uint32           // LO register for division quotient and multiplication low result
	currentPC        uint32           // address of instruction currently being executed. used for setting EPC in exceptions
	branching        bool             // set by current instruction if branch occured and next instruction will be in delay slot
	instrInDelaySlot bool             // set if the current instruction executes in the delay slot
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
	cpu.nextPC = cpu.pc + 4
	cpu.branching = false
	cpu.instrInDelaySlot = false
	cpu.regs = Registers{zero: 0} // TODO - probably reset to garbage but idc
	cpu.outRegs = cpu.regs
	cpu.loadReg = LoadRegPair{0, 0}
	cpu.copZeroRegs = CopZeroRegisters{}
	cpu.nextInstruction = Instruction(0x0) // NOP
	cpu.hi = 0xbeaf
	cpu.lo = 0xfeab
	log.Info("Reset CPU state")
}

// GetReg get the register at index - just calls the reg's method
func (cpu *CPU) GetReg(index RegIndex) uint32 {
	return cpu.regs.GetReg(index)
}

// SetReg sets the register at index to val - just calls the reg's method
func (cpu *CPU) SetReg(index RegIndex, val uint32) {
	cpu.outRegs.SetReg(index, val)
}

// GetCopZeroReg get cop zero reg and handle logging if need be
func (cpu *CPU) GetCopZeroReg(index RegIndex) uint32 {
	val, _ := cpu.copZeroRegs.GetReg(index)
	return val
}

// SetCopZeroReg set cop zero reg and handle logging if need be
func (cpu *CPU) SetCopZeroReg(index RegIndex, val uint32) {
	_ = cpu.copZeroRegs.SetReg(index, val)
}

// SetLoadReg set the loadReg to index and val
func (cpu *CPU) SetLoadReg(index RegIndex, val uint32)  {
	cpu.loadReg.target = index
	cpu.loadReg.val = val
}

// RunNextInstruction run the next instruction
func (cpu *CPU) RunNextInstruction() {
	instruction := Instruction(cpu.load32(cpu.pc))

	cpu.instrInDelaySlot = cpu.branching
	cpu.branching = false

	// save address of current instruction
	cpu.currentPC = cpu.pc

	// increment next pc to point to next instruction
	cpu.pc = cpu.nextPC
	cpu.nextPC += 4

	// execute pending load
	cpu.SetReg(cpu.loadReg.target, cpu.loadReg.val)

	// reset load to target 0 for next instr
	cpu.SetLoadReg(0, 0)

	cpu.decodeAndExecuteInstr(instruction)
	
	// set the regs to the outregs
	// FIXME - optimize later
	cpu.regs = cpu.outRegs
}

// decodeAndExecuteInstr decode and execute an instruction
// TODO - switch from binary to hex cuz nicer
func (cpu *CPU) decodeAndExecuteInstr(instruction Instruction) {
	if instruction.function() > 0x3f {
		UnknownOpcode.runFunc(cpu, instruction)
	}
	
	Opcodes[instruction.function()].runFunc(cpu, instruction)
}

// copZeroOpcode Coprocessor 0 opcode
func (cpu *CPU) copZeroOpcode(instruction Instruction) {
	switch instruction.copOpcode() {
	case 0b00000: // MFC0
		cpu.moveFromCopZero(instruction)
	case 0b00100: // MTC0
		cpu.moveToCopZero(instruction)
	case 0b10000: // RFE
		cpu.returnFromException(instruction)
	default:
		log.Panicf("Unknown cop zero instruction - 0x%08x, 0x%02x", instruction, instruction.copOpcode())
	}
}

// executeSubInstr decode and execute sub instruction (special)
func (cpu *CPU) executeSubInstr(instruction Instruction) {
	if instruction.subFunction() > 0x3f {
		UnknownSubOpcode.runFunc(cpu, instruction)
	}

	SubOpcodes[instruction.subFunction()].runFunc(cpu, instruction)
}

// load32 Load and return the value at given address addr
func (cpu *CPU) load32(addr uint32) uint32 {
	data, err := cpu.bus.Load32(addr)
	if err != nil {
		log.Fatalf("Load32 failed - %v", err)
	}

	return data
}

// Load16 load 16-bit halfword
func (cpu *CPU) Load16(addr uint32) uint16 {
	data, err := cpu.bus.Load16(addr)
	if err != nil {
		log.Fatalf("Load16 failed - %v", err)
	}

	return data
}

// Load8 load 8 bit val from memory
func (cpu *CPU) Load8(addr uint32) uint8 {
	val, err := cpu.bus.Load8(addr)
	if err != nil {
		log.Fatalf("Load8 failed - %v", err)
	}

	return val
}

// Store32 store given value val into address addr
func (cpu *CPU) Store32(addr, val uint32) {
	err := cpu.bus.Store32(addr, val)
	if err != nil {
		// log.Infof("%+v", cpu.regs)
		log.Panicf("Store32 Failed - %v", err)
	}
}

// Store16 store given 16bit value into memory
func (cpu *CPU) Store16(addr uint32, val uint16) {
	err := cpu.bus.Store16(addr, val)
	if err != nil {
		log.Panicf("Store16 Failed - %v", err)
	}
}

// Store8 store 8-bit byte into memory
func (cpu *CPU) Store8(addr uint32, val uint8) {
	err := cpu.bus.Store8(addr, val)
	if err != nil {
		log.Panicf("Store8 Failed - %v", err)
	}
}

// branchingSetup set some common settigns for branching
func (cpu *CPU) branchingSetup()  {
	cpu.branching = true

	if cpu.nextPC % 4 != 0 { // TODO - check whether this should be nextPC
		// PC is not correctly aligned
		cpu.Exception(LoadAddressError)
	}
}

// branch branch to the immediate value offset - basic branch used by
// other instructions
func (cpu *CPU) branch(offset uint32) {
	offset = offset << 2
	cpu.nextPC += offset

	cpu.nextPC -= 4 // have to compensate for the += 4 in the run next instr
	cpu.branching = true
	cpu.branchingSetup()
}

// jump jump
func (cpu *CPU) jump(instr Instruction) {
	immediate := instr.jumpImmediate()

	cpu.nextPC = (cpu.nextPC & 0xf0000000) | (immediate << 2)
	cpu.branching = true
	cpu.branchingSetup()
}

// exception enums
const (
	SysCall  = 0x8
	Overflow = 0xc
	LoadAddressError = 0x4
	StoreAddressError = 0x5
)

// Exception Trigger an exception
//
// TODO - recheck this later
func (cpu *CPU) Exception(cause int) {
	// exception handler address depends on 'BEV' bit:
	sr := cpu.GetCopZeroReg(12)
	var handler uint32 = 0x80000080
	if sr&(1<<22) != 0 {
		handler = 0xbfc00180
	}

	// shift bits [5:0] of SR two places to left.  These bits are
	// threee pairs of interrupt enable/ user mode bits that behave
	// like a stack 3 entires deep. Entering an exception pushes a
	// pair of zeroes by left shifting the stack which disabled
	// interrutpts and puts the CU in kernel mode. The original third
	// entry is discarded.
	mode := sr & 0x3f
	sr &= ^uint32(0x3f) // TODO - check
	sr |= (mode << 2) & 0x3f

	// write sr back
	cpu.SetCopZeroReg(12, sr)

	cpu.copZeroRegs.cause = uint32(cause) << 2
	cpu.copZeroRegs.epc = cpu.currentPC

	if cpu.instrInDelaySlot {
		// when exception occurs in delay slot 'epc' points to branch
		// instruction and bit31 of cause is set
		cpu.copZeroRegs.epc -= 4
		cpu.copZeroRegs.cause |= 1 << 31
	}

	cpu.pc = handler
	cpu.nextPC = cpu.pc + 4
}
