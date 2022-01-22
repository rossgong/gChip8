package gui

import (
	"gioui.org/io/key"
	"gongaware.org/gChip8/pkg/chip8"
)

var keymap = map[string]byte{
	"X": 0x0,
	"1": 0x1,
	"2": 0x2,
	"3": 0x3,
	"Q": 0x4,
	"W": 0x5,
	"E": 0x6,
	"A": 0x7,
	"S": 0x8,
	"D": 0x9,
	"Z": 0xa,
	"C": 0xb,
	"4": 0xc,
	"R": 0xd,
	"F": 0xe,
	"V": 0xf,
}

func handleKeys(event key.Event, input *chip8.Input) {
	eventKey, ok := keymap[event.Name]
	if ok {
		switch event.State {
		case key.Press:
			input.PressKey(eventKey)
		case key.Release:
			input.ReleaseKey(eventKey)
		}
	}
}
