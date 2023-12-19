data = {0x48, 0x45, 0x4c, 0x4c, 0x4f, 0x20, 0x46, 0x52, 0x4f, 0x4d, 0x20, 0x4c, 0x55, 0x41, 0x0d}

function arrange()
    read_ptr = 1
end

function trap(trap_code)
    set_accu(data[read_ptr])
    read_ptr = read_ptr + 1
end

function assert()
    local preexec_value = read_byte(3*4096+42)
    local test_data = preexec_value == 42
    if (not test_data) then
        return false, "Prexec value not found: " .. preexec_value
    end

    local hi_addr = read_byte(load_address + 4)
    local lo_addr = read_byte(load_address + 3)
    local input_buffer_addr = hi_addr * 256 + lo_addr
    local data_read = get_memory(input_buffer_addr, #data)

    return data_read == "48454c4c4f2046524f4d204c55410d", "Read incorrect data: " .. data_read
end

function cleanup()

end