package canal_handler

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
	"github.com/pingcap/errors"
	"go-canal/libs/canal/canal_entity"
	"go-canal/libs/log"
	"go-canal/libs/lua"
	"go-canal/libs/lua/lua_rule"
	"go-canal/utils/string_util"
	"go-canal/utils/sys_util"
	"os"
	"path/filepath"
	"time"
)

type TEventRow = []interface{}

type TEventHandler struct {
	*canal.DummyEventHandler
}

type THeaderMap = map[string]int

var EventHandler = new(TEventHandler)
var position = new(mysql.Position)
var queue = make(chan interface{}, 4096)
var stop = make(chan struct{}, 1)

func (h *TEventHandler) OnRow(e *canal.RowsEvent) error {
	log.Debugf("OnRow, action: %s, count: %d", e.Action, len(e.Rows))
	headerMap := GetHeaderMap()
	defer PutHeaderMap(headerMap)

	for index, column := range e.Table.Columns {
		headerMap[column.Name] = index
	}

	if e.Action == canal.UpdateAction {
		for i := 0; i < len(e.Rows); i++ {
			if (i+1)%2 == 0 {
				entity := canal_entity.GetRowEntity()
				entity.Action = e.Action
				entity.Schema = e.Table.Schema
				entity.Table = e.Table.Name
				entity.Timestamp = e.Header.Timestamp

				RowHeaderValueMap(entity.Row, headerMap, e.Rows[i])
				RowHeaderValueMap(entity.OldRow, headerMap, e.Rows[i-1])

				queue <- entity
			}
		}
	} else {
		for _, row := range e.Rows {
			entity := canal_entity.GetRowEntity()
			entity.Action = e.Action
			entity.Schema = e.Table.Schema
			entity.Table = e.Table.Name
			entity.Timestamp = e.Header.Timestamp

			RowHeaderValueMap(entity.Row, headerMap, row)

			queue <- entity
		}
	}
	return nil
}

func (h *TEventHandler) OnPosSynced(header *replication.EventHeader, pos mysql.Position, set mysql.GTIDSet, force bool) error {
	position = &pos
	log.Debugf("OnPosSynced, Position (Name: %s, Pos: %d) Synced", position.Name, position.Pos)
	return nil
}

func (h *TEventHandler) String() string {
	return "TEventHandler"
}

func (h *TEventHandler) Handle() {
	go func() {
		for {
			select {
			case q := <-queue:
				switch q.(type) {
				case *canal_entity.TPosEntity:
				case *canal_entity.TRowEntity:
					entity := q.(*canal_entity.TRowEntity)
					key1 := fmt.Sprintf("%s.%s", entity.Schema, entity.Table)
					key2 := fmt.Sprintf("%s.*", entity.Schema)
					key3 := "*.*"
					rule, exist := lua_rule.Rules[key1]
					if !exist {
						if rule, exist = lua_rule.Rules[key2]; !exist {
							if rule, exist = lua_rule.Rules[key3]; !exist {
								continue
							}
						}
					}
					if exist {
						if err := lua.DoScript(entity, rule); err != nil {
							log.Errorf("%v", err.Error())
						}
					}
					log.Debugf("Handle, %s: %s.%s ok", entity.Action, entity.Schema, entity.Table)
					canal_entity.PutRowEntity(entity)
				}
			case <-stop:
				h.stop()
				return
			}
		}
	}()
	go func() {
		for {
			time.Sleep(time.Minute)
			h.save(*position)
		}
	}()
}

func (h *TEventHandler) save(pos mysql.Position) {
	var positionJson []byte
	var err error
	if positionJson, err = json.Marshal(pos); err != nil {
		log.Warnf("Position json marshal failed: %s", errors.Trace(err).Error())
		return
	}
	ConfigFileMD5 := string_util.MD5(flag.Lookup("config").Value.String())
	positionFile := filepath.Join(sys_util.CurrentDirectory(), "runtime", ConfigFileMD5+".position.json")
	log.Infof("Position (Name: %s, Pos: %d) saving to %s", pos.Name, pos.Pos, positionFile)
	if err = os.WriteFile(positionFile, positionJson, os.ModePerm); err != nil {
		log.Warnf("Position save failed: %s", errors.Trace(err).Error())
		return
	}
	log.Infof("Position (Name: %s, Pos: %d) saved", pos.Name, pos.Pos)
}

func (h *TEventHandler) stop() {
	log.Warnf("CanalEventHandler stopping")
	h.save(*position)
	log.Warnf("CanalEventHandler stoped")
}

func (h *TEventHandler) Stop() {
	stop <- struct{}{}
}
