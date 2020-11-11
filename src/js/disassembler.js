

class InstructionDescriptor {
    constructor(op, args, bytes) {
        this.op = op;
        if (Array.isArray(args)) {
            this.args = args;
        } else {
            this.args = [args];
        }
        this.raw = bytes;
    }

    toString() {
        return `${this.op} ${this.args.join(', ')}`;
    }
}

class MemoryAddress {
    constructor(value) {
        this.value = value;
    }

    toString() {
        return `a${toByteHexString(this.value, 3)}`;
    }
}

class HexByte {
    constructor(value) {
        this.value = value;
    }

    toString() {
        return `0x${toByteHexString(this.value, 2)}`;
    }
}

class Nibble {
    constructor(value) {
        this.value = value;
    }

    toString() {
        return `0x${toByteHexString(this.value, 1)}`;
    }
}

class Register {
    constructor(value) {
        this.value = value;
    }

    toString() {
        return `V_${toByteHexString(this.value, 1)}`;
    }
}

function chip8Disassemble(machineCode) {
    let instruction = [];
    let decoded = [];
    machineCode.forEach(byte => {
        instruction.push(byte);
        if (instruction.length == 2) {
            let descriptor = chip8Decode(instruction);

            decoded.push(descriptor);

            instruction = [];
        }
    });

    return decoded;
}

function chip8Decode(bytes) {
    if (bytes.length != 2) {
        throw "Invalid Instruction: All Chip8 instructions are 2 bytes";
    }

    //Divide all ops by the first nibble to start by masking off the nibble
    switch (bytes[0] & 0xF0) {
        //Various system calls
        case 0x00:
            return chip8Decode_0(bytes);

        case 0x10:
            //The argument for the JMP instruction is the last three Nibbles
            return new InstructionDescriptor("JMP", new MemoryAddress(((bytes[0] & 0xF) << 8) + bytes[1]), bytes);

        case 0x20:
            //The argument for the CAL instruction is the last three Nibbles
            return new InstructionDescriptor("CALL", new MemoryAddress(((bytes[0] & 0xF) << 8) + bytes[1]), bytes);

        case 0x30:
            return new InstructionDescriptor("SE", [new Register(bytes[0] & 0x0F), new HexByte(bytes[1])], bytes);

        case 0x40:
            return new InstructionDescriptor("SNE", [new Register(bytes[0] & 0x0F), new HexByte(bytes[1])], bytes);

        case 0x50:
            //According to CowGod docs all instructions with a 5 at the start must end in 0
            if ((bytes[1] & 0x0F) != 0) {
                return new InstructionDescriptor("???0x" + bytes.map(b => toByteHexString(b)).join('') + "???", [], bytes);
                //throw invalidInstruction(bytes);
            }
            //0x5xy0 instructions ignore the last nibble so make sure to shift bits in the second byte (no need to mask with above check)
            return new InstructionDescriptor("SE", [new Register(bytes[0] & 0x0F), new Register(bytes[1] >> 4)], bytes);

        case 0x60:
            return new InstructionDescriptor("LD", [new Register(bytes[0] & 0x0F), new HexByte(bytes[1])], bytes);

        case 0x70:
            return new InstructionDescriptor("ADD", [new Register(bytes[0] & 0x0F), new HexByte(bytes[1])], bytes);

        //various maths
        case 0x80:
            return chip8Decode_8(bytes);

        case 0x90:
            //According to CowGod docs all instructions with a 9 at the start must end in 0
            if ((bytes[1] & 0x0F) != 0) {
                return new InstructionDescriptor("???0x" + bytes.map(b => toByteHexString(b)).join('') + "???", [], bytes);
                //throw invalidInstruction(bytes);
            }
            //0x9xy0 instructions ignore the last nibble so make sure to shift bits in the second byte (no need to mask with above check)
            return new InstructionDescriptor("SNE", [new Register(bytes[0] & 0x0F), new Register(bytes[1] >> 4)], bytes);

        case 0xA0:
            //The argument for this LD instruction is the last three Nibbles
            return new InstructionDescriptor("LD", ["I", new MemoryAddress(((bytes[0] & 0xF) << 8) + bytes[1])], bytes);

        case 0xB0:
            //The argument for this JP instruction is the last three Nibbles
            return new InstructionDescriptor("JP", [new Register(0), new MemoryAddress(((bytes[0] & 0xF) << 8) + bytes[1])], bytes);

        case 0xC0:
            return new InstructionDescriptor("RND", [new Register(bytes[0] & 0x0F), new HexByte(bytes[1])], bytes);

        case 0xD0:
            return new InstructionDescriptor("DRW", [new Register(bytes[0] & 0x0F), new Register(bytes[1] >> 4), new Nibble(bytes[1] & 0x0F)], bytes);

        //various skips
        case 0xE0:
            return chip8Decode_E(bytes);

        //MISC
        case 0xF0:
            return chip8Decode_F(bytes);

        default:
            throw "Invalid Instruction: First Byte larger than 0xFF or isn't a number";
    }
}

function invalidInstruction(bytes) {
    var byteStrings = bytes.map(byte => {
        return toByteHexString(byte, 2);
    });

    return `Invalid instruction: 0x${byteStrings.join('')}`;
}

function toByteHexString(byte, nibbleLength = 2) {
    return byte.toString(16).toUpperCase().padStart(nibbleLength, '0');
}

