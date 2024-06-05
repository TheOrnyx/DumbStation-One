package cpu

type Instruction struct {
	
}

// DecodeInstruction decode an instruction from an instruction word
// and return a struct instance based on this
func DecodeInstruction(instr uint32) *Instruction {
	return new(Instruction)
}
