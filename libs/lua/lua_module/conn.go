package lua_module

import (
	"fmt"
	"github.com/go-mysql-org/go-mysql/client"
	lua "github.com/yuin/gopher-lua"
	"go-canal/libs/log"
)

type tPool struct {
	addr     string
	user     string
	password string
	db       string
	*client.Pool
}

func (p *tPool) New() {
	p.Pool = client.NewPool(
		log.Debugf,
		60,
		600,
		5,
		p.addr,
		p.user,
		p.password,
		p.db,
	)
}

func lock() {
	locker <- true
}

func unlock() {
	<-locker
}

func getPoolParams(dsn *lua.LTable) (string, string, string, string, string) {
	addr := dsn.RawGetString("addr").String()
	user := dsn.RawGetString("user").String()
	if user == "nil" {
		user = ""
	}
	password := dsn.RawGetString("password").String()
	if password == "nil" {
		password = ""
	}
	db := dsn.RawGetString("db").String()
	if db == "nil" {
		db = ""
	}
	key := fmt.Sprintf("%s_%s", addr, user)
	return addr, user, password, db, key
}

func getPool(dsn *lua.LTable) *tPool {
	lock()
	defer unlock()

	addr, user, password, db, key := getPoolParams(dsn)

	pool, exist := pools[key]
	if !exist {
		pool = &tPool{
			addr:     addr,
			user:     user,
			password: password,
			db:       db,
		}
		pool.New()
		pools[key] = pool
	}

	return pool
}

var (
	locker = make(chan bool, 1)
	pools  = make(map[string]*tPool, 1)
)

func getConn(dsn *lua.LTable) (*client.Conn, error) {
	return getPool(dsn).GetConn(nil)
}

func putConn(dsn *lua.LTable, conn *client.Conn) {
	getPool(dsn).PutConn(conn)
}
