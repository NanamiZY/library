package router

import (
	"books/APPV1/logic"
	"books/APPV1/model"
	"books/APPV1/tools"
	"github.com/gin-gonic/gin"
)

func adminRouter(r *gin.Engine) {
	base := r.Group("/admin")
	base.Use(AdminCheck())
	//管理员登出
	base.GET("/logout", logic.Logout)
	user := base.Group("/users")
	{
		//搜索某一用户
		user.GET("", logic.SearchUser)
		//修改用户信息
		user.PUT("/:id", logic.UpdateUserByAdmin)
		//添加用户信息
		user.POST("", logic.AddUserByAdmin)
		//删除用户信息
		user.DELETE("/:id", logic.DeleteUser)
		//获取用户已归还或者未归还的所有记录
		user.GET("/:id/records/:status", logic.GetUserBook)
	}
	//书的所有资源
	book := base.Group("/books")
	{
		//查询图书基本信息
		book.GET("/:id", logic.GetBookByAdmin)
		//添加图书
		book.POST("", logic.AddBook)
		//更新图书信息
		book.PUT("/:id", logic.UpdateBook)
		//删除图书
		book.DELETE("/:id", logic.DeleteBook)
	}
	category := base.Group("/categories")
	{
		//获取分类详细信息
		category.GET("/:id", logic.GetCategory)
		//添加分类
		category.POST("", logic.AddCategory)
		//更新分类
		category.PUT("/:id", logic.UpdateCategory)
		//删除分类
		category.DELETE("/:id", logic.DeleteCategory)
	}
	//记录表的资源
	record := base.Group("/records")
	{
		//所有归还或者未归还的记录
		record.GET("/:status", logic.GetStatusRecordsByAdmin)
	}
}
func AdminCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		data := model.GetSession(c)
		id, ok1 := data["id"]
		name, ok2 := data["name"]
		idInt64, _ := id.(int64)
		if !ok1 || !ok2 || idInt64 <= 0 || name == "" {
			c.AbortWithStatusJSON(401, tools.HttpCode{
				Code:    tools.NotLogin,
				Message: "验签失败",
			})
			return
		}
		c.Set("name", name)
		c.Set("userId", idInt64)
		c.Next()
	}
}
