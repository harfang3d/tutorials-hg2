# Creating a Lua script VM for a scene object and communicating with a Script component

import harfang as hg

scene = hg.Scene()
script = scene.CreateScript('example')

lua_vm = hg.SceneLuaVM()
lua_vm.CreateScriptFromSource(scene, script, """
a = 4

function CallToPrintA() print('CallToPrintA: '..a) end
function CallToPrintV(v) print('CallToPrint: '..v) end
function CallToPrintScriptPath(s) print('CallToPrintScriptPath: '..s:GetPath()) end

function CallToReturnValue() return 'String returned from scene VM to host VM' end
""")

a = lua_vm.GetScriptValue(script, 'a')
print('GetScriptValue returned a=' + str(lua_vm.Unpack(a)))

lua_vm.SetScriptValue(script, 'a', lua_vm.Pack(24))

a = lua_vm.GetScriptValue(script, 'a')
print('GetScriptValue returned a=' + str(lua_vm.Unpack(a)))

assert lua_vm.Call(script, 'CallToPrintA', [])[0] == True
assert lua_vm.Call(script, 'CallToPrintV', [lua_vm.Pack(8)])[0] == True
assert lua_vm.Call(script, 'CallToPrintScriptPath', [lua_vm.Pack(script)])[0] == True

assert lua_vm.Call(script, 'InvalidCall', [])[0] == False

success, rvalues = lua_vm.Call(script, 'CallToReturnValue', [])
assert success == True

print('CallToReturnValue return value=' + str(lua_vm.Unpack(rvalues[0])))
