# `6502profiler`

## Overview

This software is in essence an emulator for the MOS 6502, 6510 (and in future versions the 65C02) microprocessors. In contrast to the plethora of 
emulators that already exist for these microprocessors it does not aim to emulate an existing retro computer with all its features like graphics 
and sound. It is rather intended to be a development tool for optimizing and verifying the implementation of pure algorithms on these old computers. 
To state this again: No graphics or sound capabilities of any of the old computers are emulated and therefore routines like that can not be optimized
by using `6502profiler`.

`6502profiler` reads a binary as for instance created by the `ACME` macro assembler and executes it inside the emulator. While running the program
it counts the nunber of clock cycles that are used up during execution. On top of that it can be used to identify "hot spots" in the program because
it keeps track of how many times each memory cell is accessed (i.e. read and/or written). 

On top of that `6502profiler` offers the possibility to implement tests for assembler subroutines where arranging the test data and evaluating the
results is offloaded to a Lua script.

**Caution: This is work in pogress, things will change and maybe even break.**

## Emulator configuration

The config is stored in a JSON file and can be used through the `-c` option. The config file is structured as follows

```
{
    "Model": 1,
    "MemSpec": "Linear16K",
    "IoMask": 45,
    "IoAddrConfig": {
        "221": "file:output.bin"   
    },
    "AcmeBinary": "acme",
    "AcmeSrcDir": "./testprg",
    "AcmeBinDir": "./testprg/tests/bin",
    "AcmeTestDir": "./testprg/tests"
}
```

`Model` can be 0 or 1. The number 0 encodes a standard 6502/6510 and 1 stands for a 65C02. At the moment `MemSpec` can be 
`Linear16K`, `Linear32K` or `Linear64K` and denotes a contiguous chunk of memory starting at address 0 with a length of 16, 
32 or 64 kilobyte. `IoMask` and `IoAddrConfig` can be used to configure special I/O adresses that allow to exfiltrate data 
from the emulator by means of write instructions to a special I/O address. 

The value in `IoMask` spcifies the hi byte of all such special addresses and each entry in `IoAddrConfig` specifies the 
corresponding lo byte of one special address as well as a name of a file to which bytes are written each time data is stored 
in that address via `sta`, `stx`, `sty` or instructions that modify data in place as for instance `inc`. In the example above 
the resulting special address is `$2ddd` ($2d=45, $dd=221). Entries for such file store addresses start with `file:` and the 
remaining string specifies the file name.

The `AcmeBinary`entry defines the path to the `acme` program binary. `AcmeSrcDir` has to describe the path to the directory where 
the assembler source files (which do not implement the tests themselves) are stored. `AcmeTestDir` holds the directory where
the test case files, the assembler source for the test drivers and the test scripts are located. Assembled test drivers are 
stored in the directory referenced by `AcmeBinDir`.

# How to use `6502profiler`

`6502profiler` has a command line interface. The first parameter is a so called command. Currently the `profile` and the 
`verify` commands are implemented.

## The `profile` command

This command is intended to collect and evaluate data about the performance of the program under test. It understands the following 
command line options:

```
./6502profiler profile -h
Usage of 6502profiler profile:
  -c string
    	Config file name
  -label string
    	Path to the label file generated by the ACME assembler
  -out string
    	Path to the out file that holds the generated data
  -prcnt uint
    	Percentage used to determine cut off value (default 10)
  -prg string
    	Path to the program to run
  -strategy string
    	Strategy to determine cutoff value (default "median")
```

The most important option is the `-prg` option which can be used to specify the binary to run. It is expected that the first two bytes
of the binary contain the load address in the usual form (lo byte first). It is also expected that exeution of the program will start at
this address. The final instruction in a program that is run by `6502profiler` has to be `BRK` and not `RTS`.

If the `-out` option is specified `6502profiler` outputs statistical data about the current program execution. The output contains two
types of lines. Label lines and address lines. The following example illustrates a label line followed by three address lines.

```
SQ_TAB_LSB
     0803: 00 50350
     0804: 01 1056087
###  0805: 8A 5591244
``` 

Label lines are created by the data contained in the symbol list file generated by the `ACME` macro assember when the `-l` option is
used with `ACME`. The path to this file has to be provided through the `-label` option of `6502profiler`. Label lines serve as a basic
link between the output of `6502profiler` and the source code of the program that is evaluated. An address line contains an address 
followed by a colon. This address is followed by the byte stored at this memory location at the end of the execution of the program which in 
turn is follwed by the number of times the address has been accessed (read and written) by the running program. When an address line 
starts with `###` it has been accessed "more often" than is usual during program execution. 

The condition what "more often" actually is, is determined by the options `-strategy` and `-prcnt`. `-prcnt` has to be a number between
0 and 100. The value 25 for instance means that all addresses in the output are flagged which are in the top 25% percent of all access
numbers. The `-strategy` option determines what constitiutes the overall set of numbers. `median` sorts all access values and uses
the lowest value in the top `n%` (n being the value of the `-prcnt` option) as a threshold. Any other value for the `-strategy` option
sorts the access values after removing all duplicate values. In first experiments no significant differences between the two strategies
have been found. If `-out` is used, then `-label` also has to be specified. `-prcnt` and `-strategy` are optional. The default values 
are 10 and `median`.

The report file created by `ACME` when specifying the `-r` option can be used to more precisely link the output of `6502profiler`
to the assembly source code. This may serve as an example:

