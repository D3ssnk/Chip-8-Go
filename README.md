# Chip-8-Go

A custom CHIP-8 emulator built entirely in Go. This project simulates the classic CHIP-8 system architecture, including CPU cycles, memory management, and display rendering, providing a runtime environment for classic CHIP-8 ROMs.

## Demo

#### Keypad Test
https://github.com/user-attachments/assets/0d4911c7-a8c0-4788-85cd-8c59ffc3385d

#### Breakout
https://github.com/user-attachments/assets/caabc610-ff38-4982-896f-7b062c359d97

#### IBM Logo
<img width="752" height="464" alt="Screenshot 2026-06-23 at 11 20 31" src="https://github.com/user-attachments/assets/3a16ac3d-0023-44da-a26b-931c14f300cb" />

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

## Controls

The original CHIP-8 used a 16-key hexadecimal keypad. This emulator maps those classic keys to the left side of a modern keyboard.

| Modern Keyboard | CHIP-8 Keypad |
| --- | --- |
| `1` `2` `3` `4` | `1` `2` `3` `C` |
| `Q` `W` `E` `R` | `4` `5` `6` `D` |
| `A` `S` `D` `F` | `7` `8` `9` `E` |
| `Z` `X` `C` `V` | `A` `0` `B` `F` |

## Running Tests

To run the test suite for the CPU logic, memory operations, and other core components, execute the following command from the root of the project:

```bash
go test ./...

```

## Architecture Overview

The core of the emulator revolves around the machine state, typically encapsulated in a central `CPU` struct. Key components include:

* **Memory:** 4096 bytes (`0x000` to `0xFFF`). ROMs are loaded into memory starting at `0x200`.
* **Registers:** 16 8-bit data registers for general operations, plus specialized registers for addresses and timers.
* **Stack:** Holds up to 16 16-bit return addresses for subroutine calls.

The standard fetch-decode-execute cycle reads 16-bit opcodes from memory, decodes the instruction using bitwise masking, and routes execution to the appropriate handler.
