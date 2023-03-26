function push(data)
    sp = get_sp()
    write_byte(0x100+sp, data)
    if sp == 0 then
        set_sp(0xFF)
    else
        set_sp(sp - 1)
    end
end

function arrange()
    set_sp(0xFF)
    push(2)
    push(3)
    push(4)
end

function assert()
    mem = get_memory(load_address+3, 3)

    if mem == "040302" then
        return true, ""
    end

    return false, string.format("Unexpected value: %s", mem)
end
