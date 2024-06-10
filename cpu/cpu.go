package cpu

import (
	"github.com/TheOrnyx/psx-go/log"
	"github.com/TheOrnyx/psx-go/memory"
)

// The CPU struct
// NOTES -
// All registers are 32-bits wide (uint32)
type CPU struct {
	pc              uint32           // The program counter
	regs            Registers        // the rest of the registers
	outRegs         Registers        // second set of registers used to emulate load delay slot - this sucks
	loadReg         LoadRegPair      // the pair to use for loading
	copZeroRegs     CopZeroRegisters // Coprocessor zero's registers
	bus             *memory.Bus      // the memory bus
	nextInstruction Instruction      // the next instruction, used to simulate branch delay shot
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
	cpu.outRegs = cpu.regs
	cpu.loadReg = LoadRegPair{0, 0}
	cpu.copZeroRegs = CopZeroRegisters{}
	cpu.nextInstruction = Instruction(0x0) // NOP
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

// RunNextInstruction run the next instruction
func (cpu *CPU) RunNextInstruction() {
	instruction := cpu.nextInstruction
	cpu.nextInstruction = Instruction(cpu.load32(cpu.pc))
	cpu.pc += 4 // increment pc by 4 bytes

	// set for the delay slot
	cpu.SetReg(cpu.loadReg.target, cpu.loadReg.val)
	cpu.loadReg = LoadRegPair{0, 0} // TODO - check this is not inefficient compared to setting vals
	
	cpu.decodeAndExecuteInstr(instruction)

	// set the regs to the outregs
	// FIXME - optimize later
	cpu.regs = cpu.outRegs
}

// decodeAndExecuteInstr decode and execute an instruction
// TODO - switch from binary to hex cuz nicer
func (cpu *CPU) decodeAndExecuteInstr(instruction Instruction) {
	switch instruction.function() {
	case 0x00: // SPECIAL
		cpu.executeSubInstr(instruction)
	case 0x02: // J
		cpu.jump(instruction)
	case 0x03: // JAL
		cpu.jumpAndLink(instruction)
	case 0x04: // BEQ
		cpu.branchIfEqual(instruction)
	case 0x05: // BNE
		cpu.branchNotEqual(instruction)
	case 0x08: // ADDI
		cpu.addImmediate(instruction)
	case 0x0c: // ANDI
		cpu.andImmediate(instruction)
	case 0x10: // COP0
		cpu.copZeroOpcode(instruction)
	case 0x0f: // LUID
		cpu.loadUpperImmediate(instruction)
	case 0x0d: // ORI
		cpu.orImmediate(instruction)
	case 0x20: // LB
		cpu.loadByte(instruction)
	case 0x23: // LW
		cpu.loadWord(instruction)
	case 0x28: // SB
		cpu.storeByte(instruction)
	case 0x29: // SH
		cpu.storeHalfWord(instruction)
	case 0x2b: // SW
		cpu.storeWord(instruction)
	case 0x09: // ADDIU
		cpu.addImmediateUnsigned(instruction)
	default:
		log.Panicf("Unknown instruction - 0x%08x, 0x%02x", instruction, instruction.function())
	}
}

// copZeroOpcode Coprocessor 0 opcode
func (cpu *CPU) copZeroOpcode(instruction Instruction) {
	switch instruction.copOpcode() {
	case 0x00: // MFC0
		cpu.moveFromCopZero(instruction)
	case 0b00100: // MTC0
		cpu.moveToCopZero(instruction)
	default:
		log.Panicf("Unknown cop zero instruction - 0x%08x, 0x%02x", instruction, instruction.copOpcode())
	}
}

// executeSubInstr decode and execute sub instruction (special)
//
// NOTE - This is a seperate function cuz I didn't wanna have a like nested switch
func (cpu *CPU) executeSubInstr(instruction Instruction) {
	switch instruction.subFunction() {
	case 0x00: // SLL
		cpu.shiftLeftLogical(instruction)
	case 0x08: // JR
		cpu.jumpRegister(instruction)
	case 0x20: // ADD
		cpu.add(instruction)
	case 0x21: // ADDU
		cpu.addUnsigned(instruction)
	case 0x24: // AND
		cpu.and(instruction)
	case 0x25: // OR
		cpu.or(instruction)
	case 0x2b: // SLTU
		cpu.setOnLessThanUnsigned(instruction)
	default:
		log.Panicf("Unknown sub instruction - 0x%08x, 0x%02x", instruction, instruction.subFunction())
	}
}

// load32 Load and return the value at given address addr
func (cpu *CPU) load32(addr uint32) uint32 {
	data, err := cpu.bus.Load32(addr)
	if err != nil {
		log.Fatalf("Load32 failed - %v", err)
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
func (cpu *CPU) Store16(addr uint32, val uint16)  {
	err := cpu.bus.Store16(addr, val)
	if err != nil {
		log.Panicf("Store16 Failed - %v", err)
	}
}

// Store8 store 8-bit byte into memory
func (cpu *CPU) Store8(addr uint32, val uint8)  {
	err := cpu.bus.Store8(addr, val)
	if err != nil {
		log.Panicf("Store8 Failed - %v", err)
	}
}

// branch branch to the immediate value offset - basic branch used by
// other instructions
func (cpu *CPU) branch(offset uint32) {
	offset = offset << 2
	cpu.pc += offset

	cpu.pc -= 4 // have to compensate for the += 4 in the run next instr
}
