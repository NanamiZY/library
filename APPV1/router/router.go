package router

import (
	"books/APPV1/logic"
	"books/APPV1/model"
	_ "books/docs"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"time"
)

func New() *gin.Engine {
	r := gin.Default()
	userRouter(r)
	adminRouter(r)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//获取验证码
	r.GET("/getCode", logic.GetCode)
	//用户登录
	r.POST("/userLogin", logic.UserLogin)
	//用户注册
	r.POST("/users", logic.AddUser)
	//管理员登陆
	r.POST("/adminLogin", logic.AdminLogin)
	//游客可以浏览书籍和分类
	r.GET("/books", logic.SearchBook)
	r.GET("/categories", logic.SearchCategory)
	//获取某一分类下的所有图书
	r.GET("/books/:category_id", logic.GetCategoryBooks)
	//获取手机验证码
	r.GET("/getPhoneCode/:phone", logic.GetPhoneCode)
	//手机号登录
	r.POST("/phoneLogin", logic.PhoneLogin)

	go func() {
		// 定义一个定时任务，每隔10分钟执行一次
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				// 在这里调用 checkOverdueBooks 函数执行检查过期图书的操作
				checkOverdueBooks()
			}
		}
	}()
	// 每天0点定时执行一次检查即将过期图书的操作
	go func() {
		for {
			now := time.Now()
			// 计算距离下一个0点的时间
			next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
			duration := next.Sub(now)
			// 等待一段时间后执行检查过期图书的操作
			time.Sleep(duration)
			checkWillBooks()
		}
	}()
	return r
}

func checkOverdueBooks() {
	fmt.Println("正在检查过期图书...")
	overID := model.CheckRecord()
	model.BanUser(overID)
}
func checkWillBooks() {
	model.CheckWillRecord()
}
