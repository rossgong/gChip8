package chip8

import (
	"fmt"
	"math/rand"
)

//Defines the alu functions of the chip8 cpu
//These for the most part map to opCode functions

//CLS opcode
/*
TODO: Implement GPU
func clearDisplay(display *GPU) Operation {
	return GPU.clear()
}
*/
func clearDisplay() (Operation, error) {
	return nil, fmt.Errorf("cls error: display not implemented")
}

//Return from subroutine
func subroutineReturn(cpu *CPU) (Operation, error) {
	if cpu.stackPointer > 0 {
		return func() {
			cpu.programCounter = cpu.stack[cpu.stackPointer]
			cpu.stackPointer--
		}, nil
	} else {
		return nil, fmt.Errorf("ret error: stack empty, nothing to return too")
	}
}

func jump(programCounter *Address, address Address) Operation {
	return func() {
		*programCounter = address
	}
}

func jumpOffset(programCounter *Address, offset byte, address Address) Operation {
	return jump(programCounter, address+Address(offset))
}

func subroutineCall(cpu *CPU, subroutine Address) (Operation, error) {
	if cpu.stackPointer < maxSubroutineLevel {
		return func() {
			cpu.stack[cpu.stackPointer] = cpu.programCounter
			cpu.programCounter = subroutine
			cpu.stackPointer++
		}, nil
	} else {
		return nil, fmt.Errorf("call error: stack overflow")
	}
}

func skipInstructionIfTrue(programCounter *Address, condition bool) Operation {
	return func() {
		if condition {
			*programCounter += instructionSize
		}
	}
}

//Load value into register number [0x0-0xF]
func loadRegister(vX *byte, value byte) Operation {
	return func() {
		*vX = value
	}
}

func or(vX *byte, vY byte) Operation {
	return func() {
		*vX |= vY
	}
}

func and(vX *byte, vY byte) Operation {
	return func() {
		*vX &= vY
	}
}

func xor(vX *byte, vY byte) Operation {
	return func() {
		*vX ^= vY
	}
}

func add(status *byte, vX *byte, vY byte) Operation {
	return func() {
		temp := *vX
		*vX += vY
		if temp < *vX { //overflow
			*status = 1
		} else {
			*status = 0
		}

	}
}

func subtract(status *byte, vX *byte, vY byte) Operation {
	return func() {
		temp := *vX
		*vX -= vY
		if temp > vY { //NOT borrow
			*status = 1
		} else {
			*status = 0
		}
	}
}

func shiftRight(status *byte, vX *byte) Operation {
	return func() {
		*status = *vX & 1 //Shift right bit unto status
		*vX >>= 1
	}
}

func subtractN(status *byte, vX *byte, vY byte) Operation {
	return func() {
		temp := *vX
		*vX = vY - *vX
		if vY > temp { //NOT borrow
			*status = 1
		} else {
			*status = 0
		}
	}
}

func shiftLeft(status *byte, vX *byte) Operation {
	return func() {
		*status = *vX >> 7 //Shift Left bit unto status
		*vX <<= 1
	}
}

//Load value into I
func loadAddress(registerI *Address, value Address) Operation {
	return func() {
		*registerI = value
	}
}

func randByteMasked(rand *rand.Rand, vX *byte, mask byte) Operation {
	return func() {
		rand := byte(rand.Uint64() >> 56) //Shift random uin64 56 places in order to have only 8 bits of random
		*vX = rand & mask
	}
}

/*
TODO: Implement GPU
*/
func draw(cpu *CPU, vX *byte, vY byte, nibble uint8) Operation {
	return nil
}

func loadKeyPress(cpu *CPU) Operation {
	return nil
}

func addI(registerI *Address, vX byte) Operation {
	return func() {
		*registerI += Address(vX)
	}
}

func loadDigit(registerI *Address, vX byte) Operation {
	return func() {
		*registerI = (Address(vX) * 5) + digitSpriteLocation
	}
}

func storeBCD(registerI Address, vX byte, memory *Memory) Operation {
	return func() {
		//i < 3 is so it includes the 0 and detects when the value overflows janky but works
		for i := Address(2); i < 3; i-- {
			//Get the ones place and then "shift" one place over
			memory[registerI+i] = vX % 10
			vX /= 10
		}
	}
}

func storeRegisters(registers *[registerCount]byte, registerI Address, memory *Memory) Operation {
	return func() {
		for i, value := range registers {
			memory[registerI+Address(i)] = value
		}
	}
}

func loadRegisters(registers *[registerCount]byte, registerI Address, memory *Memory) Operation {
	return func() {
		for i, _ := range registers {
			registers[i] = memory[registerI+Address(i)]
		}
	}
}
