test_data = {
    [0] = {val1 = 17, val2= 34},
    [1] = {val1 = 234, val2= 178},
    [2] = {val1 = 254, val2= 255},
    [3] = {val1 = 0, val2= 255}    
}

function num_iterations() 
    return 4
end

iter_count = 0

function arrange()
    set_pc(load_address)
    set_accu(test_data[iter_count].val1)
    set_xreg(test_data[iter_count].val2)
end

function assert()
    in1 = test_data[iter_count].val1
    in2 = test_data[iter_count].val2
    res = (get_accu() * 256) + get_xreg()

    iter_count = iter_count + 1
    return res == in1 * in2, "Unexpected value"
end