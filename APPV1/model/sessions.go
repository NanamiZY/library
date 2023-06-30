package model

import (
	"fmt"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

var sessionName = "session-name"
var store, _ = redis.NewStoreWithDB(10, "tcp", "localhost:6379", "", "1", []byte("咕噜咕噜"))

func GetSession(c *gin.Context) map[interface{}]interface{} {
	id, _ := c.Cookie("id")
	str := fmt.Sprintf("%s%s", sessionName, id)
	session, _ := store.Get(c.Request, str)
	fmt.Printf("session:%+v\n", session.Values)
	return session.Values
}

func SetSession(c *gin.Context, name string, id int64) error {
	str := fmt.Sprintf("%s%s", sessionName, string(id))
	fmt.Printf(str)
	session, _ := store.Get(c.Request, sessionName)
	session.Values["name"] = name
	session.Values["id"] = id
	return session.Save(c.Request, c.Writer)
}

func FlushSession(c *gin.Context) error {
	id, _ := c.Cookie("id")
	str := fmt.Sprintf("%s%s", sessionName, id)
	session, _ := store.Get(c.Request, str)
	fmt.Printf("session : %+v\n", session.Values)
	session.Values["name"] = ""
	session.Values["id"] = 0
	return session.Save(c.Request, c.Writer)
}
