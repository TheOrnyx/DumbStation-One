package cpu

// The general purpose registers (basically every register except PC)
type Registers struct {
	//             Name | Alias    | Common Usage
	//------------------+----------+----------------------------------------------------
	zero uint32 // (R0) | zero     | Constant as always 0 (not real register apparently)
	at   uint32 // R1   | at       | Assembler temporary
	v0   uint32 // R2   | v0       | Subroutine return value
	v1   uint32 // R3   | v1       | Subroutine return value
	a0   uint32 // R4   | a0       | Subroutine argument (subroutine may change)
	a1   uint32 // R5   | a1       | Subroutine argument (subroutine may change)
	a2   uint32 // R6   | a2       | Subroutine argument (subroutine may change)
	a3   uint32 // R7   | a3       | Subroutine argument (subroutine may change)
	
	t0   uint32 // R8   | t0       | Temporary (subroutine may change)
	t1   uint32 // R9   | t1       | Temporary (subroutine may change)
	t2   uint32 // R10  | t2       | Temporary (subroutine may change)
	t3   uint32 // R11  | t3       | Temporary (subroutine may change)
	t4   uint32 // R12  | t4       | Temporary (subroutine may change)
	t5   uint32 // R13  | t5       | Temporary (subroutine may change)
	t6   uint32 // R14  | t6       | Temporary (subroutine may change)
	t7   uint32 // R15  | t7       | Temporary (subroutine may change)
	s0   uint32 // R16  | s0       | Static vars (must be saved by subs)
	s1   uint32 // R17  | s1       | Static vars (must be saved by subs)
	s2   uint32 // R18  | s2       | Static vars (must be saved by subs)
	s3   uint32 // R19  | s3       | Static vars (must be saved by subs)
	s4   uint32 // R20  | s4       | Static vars (must be saved by subs)
	s5   uint32 // R21  | s5       | Static vars (must be saved by subs)
	s6   uint32 // R22  | s6       | Static vars (must be saved by subs)
	s7   uint32 // R23  | s7       | Static vars (must be saved by subs)
	t8   uint32 // R24  | t8       | Temporary (subroutine may change)
	t9   uint32 // R25  | t9       | Temporary (subroutine may change)
	k0   uint32 // R26  | k0       | Reserved for Kernel (destroyed by some IRQ handlers)
	k1   uint32 // R27  | k1       | Reserved for Kernel (destroyed by some IRQ handlers)
	gp   uint32 // R28  | gp       | Global pointer (rearely used)
	sp   uint32 // R29  | sp       | Stack Pointer
	fp   uint32 // R30  | fp(s8)   | Frame pointer, or 9th static var, must be saved
	ra   uint32 // R31  | ra       | Return address (used so by JAL,BLTZAL,BGEZAL opcodes)
	hi   uint32 // -    | hi       | Multiply/divide result (subroutine may change)
	lo   uint32 // -    | lo       | Multiply/divide result (subroutine may change)
}
