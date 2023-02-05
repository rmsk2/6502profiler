# 6502profiler

This software is in essence an emulator for the MOS 6502, 6510 (and in future versions the 65C02) microprocessors. In contrast to the plethora of 
emulators that already exist for these microprocessors it does not aim to emulate an existing retro computer with all its features like graphics 
and sound. It is rather intended to be a development tool for optimizing (and in the future verify) the implementation of pure algorithms on these
old computers. To state this again: No graphics or sound capabilities of any of the old computers are emulated and therefore routines like that can 
not be optimized using `6502profiler`.

`6502profiler` reads a binary as for instance created by the `ACME` macro assembler and executes it inside the emulator. While running the program
it counts the nunber of clock cycles that are used up during execution. On top of that it can be used to identify "hot spots" in the program because
it keeps track of how many times each memory cell is accessed (i.e. read and/or written). 

# Upcoming

- Implement a feature that allows to test assembly code where the verification is done in a scripting language (the plan is to use Lua)
- Implement the additional addressing modes and instructions of the 65C02 processor
- Implement the memory model used by the Commander X16
