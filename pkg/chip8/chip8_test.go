package chip8

import (
	"testing"
)

func TestJump(t *testing.T) {
	program := []byte{0x12, 0x30}
	system := createNewSystem(program)

	err := system.cpu.cycle()
	if err != nil {
		t.Error(err)
	} else {
		if system.cpu.programCounter != 0x230 {

			t.Errorf("FAIL pc=0x%.3X (expected 0x230)", system.cpu.programCounter)
		}
	}
}

func TestCall(t *testing.T) {
	program := []byte{0x22, 0x30}
	system := createNewSystem(program)

	err := system.cpu.cycle()
	if err != nil {
		t.Error(err)
	} else {
		if system.cpu.programCounter != 0x230 || system.cpu.stack[0] != 0x202 || system.cpu.stackPointer != 1 {

			t.Errorf("FAIL\npc=0x%.3X (expected 0x230)\nstoredStack=0x%.3X (expected 0x202)\nSP = %v (expected 1)", system.cpu.programCounter, system.cpu.stack[0], system.cpu.stackPointer)
		}
	}
}

func TestSkipIfEqualValue(t *testing.T) {
	program := []byte{0x32, 0x02, 0, 0, 0x32, 03}
	system := createNewSystem(program)

	t.Run("Equals", func(t *testing.T) {
		err := system.cpu.cycle()
		if err != nil {
			t.Error(err)
		} else {
			if system.cpu.programCounter != 0x204 {

				t.Errorf("FAIL on equals pc=0x%.3X (expected 0x204)", system.cpu.programCounter)
			}
		}
	})

	t.Run("Not equals", func(t *testing.T) {
		err := system.cpu.cycle()
		if err != nil {
			t.Error(err)
		} else {
			if system.cpu.programCounter != 0x206 {

				t.Errorf("FAIL on not equals pc=0x%.3X (expected 0x206)", system.cpu.programCounter)
			}
		}
	})
}

func TestSkipIfNotEqualValue(t *testing.T) {
	program := []byte{0x42, 0x02, 0x42, 03}
	system := createNewSystem(program)

	t.Run("Equals", func(t *testing.T) {
		err := system.cpu.cycle()
		if err != nil {
			t.Error(err)
		} else {
			if system.cpu.programCounter != 0x202 {
				t.Errorf("FAIL on equals pc=0x%.3X (expected 0x202)", system.cpu.programCounter)
			}
		}
	})

	t.Run("Not equals", func(t *testing.T) {
		err := system.cpu.cycle()
		if err != nil {
			t.Error(err)
		} else {
			if system.cpu.programCounter != 0x206 {

				t.Errorf("FAIL on not equals pc=0x%.3X (expected 0x206)", system.cpu.programCounter)
			}
		}
	})
}

func TestSkipIfEqualRegister(t *testing.T) {
	program := []byte{0x52, 0x20, 0, 0, 0x52, 0x30}
	system := createNewSystem(program)

	t.Run("Equals", func(t *testing.T) {
		err := system.cpu.cycle()
		if err != nil {
			t.Error(err)
		} else {
			if system.cpu.programCounter != 0x204 {
				t.Errorf("FAIL on equals pc=0x%.3X (expected 0x204)", system.cpu.programCounter)
			}
		}
	})

	t.Run("Not equals", func(t *testing.T) {
		err := system.cpu.cycle()
		if err != nil {
			t.Error(err)
		} else {
			if system.cpu.programCounter != 0x206 {

				t.Errorf("FAIL on not equals pc=0x%.3X (expected 0x206)", system.cpu.programCounter)
			}
		}
	})
}

func TestLoad(t *testing.T) {
	program := []byte{0x62, 0xFF}
	system := createNewSystem(program)

	err := system.cpu.cycle()
	if err != nil {
		t.Error(err)
	} else {
		if system.cpu.Registers[2] != 0xFF {
			t.Errorf("FAIL v2=0x%.3X (expected 0xFF)", system.cpu.Registers[2])
		}
	}
}

func TestAdd(t *testing.T) {
	program := []byte{0x72, 0x1E, 0x73, 0xFF}
	system := createNewSystem(program)

	t.Run("No Overflow", func(t *testing.T) {
		err := system.cpu.cycle()
		if err != nil {
			t.Error(err)
		} else {
			if system.cpu.Registers[2] != 0x20 || system.cpu.Registers[statusRegister] != 0 {
				t.Errorf("FAIL v2=0x%.3X (expected 0x20) status=%v (expected 0)", system.cpu.Registers[2], system.cpu.Registers[statusRegister])
			}
		}
	})

	t.Run("Overflow", func(t *testing.T) {
		err := system.cpu.cycle()
		if err != nil {
			t.Error(err)
		} else {
			if system.cpu.Registers[3] != 0x02 || system.cpu.Registers[statusRegister] != 1 {
				t.Errorf("FAIL v2=0x%.3X (expected 0x20) status=%v (expected 1)", system.cpu.Registers[2], system.cpu.Registers[statusRegister])
			}
		}
	})
}

func TestStoreBCD(t *testing.T) {
	program := []byte{0xF3, 0x33}
	system := createNewSystem(program)

	err := system.cpu.cycle()
	if err != nil {
		t.Error(err)
	} else {
		if system.cpu.ram[0x300] != 0x00 || system.cpu.ram[0x301] != 0x00 || system.cpu.ram[0x302] != 0x03 {
			t.Errorf("FAIL [I]=0x%.3X (expected 0x00)\n[I+1]=0x%.3X (expected 0x00)\n[I+2]=0x%.3X (expected 0x03)\n", system.cpu.ram[0x300], system.cpu.ram[0x301], system.cpu.ram[0x302])
		}
	}
}

//Utility functions
func createNewSystem(program []byte) *Chip8 {
	system := New(make(chan<- DotGrid), make(<-chan Input), make(<-chan bool))

	for i := range system.cpu.Registers {
		system.cpu.Registers[i] = byte(i)
	}
	system.cpu.RegisterI = 0x300
	system.LoadProgram(program)

	return system
}
