trap_nr = 123

function arrange()
    write_byte(load_address+3, trap_nr)
end

function trap(trap_code)
    set_yreg(trap_code)
end

function assert()
    preexec_value = read_byte(3*4096+42)
    test_data = preexec_value == 42
    if (not test_data) then
        return false, "Prexec value not found: " .. preexec_value
    end

    x_reg = get_xreg()
    y_reg = get_yreg()
    accu = get_accu()
    return (x_reg == trap_nr) and (accu == 0x42) and (y_reg == trap_nr+1) and test_data, "Registers contain unexpected values"
end

function cleanup()

end