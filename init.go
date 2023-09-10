package main

import (
	"flag"
	"github.com/pingcap/errors"
	"go-canal/libs/canal"
	"go-canal/libs/log"
	"go-canal/libs/lua/lua_rule"
	"go-canal/libs/service"
	"os"
)

func init() {
	flag.StringVar(&configFile, "config", "config.yml", "application config file_util")
	flag.Parse()
	if err := InitConfigFile(configFile); err != nil {
		log.Panic(errors.Trace(err).Error())
		os.Exit(1)
	}
	lua_rule.Init()
	canal.Init()
	service.Init()
}
