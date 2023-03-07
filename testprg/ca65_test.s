start:
    lda #<GEO_RAM
    sta $80
    lda #>GEO_RAM
    sta $81

    lda #0
    sta TRACK_REG
    sta SECTOR_REG

    ldy #5
    lda #42
    sta ($80), y

    inc SECTOR_REG
    lda #43
    sta ($80), y
    
    lda #0
    sta SECTOR_REG
    lda ($80), y
    sta OUT_1

    inc SECTOR_REG
    lda ($80), y
    sta OUT_2
    brk