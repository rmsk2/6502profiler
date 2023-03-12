function arrange()
    set_pc(load_address+1)
end


function assert() 
    err_msg = ""
    res_byte = read_byte(load_address)

    if res_byte ~= 10 then
        err_msg = "Unexpected memory contents"
    end

    return res_byte == 10, err_msg
end