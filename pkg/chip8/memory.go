package chip8

import "fmt"

const (
	RamSize             = 0x1000
	programStart        = 0x200
	digitSpriteLocation = 0x0 //Address where the digit sprites start
)

type memory [RamSize]byte

func (ram *memory) loadProgam(program []byte) error {
	if len(program)-200 > RamSize {
		return fmt.Errorf("ram error: program is too large and cannot be loaded into memory")
	}
	for i, programByte := range program {
		ram[i+programStart] = programByte
	}
	return nil
}

func (memory *memory) getSprite(address Address, size byte) []byte {
	sprite := make([]byte, size)

	for i, _ := range sprite {
		sprite[i] = memory[address+Address(i)]
	}

	return sprite
}

func (ram *memory) loadFont() {
	for i, fontByte := range font {
		ram[i+digitSpriteLocation] = fontByte
	}
}
