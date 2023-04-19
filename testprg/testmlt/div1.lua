function arrange()
end

function assert() 
    res_mem = get_memory(0xDE14, 4)
    return res_mem == "01004444", res_mem
end