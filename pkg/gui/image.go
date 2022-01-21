package gui

import (
	"image"
	"image/color"
	"image/draw"

	"gongaware.org/gChip8/pkg/chip8"
)

const (
	defaultScale = 5
)

//Default color definitions
var defaultOnColor color.Color = color.White
var defaultOffColor color.Color = color.Black

//pixel mask
type pixelMask struct {
	dots chip8.DotGrid

	scale int
}

func (mask pixelMask) At(x, y int) color.Color {
	if mask.dots[y/mask.scale][x/mask.scale] {
		return color.Transparent
	} else {
		return color.Opaque
	}
}

func (mask pixelMask) ColorModel() color.Model {
	return color.AlphaModel
}

func (mask pixelMask) Bounds() image.Rectangle {
	return image.Rect(0, 0, len(mask.dots[0])*mask.scale, len(mask.dots)*mask.scale)
}

func CreateImageFromDisplay(display *chip8.Display) *image.RGBA {
	//these can be replaced later as arguments for more custom images
	onImage, offImage := image.Uniform{defaultOnColor}, image.Uniform{defaultOffColor}
	scale := defaultScale

	//Create bounds for image and multiply by scale
	x, y := display.GetSize()

	//Initialize display as all off
	result := image.NewRGBA(image.Rect(0, 0, x*scale, y*scale))
	draw.Draw(result, result.Bounds(), &onImage, image.Point{}, draw.Src)

	mask := pixelMask{display.ToBoolArray(), scale}

	draw.DrawMask(result, result.Bounds(), &offImage, image.Point{}, mask, image.Point{}, draw.Src)

	return result
}
