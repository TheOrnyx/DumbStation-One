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
func (cpu *CPU) loadUpperImmediate(instr Instruction) {
	immediate := instr.immediate16()
	targetReg := instr.targetReg()

	cpu.SetReg(targetReg, immediate<<16)
}

// or bitwise or between two registers
func (cpu *CPU) or(instr Instruction) {
	sourceReg := instr.sourceReg()
	targetReg := instr.targetReg()

	val := cpu.GetReg(sourceReg) | cpu.GetReg(targetReg)
	cpu.SetReg(instr.destReg(), val)
}

// orImmediate bitwise or immediate
func (cpu *CPU) orImmediate(instr Instruction) {
	immediate := instr.immediate16()
	sourceReg := instr.sourceReg()
	val := cpu.GetReg(sourceReg) | immediate

	cpu.SetReg(instr.targetReg(), val)
}

// and bitwise and
func (cpu *CPU) and(instr Instruction) {
	sourceReg := instr.sourceReg()
	targetReg := instr.targetReg()

	val := cpu.GetReg(sourceReg) & cpu.GetReg(targetReg)
	cpu.SetReg(instr.destReg(), val)
}

// andImmediate bitwise and immediate
func (cpu *CPU) andImmediate(instr Instruction) {
	immediate := instr.immediate16()
	sourceReg := instr.sourceReg()
	val := cpu.GetReg(sourceReg) & immediate

	cpu.SetReg(instr.targetReg(), val)
}

// nor bitwise not or
func (cpu *CPU) nor(instr Instruction)  {
	source := cpu.GetReg(instr.sourceReg())
	target := cpu.GetReg(instr.targetReg())
	val := ^(source | target)

	cpu.SetReg(instr.destReg(), val)
}

// xor bitwise xor between two registers
func (cpu *CPU) xor(instr Instruction) {
	sourceReg := instr.sourceReg()
	targetReg := instr.targetReg()

	val := cpu.GetReg(sourceReg) ^ cpu.GetReg(targetReg)
	cpu.SetReg(instr.destReg(), val)
}

// xorImmediate bitwise xor immediate
func (cpu *CPU) xorImmediate(instr Instruction) {
	immediate := instr.immediate16()
	sourceReg := instr.sourceReg()
	val := cpu.GetReg(sourceReg) ^ immediate

	cpu.SetReg(instr.targetReg(), val)
}

// storeWord Store Word
func (cpu *CPU) storeWord(instr Instruction) {
	if cpu.GetCopZeroReg(REG_SR)&0x10000 != 0 {
		log.Info("Ignoring storeWord while the cache is isolated")
		return
	}

	targetReg := instr.targetReg()
	sourceReg := instr.sourceReg()
	val := cpu.GetReg(targetReg)
	addr := cpu.GetReg(sourceReg) + instr.immediate16Se()

	if addr % 4 == 0 {
		cpu.Store32(addr, val)
	} else {
		cpu.Exception(StoreAddressError)
	}
}

// storeWordLeft store word left with unaligned
func (cpu *CPU) storeWordLeft(instr Instruction)  {
	immediate := instr.immediate16Se()
	targetReg := instr.targetReg()
	source := cpu.GetReg(instr.sourceReg())

	addr := source + immediate
	val := cpu.GetReg(targetReg)

	alignedAddr := addr & (^uint32(3))
	// load current value at target
	curMem := cpu.load32(alignedAddr)

	var mem uint32
	switch addr & 3 {
	case 0:
		mem = (curMem & 0xffffff00) | (val >> 24)
	case 1:
		mem = (curMem & 0xffff0000) | (val >> 16)
	case 2:
		mem = (curMem & 0xff000000) | (val >> 8)
	case 3:
		mem = (curMem & 0x00000000) | (val >> 0)
	default:
		log.Panicf("storeWordLeft failed - This shouldn't happen: 0x%08x", addr & 3)
	}

	cpu.Store32(alignedAddr, mem)
}

