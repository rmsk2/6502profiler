test_addr = 65536 + 10251

function arrange()
    write_byte_long(test_addr, 17)
end

function assert() 
    res = (read_byte(load_address+3) == 42) and ((read_byte(load_address+4) == 43))
    error_msg = ""

    if not res then
        error_msg = "Unexpected value"
        return res, error_msg
    end

    if read_byte_long(65536 + 5) ~= 42 then
        return false, "Wrong value when accessing linear address space"
    end

    if read_byte_long(65536 + 5 + 256) ~= 43 then
        return false, "Wrong value when accessing linear address space"
    end    

    if read_byte_long(test_addr) ~= 17 then
        return false, "Wrong value when accessing linear address space"
    end

    return true, ""
end