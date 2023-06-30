package model

import "fmt"

func GetMessage(id int64) []*Message {
	ret := make([]*Message, 0)
	sql := "select * from messages where user_id=? ORDER BY create_time DESC"
	err := GlobalConn.Raw(sql, id).Scan(&ret).Error
	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
	}
	sql = "update messages set status=1 where user_id=?"
	err = GlobalConn.Exec(sql, id).Error
	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
	}
	return ret
}
func GetCount(id int64) int64 {
	var count int64
	sql := "select count(*) from messages where user_id=? and status=0"
	err := GlobalConn.Raw(sql, id).Scan(&count).Error
	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
	}
	return count
}
