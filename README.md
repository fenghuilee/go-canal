# go-canal
MySQL/MariaDB database binlog incremental subscription & consumer components Canal's golang client

# build
`go build`

# run
`go run`

# command
`
go-canal -config config.yml
`

# config.yml
```#日志相关配置
#logger相关配置
logger:
  level: info #日志级别；支持：trace|debug|info|warn|error|fatal|panic，默认info

#canal相关配置
canal:
  addr: 127.0.0.1:3306
  user: user
  password: password
  charset: utf8mb4
  flavor: mysql
  include_table_regex:
    - "test\\.*"
  exclude_table_regex:
    - "test_ignore\\.*"

#规则配置
rules:
  - lua_path: lua/test.lua # *.* 通用处理规则
  - schema: test           # test.* 精确库名，模糊表名处理规则
    lua_path: lua/test.lua
  - table: canal_test      # *.canal_test 模糊库名，精确表名处理规则
    lua_path: lua/test.lua
  - schema: test           # test.canal_test 精确表名匹配处理规则
    table: canal_test
    lua_path: lua/test.lua
```
