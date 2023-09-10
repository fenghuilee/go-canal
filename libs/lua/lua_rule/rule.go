package lua_rule

import (
	"fmt"
	"github.com/pingcap/errors"
	lua "github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/ast"
	"github.com/yuin/gopher-lua/parse"
	"go-canal/libs/log"
	"go-canal/utils/file_util"
	"go-canal/utils/sys_util"
	"os"
	"path/filepath"
	"strings"
)

type TRule struct {
	Schema      string `yaml:"schema"`
	Table       string `yaml:"table"`
	LuaPath     string `yaml:"lua_path"`
	LuaProto    *lua.FunctionProto
	LuaFunction *lua.LFunction
}

type TRules map[string]*TRule

var Config = &struct {
	Rules []*TRule `yaml:"rules"`
}{
	Rules: make([]*TRule, 0),
}
var Rules = make(TRules)

var luaProtos map[string]*lua.FunctionProto

func Init() {
	c := 0
	for _, config := range Config.Rules {
		if config.Schema == "" {
			config.Schema = "*"
		}
		if config.Table == "" {
			config.Table = "*"
		}
		c += len(strings.Split(config.Table, ","))
	}
	luaProtos = make(map[string]*lua.FunctionProto, len(Config.Rules))
	for _, config := range Config.Rules {
		if err := config.Compile(); err != nil {
			panic(errors.Trace(err))
			os.Exit(1)
		}
		config.LuaProto = luaProtos[config.LuaPath]
		tables := strings.Split(config.Table, ",")
		for index, table := range tables {
			rule := &TRule{
				Schema:      config.Schema,
				Table:       table,
				LuaPath:     config.LuaPath,
				LuaProto:    config.LuaProto,
				LuaFunction: config.LuaFunction,
			}
			key := fmt.Sprintf("%s.%s", config.Schema, table)
			Rules[key] = rule
			log.Infof("Rule #%d: %s.%s -> %s", index+1, rule.Schema, rule.Table, rule.LuaPath)
		}
	}
}

func (r *TRule) Compile() error {
	if r.LuaPath == "" {
		return errors.New("empty lua_lua path not allowed")
	}
	if _, exist := luaProtos[r.LuaPath]; exist {
		return nil
	}
	LuaPath := r.LuaPath
	if !strings.HasPrefix(r.LuaPath, "/") {
		LuaPath = filepath.Join(sys_util.CurrentDirectory(), r.LuaPath)
	}
	if !file_util.IsExist(r.LuaPath) {
		return errors.New(r.LuaPath + " lua_lua script file_util not found")
	}
	var data []byte
	var err error
	if data, err = os.ReadFile(LuaPath); err != nil {
		return errors.Trace(err)
	}
	script := string(data)

	if script == "" {
		return errors.New(r.LuaPath + " lua_lua script file_util is empty")
	}

	reader := strings.NewReader(script)
	var chunk []ast.Stmt
	if chunk, err = parse.Parse(reader, script); err != nil {
		return errors.Trace(err)
	}

	var proto *lua.FunctionProto
	if proto, err = lua.Compile(chunk, script); err != nil {
		return errors.Trace(err)
	}
	luaProtos[r.LuaPath] = proto

	return nil
}
