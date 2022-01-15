package chip8

import "math/rand"

type (
	Register  uint8
	Address   uint16
	Operation func()
)

const (
	registerCount      = 16
	maxSubroutineLevel = 16
	instructionSize    = 2   //bytes
	statusRegister     = 0xF //The F register is used for any status flags
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
