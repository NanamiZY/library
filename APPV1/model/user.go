package model

import "fmt"

func AddUser(user *User) error {
	sql := "insert into `users`(user_name,password,name,sex,phone,status) values(?,?,?,?,?,?)"
	err := GlobalConn.Exec(sql, user.UserName, user.Password, user.Name, user.Sex, user.Phone, user.Status).Error
	return err
}
func GetMyUser(id int64) *User {
	user := &User{}
	sql := "select * from `users` where id=?"
	err := GlobalConn.Raw(sql, id).Scan(user).Error
	if err != nil {
		fmt.Printf("err:%v", err.Error())
	}
	return user
}
func UpdateUser(user *User) error {
	sql := "update `users` set password=?,name=?,sex=?,phone=?,status=? where id=?"
	err := GlobalConn.Exec(sql, user.Password, user.Name, user.Sex, user.Phone, user.Status, user.Id)
	return err.Error
}
func GetUsers(search string) []*User {
	ret := make([]*User, 0)
	if search == "" {
		sql := "select * from `users` where id>0"
		err := GlobalConn.Raw(sql).Scan(&ret).Error
		if err != nil {
			fmt.Printf("err:%v", err.Error())
		}
		return ret
	}
	sql := "select * from users where user_name like ? or name like ?"
	err := GlobalConn.Raw(sql, "%"+search+"%", "%"+search+"%").Scan(&ret).Error
	if err != nil {
		fmt.Printf("err:%v", err.Error())
	}
	return ret
}
func DeleteUser(id int64) error {
	sql := "delete from users where id=?"
	return GlobalConn.Exec(sql, id).Error
}
func BanUser(id []*int64) {
	sql := "update users set status=1 where status=0 and id=?"
	for i := 0; i < len(id); i++ {
		GlobalConn.Exec(sql, id[i])
	}
}
func RestoreUser(id int64) {
	sql := "update users set status=0 where id=?"
	GlobalConn.Exec(sql, id)
}
