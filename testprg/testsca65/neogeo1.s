
TRACK_REG = $DFFE
SECTOR_REG = $DFFF
GEO_RAM = $DE00

jmp start

OUT_1:
.byte 32
OUT_2:
.byte 33

.include "ca65_test.s"