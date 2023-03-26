require(test_dir .. "tools")

test_vector = "10203040"

function arrange()
    set_memory(load_address+3, test_vector)
end

function assert()
    d = get_memory(load_address+7, 4)
    data_ok = (d == test_vector)
    negative_is_set = is_flag_set("N")
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