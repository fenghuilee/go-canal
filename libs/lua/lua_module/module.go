package lua_module

import (
	lua "github.com/yuin/gopher-lua"
	. "go-canal/libs/lua/lua_lua"
)

func rawRow(L *lua.LState) int {
	row := L.GetGlobal(GlobalROW)
	L.Push(row)
	return 1
}

func rawOldRow(L *lua.LState) int {
	row := L.GetGlobal(GlobalORW)
	L.Push(row)
	return 1
}

func rawAction(L *lua.LState) int {
	act := L.GetGlobal(GlobalACT)
	L.Push(act)
	return 1
}

func rawTable(L *lua.LState) int {
	tbl := L.GetGlobal(GlobalTBL)
	L.Push(tbl)
	return 1
}

func rawSchema(L *lua.LState) int {
	scm := L.GetGlobal(GlobalSCM)
	L.Push(scm)
	return 1
}

func Preload(L *lua.LState) {
	L.PreloadModule("canal", canalModule)
	L.PreloadModule("sql", sqlModule)
}
