package model

import (
	"time"
)

// 图书
type Book struct {
	Id          int64  `json:"id" form:"id" gorm:"primaryKey"`
	BN          string `json:"bn" form:"bn" gorm:"type:varchar(100);index"`
	Name        string `json:"name" form:"name" gorm:"type:varchar(100);uniqueIndex"`
	Description string `json:"description" form:"description" gorm:"type:varchar(15000)"`
	ImgUrl      string `json:"imgUrl" form:"imgUrl"`
	Count       int    `json:"count" form:"count"`
	CategoryId  uint64 `json:"categoryId" form:"categoryId"`
}

// 分类
type Category struct {
	Id   int64
	Name string `json:"name" form:"name" gorm:"type:varchar(100)"`
}

// 用户
type User struct {
	Id       int64  `json:"id" form:"id"`
	UserName string `json:"userName" form:"userName" gorm:"type:varchar(100)"`
	Password string `json:"password" form:"password" gorm:"type:varchar(100)"`
	Name     string `json:"name" form:"name" gorm:"type:varchar(100)"`
	Sex      string `json:"sex" form:"sex" gorm:"type:varchar(100)"`
	Phone    string `json:"phone" form:"phone" gorm:"type:varchar(100)"`
	Status   int    `json:"status" form:"status"` //`json:""默认正常0 封禁1
}

// 管理员
type Librarian struct {
	Id       int64  `json:"id" form:"id"`
	UserName string `json:"userName" form:"userName" gorm:"type:varchar(100)"`
	Password string `json:"password" form:"password" gorm:"type:varchar(100)"`
	Name     string `json:"name" form:"name" gorm:"type:varchar(100)"`
	Sex      string `json:"sex" form:"sex" gorm:"type:varchar(100)"`
	Phone    string `json:"phone" form:"phone" gorm:"type:varchar(100)"`
}

// 借还记录表
type Record struct {
	Id         int64
	UserId     int64
	BookId     int64
	Status     int `json:"status"` //已归还1 未归还0
	StartTime  time.Time
	OverTime   time.Time
	ReturnTime time.Time
}

type Message struct {
	Id         int64     `json:"id" form:"id"`
	UserId     int64     `json:"userId" form:"id"'`
	Message    string    `json:"message" form:"message"`
	Status     int       `json:"status" form:"status"` //0未读,1已读
	CreateTime time.Time `json:"createTime" form:"createTime"`
}
