-- Creating a Lua script VM for a scene object and communicating with a Script component

hg = require("harfang")

scene = hg.Scene()
script = scene:CreateScript('example')

lua_vm = hg.SceneLuaVM()
lua_vm:CreateScriptFromSource(scene, script, "\
a = 4\
\
function CallToPrintA() print('CallToPrintA: '..a) end\
function CallToPrintV(v) print('CallToPrint: '..v) end\
function CallToPrintScriptPath(s) print('CallToPrintScriptPath: '..s:GetPath()) end\
\
function CallToReturnValue() return 'String returned from scene VM to host VM' end\
")

a = lua_vm:GetScriptValue(script, 'a')
print('GetScriptValue returned a='..a)

lua_vm:SetScriptValue(script, 'a', 24)

a = lua_vm:GetScriptValue(script, 'a')
print('GetScriptValue returned a='..a)

assert(lua_vm:Call(script, 'CallToPrintA', hg.LuaObjectList()) == true)
assert(lua_vm:Call(script, 'CallToPrintV', {8}) == true)
assert(lua_vm:Call(script, 'CallToPrintScriptPath', {script}) == true)

assert(lua_vm:Call(script, 'InvalidCall', {}) == false)

success, rvalues = lua_vm:Call(script, 'CallToReturnValue', {})
assert(success == true)

print('CallToReturnValue return value='..rvalues[1])
