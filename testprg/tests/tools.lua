function push(data)
    sp = get_sp()
    write_byte(0x100 + sp, data)
    if sp == 0 then
        set_sp(0xFF)
    else
        set_sp(sp - 1)
    end
end

function pull()
    sp = get_sp()
    if sp == 0xFF then
        sp = 0x00
    else
        sp = sp + 1
    end

    set_sp(sp)
    return read_byte(0x100 + sp)
end

function restart()
    set_pc(load_address)
end

function is_flag_set(f)
    return string.find(get_flags(), f, 0, true) ~= nil
end