//methods to simplify the above switch
function chip8Decode_0(bytes) {
    //All of these should start with 0x00
    if (bytes[0] != 00) {
        return new InstructionDescriptor("IGNORE SYS CALL", bytes);
    }

    switch (bytes[1]) {
        //Clear display
        case 0xE0:
            return new InstructionDescriptor("CLS", [], bytes);

        //Subroutine return
        case 0xEE:
            return new InstructionDescriptor("RET", [], bytes);

        //**Super Chip-48 only
        case 0xFB:
            return new InstructionDescriptor("SCR*", [], bytes);
        case 0xFC:
            return new InstructionDescriptor("SCL*", [], bytes);
        case 0xFD:
            return new InstructionDescriptor("EXIT*", [], bytes);
        case 0xFE:
            return new InstructionDescriptor("LOW*", [], bytes);
        case 0xFF:
            return new InstructionDescriptor("HIGH*", [], bytes);

        default:
            //Check for 0xC in the first nibble for SCD instruction by masking
            if ((bytes[1] & 0xF0) == 0xC0) {
                //SCD instruction takes the last nibble as an argument
                return new InstructionDescriptor("SCD*", [], new Nibble(bytes[1] & 0x0F), bytes);
            } else {
                return new InstructionDescriptor("IGNORE SYS CALL", bytes);
            }
    }
}

function chip8Decode_8(bytes) {
    //0x8xxx instruction change base on the last nibble of the instruction so mask accordingly
    //0x8xy0 instructions ignore the last nibble after decode so make sure to shift bits in the second byte (no need to mask?)
    switch (bytes[1] & 0x0F) {
        case 0x0:
            return new InstructionDescriptor("LD", [new Register(bytes[0] & 0x0F), new Register(bytes[1] >> 4)], bytes);

        case 0x1:
            return new InstructionDescriptor("OR", [new Register(bytes[0] & 0x0F), new Register(bytes[1] >> 4)], bytes);

        case 0x2:
            return new InstructionDescriptor("AND", [new Register(bytes[0] & 0x0F), new Register(bytes[1] >> 4)], bytes);

        case 0x3:
            return new InstructionDescriptor("XOR", [new Register(bytes[0] & 0x0F), new Register(bytes[1] >> 4)], bytes);

        case 0x4:
            return new InstructionDescriptor("ADD", [new Register(bytes[0] & 0x0F), new Register(bytes[1] >> 4)], bytes);

        case 0x5:
            return new InstructionDescriptor("SUB", [new Register(bytes[0] & 0x0F), new Register(bytes[1] >> 4)], bytes);

        case 0x6:
            return new InstructionDescriptor("SHR", [new Register(bytes[0] & 0x0F), new Register(bytes[1] >> 4)], bytes);

        case 0x7:
            return new InstructionDescriptor("SUBN", [new Register(bytes[0] & 0x0F), new Register(bytes[1] >> 4)], bytes);

        case 0xE:
            return new InstructionDescriptor("SHL", [new Register(bytes[0] & 0x0F), new Register(bytes[1] >> 4)], bytes);

        default:
            return new InstructionDescriptor("???0x" + bytes.map(b => toByteHexString(b)).join('') + "???", [], bytes);
        //throw invalidInstruction(bytes);
    }
}

function chip8Decode_E(bytes) {
    //All 0xExxx instruction are based on the last byte
    switch (bytes[1]) {
        case 0x9E:
            return new InstructionDescriptor("SKP", new Register(bytes[0] & 0x0F), bytes);
        case 0xA1:
            return new InstructionDescriptor("SKNP", new Register(bytes[0] & 0x0F), bytes);

        default:
            return new InstructionDescriptor("???0x" + bytes.map(b => toByteHexString(b)).join('') + "???", [], bytes);
        //throw invalidInstruction(bytes);
    }

}

function chip8Decode_F(bytes) {
    //All 0xFxxx instruction are based on the last byte
    switch (bytes[1]) {
        case 0x07:
            return new InstructionDescriptor("LD", [new Register(bytes[0] & 0x0F), 'DT'], bytes);

        case 0x0A:
            return new InstructionDescriptor("LD", [new Register(bytes[0] & 0x0F), 'K'], bytes);

        case 0x15:
            return new InstructionDescriptor("LD", ['DT', new Register(bytes[0] & 0x0F)], bytes);

        case 0x18:
            return new InstructionDescriptor("LD", ['ST', new Register(bytes[0] & 0x0F)], bytes);

        case 0x1E:
            return new InstructionDescriptor("ADD", ['I', new Register(bytes[0] & 0x0F)], bytes);

        case 0x29:
            return new InstructionDescriptor("LD", ['F', new Register(bytes[0] & 0x0F)], bytes);

        case 0x30:
            return new InstructionDescriptor("LD*", ['HF', new Register(bytes[0] & 0x0F)], bytes);

        case 0x33:
            return new InstructionDescriptor("LD", ['B', new Register(bytes[0] & 0x0F)], bytes);

        case 0x55:
            return new InstructionDescriptor("LD", ['I', new Register(bytes[0] & 0x0F)], bytes);

        case 0x65:
            return new InstructionDescriptor("LD", [new Register(bytes[0] & 0x0F), 'I'], bytes);

        case 0x75:
            return new InstructionDescriptor("LD*", ['R', new Register(bytes[0] & 0x0F)], bytes);

        case 0x85:
            return new InstructionDescriptor("LD*", [new Register(bytes[0] & 0x0F), 'R'], bytes);

        default:
            return new InstructionDescriptor("???0x" + bytes.map(b => toByteHexString(b)).join('') + "???", [], bytes);
        //throw invalidInstruction(bytes);
    }
}

