package main

import (
	"io/ioutil"
	"log"
	"os"

	"gioui.org/app"
	"gongaware.org/gChip8/pkg/chip8"
	"gongaware.org/gChip8/pkg/gui"
)

func main() {
	program, err := loadRomFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	errChannel := make(chan error)

	system, displayChannel, inputChannel, powerChannel := chip8.New()
	system.LoadProgram(program)
	_ = powerChannel //Get rid of error

	go func() {
		errChannel <- system.Run()
	}()

	go func() {
		window := gui.New(displayChannel, inputChannel)
		err := window.Run()
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loadRomFile(filename string) ([]byte, error) {
	romFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	program, err := ioutil.ReadAll(romFile)
	if err != nil {
		return nil, err
	}

	return program, nil
}
