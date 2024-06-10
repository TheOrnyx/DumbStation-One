package cpu

import (
	"github.com/TheOrnyx/psx-go/log"
	"github.com/TheOrnyx/psx-go/utils"
)

type Instruction uint32

// function return the primary opcode field of the instruction (bits 26..31)
func (i Instruction) function() uint32 {
	return uint32(i) >> 26
}

// subFunction get the special function bits for the instruction from bits [5:0]
//
// called when function() is 0b000000
func (i Instruction) subFunction() uint32 {
	return uint32(i) & 0x3f
}

// targetReg get the target register of the current instruction
func (i Instruction) targetReg() RegIndex {
	return RegIndex((i >> 16) & 0x1f)
}

// sourceReg get the source reg of the current instruction (bits 21..25)
func (i Instruction) sourceReg() RegIndex {
	return RegIndex((i >> 21) & 0x1f)
}

// destReg get the index of the destination register from bits [15:11]
func (i Instruction) destReg() RegIndex {
	return RegIndex((i >> 11) & 0x1f)
}

// immediate16 return the immediate 16 bit value for the instruction
func (i Instruction) immediate16() uint32 {
	return uint32(i) & 0xffff
}

// immediate16Se return immediate value in bits [16:0] as a
// sign-extended 32 bit val
func (i Instruction) immediate16Se() uint32 {
	val := int16(i & 0xffff)
	return uint32(val)
}

// shiftImmediate return the immediate 5-bit value for the shift in bits [10:6]
func (i Instruction) shiftImmediate() uint32 {
	return (uint32(i) >> 6) & 0x1f
}

// jumpImmediate the jump target stored in bits [25:0]
func (i Instruction) jumpImmediate() uint32 {
	return uint32(i) & 0x3ffffff
}

// copOpcode the coprocessor opcode from bits [25:21]
//
// TODO - check how we know this is right??
func (i Instruction) copOpcode() uint32 {
	return (uint32(i) >> 21) & 0x1f
}

/////////////////////////////////////
// The CPU instructions themselves //
/////////////////////////////////////

// loadUpperImmediate Load upper immediate - loadUpperImmediate
func (cpu *CPU) loadUpperImmediate(instr Instruction)  {
	immediate := instr.immediate16()
	targetReg := instr.targetReg()

	cpu.SetReg(targetReg, immediate << 16)
}

// or bitwise or between two registers
func (cpu *CPU) or(instr Instruction)  {
	sourceReg := instr.sourceReg()
	targetReg := instr.targetReg()
	
	val := cpu.GetReg(sourceReg) | cpu.GetReg(targetReg)
	cpu.SetReg(instr.destReg(), val)
}

// orImmediate bitwise or immediate
func (cpu *CPU) orImmediate(instr Instruction)  {
	immediate := instr.immediate16()
	sourceReg := instr.sourceReg()
	val := cpu.GetReg(sourceReg) | immediate
	
	cpu.SetReg(instr.targetReg(), val)
}

// and bitwise and
func (cpu *CPU) and(instr Instruction)  {
	sourceReg := instr.sourceReg()
	targetReg := instr.targetReg()
	
	val := cpu.GetReg(sourceReg) & cpu.GetReg(targetReg)
	cpu.SetReg(instr.destReg(), val)
}

// andImmediate bitwise and immediate
func (cpu *CPU) andImmediate(instr Instruction)  {
	immediate := instr.immediate16()
	sourceReg := instr.sourceReg()
	val := cpu.GetReg(sourceReg) & immediate
	
	cpu.SetReg(instr.targetReg(), val)
}

// storeWord Store Word
func (cpu *CPU) storeWord(instr Instruction)  {
	if cpu.GetCopZeroReg(REG_SR) & 0x10000 != 0 {
		log.Info("Ignoring storeWord while the cache is isolated")
		return
	}
	
	targetReg := instr.targetReg()
	sourceReg := instr.sourceReg()

	addr := cpu.GetReg(sourceReg) + instr.immediate16Se()
	val := cpu.GetReg(targetReg)
	cpu.Store32(addr, val)
}

// storeHalfWord store half word into memory
func (cpu *CPU) storeHalfWord(instr Instruction)  {
	if cpu.GetCopZeroReg(REG_SR) & 0x10000 != 0 {
		log.Info("Ignoring storeHalfWord while the cache is isolated")
		return
	}

	immediate := instr.immediate16Se()
	targetReg := cpu.GetReg(instr.targetReg())
	sourceReg := cpu.GetReg(instr.sourceReg())
	addr := sourceReg + immediate
	
	cpu.Store16(addr, uint16(targetReg))
}

// storeByte store byte
func (cpu *CPU) storeByte(instr Instruction)  {
	if cpu.GetCopZeroReg(REG_SR) & 0x10000 != 0 {
		log.Info("Ignoring storeByte while the cache is isolated")
		return
	}

	immediate := instr.immediate16Se()
	targetReg := cpu.GetReg(instr.targetReg())
	sourceReg := cpu.GetReg(instr.sourceReg())
	addr := sourceReg + immediate
	
	cpu.Store8(addr, uint8(targetReg))
}

