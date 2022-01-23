package chip8

import (
	"math"
	"time"
)

const (
	defaultFrequency = 10000.0 //hz
	counterFrequency = 60.0    //hz

	channelBuffer
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

func New() (*Chip8, <-chan Display, chan<- Input, chan<- bool) {
	system := Chip8{}
	system.ram.loadFont()
	system.cpu.initialize(&system.ram, &system.input, &system.display)
	system.frequency = defaultFrequency
	system.cyclesPerFrame = int(math.Floor(defaultFrequency / counterFrequency))

	displayChan, inputChan, powerChan := make(chan Display, channelBuffer), make(chan Input, channelBuffer), make(chan bool, channelBuffer)

	system.displayChannel = displayChan
	system.inputChannel = inputChan
	system.powerChannel = powerChan

	return &system, displayChan, inputChan, powerChan
}

func (system *Chip8) LoadProgram(program []byte) {
	system.ram.loadProgam(program)
}

func (system *Chip8) Run() error {
	system.IsRunning = true

	delayTicker := time.NewTicker(time.Second / counterFrequency)
	for system.IsRunning {
		select {
		case <-delayTicker.C:
			if system.cpu.SoundRegister > 0 {
				system.cpu.SoundRegister--
			}
			if system.cpu.DelayRegister > 0 {
				system.cpu.DelayRegister--
			}

			if system.display.hasChanged {
				system.displayChannel <- system.display
				system.display.hasChanged = false
			}
		case system.input = <-system.inputChannel:
			// fmt.Printf("%.16b\n", system.input)
		default:
			err := system.cpu.cycle()
			if err != nil {
				return err
			}
		}
	}

	return nil
}
