package cpu

import "github.com/TheOrnyx/psx-go/log"

type RegIndex uint32 // register index type

// The general purpose registers (and hi and lo) (basically every register except PC)
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

// GetReg get register from given value
// TODO - this is really gross, find a nicer way
// FIXME - I can replcae this with reflection if needed!!
func (reg *Registers) GetReg(index RegIndex) uint32 {
	switch index {
	case  0: return reg.zero
	case  1: return reg.at
	case  2: return reg.v0
	case  3: return reg.v1
	case  4: return reg.a0
	case  5: return reg.a1
	case  6: return reg.a2
	case  7: return reg.a3
	case  8: return reg.t0
	case  9: return reg.t1
	case 10: return reg.t2
	case 11: return reg.t3
	case 12: return reg.t4
	case 13: return reg.t5
	case 14: return reg.t6
	case 15: return reg.t7
	case 16: return reg.s0
	case 17: return reg.s1
	case 18: return reg.s2
	case 19: return reg.s3
	case 20: return reg.s4
	case 21: return reg.s5
	case 22: return reg.s6
	case 23: return reg.s7
	case 24: return reg.t8
	case 25: return reg.t9
	case 26: return reg.k0
	case 27: return reg.k1
	case 28: return reg.gp
	case 29: return reg.sp
	case 30: return reg.fp
	case 31: return reg.ra

	default:
		log.Warnf("Unknown register index %v\n", index)
		return reg.zero
	}
}

// SetReg set register at index to given value val
// TODO - fix this one too, grossss
func (reg *Registers) SetReg(index RegIndex, val uint32)  {
	switch index {
	case  0: // Do nothing here :P
	case  1: reg.at = val
	case  2: reg.v0 = val
	case  3: reg.v1 = val
	case  4: reg.a0 = val
	case  5: reg.a1 = val
	case  6: reg.a2 = val
	case  7: reg.a3 = val
	case  8: reg.t0 = val
	case  9: reg.t1 = val
	case 10: reg.t2 = val
	case 11: reg.t3 = val
	case 12: reg.t4 = val
	case 13: reg.t5 = val
	case 14: reg.t6 = val
	case 15: reg.t7 = val
	case 16: reg.s0 = val
	case 17: reg.s1 = val
	case 18: reg.s2 = val
	case 19: reg.s3 = val
	case 20: reg.s4 = val
	case 21: reg.s5 = val
	case 22: reg.s6 = val
	case 23: reg.s7 = val
	case 24: reg.t8 = val
	case 25: reg.t9 = val
	case 26: reg.k0 = val
	case 27: reg.k1 = val
	case 28: reg.gp = val
	case 29: reg.sp = val
	case 30: reg.fp = val
	case 31: reg.ra = val
	default:
		log.Warnf("Unknown register index %v\n", index)
	}

	reg.zero = 0
}
