package model

import (
	"fmt"
	"github.com/go-redis/redis"
	"math/rand"
	"time"
)

func NewClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis地址
		Password: "",               // Redis密码，如果没有则为空字符串
		DB:       2,                // Redis数据库编号，0表示默认的数据库
	})
}
func GetBooks(size int, id int) []*Book {
	clint := NewClient()
	key := "books:" + string(id)
	data, err := clint.Get(key).Bytes()
	if err == redis.Nil {
		// 如果Redis中没有缓存数据，则查询MySQL数据库，并将结果压缩存放到Redis中
		var books []*Book
		sql := "select * from books where id>=? limit ?" //分页查询
		GlobalConn.Raw(sql, id, size).Scan(&books)
		//压缩
		Ybyte := Compress(books)
		data = Ybyte
		//防止击穿这种方法可以用来避免大量的键同时过期，从而减轻 Redis 的负载压力。同时，随机过期时间也可以使键的过期时间更加均匀，
		//避免数据集中存储在某个时间段内过期，造成 Redis 的短时间内负载过高。
		num := rand.Intn(5)
		err = clint.Set(key, Ybyte, time.Duration(num)*time.Second).Err()
		if err != nil {
			return nil
		}
	}
	return Decompression(data)
	//sql := "select * from books where bn like ? or name like ?"
	//err := GlobalConn.Raw(sql, "%"+search+"%", "%"+search+"%").Scan(&ret).Error
	//if err != nil {
	//	fmt.Printf("err:%v", err.Error())
	//}
	//return ret
}
func GetBook(id int64) *Book {
	ret := &Book{}
	sql := "select * from `books` where id=?"
	err := GlobalConn.Raw(sql, id).Scan(ret).Error
	if err != nil {
		fmt.Printf("err:%v", err.Error())
	}
	return ret
}
func GetCategoryBook(id int64) []*Book {
	ret := make([]*Book, 0)
	sql := "select * from books where category_id=?"
	err := GlobalConn.Raw(sql, id).Scan(&ret).Error
	if err != nil {
		fmt.Printf("err:%v", err.Error())
	}
	return ret
}

// 消息队列可实现超卖
func BorrowBook(record *Record) {
	var err error
	tx := GlobalConn.Begin()
	sql := "select * from books where id=? for update"
	tx.Raw(sql, record.Id)
	sql = "update books set count=count-1 where id=?"
	err = tx.Exec(sql, record.BookId).Error
	if err != nil {
		fmt.Printf("err:%v", err.Error())
		tx.Rollback()
		return
	}
	sql = "insert into records (user_id,book_id,status,start_time,over_time) values (?,?,?,?,?)"
	err = tx.Exec(sql, record.UserId, record.BookId, record.Status, record.StartTime, record.OverTime).Error
	if err != nil {
		fmt.Printf("err:%v", err.Error())
		tx.Rollback()
		return
	}
	tx.Commit()
}
func ReturnBook(record *Record) {
	var err error
	tx := GlobalConn.Begin()
	sql := "update books set count=count+1 where id=?"
	err = tx.Exec(sql, record.BookId).Error
	if err != nil {
		fmt.Printf("err:%v", err.Error())
		tx.Rollback()
		return
	}
	sql = "update records set status=?,return_time=? where user_id=? and book_id=?"
	err = tx.Exec(sql, record.Status, record.ReturnTime, record.UserId, record.BookId).Error
	if err != nil {
		fmt.Printf("err:%v", err.Error())
		tx.Rollback()
		return
	}
	tx.Commit()
}

func AddBook(book *Book) error {
	sql := "insert into books (bn,name,description,count,category_id) values(?,?,?,?,?) "
	return GlobalConn.Exec(sql, book.BN, book.Name, book.Description, book.Count, book.CategoryId).Error
}
func UpdateBook(book *Book) error {
	sql := "update books set bn=?,name=?,description=?,count=?,category_id=? where id=?"
	return GlobalConn.Exec(sql, book.BN, book.Name, book.Description, book.Count, book.CategoryId, book.Id).Error
}
func DeleteBook(id int64) error {
	sql := "delete from books where id=?"
	return GlobalConn.Exec(sql, id).Error
}
func Determine(bookId int64, userId int64) []*Record {
	ret := make([]*Record, 0)
	sql := "select * from records where book_id=? and status=0 and user_id=?"
	GlobalConn.Raw(sql, bookId, userId).Scan(&ret)
	return ret
}
