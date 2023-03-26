require (test_dir .. "tools")

function arrange()
    set_sp(0xFF)
    push(2)
    push(3)
    push(4)
end

function assert()
    mem = get_memory(load_address+3, 3)

    if mem == "040302" then
        return true, ""
    end

    return false, string.format("Unexpected value: %s", mem)
end
