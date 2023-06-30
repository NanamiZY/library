package model

import (
	"fmt"
	"strconv"
	"time"
)

func GetAllRecord(id int64) []*Record {
	record := make([]*Record, 0)
	sql := "select * from `records` where user_id=?"
	err := GlobalConn.Raw(sql, id).Scan(&record).Error
	if err != nil {
		fmt.Printf("err:%v", err.Error())
	}
	return record
}
func GetStatusRecord(id, record int64) []*Record {
	ret := make([]*Record, 0)
	sql := "select * from `records` where user_id=? and status=?"
	err := GlobalConn.Raw(sql, id, record).Scan(&ret).Error
	if err != nil {
		fmt.Printf("err:%v", err.Error())
	}
	return ret
}
func GetAllStatusRecord(status int64) []*Record {
	record := make([]*Record, 0)
	sql := "select * from records where status=?"
	err := GlobalConn.Raw(sql, status).Scan(&record).Error
	if err != nil {
		fmt.Printf("err:%v", err.Error())
	}
	return record
}
func CheckRecord() []*int64 {
	ret := make([]*int64, 0)
	sql := "select distinct user_id from records where status=0 and over_time<=?  "
	err := GlobalConn.Raw(sql, time.Now()).Scan(&ret).Error
	if err != nil {
		fmt.Printf("err:%v", err.Error())
	}
	return ret
}
func CheckWillRecord() {
	ret := make([]*Record, 0)
	sql := "select * from records where status=0 and over_time>? and DATE_SUB(over_time, INTERVAL 1 DAY)<=?"
	err := GlobalConn.Raw(sql, time.Now(), time.Now()).Scan(&ret).Error
	if err != nil {
		fmt.Printf("err:%v", err.Error())
	}
	sql = "insert into messages (message,create_time,user_id,status) values (?,?,?,?)"
	if len(ret) > 0 {
		for i := 0; i < len(ret); i++ {
			if ret[i].Status == 1 {
				continue
			}
			GlobalConn.Exec(sql, "请尽快归还编号为 "+strconv.FormatInt(ret[i].BookId, 10)+" 的书", time.Now(), ret[i].UserId, 0)
		}
	}
}
