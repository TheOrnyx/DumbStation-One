package cpu

import "github.com/TheOrnyx/psx-go/log"

// struct for holding individual opcodes along with their respective run methods
type OpCode struct {
	fieldVal uint32
	name     string
	runFunc  func(cpu *CPU, instr Instruction)
}

var UnknownOpcode OpCode = OpCode{fieldVal: 0xFFF, name: "Unknown Opcode",
	runFunc: func(cpu *CPU, instr Instruction) {
		log.Panicf("Unknown Opcode - 0x%08x, 0x%02x", instr, instr.function())
	}}

var UnknownSubOpcode OpCode = OpCode{fieldVal: 0xFFF, name: "Unknown Sub Opcode",
	runFunc: func(cpu *CPU, instr Instruction) {
		log.Panicf("Unknown SubOpcode - 0x%08x, 0x%02x", instr, instr.subFunction())
	}}

var UnknownCopZeroOpcode OpCode = OpCode{fieldVal: 0xFFF, name: "Unknown Cop Zero Opcode",
	runFunc: func(cpu *CPU, instr Instruction) {
		log.Panicf("Unknown SubOpcode - 0x%08x, 0x%02x", instr, instr.copOpcode())
	}}

var Opcodes [0x40]OpCode = [0x40]OpCode{
	{0x00, "SPECIAL", func(cpu *CPU, instr Instruction) { cpu.executeSubInstr(instr) }},
	{0x01, "BcondZ", func(cpu *CPU, instr Instruction) { cpu.branchVarious(instr) }},
	{0x02, "J", func(cpu *CPU, instr Instruction) { cpu.jump(instr) }},
	{0x03, "JAL", func(cpu *CPU, instr Instruction) { cpu.jumpAndLink(instr) }},
	{0x04, "BEQ", func(cpu *CPU, instr Instruction) { cpu.branchIfEqual(instr) }},
	{0x05, "BNE", func(cpu *CPU, instr Instruction) { cpu.branchNotEqual(instr) }},
	{0x06, "BLEZ", func(cpu *CPU, instr Instruction) { cpu.branchLessOrEqualZero(instr) }},
	{0x07, "BGTZ", func(cpu *CPU, instr Instruction) { cpu.branchGreaterThanZero(instr) }},
	{0x08, "ADDI", func(cpu *CPU, instr Instruction) { cpu.addImmediate(instr) }},
	{0x09, "ADDIU", func(cpu *CPU, instr Instruction) { cpu.addImmediateUnsigned(instr) }},
	{0x0a, "SLTI", func(cpu *CPU, instr Instruction) { cpu.setIfLessThanImm(instr) }},
	{0x0b, "SLTIU", func(cpu *CPU, instr Instruction) { cpu.setIfLessThanImmUnsigned(instr) }},
	{0x0c, "ANDI", func(cpu *CPU, instr Instruction) { cpu.andImmediate(instr) }},
	{0x0d, "ORI", func(cpu *CPU, instr Instruction) { cpu.orImmediate(instr) }},
	{0x0e, "XORI", func(cpu *CPU, instr Instruction) { cpu.xorImmediate(instr) }},
	{0x0f, "LUI", func(cpu *CPU, instr Instruction) { cpu.loadUpperImmediate(instr) }},
	{0x10, "COP0", func(cpu *CPU, instr Instruction) { cpu.copZeroOpcode(instr) }},
	{0x11, "COP1", func(cpu *CPU, instr Instruction) { cpu.copOne(instr) }},
	{0x12, "COP2", func(cpu *CPU, instr Instruction) { cpu.copTwoOpcode(instr) }},
	{0x13, "COP3", func(cpu *CPU, instr Instruction) { cpu.copThree(instr) }},
	UnknownOpcode, // 0x14
	UnknownOpcode, // 0x15
	UnknownOpcode, // 0x16
	UnknownOpcode, // 0x17
	UnknownOpcode, // 0x18
	UnknownOpcode, // 0x19
	UnknownOpcode, // 0x1a
	UnknownOpcode, // 0x1b
	UnknownOpcode, // 0x1c
	UnknownOpcode, // 0x1d
	UnknownOpcode, // 0x1e
	UnknownOpcode, // 0x1f
	{0x20, "LB", func(cpu *CPU, instr Instruction) { cpu.loadByte(instr) }},
	{0x21, "LH", func(cpu *CPU, instr Instruction) { cpu.loadHalfWord(instr) }},
	{0x22, "LWL", func(cpu *CPU, instr Instruction) { cpu.loadWordLeft(instr) }},
	{0x23, "LW", func(cpu *CPU, instr Instruction) { cpu.loadWord(instr) }},
	{0x24, "LBU", func(cpu *CPU, instr Instruction) { cpu.loadByteUnsigned(instr) }},
	{0x25, "LHU", func(cpu *CPU, instr Instruction) { cpu.loadHalfWordUnsigned(instr) }},
	{0x26, "LWR", func(cpu *CPU, instr Instruction) { cpu.loadWordRight(instr) }},
	UnknownOpcode, // 0x27
	{0x28, "SB", func(cpu *CPU, instr Instruction) { cpu.storeByte(instr) }},
	{0x29, "SH", func(cpu *CPU, instr Instruction) { cpu.storeHalfWord(instr) }},
	{0x2a, "SWL", func(cpu *CPU, instr Instruction) { cpu.storeWordLeft(instr) }},
	{0x2b, "SW", func(cpu *CPU, instr Instruction) { cpu.storeWord(instr) }},
	UnknownOpcode, // 0x2c
	UnknownOpcode, // 0x2d
	{0x2e, "SWR", func(cpu *CPU, instr Instruction) { cpu.storeWordRight(instr) }},
	UnknownOpcode, // 0x2f
	{0x30, "LWC0", func(cpu *CPU, instr Instruction) { cpu.loadWordInCopZero(instr) }},
	{0x31, "LWC1", func(cpu *CPU, instr Instruction) { cpu.loadWordInCopOne(instr) }},
	{0x32, "LWC2", func(cpu *CPU, instr Instruction) { cpu.loadWordInCopTwo(instr) }},
	{0x33, "LWC3", func(cpu *CPU, instr Instruction) { cpu.loadWordInCopThree(instr) }},
	UnknownOpcode, // 0x34
	UnknownOpcode, // 0x35
	UnknownOpcode, // 0x36
	UnknownOpcode, // 0x37
	{0x38, "SWC0", func(cpu *CPU, instr Instruction) { cpu.storeWordInCopZero(instr) }},
	{0x39, "SWC1", func(cpu *CPU, instr Instruction) { cpu.storeWordInCopOne(instr) }},
	{0x3a, "SWC2", func(cpu *CPU, instr Instruction) { cpu.storeWordInCopTwo(instr) }},
	{0x3b, "SWC3", func(cpu *CPU, instr Instruction) { cpu.storeWordInCopThree(instr) }},
	UnknownOpcode, // 0x3c
	UnknownOpcode, // 0x3d
	UnknownOpcode, // 0x3e
	UnknownOpcode, // 0x3f
}