// loadWord load word
func (cpu *CPU) loadWord(instr Instruction)  {
	if cpu.GetCopZeroReg(REG_SR) & 0x10000 != 0 {
		log.Info("Ignoring loadWord while the cache is isolated")
		return
	}

	immediate := instr.immediate16Se()
	targetReg := instr.targetReg()
	sourceReg := instr.sourceReg()
	addr := cpu.GetReg(sourceReg) + immediate

	val := cpu.load32(addr)
	cpu.loadReg = LoadRegPair{targetReg, val} // TODO - once again check performance
}

// loadByte load signed byte
func (cpu *CPU) loadByte(instr Instruction)  {
	immediate := instr.immediate16Se()
	targetReg := instr.targetReg()
	sourceReg := cpu.GetReg(instr.sourceReg())

	addr := sourceReg + immediate
	val := int8(cpu.Load8(addr))

	cpu.loadReg = LoadRegPair{target: targetReg, val: uint32(val)}
}

// shiftLeftLogical Shift left logical
func (cpu *CPU) shiftLeftLogical(instr Instruction)  {
	immediate := instr.shiftImmediate()
	targetReg := instr.targetReg()

	val := cpu.GetReg(targetReg) << immediate
	cpu.SetReg(instr.destReg(), val)
}

// addImmediateUnsigned add immediate unsigned
//
// This one apparently doesn't generate an exception on overflow while ADDI does?
func (cpu *CPU) addImmediateUnsigned(instr Instruction)  {
	immediate := instr.immediate16Se()
	sourceReg := instr.sourceReg()

	val := cpu.GetReg(sourceReg) + immediate
	cpu.SetReg(instr.targetReg(), val)
}

// addImmediate add immediate - generates exception when overflows
func (cpu *CPU) addImmediate(instr Instruction)  {
	immediate := instr.immediate16Se()
	sourceReg := cpu.GetReg(instr.sourceReg())

	result, overflowed := utils.AddSigned16(sourceReg, immediate)
	if overflowed {
		log.Panicf("ADDI overflowed with imm:0x%08x reg:0x%08x", immediate, sourceReg)
	}

	cpu.SetReg(instr.targetReg(), result)
}

// add add and generate exception on overflow
func (cpu *CPU) add(instr Instruction)  {
	sourceReg := cpu.GetReg(instr.sourceReg())
	targetReg := cpu.GetReg(instr.targetReg())

	val, overflowed := utils.AddSigned16(sourceReg, targetReg)
	if overflowed {
		log.Panicf("ADD overflowed with source:0x%08x target:0x%08x", sourceReg, targetReg)
	}

	cpu.SetReg(instr.destReg(), val)
}

// addUnsigned Add two unsigned registers together and store in target
func (cpu *CPU) addUnsigned(instr Instruction)  {
	sourceReg := cpu.GetReg(instr.sourceReg())
	targetReg := cpu.GetReg(instr.targetReg())
	val := sourceReg + targetReg

	cpu.SetReg(instr.destReg(), val)
}

// jump jump
func (cpu *CPU) jump(instr Instruction)  {
	immediate := instr.jumpImmediate()

	cpu.pc = (cpu.pc & 0xf0000000) | (immediate << 2)
}

// jumpAndLink jump and link by storing return address in $ra
func (cpu *CPU) jumpAndLink(instr Instruction)  {
	returnAddr := cpu.pc

	cpu.SetReg(RegIndex(31), returnAddr)
	cpu.jump(instr)
}

// jumpRegister jump to val in register
func (cpu *CPU) jumpRegister(instr Instruction)  {
	sourceReg := cpu.GetReg(instr.sourceReg())
	cpu.pc = sourceReg
}

// branchNotEqual branch if not equal
func (cpu *CPU) branchNotEqual(instr Instruction)  {
	offset := instr.immediate16Se()
	sourceReg := instr.sourceReg()
	targetReg := instr.targetReg()

	if cpu.GetReg(sourceReg) != cpu.GetReg(targetReg) {
		cpu.branch(offset)
	}
}

// branchIfEqual branch if equal
func (cpu *CPU) branchIfEqual(instr Instruction)  {
	offset := instr.immediate16Se()
	sourceReg := instr.sourceReg()
	targetReg := instr.targetReg()

	if cpu.GetReg(sourceReg) == cpu.GetReg(targetReg) {
		cpu.branch(offset)
	}
}

// setOnLessThanUnsigned set the dest to 1 when source is less than target
func (cpu *CPU) setOnLessThanUnsigned(instr Instruction)  {
	sourceReg := cpu.GetReg(instr.sourceReg())
	targetReg := cpu.GetReg(instr.targetReg())
	var val uint32 = 0

	if sourceReg < targetReg {
		val = 1
	}

	cpu.SetReg(instr.destReg(), val)
}

/////////////////////////////////
// Coprocessor Instructions ✨ //
/////////////////////////////////

// moveToCopZero move register contents to coprocessor zero
func (cpu *CPU) moveToCopZero(instr Instruction)  {
	cpuReg := instr.targetReg()
	copReg := instr.destReg()

	val := cpu.GetReg(cpuReg)
	cpu.SetCopZeroReg(copReg, val)
}

// moveFromCopZero move from Coprocessor 0
func (cpu *CPU) moveFromCopZero(instr Instruction)  {
	cpuReg := instr.targetReg()
	copReg := instr.destReg()

	val := cpu.GetCopZeroReg(copReg)

	cpu.loadReg = LoadRegPair{cpuReg, val}
}
