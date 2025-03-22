package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	dsn string
)

// DB 外部调用的单例
var DB *gorm.DB

func GetDBv2() (*gorm.DB, error) {
	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	sqlDB, err := db.DB()
	// 设置连接池最大空闲连接数
	sqlDB.SetMaxIdleConns(150)
	// 设置连接池最大打开连接数
	sqlDB.SetMaxOpenConns(250)
	// 设置连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(30 * 60)
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		return nil, err
	}
	return db, nil
}

func GetDB() (*gorm.DB, error) {
	// 连接数据库
	if DB != nil {
		return DB, nil
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	sqlDB, err := db.DB()
	// 设置连接池最大空闲连接数
	sqlDB.SetMaxIdleConns(150)
	// 设置连接池最大打开连接数
	sqlDB.SetMaxOpenConns(250)
	// 设置连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(30 * 60)
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		return nil, err
	}
	return db, nil
}

func init() {
	//config := conf.LoadConfig()
	//dsn = config.Database.URL
	dsn = "root:123456@tcp(117.72.72.114:13306)/eshop_cart?charset=utf8mb4&parseTime=True&loc=Local"
}

func Init() {
	var err error
	DB, err = GetDB()
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
	}
}
