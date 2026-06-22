# Chip-8-Go

A custom CHIP-8 emulator built entirely in Go. This project simulates the classic CHIP-8 system architecture, including CPU cycles, memory management, and display rendering, providing a runtime environment for classic CHIP-8 ROMs.

## Features

* **CPU Emulation:** Accurate instruction set execution, including arithmetic, logical, and control flow opcodes.
* **Memory Management:** Robust file-loading logic to read binary ROMs and map them safely into the 4KB address space.
* **Register Handling:** Full implementation of the 16 general-purpose registers (V0-VF), the index register (I), and the program counter (PC).
* **Timers:** Implementation of both delay and sound timers operating at the standard 60Hz decrement rate.
* **Display & Input:** *(Note: Update this section with the specific rendering/input library you decide to use, such as Ebitengine, SDL2, or Raylib)*.

## Prerequisites

* [Go](https://go.dev/dl/) 1.18 or higher

## Installation

Clone the repository and build the binary:

```bash
git clone [https://github.com/D3ssnk/Chip-8-Go.git](https://github.com/D3ssnk/Chip-8-Go.git)
cd Chip-8-Go
go build -o chip8

```

## Usage

To run the emulator, pass the path to a valid CHIP-8 ROM file as a command-line argument:

```bash
./chip8 path/to/rom.ch8

```

## Architecture Overview

The core of the emulator revolves around the machine state, typically encapsulated in a central `CPU` struct. Key components include:

* **Memory:** 4096 bytes (`0x000` to `0xFFF`). ROMs are loaded into memory starting at `0x200`.
* **Registers:** 16 8-bit data registers for general operations, plus specialized registers for addresses and timers.
* **Stack:** Holds up to 16 16-bit return addresses for subroutine calls.

The standard fetch-decode-execute cycle reads 16-bit opcodes from memory, decodes the instruction using bitwise masking, and routes execution to the appropriate handler.

## Roadmap

* Complete remaining opcode implementations.
* Finalize display rendering and keyboard mapping.
* Add unit tests for core CPU instructions and register operations.

## License

[MIT](https://www.google.com/search?q=LICENSE) *(Update if using a different license)*
