package chip8

import (
	"fmt"
	"math/rand"
)

type (
	Address     uint16
	Instruction uint16
	Operation   func()
)

const (
	registerCount      = 16
	maxSubroutineLevel = 16
	instructionSize    = 2   //bytes
	statusRegister     = 0xF //The F register is used for any status flags
)

type CPU struct {
	//Accessible registers
	Registers     [registerCount]byte
	DelayRegister byte
	SoundRegister byte
	RegisterI     Address //Register used for addresses

	//Internal data
	programCounter Address
	stackPointer   byte
	stack          [maxSubroutineLevel]Address
	ram            *Memory

	execute Operation

	random rand.Rand
}

func (cpu *CPU) Cycle() error {
	//fetch
	opcode, err := cpu.fetch()
	if err == nil {
		return err
	}

	//decode
	cpu.execute, err = cpu.decode(opcode)
	if err == nil {
		return err
	}
	//execute
	cpu.execute()
	return nil
}

func (cpu *CPU) fetch() (Instruction, error) {
	address := cpu.programCounter
	cpu.programCounter += 2
	_ = address
	return 0, fmt.Errorf("fetch error not implemented")
}

func (cpu *CPU) decode(opcode Instruction) (Operation, error) {
	startingNibble := opcode & 0xF000 //Mask to solo the first nibble

	switch startingNibble {
	case 0x0000:
		lastNibble := byte(opcode & 0x00FF) //Mask to solo the last nibble
		return decode0(cpu, lastNibble)
	case 0x1000: //JP instruction
		return jump(&cpu.programCounter, maskAddress(opcode)), nil
	case 0x2000: //CALL instruction
		return subroutineCall(cpu, maskAddress(opcode))
	case 0x3000: //SE Skip if register equals byte
		return skipInstructionIfTrue(&cpu.programCounter,
			cpu.Registers[maskXRegister(opcode)] == maskEndingByte(opcode)), nil
	case 0x4000: //SNE Skip if register does not equal byte
		return skipInstructionIfTrue(&cpu.programCounter,
			cpu.Registers[maskXRegister(opcode)] != maskEndingByte(opcode)), nil
	case 0x5000: //SE Skip if register equals register
		return skipInstructionIfTrue(&cpu.programCounter,
			cpu.Registers[maskXRegister(opcode)] == cpu.Registers[maskYRegister(opcode)]), nil
	case 0x6000: //LD Load byte into register
		return loadRegister(&cpu.Registers[maskXRegister(opcode)], maskEndingByte(opcode)), nil
	case 0x7000: //ADD Adds byte into register
		xRegister := maskXRegister(opcode)
		return add(&cpu.Registers[statusRegister], &cpu.Registers[xRegister], cpu.Registers[xRegister]), nil
	case 0x8000: //Various math functions requires both registers
		xRegister := &cpu.Registers[maskXRegister(opcode)]
		yRegister := cpu.Registers[maskYRegister(opcode)]
		lastNibble := opcode & 0x000F //Mask to solo the last nibble
		op := decode8(&cpu.Registers[statusRegister], xRegister, yRegister, byte(lastNibble))

		//This is needed as not all 0x8xxx opcodes are valid
		if op != nil {
			return op, nil
		}
	case 0x9000: //SNE Skip if register X is not equal to register Y. Opcode must end in 0?
		if opcode&0x000F == 0 {
			return skipInstructionIfTrue(&cpu.programCounter,
				cpu.Registers[maskXRegister(opcode)] != cpu.Registers[maskYRegister(opcode)]), nil
		}
	case 0xA000: //LD Load address into I
		return loadAddress(&cpu.RegisterI, maskAddress(opcode)), nil
	case 0xB000: //JP Offset opcode address with register 0 and jump there
		return jumpOffset(&cpu.programCounter, cpu.Registers[0], maskAddress(opcode)), nil
	case 0xC000: //RND load a register x with a random byte AND a byte mask
		return randByteMasked(&cpu.random, &cpu.Registers[maskXRegister(opcode)], maskEndingByte(opcode)), nil
	case 0xD000: //DRW
		break
	case 0xE000: //Keyboard functions
		break
	case 0xF000:
		xRegister := &cpu.Registers[maskXRegister(opcode)]
		lastByte := byte(opcode & 0x00FF) //Mask to solo the last byte
		op := decodeF(cpu, xRegister, lastByte)

		//This is needed as not all 0xFxxx opcodes are valid
		if op != nil {
			return op, nil
		}
	}
	return nil, fmt.Errorf("decode error: 0x%X not implemented/supported", opcode)
}

