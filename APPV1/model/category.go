package model

import (
	"fmt"
)

func GetCategories() []*Category {
	ret := make([]*Category, 0)
	sql := "select * from `categories` where id>0"
	err := GlobalConn.Raw(sql).Scan(&ret).Error
	if err != nil {
		fmt.Printf("err:%v", err.Error())
	}
	return ret
}
func GetCategory(id int64) *Category {
	cate := &Category{}
	sql := "select * from categories where id=?"
	err := GlobalConn.Raw(sql, id).Scan(cate).Error
	if err != nil {
		fmt.Printf("err:%v", err.Error())
	}
	return cate
}
func AddCategory(cate *Category) error {
	sql := "insert into categories (name) values (?)"
	return GlobalConn.Exec(sql, cate.Name).Error
}
func UpdateCategory(cate *Category) error {
	sql := "update categories set name=? where id=?"
	return GlobalConn.Exec(sql, cate.Name, cate.Id).Error
}
func DeleteCategory(id int64) error {
	sql := "delete from categories where id=?"
	return GlobalConn.Exec(sql, id).Error
}
