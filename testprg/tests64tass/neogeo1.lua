function arrange()
end

function assert() 
    res = (read_byte(load_address+3) == 42) and ((read_byte(load_address+4) == 43))
    error_msg = ""

    if not res then
        error_msg = "Unexpected value"
    end

    return res, error_msg
end