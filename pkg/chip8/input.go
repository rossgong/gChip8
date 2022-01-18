package chip8

const (
	numKeys = 16
)

type Input uint16

func (input Input) checkKey(key byte) bool {
	return input&(0x1<<key) > 0 //Check key by moving digit to the correct bit and then masking
}

func (input *Input) PressKey(key byte) {
	*input |= (0x1 << key)
}

func (input *Input) ReleaseKey(key byte) {
	*input &= ^(0x1 << key) //NOTAND always turns off
}
