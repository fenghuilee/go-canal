package lua

import (
	"github.com/pingcap/errors"
	libs "github.com/vadv/gopher-lua-libs"
	lua "github.com/yuin/gopher-lua"
	"go-canal/libs/canal/canal_entity"
	. "go-canal/libs/lua/lua_lua"
	"go-canal/libs/lua/lua_module"
	"go-canal/libs/lua/lua_rule"
	"sync"
)

var pool = sync.Pool{
	New: func() interface{} {
		L := lua.NewState()
		libs.Preload(L)
		lua_module.Preload(L)
		return L
	},
}

func Get() *lua.LState {
	return pool.Get().(*lua.LState)
}

func Put(L *lua.LState) {
	pool.Put(L)
}

func DoScript(entity *canal_entity.TRowEntity, rule *lua_rule.TRule) error {
	L := Get()
	defer Put(L)

	L.SetGlobal(GlobalROW, MapToTable(L, entity.Row))
	L.SetGlobal(GlobalORW, MapToTable(L, entity.OldRow))
	L.SetGlobal(GlobalACT, lua.LString(entity.Action))
	L.SetGlobal(GlobalTBL, lua.LString(entity.Table))
	L.SetGlobal(GlobalSCM, lua.LString(entity.Schema))

	funcFromProto := L.NewFunctionFromProto(rule.LuaProto)
	L.Push(funcFromProto)

	return errors.Trace(L.PCall(0, lua.MultRet, nil))
}
