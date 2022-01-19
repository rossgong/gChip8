package chip8

import (
	"fmt"
	"math/rand"
)

//Defines the alu functions of the chip8 cpu
//These for the most part map to opCode functions

//CLS opcode
func clearDisplay(display *Display) Operation {
	return func() {
		*display = Display{}
	}
}

//Return from subroutine
func subroutineReturn(cpu *cpu) (Operation, error) {
	if cpu.stackPointer > 0 {
		return func() {
			cpu.stackPointer--
			cpu.programCounter = cpu.stack[cpu.stackPointer]
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

func subroutineCall(cpu *cpu, subroutine Address) (Operation, error) {
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
		if temp > *vX { //overflow
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

func draw(display *Display, sprite []byte, vX byte, vY byte, status *byte) Operation {
	return func() {
		if display.drawSprite(sprite, vX, vY) {
			*status = 1
		} else {
			*status = 0
		}
	}
}

func loadKeyPress(cpu *cpu, vX *byte, keys *Input) Operation {
	cpu.isWaitingForInput = true
	initialKeys := *keys
	return func() {
		if initialKeys >= *keys { //We are waiting for a release which would be when a bit is unset
			initialKeys = *keys //Reset keys to check against as additional keys could be pressed
		} else {
			keysReleased := ^(^initialKeys | *keys) //Bitmagic or NOTImplication
			for i := byte(0); i < numKeys; i++ {
				if keysReleased.checkKey(0) { //Check first key
					*vX = i
					cpu.isWaitingForInput = false
					break
				} else { //If not shift and then check again
					keysReleased >>= 1
				}
			}
		}
	}
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

func storeBCD(registerI Address, vX byte, memory *memory) Operation {
	return func() {
		memory[registerI] = vX / 100
		vX /= 10
		memory[registerI+1] = vX / 10
		memory[registerI+2] = vX % 10
	}
}

func storeRegisters(registers *[registerCount]byte, registerI Address, vX byte, memory *memory) Operation {
	return func() {
		for offset := Address(0); offset <= Address(vX); offset++ {
			memory[registerI+offset] = registers[offset]
		}
	}
}

func loadRegisters(registers *[registerCount]byte, registerI Address, vX byte, memory *memory) Operation {
	return func() {
		for offset := Address(0); offset <= Address(vX); offset++ {
			registers[offset] = memory[registerI+offset]
		}
	}
}
