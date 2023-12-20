data = {0x48, 0x45, 0x4c, 0x4c, 0x4f, 0x20, 0x46, 0x52, 0x4f, 0x4d, 0x20, 0x4c, 0x55, 0x41, 0x0d}

function arrange()
    read_ptr = 1
end

function trap(trap_code)
    set_accu(data[read_ptr])
    read_ptr = read_ptr + 1
end

function de_ref(ptr_addr)
    local hi_addr = read_byte(ptr_addr + 1)
    local lo_addr = read_byte(ptr_addr)
    
    return hi_addr * 256 + lo_addr
end

function assert()
    local preexec_value = read_byte(3*4096+42)
    local test_data = preexec_value == 42
    if (not test_data) then
        return false, "Prexec value not found: " .. preexec_value
    end

    local data_read = get_memory(de_ref(load_address + 3), #data)

    return data_read == "48454c4c4f2046524f4d204c55410d", "Read incorrect data: " .. data_read
end

function cleanup()

end