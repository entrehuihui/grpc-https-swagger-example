package databases

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
)

// DB 连接DB
var db *gorm.DB

// GetDB 获取DB
func GetDB() *gorm.DB {
	return db
}

// InitMysql 初始化mysql连接
func InitMysql(user, password, host, dbname string) (err error) {
	dbparameter := user + ":" + password + "@tcp(" + host + ")/" + dbname + "?charset=utf8&parseTime=True&loc=Local"
loop:
	// 循环重连数据库
	db, err = gorm.Open("mysql", dbparameter)
	if err != nil {
		log.Println("mysql db isn't connect error :", err)
		time.Sleep(3 * time.Second)
		goto loop
	}

	// 开启打印日志
	db.LogMode(true)
	// 设置连接池为100
	db.DB().SetMaxIdleConns(1000)
	// 迁移数据表
	err = db.AutoMigrate(&User{}, &Grade{}, &Roles{}, &Application{}, &Appauthority{}, &Authority{}).Error
	if err != nil {
		log.Fatal(err)
	}
	return
}
