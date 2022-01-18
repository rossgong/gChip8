package chip8

const (
	defaultFrequency = 1000.0 //hz
	counterFrequency = 60.0   //hz
)

type Chip8 struct {
	cpu cpu
	ram memory

	IsRunning bool
	frequency float64
}

func New() *Chip8 {
	system := Chip8{}
	system.ram.loadFont()
	system.cpu.initialize(&system.ram)
	system.frequency = defaultFrequency
	return &system
}

func (system *Chip8) LoadProgram(program []byte) {
	system.ram.loadProgam(program)
}

func (system *Chip8) Run() {
	system.cpu.cycle()
}
