function arrange()
    set_accu(34)
    set_xreg(45)
end

function assert()
    result = 256 * get_accu() + get_xreg()
    error_msg = ""

    res = (result == 1530)

    if not res then
        error_msg = string.format("Wrong result: %d", result)
    end

    return res, error_msg
end