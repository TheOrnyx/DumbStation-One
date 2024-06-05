package cpu

type Instruction uint32

// function return the primary opcode field of the instruction (bits 26..31)
func (i Instruction) function() uint32 {
	return uint32(i) >> 24
}

// targetReg get the target register of the current instruction
func (i Instruction) targetReg() uint32 {
	return (uint32(i) >> 16) & 0x1f
}

// sourceReg get the source reg of the current instruction (bits 21..25)
func (i Instruction) sourceReg() uint32 {
	return (uint32(i) >> 21) & 0x1f
}

// immediate16 return the immediate 16 bit value for the instruction
func (i Instruction) immediate16() uint32 {
	return uint32(i) & 0xffff
}

/////////////////////////////////////
// The CPU instructions themselves //
/////////////////////////////////////

// lui Load upper immediate - lui
func (cpu *CPU) lui(instr Instruction)  {
	immediate := instr.immediate16()
	targetReg := instr.targetReg()

	cpu.regs.SetReg(targetReg, immediate << 16)
}

// ori bitwise or immediate
func (cpu *CPU) ori(instr Instruction)  {
	immediate := instr.immediate16()
	sourceReg := instr.sourceReg()
	val := cpu.regs.GetReg(sourceReg) | immediate
	
	cpu.regs.SetReg(instr.targetReg(), val)
}
