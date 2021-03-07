package dao

import "github.com/jinzhu/gorm"

var (
	// DB mysql db
	DB *gorm.DB
)

// InitMySQL 初始化数据库
func InitMySQL() (err error) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/db2?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		return err
	}
	return DB.DB().Ping()
}

// Close 关闭数据库
func Close() {
	DB.Close()
}