// storeWordRight store word right unaligned
func (cpu *CPU) storeWordRight(instr Instruction)  {
	immediate := instr.immediate16Se()
	targetReg := instr.targetReg()
	source := cpu.GetReg(instr.sourceReg())

	addr := source + immediate
	val := cpu.GetReg(targetReg)

	alignedAddr := addr & (^uint32(3))
	// load current value at target
	curMem := cpu.load32(alignedAddr)

	var mem uint32
	switch addr & 3 {
	case 0:
		mem = (curMem & 0x00000000) | (val << 0)
	case 1:
		mem = (curMem & 0x000000ff) | (val << 8)
	case 2:
		mem = (curMem & 0x0000ffff) | (val << 16)
	case 3:
		mem = (curMem & 0x00ffffff) | (val << 24)
	default:
		log.Panicf("storeWordLeft failed - This shouldn't happen: 0x%08x", addr & 3)
	}

	cpu.Store32(alignedAddr, mem)
}

// storeHalfWord store half word into memory
func (cpu *CPU) storeHalfWord(instr Instruction) {
	if cpu.GetCopZeroReg(REG_SR)&0x10000 != 0 {
		log.Info("Ignoring storeHalfWord while the cache is isolated")
		return
	}

	immediate := instr.immediate16Se()
	targetReg := cpu.GetReg(instr.targetReg())
	sourceReg := cpu.GetReg(instr.sourceReg())
	addr := sourceReg + immediate

	if addr % 2 == 0 {
		cpu.Store16(addr, uint16(targetReg))
	} else {
		cpu.Exception(StoreAddressError)
	}
}

