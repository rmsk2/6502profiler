function arrange()
    set_memory("10203040", load_address+3)
end

function assert()
    d = get_memory(load_address+7, 4)
    fl = get_flags()
    data_ok = (d == "10203040")
    negative_is_set = (string.find(fl, "N", 0, true) ~= nil)

    res = data_ok and negative_is_set
    if not res then
        print()
    end

    if not data_ok then
        print(string.format("data wrong '%s'", d))
    end

    if not negative_is_set then
        print(string.format("negative flag not set: %s", fl))
    end

    return res
end