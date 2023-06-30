package model

import "fmt"

func GetUser(userName, pwd string) *User {
	ret := &User{}
	sql := "select * from users where user_name=? and password=? limit 1"
	err := GlobalConn.Raw(sql, userName, pwd).Scan(ret).Error
	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
	}
	return ret
}
func GetAdmin(userName, pwd string) *Librarian {
	ret := &Librarian{}
	sql := "select * from `librarians` where user_name=? and password=? limit 1"
	err := GlobalConn.Raw(sql, userName, pwd).Scan(ret).Error
	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
	}
	return ret
}
func GetUserByPhone(phone string) *User {
	ret := &User{}
	sql := "select * from users where phone=? limit 1"
	err := GlobalConn.Raw(sql, phone).Scan(ret).Error
	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
	}
	return ret
}
