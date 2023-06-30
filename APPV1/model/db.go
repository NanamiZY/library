package model

import (
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var GlobalConn *gorm.DB
var mysqlLogger logger.Interface

func New() {
	mysqlLogger = logger.Default.LogMode(logger.Info)
	mysqlLogger.Info(context.Background(), "连接数据库···")
	//parseTime=True&loc=Local MySQL 默认时间是格林尼治时间，与我们差八小时，需要定位到我们当地时间
	my := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", "root", "886dsbzy", "127.0.0.1:3306", "library")
	conn, err := gorm.Open(mysql.Open(my), &gorm.Config{})
	if err != nil {
		fmt.Printf("err:%s\n", err)
		panic(err)
	}
	GlobalConn = conn
	GlobalConn.AutoMigrate(&Book{}, &Category{}, &User{}, &Librarian{}, &Record{}, &Message{})
	MysqlLogger()

}

func MysqlLogger() {
	GlobalConn = GlobalConn.Session(&gorm.Session{
		Logger: mysqlLogger,
	})
}

//create table `users`(
//	`id` bigint not null AUTO_INCREMENT,
//    `name` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL,
//	`pwd` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL,
//	PRIMARY KEY (`id`)
//	)
//create table `books`(
//`id` bigint not null AUTO_INCREMENT,
//`book_name` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL,
//`category_id` bigint DEFAULT NULL,
//`count` bigint DEFAULT NULL ,
//PRIMARY KEY (`id`)
//)
//create table `borrow`(
//	`id` bigint not null AUTO_INCREMENT,
//	`user_id` bigint DEFAULT NULL,
//	`book_id` bigint DEFAULT NULL,
//	`borrow_time` datetime DEFAULT NULL,
//	`return_status` int DEFAULT 0 COMMENT '0归还,1未归还',
//	PRIMARY KEY (`id`)
//	)
//create table `return`(
//	`id` bigint not null AUTO_INCREMENT,
//	`user_id` bigint DEFAULT NULL,
//	`book_id` bigint DEFAULT NULL,
//	`return_time` datetime DEFAULT NULL,
//	PRIMARY KEY (`id`)
//)
//create table `category`(
//	`category_id` bigint not null AUTO_INCREMENT,
//	`category_name` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL,
//	PRIMARY KEY (`category_id`)
//	)
//create table `admin`(
//	`id` bigint not null AUTO_INCREMENT,
//    `name` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL,
//	`pwd` varchar(50) COLLATE utf8mb4_bin DEFAULT NULL,
//	PRIMARY KEY (`id`)
//	)
