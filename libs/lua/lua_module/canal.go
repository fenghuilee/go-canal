package lua_module

import (
	lua "github.com/yuin/gopher-lua"
)

func canalModule(L *lua.LState) int {
	t := L.NewTable()
	L.SetFuncs(t, canalModuleApi)
	L.Push(t)
	return 1
}

var canalModuleApi = map[string]lua.LGFunction{
	"rawRow":    rawRow,
	"rawOldRow": rawOldRow,
	"rawAction": rawAction,
	"rawTable":  rawTable,
	"rawSchema": rawSchema,
}
