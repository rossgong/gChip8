package chip8

type (
	Register uint8
	Address  uint16
)

const (
	registerCount  = 16
	statusRegister = 0xF //The F register is used for any status flags

	maxSubroutineLevel = 16
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
}
