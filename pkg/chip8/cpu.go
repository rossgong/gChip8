package chip8

import (
	"fmt"
	"math/rand"
)

type (
	Register  uint8
	Address   uint16
	Operation func()
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
	Registers     [registerCount]Register
	DelayRegister Register
	SoundRegister Register
	RegisterI     Address //Register used for addresses

	//Internal data
	programCounter Address
	stackPointer   Register
	stack          [maxSubroutineLevel]Address

	execute Operation

	randomSource rand.Rand
}

//Error functions
func checkOutOfBounds(fnString string, indicies ...uint8) error {
	hasError := false
	fnString += " ("
	for i, reg := range indicies {
		if reg >= registerCount {
			if hasError { //If there is already an invlid register
				fnString += fmt.Sprintf("&&|V%s", rune(byte('x')+byte(i)))
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
