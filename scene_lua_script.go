package main

import (
	hg "github.com/harfang3d/harfang-go/v3"
)

func main() {
	scene := hg.NewScene()
	script := scene.CreateScriptWithPath("example")

	lua_vm := hg.NewSceneLuaVM()
	lua_vm.CreateScriptFromSource(scene, script, `
a = 4

function CallToPrintA() print("CallToPrintA: "..a) end
function CallToPrintV(v) print("CallToPrint: "..v) end
function CallToPrintScriptPath(s) print("CallToPrintScriptPath: "..s:GetPath()) end

function CallToReturnValue() return "String returned from scene VM to host VM" end
`)

	//a := lua_vm.GetScriptValue(script, "a")
	//fmt.Println("GetScriptValue returned a=" + str(lua_vm.Unpack(a)))

	//lua_vm.SetScriptValue(script, "a", lua_vm.Pack(24))

	//a = lua_vm.GetScriptValue(script, "a")
	//	fmt.Println("GetScriptValue returned a=" + str(lua_vm.Unpack(a)))

	if ok, _ := lua_vm.CallWithSliceOfArgs(script, "CallToPrintA", hg.GoSliceOfLuaObject{}); !ok {
		panic("CallToPrintA")
	}
	/*if ok, _ := lua_vm.CallWithSliceOfArgs(script, "CallToPrintV", hg.GoSliceOfLuaObject{lua_vm.Pack(8)}); !ok {
		panic("CallToPrintA")
	}
	if ok, _ := lua_vm.CallWithSliceOfArgs(script, "CallToPrintScriptPath", hg.GoSliceOfLuaObject{lua_vm.Pack(script)}); !ok {
		panic("CallToPrintA")
	}*/

	if ok, _ := lua_vm.CallWithSliceOfArgs(script, "InvalidCall", hg.GoSliceOfLuaObject{}); ok {
		panic("CallToPrintA")
	}

	success, _ /*rvalues*/ := lua_vm.CallWithSliceOfArgs(script, "CallToReturnValue", hg.GoSliceOfLuaObject{})
	if !success {
		panic("CallToPrintA")
	}

	//fmt.Println("CallToReturnValue return value=" + str(lua_vm.Unpack(rvalues[0])))
}
