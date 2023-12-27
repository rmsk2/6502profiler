require("io")
require("string")

out_file_name = "apfel.bin"

function trap(trap_code)
    out_file:write(string.char(trap_code))
end

function cleanup()
    out_file:close()
end

function init()
    out_file = io.open(out_file_name, "wb+")
end

init()