function arrange()
    setmemory("10203040", 2048+3)
end

function assert()
    d = getmemory(2048+7, 4)
    return d == "10203040"
end