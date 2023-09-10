package lua_module

import (
	"github.com/go-mysql-org/go-mysql/client"
	"github.com/go-mysql-org/go-mysql/mysql"
	lua "github.com/yuin/gopher-lua"
	"go-canal/libs/log"
	"go-canal/libs/lua/lua_lua"
)

func sqlModule(L *lua.LState) int {
	t := L.NewTable()
	L.SetFuncs(t, sqlModuleApi)
	L.Push(t)
	return 1
}

var sqlModuleApi = map[string]lua.LGFunction{
	"query":   query,
	"execute": execute,
}

func query(L *lua.LState) int {
	dsn := L.CheckTable(1)
	sql := L.CheckString(2)

	var conn *client.Conn
	var err error
	if conn, err = getConn(dsn); err != nil {
		log.Error(err.Error())
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	defer putConn(dsn, conn)

	var result *mysql.Result
	if result, err = conn.Execute(sql); err != nil {
		log.Error(err.Error())
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	rows := L.NewTable()
	rowNumber := result.RowNumber()
	if rowNumber > 0 {
		for i := 0; i < rowNumber; i++ {
			row := L.NewTable()
			for field, index := range result.FieldNames {
				value, err := result.GetValue(i, index)
				if err != nil {
					log.Error(err.Error())
					L.Push(lua.LNil)
					L.Push(lua.LString(err.Error()))
					return 2
				}
				L.SetTable(row, lua.LString(field), lua_lua.InterfaceToLValue(value))
			}
			L.SetTable(rows, lua.LNumber(i+1), row)
		}
	}

	L.Push(rows)
	L.Push(lua.LNil)
	return 2
}

func execute(L *lua.LState) int {
	dsn := L.CheckTable(1)
	sql := L.CheckString(2)

	var conn *client.Conn
	var err error
	if conn, err = getConn(dsn); err != nil {
		log.Error(err.Error())
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	defer putConn(dsn, conn)

	var result *mysql.Result
	if result, err = conn.Execute(sql); err != nil {
		log.Error(err.Error())
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	response := L.NewTable()
	L.SetTable(response, lua.LString("Status"), lua.LNumber(result.Status))
	L.SetTable(response, lua.LString("InsertId"), lua.LNumber(result.InsertId))
	L.SetTable(response, lua.LString("AffectedRows"), lua.LNumber(result.AffectedRows))

	L.Push(response)
	L.Push(lua.LNil)
	return 2
}