```
     5  0800 a922               lda #34
     6  0802 a22d               ldx #45
     7  0804 20080a             jsr mul16BitFast
     8  0807 00                 brk
     9                          
    10                          ; xy = (x^2 + y^2 - (x-y)^2)/2
    11                          ; The following tables contain the LSB and MSB of i^2 where i=0, ..., 255
    12                          SQ_TAB_LSB
    13  0808 0001040910192431...!byte $00, $01, $04, $09, $10, $19, $24, $31, $40, $51, $64, $79, $90, $A9, $C4, $E1
```

The first number is the line number in the source file. The second number is the address to which the machine language
instruction has been written in the output.

## The `verify` command

This command is intended to facilitate the testing of assembler subroutines. The syntax for using the `verify` command is

```
./6502profiler verify -h
Usage of 6502profiler verify:
  -c string
    	Config file name
  -t string
    	Test case file
```

The name of the test case file is interpreted relative to the `AcmeTestDir` configuration entry.

The general idea is to have a source file which contains the subroutine to test in one directory (the source directory) 
and an additional separate test driver program in a test directory which calls the routines that are to be tested in an appropriate
fashion. The test driver includes the files which contain the subroutines to test (using the `!source` pseudo opcode) from the 
source directory. Then the test driver is assembled (or compiled) into the test binary directory.

The `verify` command then loads the test driver binary and a corresponding Lua test script. This script has to define at least
two functions `arrange` and `assert`. Before running the test driver in the emulator the `verify` command calls the `arrange`
function in the Lua script which can modify the emulator state before the test driver is run (for instance to arrange test data). 
Then the test driver is run by the emulator and when that is done the `assert`function of the test script evaluates whether
the program returned the expected result. The test is successfull if the `assert` script returns `true`.

The source files for the test driver and the test script have to be referenced in a JSON test case file which has the following
format:

```
{
    "Name": "Simple loop test",
    "TestDriverSource": "test1.a",
    "TestScript": "test1.lua"
}

```

The file names are interpreted relative to the `AcmeTestDir` configuration entry. Here an example for a test driver and a test
script. Let's say we want to test the subroutine `simpleLoop` defined in `test_loop.a` in the source directory:

```
.DATA_IN
!byte 0x40,0x30,0x20,0x10
.DATA_OUT
!byte 0,0,0,0

simpleLoop
    ldy #3
    lda #<.DATA_IN
    sta $12
    lda #>.DATA_IN
    sta $13
.loop
    lda ($12), y
    sta .DATA_OUT,y
    dey
    bpl .loop
    rts
```

We then write the test driver and store it as `test1.a` in the test directory.

```
* = $0800

jmp testStart

!source "test_loop.a"

testStart
    jsr simpleLoop
    brk
```

It is assumed that the test driver starts its execution at the load address. Finally the corresponding test script is implemented
and also stored (as `test1.lua` ) in the test directory.

```
function arrange()
    set_memory("10203040", load_address+3)
end

function assert()
    d = get_memory(load_address+7, 4)
    fl = get_flags()
    data_ok = (d == "10203040")
    negative_is_set = (string.find(fl, "N", 0, true) ~= nil)

    res = data_ok and negative_is_set
    if not res then
        print()
    end

    if not data_ok then
        print(string.format("data wrong '%s'", d))
    end

    if not negative_is_set then
        print(string.format("negative flag not set: %s", fl))
    end

    return res
end
```

The `set_memory` and `get_memory` functions can be used to get and set emulator memory. Memory contents is always represented as a
hex string. On top of that the load address and the length of the test driver can be referenced in Lua by the variables `load_address`
and `prog_len`. The test can be run by `./6502profiler verify -c config.json -t test1.json`. The following functions can be used in 
Lua to query and manipulate the processor state:

|Function Name| Description |
|-|-|
| `set_memory(hex_data, address)` | Store the data given in `hex_data` at address `address`| 
| `get_memory(address, length)` | Return `length` bytes from the emulator beginning with the byte at address `address`| 
| `get_flags()` | Returns an eight character string that contains the letters `NVBDIZC-`. A letter is in the string if the corresponding flag is set|
| `get_flags(flag_data)` | Sets the value of the flag register. If `flag_data` contains any of the letters described above the corresponding flag is set. Using `""` clears all flags |

# Performance

I have used `6502profiler` to further optimize the calculation routines for my [C64](https://github.com/rmsk2/c64_mandelbrot) 
and [Commander X16](https://github.com/rmsk2/X16_mandelbrot) Mandelbrot set viewers. A C64 needs about 75 minutes to create 
the default visualization in hires mode using a program of 1827 bytes length. `6502profiler` executes this program in about
a minute. The corresponding assembler source code can be found in `testprg/fixed_test.a` and `testprg/fixed_point.a`

# Limitations

Currently all 6502/6510 addressing modes and all but one instruction are emulated. The missing instruction is `RTI` as I do not
see any use for this instruction on the purely logical level on which `6502profiler` operates. On top of that only a few
65C02 spedific addressing modes and instructions have been implemented up to this point. 

The `verify` command currently can only run one test case. This will change very soon ... .

# Building `6502profiler`

The software is written in Go and therefore it can be built by the usual `go build` command. Tests are provided for all
6502 instructions and can be executed through `go test ./...`.

# Upcoming

- Implement a feature that allows to test assembly code where the verification is done in a scripting language (the plan is to use Lua)
- Implement the additional addressing modes and instructions of the 65C02 processor
- Implement the memory model used by the Commander X16
- Maybe implement a single stepping mode
