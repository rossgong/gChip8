package chip8

const (
	defaultWidth     = 64
	defaultHeight    = 32
	defaultByteWidth = defaultWidth / 8
)

type pixelsArray [defaultHeight][defaultByteWidth]byte
type DotGrid [defaultHeight][defaultWidth]bool

//Struct instead of alias for future superchip resolution support
type Display struct {
	pixels pixelsArray

	hasChanged bool
}

//Returns collison
func (display *Display) drawSprite(sprite []byte, x byte, y byte) bool {
	hasCollided := false
	bitOffset := x % 8            //This is the offset the the first byte needs to be shifts right
	startingXByte := (x % 64) / 8 //First byte that needs to be XORed

	// fmt.Printf("draw(%v,%v)*%v\n", x, y, len(sprite))
	for i, spriteLine := range sprite {
		display.pixels[y+byte(i)][startingXByte] ^= (spriteLine >> bitOffset)
		if bitOffset > 0 && startingXByte+1 < byte(len(display.pixels[0])) {
			display.pixels[y+byte(i)][startingXByte+1] ^= (spriteLine << (8 - bitOffset)) //Shift right for the second byte
		}
	}
	display.hasChanged = true

	return hasCollided
}

func (display *Display) clearScreen() {
	display.pixels = pixelsArray{}
	display.hasChanged = true
}

func (display *Display) ToBoolArray() DotGrid {
	result := [defaultHeight][defaultWidth]bool{}

	for y := 0; y < defaultHeight; y++ {
		result[y] = rowToBoolArray(&display.pixels[y])
	}

	return result
}

func (display Display) HasChanged() bool {
	return display.hasChanged
}

func (display Display) GetSize() (maxX, maxY int) {
	return defaultWidth, defaultHeight
}

func rowToBoolArray(row *[defaultByteWidth]byte) [defaultWidth]bool {
	result := [defaultWidth]bool{}

	currentByte := byte(0)
	for bit := range result {
		if bit%8 == 0 {
			currentByte = row[bit/8]
		}

		result[bit] = (currentByte & 0x80) > 0 //Mask off last bit
		currentByte <<= 1                      //Shift left to get the next bit
	}

	return result
}
