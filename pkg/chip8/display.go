package chip8

const (
	defaultWidth     = 64
	defaultHeight    = 32
	defaultByteWidth = defaultWidth / 8
)

type Display [defaultHeight][defaultByteWidth]byte

//Returns collison
func (display *Display) drawSprite(sprite []byte, x byte, y byte) bool {
	hasCollided := false
	bitOffset := x % 8     //This is the offset the the first byte needs to be shifts right
	startingXByte := x / 8 //First byte that needs to be XORed

	for i, spriteLine := range sprite {
		display[y+byte(i)][startingXByte] ^= (spriteLine >> bitOffset)
		if bitOffset > 0 {
			display[y+byte(i)][startingXByte+1] ^= (spriteLine << (8 - bitOffset)) //Shift right for the second byte
		}
	}

	return hasCollided
}

func (display *Display) ToBoolArray() [defaultHeight][defaultWidth]bool {
	result := [defaultHeight][defaultWidth]bool{}

	for y := 0; y < defaultHeight; y++ {
		result[y] = rowToBoolArray(&display[y])
	}

	return result
}

func rowToBoolArray(row *[defaultByteWidth]byte) [defaultWidth]bool {
	result := [defaultWidth]bool{}

	currentByte := byte(0)
	for bit, _ := range result {
		if bit%8 == 0 {
			currentByte = row[bit/8]
		}

		result[bit] = (currentByte & 0x10) > 0 //Mask off last bit
		currentByte <<= 1                      //Shift left to get the next bit
	}

	return result
}
