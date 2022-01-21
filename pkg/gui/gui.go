package gui

import (
	"fmt"
	"image"
	"image/color"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gongaware.org/gChip8/pkg/chip8"
)

type GChipGUI struct {
	window *app.Window

	currentFrame  image.Image
	bufferedFrame image.Image
	currentOps    *op.Ops

	frameBuffered bool

	//channel to engine
	// inputChannel <-chan chip8.Input
	displayChannel <-chan chip8.Display
}

func New(dispChan <-chan chip8.Display) GChipGUI {
	result := GChipGUI{}
	result.window = app.NewWindow(
		app.Title("gChip8"),
		app.Size(unit.Dp(500), unit.Dp(500)),
	)

	result.displayChannel = dispChan
	result.currentOps = new(op.Ops)
	result.currentFrame = &image.Uniform{color.Black}

	return result
}

func (gui *GChipGUI) Run() error {
	for {
		select {
		case event := <-gui.window.Events():
			err := gui.handleWindowEvent(event)
			if err != nil {
				return err
			}
		case display := <-gui.displayChannel:
			gui.bufferedFrame = CreateImageFromDisplay(&display)
			gui.frameBuffered = true
			gui.window.Invalidate()
		}
	}
}

func (gui *GChipGUI) handleWindowEvent(event event.Event) error {
	switch event := event.(type) {
	case system.DestroyEvent:
		return fmt.Errorf("closed")
	case system.FrameEvent:
		gtx := layout.NewContext(gui.currentOps, event)

		if gui.frameBuffered {
			gui.currentFrame = gui.bufferedFrame
			gui.frameBuffered = false
		}

		imageOp := paint.NewImageOp(gui.currentFrame)
		imageOp.Add(gtx.Ops)
		op.Affine(f32.Affine2D{}.Scale(f32.Pt(0, 0), f32.Pt(4, 4)))
		paint.PaintOp{}.Add(gtx.Ops)

		event.Frame(gtx.Ops)

	}
	return nil
}
