package todo

import (
    //"github.com/QLeelulu/goku"
	"goku-master/goku-master"
    "path"
    "runtime"
    "time"
)

var (
    DATABASE_Driver string = "mymysql"
    // mysql: "user:password@/dbname?charset=utf8&keepalive=1"
    // mymysql: tcp:localhost:3306*dbname/user/pwd
	DATABASE_DSN string = "tcp:192.168.5.79:3306*jxcrm/test_db_online/QwerAs131126dfTyui1209"
)

var Config *goku.ServerConfig = &goku.ServerConfig{
    Addr:           ":8080",
    ReadTimeout:    10 * time.Second,
    WriteTimeout:   10 * time.Second,
    MaxHeaderBytes: 1 << 20,
    //RootDir:        _, filename, _, _ := runtime.Caller(1),
    StaticPath: "static", // static content 
    ViewPath:   "views",

    LogLevel: goku.LOG_LEVEL_LOG,
    Debug:    true,
}

func init() {
    // WTF, i just want to set the RootDir as current dir.
    _, filename, _, _ := runtime.Caller(1)
    Config.RootDir = path.Dir(filename)

    goku.SetGlobalViewData("SiteName", "Todo - by {goku}")
}
