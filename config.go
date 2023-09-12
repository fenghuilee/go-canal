package main

import (
	"github.com/pingcap/errors"
	"github.com/sirupsen/logrus"
	"go-canal/libs/canal"
	"go-canal/libs/log"
	"go-canal/libs/lua/lua_rule"
	"go-canal/utils/file_util"
	"go-canal/utils/sys_util"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

func InitConfigFile() error {
	if !strings.HasPrefix(configFile, "/") {
		configFile = filepath.Join(sys_util.CurrentDirectory(), configFile)
	}
	if !file_util.IsExist(configFile) {
		return errors.New("config file not found")
	}
	var data []byte
	var err error
	if data, err = os.ReadFile(configFile); err != nil {
		return errors.Trace(err)
	}

	if err = yaml.Unmarshal(data, log.Config); err != nil {
		return errors.Trace(err)
	}

	if err = yaml.Unmarshal(data, canal.Config); err != nil {
		return errors.Trace(err)
	}

	if err = yaml.Unmarshal(data, lua_rule.Config); err != nil {
		return errors.Trace(err)
	}

	InitConfig()

	return nil
}

func InitConfig() {
	if log.Config.Logger.Level == "" {
		log.Config.Logger.Level = "info"
	}
	log.SetLevel(log.Config.Logger.Level)
	log.Logger.Formatter.(*logrus.TextFormatter).ForceColors = isColor
	log.Infof("Canal (%s, %s)", canal.Config.Canal.Addr, canal.Config.Canal.User)
}