// storeByte store byte
func (cpu *CPU) storeByte(instr Instruction) {
	if cpu.GetCopZeroReg(REG_SR)&0x10000 != 0 {
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
func (cpu *CPU) loadWord(instr Instruction) {
	if cpu.GetCopZeroReg(REG_SR)&0x10000 != 0 {
		log.Info("Ignoring loadWord while the cache is isolated")
		return
	}

	immediate := instr.immediate16Se()
	sourceReg := instr.sourceReg()
	addr := cpu.GetReg(sourceReg) + immediate

	if addr % 4 == 0 {
		val := cpu.load32(addr)
		cpu.SetLoadReg(instr.targetReg(), val)
	} else {
		cpu.Exception(LoadAddressError)
	}
}

// loadHalfWordUnsigned load half word unsigned from memory
func (cpu *CPU) loadHalfWordUnsigned(instr Instruction)  {
	immediate := instr.immediate16Se()
	source := cpu.GetReg(instr.sourceReg())
	addr := source + immediate

	if addr % 2 == 0 {
		val := cpu.Load16(addr)
		cpu.SetLoadReg(instr.targetReg(), uint32(val))
	} else {
		cpu.Exception(LoadAddressError)
	}
}

// loadHalfWord load halfword (signed)
func (cpu *CPU) loadHalfWord(instr Instruction)  {
	immediate := instr.immediate16Se()
	source := cpu.GetReg(instr.sourceReg())
	addr := source + immediate

	val := int16(cpu.Load16(addr))

	cpu.SetLoadReg(instr.targetReg(), uint32(val))
}

// loadWordLeft load word left - can load unaligned?
func (cpu *CPU) loadWordLeft(instr Instruction)  {
	immediate := instr.immediate16Se()
	source := cpu.GetReg(instr.sourceReg())

	addr := source + immediate

	curVal := cpu.outRegs.GetReg(instr.targetReg()) // TODO - check this is right?

	alignedAddr := addr & (^uint32(3))
	alignedWord := cpu.load32(alignedAddr)

	// based on alignment we fetch 1,2,3 or 4 most significant bytes
	var val uint32
	switch addr & 3 {
	case 0:
		val = (curVal & 0x00ffffff) | (alignedWord << 24)
	case 1:
		val = (curVal & 0x0000ffff) | (alignedWord << 16)
	case 2:
		val = (curVal & 0x000000ff) | (alignedWord << 8)
	case 3:
		val = (curVal & 0x00000000) | (alignedWord << 0)
	default:
		log.Panicf("LoadWordLeft failed - This shouldn't be reached: 0x%08x", addr & 3)
	}

	cpu.SetLoadReg(instr.targetReg(), val)
}

// loadWordRight load word right - can be unaligned
func (cpu *CPU) loadWordRight(instr Instruction)  {
	immediate := instr.immediate16Se()
	targetReg := instr.targetReg()
	source := cpu.GetReg(instr.sourceReg())

	addr := source + immediate

	curVal := cpu.outRegs.GetReg(targetReg)

	alignedAddr := addr & (^uint32(3))
	alignedWord := cpu.load32(alignedAddr)

	// based on address alignment we fetch 1,2,3 or 4 least sginificatn bytes
	var val uint32
	switch addr & 3 {
	case 0:
		val = (curVal & 0x00000000) | (alignedWord >> 0)
	case 1:
		val = (curVal & 0xff000000) | (alignedWord >> 8)
	case 2:
		val = (curVal & 0xffff0000) | (alignedWord >> 16)
	case 3:
		val = (curVal & 0xffffff00) | (alignedWord >> 24)
	case 4:
		log.Panicf("Loadwordright failed - shouldn't be reached: 0x%08x", addr & 3)
	}

	cpu.SetLoadReg(targetReg, val)
}

// loadByte load signed byte
func (cpu *CPU) loadByte(instr Instruction) {
	immediate := instr.immediate16Se()
	targetReg := instr.targetReg()
	sourceReg := cpu.GetReg(instr.sourceReg())

	addr := sourceReg + immediate
	val := int8(cpu.Load8(addr))

	cpu.loadReg = LoadRegPair{target: targetReg, val: uint32(val)}
}

// loadByteUnsigned load bytes unsigned (lb without sign extension,
// high 2 bits of target are set to 0)
func (cpu *CPU) loadByteUnsigned(instr Instruction) {
	immediate := instr.immediate16Se()
	targetReg := instr.targetReg()
	sourceReg := cpu.GetReg(instr.sourceReg())

	addr := sourceReg + immediate
	val := cpu.Load8(addr)

	cpu.loadReg = LoadRegPair{target: targetReg, val: uint32(val)}
}

// shiftLeftLogical Shift left logical
func (cpu *CPU) shiftLeftLogical(instr Instruction) {
	immediate := instr.shiftImmediate()
	targetReg := instr.targetReg()

	val := cpu.GetReg(targetReg) << immediate
	cpu.SetReg(instr.destReg(), val)
}

// shiftLeftLogicalVar shift left logical variable, basically regular
// but shift amount is in regiter
func (cpu *CPU) shiftLeftLogicalVar(instr Instruction)  {
	target := cpu.GetReg(instr.targetReg())
	shift := cpu.GetReg(instr.sourceReg())

	// shift amount is truncated to 5 bits
	val := target << (shift & 0x1f)

	cpu.SetReg(instr.destReg(), val)
}

// shiftRightLogical shift right logical
func (cpu *CPU) shiftRightLogical(instr Instruction)  {
	shift := instr.shiftImmediate()
	target := cpu.GetReg(instr.targetReg())
	val := target >> shift

	cpu.SetReg(instr.destReg(), val)
}

// shiftRightLogicalVar shift right logical variable, basically regular
// but shift amount is in register
func (cpu *CPU) shiftRightLogicalVar(instr Instruction)  {
	target := cpu.GetReg(instr.targetReg())
	shift := cpu.GetReg(instr.sourceReg())

	// shift amount is truncated to 5 bits
	val := target >> (shift & 0x1f)

	cpu.SetReg(instr.destReg(), val)
}


// shiftRightArithmetic shift right arithmetic
func (cpu *CPU) shiftRightArithmetic(instr Instruction) {
	shift := instr.shiftImmediate()
	target := int32(cpu.GetReg(instr.targetReg()))

	val := target >> int32(shift) // TODO - check shift should be cast to int32

	cpu.SetReg(instr.destReg(), uint32(val))
}

// shiftRightArithmeticVar shift right arithmetic with variable
func (cpu *CPU) shiftRightArithmeticVar(instr Instruction)  {
	target := cpu.GetReg(instr.destReg())
	shift := cpu.GetReg(instr.sourceReg())

	// shift amount truncated to 5 bits again
	val := int32(target) >> (shift & 0x1f)

	cpu.SetReg(instr.destReg(), uint32(val))
}

// addImmediateUnsigned add immediate unsigned
//
// This one apparently doesn't generate an exception on overflow while ADDI does?
func (cpu *CPU) addImmediateUnsigned(instr Instruction) {
	immediate := instr.immediate16Se()
	sourceReg := instr.sourceReg()

	val := cpu.GetReg(sourceReg) + immediate
	cpu.SetReg(instr.targetReg(), val)
}

// addImmediate add immediate - generates exception when overflows
func (cpu *CPU) addImmediate(instr Instruction) {
	immediate := instr.immediate16Se()
	sourceReg := cpu.GetReg(instr.sourceReg())

	result, overflowed := utils.AddSigned16(sourceReg, immediate)
	if overflowed {
		cpu.Exception(Overflow)
	}

	cpu.SetReg(instr.targetReg(), result)
}

// add add and generate exception on overflow
func (cpu *CPU) add(instr Instruction) {
	sourceReg := cpu.GetReg(instr.sourceReg())
	targetReg := cpu.GetReg(instr.targetReg())

	val, overflowed := utils.AddSigned16(sourceReg, targetReg)
	if overflowed {
		cpu.Exception(Overflow)
	}

	cpu.SetReg(instr.destReg(), val)
}

// addUnsigned Add two unsigned registers together and store in target
func (cpu *CPU) addUnsigned(instr Instruction) {
	sourceReg := cpu.GetReg(instr.sourceReg())
	targetReg := cpu.GetReg(instr.targetReg())
	val := sourceReg + targetReg

	cpu.SetReg(instr.destReg(), val)
}

// sub sub (signed) and check for overflows
func (cpu *CPU) sub(instr Instruction)  {
	sourceReg := cpu.GetReg(instr.sourceReg())
	targetReg := cpu.GetReg(instr.targetReg())

	val, overflowed := utils.SubSigned16(sourceReg, targetReg) // TODO - check
	if overflowed {
		cpu.Exception(Overflow)
	}

	cpu.SetReg(instr.destReg(), val)
}

// subUnsigned subtract unsigned
func (cpu *CPU) subUnsigned(instr Instruction) {
	source := cpu.GetReg(instr.sourceReg())
	target := cpu.GetReg(instr.targetReg())

	cpu.SetReg(instr.destReg(), source-target)
}

// jumpAndLink jump and link by storing return address in $ra
func (cpu *CPU) jumpAndLink(instr Instruction) {
	returnAddr := cpu.nextPC

	cpu.SetReg(RegIndex(31), returnAddr)
	cpu.jump(instr)
}

// jumpAndLinkReg jump and link by storing return address in register
func (cpu *CPU) jumpAndLinkReg(instr Instruction) {
	jumpLoc := cpu.GetReg(instr.sourceReg())
	returnAddr := cpu.nextPC

	cpu.SetReg(instr.destReg(), returnAddr)

	cpu.nextPC = jumpLoc
	cpu.branching = true
	cpu.branchingSetup()
}

// jumpRegister jump to val in register
func (cpu *CPU) jumpRegister(instr Instruction) {
	sourceReg := cpu.GetReg(instr.sourceReg())
	cpu.nextPC = sourceReg
	cpu.branching = true
	cpu.branchingSetup()
}

// branchNotEqual branch if not equal
func (cpu *CPU) branchNotEqual(instr Instruction) {
	offset := instr.immediate16Se()
	sourceReg := instr.sourceReg()
	targetReg := instr.targetReg()

	if cpu.GetReg(sourceReg) != cpu.GetReg(targetReg) {
		cpu.branch(offset)
	}
}

// branchIfEqual branch if equal
func (cpu *CPU) branchIfEqual(instr Instruction) {
	offset := instr.immediate16Se()
	sourceReg := instr.sourceReg()
	targetReg := instr.targetReg()

	if cpu.GetReg(sourceReg) == cpu.GetReg(targetReg) {
		cpu.branch(offset)
	}
}

// branchGreaterThanZero branch if greater than zero (signed)
func (cpu *CPU) branchGreaterThanZero(instr Instruction) {
	immediate := instr.immediate16Se()
	source := int32(cpu.GetReg(instr.sourceReg()))

	if source > 0 {
		cpu.branch(immediate)
	}
}

// branchLessOrEqualZero branch if less than or equal to zero (signed)
func (cpu *CPU) branchLessOrEqualZero(instr Instruction) {
	immediate := instr.immediate16Se()
	source := int32(cpu.GetReg(instr.sourceReg()))

	if source <= 0 {
		cpu.branch(immediate)
	}
}

// branchVarious various branch instructions: BGEZ, BLTZ, BGEZAL and BLTZAL
//
// bits 16 and 20 used to figure out which to use
func (cpu *CPU) branchVarious(instr Instruction) {
	immediate := instr.immediate16Se()
	sourceReg := instr.sourceReg()
	var test uint32 = 0x00

	inst := uint32(instr)

	isBgez := (inst >> 16) & 1
	isLink := (inst>>17)&0xf == 8

	val := int32(cpu.GetReg(sourceReg))

	if val < 0 {
		test = 0x01
	}

	// if test is 'greater than or equal to zero' we need to negate
	// comparison using a xor
	test = test ^ isBgez

	if isLink {
		returnAddr := cpu.nextPC
		// store pc in reg $ra
		cpu.SetReg(RegIndex(32), returnAddr)
	}

	if test != 0 {
		cpu.branch(immediate)
	}
}

// setIfLessThanUnsigned set the dest to 1 when source is less than target
func (cpu *CPU) setIfLessThanUnsigned(instr Instruction) {
	sourceReg := cpu.GetReg(instr.sourceReg())
	targetReg := cpu.GetReg(instr.targetReg())
	var val uint32 = 0

	if sourceReg < targetReg {
		val = 1
	}

	cpu.SetReg(instr.destReg(), val)
}

// setIfLessThanImm set if less than immediate (signed)
func (cpu *CPU) setIfLessThanImm(instr Instruction) {
	immediate := int32(instr.immediate16Se())
	source := int32(cpu.GetReg(instr.sourceReg()))

	val := source < immediate

	cpu.SetReg(instr.targetReg(), utils.BoolToUint32(val))
}

// setIfLessThanImmUnsigned set if less than immediate unsigned
func (cpu *CPU) setIfLessThanImmUnsigned(instr Instruction)  {
	immediate := instr.immediate16Se()
	source := cpu.GetReg(instr.sourceReg())

	val := source < immediate

	cpu.SetReg(instr.targetReg(), utils.BoolToUint32(val))
}

// setIfLessThan set on less than (signed)
func (cpu *CPU) setIfLessThan(instr Instruction)  {
	source := int32(cpu.GetReg(instr.sourceReg()))
	target := int32(cpu.GetReg(instr.targetReg()))

	val := source < target

	cpu.SetReg(instr.destReg(), utils.BoolToUint32(val))
}

// div divide (signed)
func (cpu *CPU) div(instr Instruction) {
	sourceReg := instr.sourceReg()
	targetReg := instr.targetReg()

	numerator := int32(cpu.GetReg(sourceReg))
	denominator := int32(cpu.GetReg(targetReg))

	if denominator == 0 {
		// divide by 0
		cpu.hi = uint32(numerator)

		if numerator >= 0 {
			cpu.lo = 0xffffffff
		} else {
			cpu.lo = 1
		}
	} else if uint32(numerator) == 0x80000000 && denominator == -1 {
		// result not representable in 32-bit signed int
		cpu.hi = 0
		cpu.lo = 0x80000000
	} else {
		cpu.hi = uint32(numerator % denominator)
		cpu.lo = uint32(numerator / denominator)
	}
}

// divUnsigned divide unsigned
func (cpu *CPU) divUnsigned(instr Instruction)  {
	sourceReg := instr.sourceReg()
	targetReg := instr.targetReg()

	numerator := cpu.GetReg(sourceReg)
	denominator := cpu.GetReg(targetReg)

	if denominator == 0 {
		// divide by 0
		cpu.hi = numerator
		cpu.lo = 0xffffffff
	} else {
		cpu.hi = numerator % denominator
		cpu.lo = numerator / denominator
	}
}

// multiply multiply (signed)
func (cpu *CPU) multiply(instr Instruction)  {
	a := int64(int32(cpu.GetReg(instr.sourceReg())))
	b := int64(int32(cpu.GetReg(instr.targetReg())))

	val := uint64(a * b)

	cpu.hi = uint32(val >> 32)
	cpu.lo = uint32(val)
}

// multiplyUnsigned multiply unsigned
func (cpu *CPU) multiplyUnsigned(instr Instruction)  {
	a := uint64(cpu.GetReg(instr.sourceReg()))
	b := uint64(cpu.GetReg(instr.targetReg()))

	val := a*b

	cpu.hi = uint32((val >> 32))
	cpu.lo = uint32(val)
}

// moveFromLO move from LO into general purpsoe register
//
// TODO - this should like stall if division isn't completed but implement later
func (cpu *CPU) moveFromLO(instr Instruction)  {
	cpu.SetReg(instr.destReg(), cpu.lo)
}

// moveToLO move to LO reg
func (cpu *CPU) moveToLO(instr Instruction)  {
	cpu.lo = cpu.GetReg(instr.sourceReg())
}

// moveFromHI move from hi into general purpose register
//
// TODO - like MFLO this should also stall but do this later
func (cpu *CPU) moveFromHI(instr Instruction)  {
	cpu.SetReg(instr.destReg(), cpu.hi)
}

// moveToHI move to HI reg
func (cpu *CPU) moveToHI(instr Instruction)  {
	cpu.hi = cpu.GetReg(instr.sourceReg())
}

// syscall system call
func (cpu *CPU) syscall(instr Instruction)  {
	cpu.Exception(SysCall)
}

// breakEx break exception
func (cpu *CPU) breakEx(instr Instruction)  {
	cpu.Exception(Break)
}

/////////////////////////////////
// Coprocessor Instructions ✨ //
/////////////////////////////////

// moveToCopZero move register contents to coprocessor zero
func (cpu *CPU) moveToCopZero(instr Instruction) {
	cpuReg := instr.targetReg()
	copReg := instr.destReg()

	val := cpu.GetReg(cpuReg)
	cpu.SetCopZeroReg(copReg, val)
}

// moveFromCopZero move from Coprocessor 0
func (cpu *CPU) moveFromCopZero(instr Instruction) {
	cpuReg := instr.targetReg()
	copReg := instr.destReg()

	val := cpu.GetCopZeroReg(copReg)

	cpu.loadReg = LoadRegPair{cpuReg, val}
}

// returnFromException return from exception
func (cpu *CPU) returnFromException(instr Instruction)  {
	// there are other instructions with same encoding so make sure we don't use
	if instr & 0x3f != 0b010000 {
		log.Panicf("Invalid cop0 instruction: %v", instr)
	}

	// restore pre-exception mode by shifting interrupt enable/user
	// mode stack back to original position
	mode := cpu.copZeroRegs.sr & 0x3f
	cpu.copZeroRegs.sr &= ^uint32(0x3f)
	cpu.copZeroRegs.sr |= mode >> 2
}

// copOne copropcessor 1 opcode (doesn't exist so throws exception)
func (cpu *CPU) copOne(instr Instruction)  {
	cpu.Exception(CoprocessorError)
}

// copThree copropcessor 3 opcode (doesn't exist so throws exception)
func (cpu *CPU) copThree(instr Instruction)  {
	cpu.Exception(CoprocessorError)
}

// loadWordInCopZero load word in cop 0 (throws exception)
func (cpu *CPU) loadWordInCopZero(instr Instruction)  {
	cpu.Exception(CoprocessorError)
}

// loadWordInCopOne load word in cop 1 (throws exception)
func (cpu *CPU) loadWordInCopOne(instr Instruction)  {
	cpu.Exception(CoprocessorError)
}

// loadWordInCopTwo load word in cop 2
func (cpu *CPU) loadWordInCopTwo(instr Instruction)  {
	log.Panicf("(Not implemented yet) unhandled LWC2: 0x%08x", instr)
}

// loadWordInCopThree load word in cop 3 (throws exception)
func (cpu *CPU) loadWordInCopThree(instr Instruction)  {
	cpu.Exception(CoprocessorError)
}



// storeWordInCopZero store word in cop 0 (throws exception)
func (cpu *CPU) storeWordInCopZero(instr Instruction)  {
	cpu.Exception(CoprocessorError)
}

// storeWordInCopOne store word in cop 1 (throws exception)
func (cpu *CPU) storeWordInCopOne(instr Instruction)  {
	cpu.Exception(CoprocessorError)
}

// storeWordInCopTwo store word in cop 2
func (cpu *CPU) storeWordInCopTwo(instr Instruction)  {
	log.Panicf("(Not implemented yet) unhandled SWC2: 0x%08x", instr)
}

// storeWordInCopThree store word in cop 3 (throws exception)
func (cpu *CPU) storeWordInCopThree(instr Instruction)  {
	cpu.Exception(CoprocessorError)
}
