function arrange()
    set_memory("10203040", load_address+3)
end

function assert()
    d = get_memory(load_address+7, 4)
    fl = get_flags()
    data_ok = (d == "10203040")
    negative_is_set = (string.find(fl, "N", 0, true) ~= nil)
    error_msg = " \n"

    res = data_ok and negative_is_set

    if not data_ok then
        error_msg = error_msg .. string.format("data wrong '%s'\n", d)
    end

    if not negative_is_set then
        error_msg = error_msg .. string.format("negative flag not set: %s\n", fl)
    end

    return res, error_msg
end