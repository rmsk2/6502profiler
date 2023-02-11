function arrange()
    setmemory("10203040", loadaddress+3)
end

function assert()
    d = getmemory(loadaddress+7, 4)
    return d == "10203040"
end