# `6502profiler`

## Overview

This software is in essence an emulator for the MOS 6502, 6510 and the 65C02 microprocessors. In contrast to the plethora of 
emulators that already exist for these microprocessors it does not aim to emulate an existing retro computer with all its features like graphics 
and sound. It is rather intended to be a development tool for optimizing and verifying the implementation of algorithms on old or new
machines which use these classic microprocessors. To state this again: No graphics or sound capabilities of any kind are emulated and 
`6502profiler` works at a purely logical level.

The two main use cases for `6502profiler` are unit testing for and performance analysis of 6502 assembly programs. `6502profiler` offers the 
possibility to implement tests where arranging the test data and evaluating the results is offloaded to a Lua script.

When used for performance analysis `6502profiler` executes an existing binary inside the emulator. While running the program the number of clock 
cycles that are used up during execution are counted. Additionally `6502profiler` can be used to identify "hot spots" in the program because it 
also keeps track of how many times each byte in memory is accessed (i.e. read and/or written). 

**Caution: This is work in progress, things will change and maybe even break.**

## How to use `6502profiler`

`6502profiler` has a command line interface. The first parameter is a so called command. Currently the following
commands are implemented.

```
The following commands are available: 
     info: Return info about program
     list: List all test cases and their descriptions
     newcase: Create a new test case skeleton
     profile: Run program, record and evaulute performance data
     verify: Run a test on an assembler program
     verifyall: Run all tests
```

`6502profiler` expects an installed assembler for most of its functionality to work. Its location in the file system can be configured through 
the `AcmeBinary` configuration entry. Currently `acme`, `64tass` and  `ca65` are supported. The type of assembler which is in use can be defined 
through the `AsmType` config entry. Use the values `acme`, `64tass` or `ca65` to set the entry.

## The `profile` command

This command is intended to collect and evaluate data about the performance of the program under test. It understands the following 
command line options:

```
./6502profiler profile -h
Usage of 6502profiler profile:
  -c string
    	Config file name
  -dump string
    	Dump memory after program has stopped. Format 'startaddr:len'
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

The most important option is the `-prg` option which is used to specify the binary to run. It is expected that the first two bytes
of the binary contain the load address in the usual form (lo byte first). It is also expected that execution of the program will start at
this address. The final instruction in a program that is run by `6502profiler` has to be `BRK`. `BRK` halts the emulator.

If the `-out` option is specified `6502profiler` outputs statistical data about the current program execution. The output contains two
types of lines. Label lines and address lines. The following example illustrates a label line followed by three address lines.

```
SQ_TAB_LSB
     0803: 00 50350
     0804: 01 1056087
###  0805: 8A 5591244
``` 

Label lines are created from the data contained in the symbol list file generated by the `acme` or `64tass` macro assember when 
called with the `-l` option. The path to this file can be provided through the `-label` option of `6502profiler`. Label lines 
serve as a basic link between the output of `6502profiler` and the source code of the program that is evaluated. `ca65` does not 
offer an easy way to just dump the addresses for which all the labels stand into a file and therefore specifying a label
file is optional.

An address line contains a 16 bit hex address followed by a colon. The address is followed by the byte stored at this memory location 
at the end of the execution of the program. This in turn is follwed by the number of times the address has been accessed (read and written) 
by the running program. When an address line starts with `###` the corresponding address has been accessed "more often" than is usual during 
program execution. The meaning of "more often" is defined by the options `-strategy` and `-prcnt`. 

`6502profiler` counts how often each byte is accessed during program execution and stores these so called access numbers so that they
can be evaluated after the program has teminated. If an access number for a byte in memory where a machine language program resides
is high then this means that the corresponding code is executed often and therefore optmizing these parts of the program has potentially 
a big effect on overall performance.

The option `-prcnt` has to be a number between 0 and 100 and denotes a percentage. Specifying any value `p` results in those addresses
in the output to be flagged with `###` that have access numbers which are in the top `p` percent of all access numbers.

The `-strategy` option determines what *all* access numbers are. When the value `median` is used all access numbers (including duplicates) 
are sorted and and the lowest value in the top `p`% (`p` being the value of the `-prcnt` option) is selected as a threshold. Any 
other value for the `-strategy` option sorts the access values after removing all duplicates and after that determines the threshold. 
In first experiments no significant differences between the two strategies have been found. `-prcnt` and `-strategy` are optional. The 
default values for these options are 10 and `median`. 