func decode0(cpu *CPU, lastByte byte) (Operation, error) {
	switch lastByte {
	case 0xE0: //CLS clear display
		return clearDisplay()
	case 0xEE: //RET Return from subroutine
		return subroutineReturn(cpu)
	}
	return nil, fmt.Errorf("decode error: 0x00%X not implemented/supported", lastByte)
}

//Function to make decode 0x8xxx not cloud up the decode function
func decode8(statusRegister *byte, xRegister *byte, yValue byte, lastByte byte) Operation {
	switch lastByte {
	case 0x0000: //LD Load register Y into register X
		return loadRegister(xRegister, yValue)
	case 0x0001: //OR Store registerX OR registerY into register X
		return or(xRegister, yValue)
	case 0x0002: //AND Store registerX AND registerY into register X
		return and(xRegister, yValue)
	case 0x0003: //XOR Store registerX XOR registerY into register X
		return xor(xRegister, yValue)
	case 0x0004: //ADD Store registerX + registerY into register X
		return add(statusRegister, xRegister, yValue)
	case 0x0005: //SUB Store registerX - registerY into register X
		return subtract(statusRegister, xRegister, yValue)
	case 0x0006: //SHR Store registerX >> 1 into register X
		return shiftRight(statusRegister, xRegister)
	case 0x0007: //SUBN Store registerY - registerX into register X
		return subtractN(statusRegister, xRegister, yValue)
	case 0x000E: //SHL Store registerX << 1 into register X
		return shiftLeft(statusRegister, xRegister)
	}
	return nil
}

func decodeF(cpu *CPU, xRegister *byte, lastByte byte) Operation {
	switch lastByte {
	case 0x07: //LD Load delay into register X
		return loadRegister(xRegister, cpu.DelayRegister)
	case 0x0A: //LD Load Keypress
		return loadKeyPress(cpu)
	case 0x15: //LD Load register X into delay
		return loadRegister(&cpu.DelayRegister, *xRegister)
	case 0x18: //LD Load register X into sound
		return loadRegister(&cpu.SoundRegister, *xRegister)
	case 0x1E: //ADD Add register X into I
		return addI(&cpu.RegisterI, *xRegister)
	case 0x29: //LD Load location of digit sprite into I
		return loadDigit(&cpu.RegisterI, *xRegister)
	case 0x33: //LD Store BCD representations of register x into I, I+1, I+2
		return storeBCD(cpu.RegisterI, *xRegister, cpu.ram)
	case 0x55: //LD Store registers starting at memory location I
		return storeRegisters(&cpu.Registers, cpu.RegisterI, cpu.ram)
	case 0x65: //LD Load registers from memory locations starting at location I
		return loadRegisters(&cpu.Registers, cpu.RegisterI, cpu.ram)
	}
	return nil
}

/*
MASKING FUNCTIONS
Faster to just cover the individual masking scenarios then create a generic function
*/

//Gets and address from an opcode with the format Xnnn with n being the address
func maskAddress(opcode Instruction) Address {
	return Address(opcode & 0x0FFF)
}

func maskXRegister(opcode Instruction) byte {
	return byte((opcode & 0x0F00) >> 8) //Mask corrrect nibble and then shift and convert
}

func maskYRegister(opcode Instruction) byte {
	return byte((opcode & 0x00F0) >> 4) //Mask corrrect nibble and then shift and convert
}

func maskEndingByte(opcode Instruction) byte {
	return byte(opcode & 0x00FF)
}

//Error functions
// func checkOutOfBounds(fnString string, indicies ...uint8) error {
// 	hasError := false
// 	fnString += " ("
// 	for i, reg := range indicies {
// 		if reg >= registerCount {
// 			if hasError { //If there is already an invlid register
// 				fnString += fmt.Sprintf("&&|V%v", rune(byte('x')+byte(i)))
// 			} else {
// 				fnString += fmt.Sprintf("|V%s", rune(byte('x')+byte(i)))
// 			}
// 			hasError = true
// 		}
// 	}
// 	if hasError {
// 		return fmt.Errorf("%s): invalid register number", fnString)
// 	} else {
// 		return nil
// 	}
// }
