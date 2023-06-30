package router

import (
	"books/APPV1/logic"
	"books/APPV1/model"
	"books/APPV1/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func userRouter(r *gin.Engine) {
	base := r.Group("/user")
	base.Use(userCheck())
	user := base.Group("/users")
	{
		//查看个人信息
		user.GET("/:id", logic.GetUser)
		//更新个人信息
		user.PUT("", logic.UpdateUser)
		//查看所有借阅记录
		user.GET("/:id/records", logic.GetRecords)
		//查看借或还的记录
		user.GET("/:id/records/:status", logic.GetStatusRecords)
		//借书
		user.POST("/records/:bookId", logic.BorrowBook)
		//还书
		user.PUT("/records/:bookId", logic.ReturnBook)
	}
	book := base.Group("/books")
	{
		//查询图书详细信息
		book.GET("/:id", logic.GetBook)
	}
	message := base.Group("/messages")
	{
		//显示一个用户的所有消息
		message.GET("/:id", logic.GetMessage)
		message.GET("/:id/count", logic.GetCount)
	}
	go func() {
		Cacheheating()
	}()
}
func userCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		fmt.Printf("auth:%+v\n", auth)
		data, err := tools.Token.VerifyToken(auth)
		if err != nil {
			c.AbortWithStatusJSON(401, tools.HttpCode{
				Code:    tools.NotLogin,
				Message: "验签失败",
			})
			return
		}
		fmt.Printf("data:%+v\n", data)
		if data.ID <= 0 || data.Name == "" {
			c.AbortWithStatusJSON(401, tools.HttpCode{
				Code:    tools.NotLogin,
				Message: "用户信息错误",
			})
			return
		}
		c.Set("userName", data.Name)
		c.Set("userId", data.ID)
		c.Next()
	}
}

func Cacheheating() {
	// 创建一个定时器，每 3秒触发一次
	ticker1 := time.NewTicker(3 * time.Second)
	defer ticker1.Stop()
	// 设置初始页码为 1
	id := 4108
	size := 100
	// 在goroutine中运行定时器
	for {
		select {
		case <-ticker1.C:
			// 调用 Loader 函数加载图书信息
			books := model.Loader(id, size)
			// 调用 Handler 函数处理图书信息
			model.Handler(id, books)
			// 修改 pageIndex 变量，准备加载下一页图书信息
			id += size
			if id == 4408 {
				id = 4108
			}
		}
	}
}
func tokenCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := c.Cookie("access_token")
		if err != nil {
			// 未登录或者 cookie 不存在
			// 处理未登录的情况，例如返回未授权的错误响应
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "未登录",
			})
			return
		}

		// 验证 token 是否有效，如果无效则需要重新登录
		data, err := tools.Token.VerifyToken(accessToken)
		if err != nil {
			// token 验证失败，处理重新登录的情况，例如返回未授权的错误响应
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "登录已过期，请重新登录",
			})
			return
		}

		// token 验证成功，将用户信息存储到上下文对象中，以供后续的请求处理函数使用
		c.Set("userId", data.ID)
		c.Set("userName", data.Name)

		// 处理登录成功的情况，执行后续的请求处理函数
		c.Next()
	}
}
