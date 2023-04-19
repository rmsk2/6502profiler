function arrange()
end

function assert() 
    res_mem = get_memory(0xDE04, 4)
    return res_mem == "158c1b25", res_mem
end