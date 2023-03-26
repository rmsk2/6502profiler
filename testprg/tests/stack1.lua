function arrange()
    write_byte(0x01FF, 2)
    write_byte(0x01FE, 3)
    write_byte(0x01FD, 4)
    set_sp(0xFC)
end

function assert()
    mem = get_memory(load_address+3, 3)

    if mem == "040302" then
        return true, ""
    end

    return false, string.format("Unexpected value: %s", mem)
end