var SubOpcodes [0x40]OpCode = [0x40]OpCode{
	{0x00, "SLL", func(cpu *CPU, instr Instruction) { cpu.shiftLeftLogical(instr) }},
	UnknownSubOpcode,
	{0x02, "SRL", func(cpu *CPU, instr Instruction) { cpu.shiftRightLogical(instr) }},
	{0x03, "SRA", func(cpu *CPU, instr Instruction) { cpu.shiftRightArithmetic(instr) }},
	{0x04, "SLLV", func(cpu *CPU, instr Instruction) { cpu.shiftLeftLogicalVar(instr) }},
	UnknownSubOpcode,
	{0x06, "SRLV", func(cpu *CPU, instr Instruction) { cpu.shiftRightLogicalVar(instr) }},
	{0x07, "SRAV", func(cpu *CPU, instr Instruction) { cpu.shiftRightArithmeticVar(instr) }},
	{0x08, "JR", func(cpu *CPU, instr Instruction) { cpu.jumpRegister(instr) }},
	{0x09, "JALR", func(cpu *CPU, instr Instruction) { cpu.jumpAndLinkReg(instr) }},
	UnknownSubOpcode, // 0x0a
	UnknownSubOpcode, // 0x0b
	{0x0c, "SYSCALL", func(cpu *CPU, instr Instruction) { cpu.syscall(instr) }},
	{0x0d, "BREAK", func(cpu *CPU, instr Instruction) { cpu.breakEx(instr) }},
	UnknownSubOpcode, // 0x0e
	UnknownSubOpcode, // 0x0f
	{0x10, "MFHI", func(cpu *CPU, instr Instruction) { cpu.moveFromHI(instr) }},
	{0x11, "MTHI", func(cpu *CPU, instr Instruction) { cpu.moveToHI(instr) }},
	{0x12, "MFLO", func(cpu *CPU, instr Instruction) { cpu.moveFromLO(instr) }},
	{0x13, "MTLO", func(cpu *CPU, instr Instruction) { cpu.moveToLO(instr) }},
	UnknownSubOpcode, // 0x14
	UnknownSubOpcode, // 0x15
	UnknownSubOpcode, // 0x16
	UnknownSubOpcode, // 0x17
	{0x18, "MULT", func(cpu *CPU, instr Instruction) { cpu.multiply(instr) }},
	{0x19, "MULTU", func(cpu *CPU, instr Instruction) { cpu.multiplyUnsigned(instr) }},
	{0x1a, "DIV", func(cpu *CPU, instr Instruction) { cpu.div(instr) }},
	{0x1b, "DIVU", func(cpu *CPU, instr Instruction) { cpu.divUnsigned(instr) }},
	UnknownSubOpcode, // 0x1c
	UnknownSubOpcode, // 0x1d
	UnknownSubOpcode, // 0x1e
	UnknownSubOpcode, // 0x1f
	{0x20, "ADD", func(cpu *CPU, instr Instruction) { cpu.add(instr) }},
	{0x21, "ADDU", func(cpu *CPU, instr Instruction) { cpu.addUnsigned(instr) }},
	{0x22, "SUB", func(cpu *CPU, instr Instruction) { cpu.sub(instr) }},
	{0x23, "SUBU", func(cpu *CPU, instr Instruction) { cpu.subUnsigned(instr) }},
	{0x24, "AND", func(cpu *CPU, instr Instruction) { cpu.and(instr) }},
	{0x25, "OR", func(cpu *CPU, instr Instruction) { cpu.or(instr) }},
	{0x26, "XOR", func(cpu *CPU, instr Instruction) { cpu.xor(instr) }},
	{0x27, "NOR", func(cpu *CPU, instr Instruction) { cpu.nor(instr) }},
	UnknownSubOpcode, // 0x28
	UnknownSubOpcode, // 0x29
	{0x2a, "SLT", func(cpu *CPU, instr Instruction) { cpu.setIfLessThan(instr) }},
	{0x2b, "SLTU", func(cpu *CPU, instr Instruction) { cpu.setIfLessThanUnsigned(instr) }},
	UnknownSubOpcode, // 0x2c
	UnknownSubOpcode, // 0x2d
	UnknownSubOpcode, // 0x2e
	UnknownSubOpcode, // 0x2f
	UnknownSubOpcode, // 0x30
	UnknownSubOpcode, // 0x31
	UnknownSubOpcode, // 0x32
	UnknownSubOpcode, // 0x33
	UnknownSubOpcode, // 0x34
	UnknownSubOpcode, // 0x35
	UnknownSubOpcode, // 0x36
	UnknownSubOpcode, // 0x37
	UnknownSubOpcode, // 0x38
	UnknownSubOpcode, // 0x39
	UnknownSubOpcode, // 0x3a
	UnknownSubOpcode, // 0x3b
	UnknownSubOpcode, // 0x3c
	UnknownSubOpcode, // 0x3d
	UnknownSubOpcode, // 0x3e
	UnknownSubOpcode, // 0x3f
}
