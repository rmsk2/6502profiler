function pop()
    local sp = get_sp()
    if sp == 0xFF then
        sp = 0x00
    else
        sp = sp + 1
    end

    set_sp(sp)
    return read_byte(0x100 + sp)
end

function trap(code)
    print("----- Running Lua code")
    print("Trap code: " .. code)
    local address_hi = pop() 
    local address_lo = pop()
    local address = address_hi * 256 + address_lo
    print("Copying data to address " .. address)
    local data_to_print = "48454c4c4f2046524f4d204c55410d0a00"
    set_memory(address, data_to_print)
    print("----- Done with Lua code")
end

function cleanup()
    print("Cleaning up Lua part")
end