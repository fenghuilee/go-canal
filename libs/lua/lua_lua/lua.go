package lua_lua

import (
	"encoding/json"
	lua "github.com/yuin/gopher-lua"
	"go-canal/libs"
)

const (
	GlobalROW = "___ROW___"
	GlobalORW = "___ORW___"
	GlobalACT = "___ACT___"
	GlobalTBL = "___TBL___"
	GlobalSCM = "___SCM___"
)

func MapToTable(l *lua.LState, kv libs.TMap) *lua.LTable {
	table := l.NewTable()
	if kv == nil {
		return table
	}
	for k, v := range kv {
		l.SetTable(table, lua.LString(k), InterfaceToLValue(v))
	}
	return table
}

func InterfaceToLValue(v interface{}) lua.LValue {
	switch v.(type) {
	case float64:
		ft := v.(float64)
		return lua.LNumber(ft)
	case float32:
		ft := v.(float32)
		return lua.LNumber(ft)
	case int:
		ft := v.(int)
		return lua.LNumber(ft)
	case uint:
		ft := v.(uint)
		return lua.LNumber(ft)
	case int8:
		ft := v.(int8)
		return lua.LNumber(ft)
	case uint8:
		ft := v.(uint8)
		return lua.LNumber(ft)
	case int16:
		ft := v.(int16)
		return lua.LNumber(ft)
	case uint16:
		ft := v.(uint16)
		return lua.LNumber(ft)
	case int32:
		ft := v.(int32)
		return lua.LNumber(ft)
	case uint32:
		ft := v.(uint32)
		return lua.LNumber(ft)
	case int64:
		ft := v.(int64)
		return lua.LNumber(ft)
	case uint64:
		ft := v.(uint64)
		return lua.LNumber(ft)
	case string:
		ft := v.(string)
		return lua.LString(ft)
	case []byte:
		ft := string(v.([]byte))
		return lua.LString(ft)
	case nil:
		return lua.LNil
	default:
		jsonValue, _ := json.Marshal(v)
		return lua.LString(jsonValue)
	}
}
