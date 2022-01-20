package chip8

import (
	"math"
	"time"
)

const (
	defaultFrequency = 1000.0 //hz
	counterFrequency = 60.0   //hz
)

type Chip8 struct {
	cpu     cpu
	ram     memory
	input   Input
	display Display

	displayChannel chan<- Display
	inputChannel   <-chan Input
	powerChannel   <-chan bool

	IsRunning      bool
	frequency      float64
	cyclesPerFrame int
}

func New(cDisp chan<- Display, cInput <-chan Input, cPower <-chan bool) *Chip8 {
	system := Chip8{}
	system.ram.loadFont()
	system.cpu.initialize(&system.ram, &system.input, &system.display)
	system.frequency = defaultFrequency
	system.cyclesPerFrame = int(math.Floor(defaultFrequency / counterFrequency))

	system.displayChannel = cDisp
	system.inputChannel = cInput
	system.powerChannel = cPower

	return &system
}

func (system *Chip8) LoadProgram(program []byte) {
	system.ram.loadProgam(program)
}

func (system *Chip8) Run() error {
	system.IsRunning = true
	batchTime := 1 / counterFrequency

	for system.IsRunning {
		batchStart := time.Now()
		for cycle := 0; cycle < system.cyclesPerFrame; cycle++ {
			err := system.cpu.cycle()
			if err != nil {
				return err
			}
		}

		for duration := time.Since(batchStart); duration.Seconds() > batchTime; duration = time.Since(batchStart) {

		}

		if system.cpu.SoundRegister > 0 {
			system.cpu.SoundRegister--
		}
		if system.cpu.DelayRegister > 0 {
			system.cpu.DelayRegister--
		}

		system.displayChannel <- system.display
	}

	return nil
}
