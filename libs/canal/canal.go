package canal

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/pingcap/errors"
	"go-canal/libs/canal/canal_handler"
	"go-canal/libs/log"
	"go-canal/libs/lua/lua_rule"
	"go-canal/utils/file_util"
	"go-canal/utils/string_util"
	"go-canal/utils/sys_util"
	"os"
	"path/filepath"
)

type TConfig struct {
	Addr              string   `yaml:"addr"`
	User              string   `yaml:"user"`
	Password          string   `yaml:"password"`
	Charset           string   `yaml:"charset"`
	Flavor            string   `yaml:"flavor"`
	IncludeTableRegex []string `yaml:"include_table_regex"`
	ExcludeTableRegex []string `yaml:"exclude_table_regex"`
}

type TCanal struct {
	*canal.Canal
	Config  *TConfig
	Handler *canal_handler.TEventHandler
}

func (c *TCanal) Run() {
	c.Handler.Handle()

	ConfigFileMD5 := string_util.MD5(flag.Lookup("config").Value.String())
	positionFile := filepath.Join(sys_util.CurrentDirectory(), "runtime", ConfigFileMD5+".position.json")
	if file_util.IsExist(positionFile) {
		var position mysql.Position
		var positionJson []byte
		var err error
		if positionJson, err = os.ReadFile(positionFile); err == nil {
			if err = json.Unmarshal(positionJson, &position); err == nil {
				log.Infof("Run from position (Name: %s, Pos: %d)", position.Name, position.Pos)
				c.Canal.RunFrom(position)
				return
			}
		}
	}

	c.Canal.Run()
}

func (c *TCanal) Close() {
	c.Handler.Stop()
	c.Canal.Close()
}

var Config = &struct {
	Canal *TConfig `yaml:"canal"`
}{
	Canal: new(TConfig),
}

var Canal *TCanal

func Init() *TCanal {
	config := canal.NewDefaultConfig()
	config.Addr = Config.Canal.Addr
	config.User = Config.Canal.User
	config.Password = Config.Canal.Password
	config.Charset = Config.Canal.Charset
	config.Flavor = Config.Canal.Flavor
	luaRulesIncludeTableRegex()
	config.IncludeTableRegex = Config.Canal.IncludeTableRegex
	config.ExcludeTableRegex = Config.Canal.ExcludeTableRegex
	config.Dump.ExecutionPath = ""
	config.Logger = log.Logger
	_canal, err := canal.NewCanal(config)
	if err != nil {
		log.Panic(errors.Trace(err).Error())
		os.Exit(1)
	}
	Canal = &TCanal{
		Canal:   _canal,
		Config:  Config.Canal,
		Handler: canal_handler.EventHandler,
	}
	Canal.SetEventHandler(Canal.Handler)
	return Canal
}

func luaRulesIncludeTableRegex() {
	includeTableRegex := ""
	for _, rule := range lua_rule.Rules {
		if rule.Schema == "*" && rule.Table == "*" {
			continue
		}
		includeTableRegex = fmt.Sprintf("%s\\.%s", rule.Schema, rule.Table)
		Config.Canal.IncludeTableRegex = append(Config.Canal.IncludeTableRegex, includeTableRegex)
	}
}
