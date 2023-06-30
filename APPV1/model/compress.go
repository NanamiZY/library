package model

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"
)

// 压缩数据
func Compress(books []*Book) []byte {
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	// 将 books 切片转换为 JSON 数据，并压缩
	jsonBytes, err := json.Marshal(&books)
	if err != nil {
		fmt.Printf("Error during  json.Marshal(books) :%+v\n", err.Error())
	}
	if _, err := gzipWriter.Write(jsonBytes); err != nil {
		fmt.Printf("Error during  gzipWriter.Write :%+v\n", err.Error())
	}
	if err := gzipWriter.Close(); err != nil {
		fmt.Printf("Error during gzipWriter.Close :%+v\n", err.Error())
	}
	return buf.Bytes()
}

// 解压缩
func Decompression(key []byte) []*Book {
	// 创建 gzip 解码器
	gzipReader, err := gzip.NewReader(bytes.NewReader(key))
	if err != nil {
		fmt.Printf("vgzip.NewReader 时出现错误！err:%+v\n", err.Error())
		return nil
	}
	defer func(gzipReader *gzip.Reader) {
		err := gzipReader.Close()
		if err != nil {
			fmt.Printf("gzipReader.Close() 时出现错误！err:%+v\n", err.Error())
			return
		}
	}(gzipReader)
	//
	books := make([]*Book, 0)
	// 分批读取解压缩后的 JSON 数据
	decoder := json.NewDecoder(gzipReader)
	batchSize := 100 // 每次读取的批次大小为 100 条记录
	for {
		var batch []*Book
		err := decoder.Decode(&batch)
		if err == io.EOF { // 已经读取完数据
			break
		} else if err != nil {
			fmt.Printf("ecoder.Decode(&batch) 时出现错误！err:%+v\n", err.Error())
			return nil
		}
		books = append(books, batch...)
		if len(books) >= batchSize { // 达到批次大小，返回结果
			break
		}
	}
	return books
}

// Loader 缓存预热
func Loader(id int, size int) []*Book {
	var books []*Book
	sql := "select * from books where id>=? limit ?" //分页查询
	GlobalConn.Raw(sql, id, size).Scan(&books)
	return books
}

// Handler 处理热点数据每3秒更新前三页数据
func Handler(id int, books []*Book) {
	clint := NewClient()
	v := Compress(books)
	id1 := strconv.Itoa(id)
	key := "books:" + id1
	err := clint.Set(key, v, 3*time.Second).Err()
	if err != nil {
		fmt.Printf("err", err)
	}
}
