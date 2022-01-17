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
	registerCount       = 16
	maxSubroutineLevel  = 16
	instructionSize     = 2   //bytes
	statusRegister      = 0xF //The F register is used for any status flags
	digitSpriteLocation = 0x0 //Address where the digit sprites start
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

	execute Operation

	randomSource rand.Rand
}

func (cpu *CPU) Cycle() error {
	//fetch
	opcode, err := cpu.fetch()
	if err == nil {
		return err
	}

	//decode
	op, err := cpu.decode(opcode)
	if err == nil {
		return err
	}
	//execute
	op()
	return nil
}

func (cpu *CPU) fetch() (Instruction, error) {
	address := cpu.programCounter
	cpu.programCounter += 2
	_ = address
	return 0, fmt.Errorf("fetch error not implemented")
}

func (cpu *CPU) decode(opcode Instruction) (Operation, error) {
	nibble1 := opcode & 0xF000 //Mask to solo the first nibble

	switch nibble1 {
	case 0x0000:
		//TODO: SYS CALLS
	case 0x1000: //JP instruction
		return jump(&cpu.programCounter, maskAddress(opcode)), nil
	case 0x2000: //CALL instruction
		return subroutineCall(cpu, maskAddress(opcode))
	case 0x3000: //SE
		return skipInstructionIfTrue(&cpu.programCounter,
			cpu.Registers[maskXRegister(opcode)] == byte(maskEndingByte(opcode))), nil
	case 0x4000:
	case 0x5000:
	case 0x6000:
	case 0x7000:
	case 0x8000:
	case 0x9000:
	case 0xA000:
	case 0xB000:
	case 0xC000:
	case 0xD000:
	case 0xE000:
	case 0xF000:
	}
	return nil, fmt.Errorf("decode error: not implemented")
}

/*
MASKING FUNCTIONS
Faster to just cover the individual masking scenarios then create a generic function
*/

//Gets and address from an opcode with the format Xnnn with n being the address
func maskAddress(opcode Instruction) Address {
	return Address(opcode & 0x0FFF)
}

func maskXRegister(opcode Instruction) uint8 {
	return uint8((opcode & 0x0F00) >> 8) //Mask corrrect nibble and then shift and convert
}

func maskYRegister(opcode Instruction) uint8 {
	return uint8((opcode & 0x00F0) >> 4) //Mask corrrect nibble and then shift and convert
}

func maskEndingByte(opcode Instruction) uint8 {
	return uint8(opcode & 0x00FF)
}

//Error functions
func checkOutOfBounds(fnString string, indicies ...uint8) error {
	hasError := false
	fnString += " ("
	for i, reg := range indicies {
		if reg >= registerCount {
			if hasError { //If there is already an invlid register
				fnString += fmt.Sprintf("&&|V%v", rune(byte('x')+byte(i)))
			} else {
				fnString += fmt.Sprintf("|V%s", rune(byte('x')+byte(i)))
			}
			hasError = true
		}
	}
	if hasError {
		return fmt.Errorf("%s): invalid register number", fnString)
	} else {
		return nil
	}
}

func checkRegistersAndReturnFunction(fnString string, function Operation, registers ...uint8) (Operation, error) {
	err := checkOutOfBounds(fnString, registers...)
	if err == nil {
		return function, nil
	} else {
		return nil, err
	}
}
