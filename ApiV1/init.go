package ApiV1

/**
 * @Author: DavinciEvans
 * @Author: zhouningsiyuan@foxmail.com
 * @Date: 2020/11/6 11:31
 * @Desc:
 */

import (
	"encoding/json"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

// DB：数据库ORM实体
var (
	dsn       string
	DB        *gorm.DB
	Config    config
	SecretKey string
)

// 数据库设置，通过 config.json读取
type config struct {
	Host        string
	User        string
	Password    string
	Dbname      string
	Port        string
	Sslmode     string
	TimeZone    string
	Development bool
	NetPort     int
	SecretKey   string
	Forge       bool
}

// init.go 初始化 ApiV1 的数据库
func init() {
	// 读取配置文件
	file, err := os.OpenFile("./config.json", os.O_RDONLY, 0755)
	if err != nil {
		fmt.Println(err.Error())
		panic("Can't open config.json!")
		return
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Config)
	if err != nil {
		fmt.Println("decode config error.")
	} else {
		fmt.Println("read config info success.")
	}
	SecretKey = Config.SecretKey

	// 创建日志文件
	dbLog, err := os.Create("./logs/db.log")
	if err != nil {
		fmt.Println("Could not open log.")
		panic(err.Error())
	}

	newLogger := logger.New(
		log.New(dbLog, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: 5 * 1000000000, // 慢 SQL 阈值，设置为 5 s
			LogLevel:      logger.Info,    // Log level
			Colorful:      false,          // 禁用彩色打印
		},
	)

	// 载入数据库
	dsn = "host=" + Config.Host + " user=" + Config.User + " password=" + Config.Password + " dbname=" + Config.Dbname + " port=" + Config.Port + " sslmode=" + Config.Sslmode + " TimeZone=" + Config.TimeZone

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		fmt.Println(err.Error())
		panic("Fail to open Databases")
	}
	fmt.Println("Open Database success.")

	// 迁移数据库
	modelsInit(Config.Development)

	// 设置是否要初始化初始数据
	if Config.Forge {
		if !Config.Development {
			Config.Forge = false
		}
		newConfig, err := os.OpenFile("./config.json", os.O_WRONLY|os.O_TRUNC, 0755)
		if err != nil {
			fmt.Println(err.Error())
			panic("Can't open config.json!")
			return
		}
		defer newConfig.Close()
		encoder := json.NewEncoder(newConfig)
		err = encoder.Encode(&Config)
		if err != nil {
			fmt.Println(err.Error())
			panic(err.Error())
		}
		Forge()
		fmt.Println("Detect it's first run, Forge data success.")
	}
}