The report file created by `ACME` when specifying its `-r` option, or the listings generated by `64tass` (`-L` option) or `ca65` (`-l` option) 
can be used to more precisely link the output of  `6502profiler` to the assembly source code. 

The `-dump` command line option can be used to print a hex dump of a portion of the emulator's memory to the screen after the program has 
finished. The start address and length of the memory to dump can be selected by the parameter of the option using the format `address:length`.
Both numbers have to be specified in decimal.

## The `verify` and `verifyall` commands

These commands are intended to facilitate the testing of assembly subroutines. You can see `6502profiler`
in action for this purpose in my 6502 arithmetic library [project](https://github.com/rmsk2/6502-Arithmetic).
The `verify` command can be used to run one specific test case and its command line syntax is as follows:

```
Usage of 6502profiler verify:
  -c string
    	Config file name
  -t string
    	Test case file
  -verbose
    	Give more information
```

The name of the test case file is interpreted relative to the directory specified by the `AcmeTestDir` configuration entry (see below). 
The `.json` suffix of the filename can be omitted. In order to run all test cases in that directory see the `verifyall` command as 
described below.

The general idea is to have a collection of source files which contain the assembly subroutines to test in one directory (the source 
directory as given in `AcmeSrcDir`) and additional separate assembly test driver programs in a test directory (named by `AcmeTestDir`) 
which call the routines that are to be tested in an appropriate fashion. The test drivers use the `!source` or `.include` pseudo opcode 
of the chosen assembler to access routines from the source directory. The test drivers are automatically assembled (or compiled) into 
the test binary directory. This directory is specified by `AcmeBinDir`.

The `verify` command loads the test driver binary and a corresponding Lua test script. This script has to define at least
two functions `arrange` and `assert`. Before running the test driver in the emulator the `verify` command calls the `arrange`
function in the Lua script which can modify the emulator state before the test driver is run (for instance to arrange test data). 
Then the test driver is run by the emulator and after it finishes the `assert` function of the test script is called to evaluate 
whether the program returned the expected results. The test is successfull if the `assert` script returns `true`.

The source files for the test driver and the Lua test script have to be referenced in a JSON test case file which has the following
format:

```
{
    "Name": "Simple loop test",
    "TestDriverSource": "test1.a",
    "TestScript": "test1.lua"
}

```

The file names in this file are interpreted relative to the directory specified by the `AcmeTestDir` configuration entry. 

Here an example for a test driver and a test script. Let's say we want to test the subroutine `simpleLoop` defined in `test_loop.a` 
in the source directory. This routine is expected to copy a four byte vector stored at the load address plus three bytes to the memory 
starting at the load address plus seven bytes. The test driver looks as follows and is stored as `test1.a` in the test directory.

```
* = $0800

jmp testStart

!source "test_loop.a"

testStart
    jsr simpleLoop
    brk
```

It is always assumed that a test driver starts its execution at the load address (here specified by the `* = $xxxx` expression). The 
corresponding Lua test script is also stored  (as `test1.lua` ) in the test directory:

```lua
test_vector = "10203040"

function arrange()
    set_memory(load_address+3, test_vector)
end

function assert()
    d = get_memory(load_address+7, 4)
    fl = get_flags()
    data_ok = (d == test_vector)
    negative_is_set = (string.find(fl, "N", 0, true) ~= nil)
    error_msg = " \n"

    res = data_ok and negative_is_set

    if not data_ok then
        error_msg = error_msg .. string.format("data wrong '%s'\n", d)
    end

    if not negative_is_set then
        error_msg = error_msg .. string.format("negative flag not set: %s\n", fl)
    end

    return res, error_msg
end
```

The `arrange` function copies the test vector into the emulator's memory before the test driver is run. After the test driver has finished the
`assert` function is called to evaluate the results. In this example it is tested whether the test vector has been copied to the correct address
and if the negative flag is set at the end of the test driver. If these conditions are not met corresponding error messages are returned. Have a
look in the `testprg/tests`directory in this repo for additional examples.

## Structure of test scripts

Test scripts have to implement an `assert` and an `arrange` function. `arrange` is expected to take no arguments and return no value. `assert` 
also takes no arguments but has to return two values. The first one is a boolean and is set to true if the test was successfull. The second 
return value is a string and should contain some helpful message in case the test has failed. The following functions can be used in Lua to 
query and manipulate the emulator's memory and processor state:

|Function Name| Description |
|-|-|
| `set_memory(address, hex_data)` | Store the data given in `hex_data` at address `address`| 
| `get_memory(address, length)` | Return `length` bytes from the emulator beginning with the byte at address `address` as a hex string| 
| `read_byte(address)`| Returns a single byte from memory at the given address|
| `write_byte(address, value)`| Writes a single byte to memory at the given address|
| `get_flags()` | Returns an eight character string that contains the letters `NVBDIZC-`. A letter is in the string if the corresponding flag is set|
| `set_flags(flag_data)` | Sets the value of the flag register. If `flag_data` contains any of the letters described above the corresponding flag is set. Using `""` clears all flags |
| `get_pc()` | Returns the program counter |
| `get_sp()`| Returns the stack pointer |
| `get_accu()` | Returns the value stored in the accumulator | 
| `set_accu(val)` | Stores `val` in the accu | 
| `get_xreg()` | Returns the value stored in the X register | 
| `set_xreg(val)` | Stores `val` in the X register | 
| `get_yreg()` | Returns the value stored in the Y register | 
| `set_yreg(val)` | Stores `val` in the Y register | 
| `get_cycles()` | Returns the number of clock cycles used for executing the test |


The `set_memory` and `get_memory` functions can be used to get and set blocks of emulator memory. These memory blocks are always 
represented as hex strings. On top of that the following three variables are injected into the Lua script from the Go host program:

|Variable Name| Description |
|-|-|
| `load_address` | Address to which the test driver has been loaded and from which it is run | 
| `prog_len` | Length in bytes of the loaded test driver | 
| `test_dir` | Path to the test dir which can be used with `require` to load additional scripts | 

Assigning a value to these variables remains local to the Lua test script and does not influence what is happening in the golang
host application.

## The `verifyall` comand

The `verifyall` command can be used to execute all test cases that are found in the `AcmeTestDir` as defined in the referenced
config file. It has the following syntax:

```
Usage of 6502profiler verifyall:
  -c string
    	Config file name
  -prexec string
    	Program to run before first test
  -verbose
    	Give more information
```

Here an example what kind of output `6502profiler verifyall -c config.json` generates

```
Executing test case '32 bit multiplication 1' ... (3180 clock cycles) OK
Executing test case '16 Bit multiplication 5' ... (98 clock cycles) OK
Executing test case '32 Bit is zero 3' ... (67 clock cycles) OK
Executing test case '32 Bit compare 5' ... (132 clock cycles) OK
Executing test case '32 bit multiplication 4' ... (3052 clock cycles) OK
Executing test case '32 Bit addition test 1' ... (166 clock cycles) OK
Executing test case '32 Bit is equal 2' ... (162 clock cycles) OK
Executing test case '32 Bit is zero 4' ... (34 clock cycles) OK
```

The `-prexec` command line option can be used to specify the source code of an assembly program that is compiled and run before the 
first test in order to perform a global test setup. The program name is interpreted relative to the `AcmeTestDir` defined in the config 
file.

## The `newcase` command

This command can be used to create a JSON test case file, a Lua script and a test driver file in the test directory. It
uses the following command line options.  

```
Usage of 6502profiler newcase:
  -c string
    	Config file name
  -d string
    	Test description
  -ext string
    	Extension to use for assembly test driver files (optional)
  -p string
    	Test case file name
  -t string
    	Full name of test driver file in test dir (optional)
```

The value of `-p` is used to generate the file names of all three files in the test directory by appending the corresponding 
file endings `.json`, `.a` and `.lua`. If `-t` is specified the test driver name in the newly created test case is set to the 
value of `-t`. This value has to include the file ending (typically `.a`) and is interpreted as a file name relative to `AcmeTestDir`.
The `-d` option is used to add a description to the test case file which is printed when the test is run. Through the option
`-ext` an alternative file extension for the assembly test driverfiles can be specified in case you do not like the default value 
of `.a`.

## The `list` command

The list command can be used to list the descriptions and the test case file names of all tests in the test directory. It has
the following syntax:

```
Usage of 6502profiler list:
  -c string
    	Config file name
```

# Emulator configuration

The config is stored in a JSON file and can be referenced through the `-c` option. The config file is structured as follows

```
{
    "Model": "6502",
    "MemSpec": "Linear64K",
    "IoMask": 45,
    "IoAddrConfig": {
        "221": "file:output.bin",
        "222": "stdout:16",
        "223": "printer:petscii"   
    },
    "PreLoad": {
        "40960": "/home/martin/data/vice_roms/C64/basic",
        "57344": "/home/martin/data/vice_roms/C64/kernal"
    },
    "AsmType": "acme",    
    "AcmeBinary": "acme",
    "AcmeSrcDir": "./testprg",
    "AcmeBinDir": "./testprg/tests/bin",
    "AcmeTestDir": "./testprg/tests"
}
```

`Model` can be `6502` or `65C02`. At the moment `MemSpec` can be `Linear16K`, `Linear32K`, `Linear48K`, `Linear64K`, 
`XSixteen512K`, `XSixteen2048K`, `GeoRam_512K` or `GeoRam_2048K`. The linear memory specifications denote a contiguous 
chunk of memory starting at address 0 with a length of 16, 32, 48 or 64 kilobytes. The `XSixteen` memory specifications 
configure the emulator to use the memory model of the Commander X16 with either 512K oder 2048K of banked RAM. `GeoRam_512K` 
and `GeoRam_2048K` can be used to emulate the memory model of the GeoRAM or NeoRAM cartridge, where the memory locations 
`$DFFE` and `$DFFF` select which 256 byte page of extended memory is banked into the address space starting at `$DE00`. 

`IoMask` and `IoAddrConfig` can be used to configure special I/O adresses that allow to exfiltrate data from the emulator by 
means of writing to a special virtual I/O address. 

The value in `IoMask` specifies the hi byte of all such special addresses and each entry in `IoAddrConfig` specifies the 
corresponding lo byte of one special address. In the example above the first resulting special address is `$2ddd` 
($2d=45, $dd=221). Each entry  in `IoAddrConfig` also has to contain a specification of what should happen each time data is 
stored in that address via `sta`, `stx`, `sty` or instructions that modify data in place as for instance `inc`. If no such 
special addresses are needed then `IoAddrConfig` should be empty.

Currently three types of special IO addresses are defined. One stores the data written to the corresponding address in a file. 
Such an address is defined by an entry that start with `file:` and the remaining string specifies the file name. The second 
type of special IO address outputs the data hex encoded to stdout. Such entries start with `sdtdout:`and the remaining part of
the entry specifies the number of bytes to be printed on one line as a decimal number. The third type prints the bytes to
stdout as characters. The value after the colon specifies the encoding to use. At the moment the only legal value is `petscii`. 

If you want to load binaries into the emulator's RAM before any program is run you can list these binaries in the `PreLoad`
property. Each entry is a key value pair where the key is the address to which the binary should be loaded and the value is
the name of the file which contains the binary to load. This can for instance be used to load ROM images. It has to be noted
though that these images are of limited use because `6502profiler` does not emulate any I/O, timing or interrupt behaviour.

The `AcmeBinary`entry defines the path to the binary of the assembler to use. If the program is in your `PATH` then the name
of the binary suffices. `AcmeSrcDir` has to describe the path to the directory where the assembler source files (which do 
not implement the tests themselves) are stored. `AcmeTestDir` holds the directory where the test case files, the assembler 
source for the test drivers and the test scripts are located. Assembled test drivers are stored in the directory referenced 
by `AcmeBinDir`. The entry `AsmType` specifies the assembler to use. Currently the values `acme`, `64tass` and `ca65` are 
allowed.

When using `ca65` the value of `AcmeBinary` only has to specify the path to the tools `ca65` and `cl65` but it must not
contain the names of the tools themselves. If for instance `ca65` and `cl65` are located in `/usr/bin` you can set `AcmeBinary`
to `/usr/bin`. If the tools are in your `PATH` then you can simply use `""`. 

# Performance

I have used `6502profiler` to further optimize the calculation routines for my [C64](https://github.com/rmsk2/c64_mandelbrot) 
and [Commander X16](https://github.com/rmsk2/X16_mandelbrot) Mandelbrot set viewers. A C64 needs about 75 minutes to create 
the default visualization in hires mode using a program of 1827 bytes length. `6502profiler` executes this program in about
a minute. The corresponding assembler source code can be found in `testprg/fixed_test.a` and `testprg/fixed_point.a`

# Limitations

Currently all 6502/6510/65C02 addressing modes and all but three instruction are emulated. The first is `RTI` as I do not
see any use for this instruction on the purely logical level on which `6502profiler` operates. Furthermore the 65C02 instructions
`STP` and `WAI` are also not implemented for the same reason.

# Building `6502profiler`

The software is written in Go and therefore it can be built by the usual `go build` command. Tests are provided for all
6502 instructions and can be executed through `go test ./...`.

# Upcoming

We will see ... .
