package logic

import (
	"books/APPV1/model"
	"books/APPV1/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type NewUser struct {
	Name        string `json:"name" form:"name"`
	OldPassword string `json:"oldPassword" form:"oldPassword" `
	NewPassword string `json:"newPassword" form:"newPassword" `
	Sex         string `json:"sex" form:"sex" `
	Phone       string `json:"phone" form:"phone" `
}

// SearchUser godoc
//
//		@Summary		搜索用户
//		@Description	若search为空,显示所有用户,否则根据search搜索用户
//		@Tags			admin
//		@Produce		json
//	    @Param search  query string false "查询条件"
//		@response		200,404,500	{object}	tools.HttpCode{data=[]model.Book}
//		@Router			/admin/users [GET]
func SearchUser(c *gin.Context) {
	search := c.Query("search")
	ret := model.GetUsers(search)
	if len(ret) > 0 {
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.OK,
			Message: "查找成功",
			Data:    ret,
		})
		return
	}
	c.JSON(404, tools.HttpCode{
		Code:    tools.NotFound,
		Message: "该用户不存在",
		Data:    nil,
	})
}

// GetUser godoc
//
//		@Summary		获取个人信息
//		@Description	用户自己获取个人的详细信息
//		@Tags			users
//		@Produce		json
//	 @CookieParam id string true "用户id"
//	 @Param Authorization header string false "Bearer 用户令牌"
//		@response		200,400,401,500	{object}	tools.HttpCode{data=model.User}
//		@Router			/user/users/{id} [GET]
func GetUser(c *gin.Context) {
	idString, _ := c.Cookie("id")
	id, _ := strconv.ParseInt(idString, 10, 64)
	ret := model.GetMyUser(id)
	ret.Password = " "
	if ret.Id > 0 {
		c.JSON(200, tools.HttpCode{
			Code:    tools.OK,
			Message: "查询成功",
			Data:    ret,
		})
		return
	}
	c.JSON(401, tools.HttpCode{
		Code:    tools.DoErr,
		Message: "用户信息错误",
		Data:    struct{}{},
	})
	return
}

// AddUser godoc
//
//	@Summary		用户注册
//	@Description	用户注册
//	@Tags			users
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			userName		formData		string		true	"用户名"
//	@Param			password		formData		string	    true	"密码"
//	@Param			name		formData		string	    true	"昵称"
//	@Param			sex		formData		string	    true	"性别"
//	@Param			phone		formData		string	    true	"电话号码"
//	@response		200,400,500	{object}	tools.HttpCode
//	@Router			/users [POST]
func AddUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(400, tools.HttpCode{
			Code:    tools.DoErr,
			Message: "绑定失败",
			Data:    struct{}{},
		})
		return
	}
	fmt.Printf("user:%v", user)
	err := model.AddUser(&user)
	if err != nil {
		c.JSON(500, tools.HttpCode{
			Code:    tools.DoErr,
			Message: "注册失败",
			Data:    struct{}{},
		})
		return
	}
	c.JSON(200, tools.HttpCode{
		Code:    tools.OK,
		Message: "注册成功",
		Data:    struct{}{},
	})
}

// AddUserByAdmin godoc
//
//	@Summary		添加用户
//	@Description	管理员添加用户
//	@Tags			admin
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			userName		formData		string		true	"用户名"
//	@Param			password		formData		string	    true	"密码"
//	@Param			name		formData		string	    true	"昵称"
//	@Param			sex		formData		string	    true	"性别"
//	@Param			phone		formData		string	    true	"电话号码"
//	@response		200,400,500	{object}	tools.HttpCode
//	@Router			/admin/users [POST]
func AddUserByAdmin(c *gin.Context) {
	var user model.User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(400, tools.HttpCode{
			Code:    tools.DoErr,
			Message: "绑定失败",
			Data:    struct{}{},
		})
		return
	}
	fmt.Printf("user:%v", user)
	err := model.AddUser(&user)
	if err != nil {
		c.JSON(500, tools.HttpCode{
			Code:    tools.DoErr,
			Message: "添加失败",
			Data:    struct{}{},
		})
		return
	}
	c.JSON(200, tools.HttpCode{
		Code:    tools.OK,
		Message: "添加成功",
		Data:    struct{}{},
	})
}

// UpdateUser godoc
//
//		@Summary		修改个人信息
//		@Description	用户修改自己的个人信息
//		@Tags			users
//	    @Accept		    multipart/form-data
//		@Produce		json
//	    @Param          Authorization header string false "Bearer 用户令牌"
//		@Param			name		formData		string		true	"昵称"
//		@Param			oldPassword		formData		string	    true	"旧密码"
//		@Param			newPassword		formData		string	    true	"新密码"
//		@Param			sex		formData		string	    true	"性别"
//		@Param			phone		formData		string	    true	"电话号码"
//		@response		200,400,500	{object}	tools.HttpCode
//		@Router			/user/users [PUT]
func UpdateUser(c *gin.Context) {
	user := &NewUser{}
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(400, tools.HttpCode{
			Code:    tools.UserInfoErr,
			Message: "绑定失败" + err.Error(),
			Data:    struct{}{},
		})
		return
	}
	idString, _ := c.Cookie("id")
	id, _ := strconv.ParseInt(idString, 10, 64)
	ret := model.GetMyUser(id)
	if user.OldPassword == ret.Password {
		ret.Password = user.NewPassword
		ret.Name = user.Name
		ret.Phone = user.Phone
		ret.Sex = user.Sex
		err := model.UpdateUser(ret)
		if err != nil {
			c.JSON(500, tools.HttpCode{
				Code:    tools.UserInfoErr,
				Message: "更新失败" + err.Error(),
				Data:    struct{}{},
			})
			return
		}
		c.JSON(200, tools.HttpCode{
			Code:    tools.OK,
			Message: "更新成功",
			Data:    struct{}{},
		})
		return
	}
	c.JSON(200, tools.HttpCode{
		Code:    tools.UserInfoErr,
		Message: "旧密码错误",
		Data:    struct{}{},
	})
}

// DeleteUser godoc
//
//			@Summary		删除用户信息
//			@Description	管理员删除用户信息
//			@Tags			admin
//			@Produce		json
//	     @Param id path int64 true "用户id"
//			@response		200,400,500	{object}	tools.HttpCode
//			@Router			/admin/users/{id} [DELETE]
func DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id <= 0 {
		c.JSON(404, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "获取用户失败",
			Data:    struct{}{},
		})
		return
	}
	if err := model.DeleteUser(id); err != nil {
		c.JSON(500, tools.HttpCode{
			Code:    tools.DoErr,
			Message: "删除失败" + err.Error(),
			Data:    struct{}{},
		})
		return
	}
	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.OK,
		Message: "删除成功",
		Data:    struct{}{},
	})
}

// UpdateUserByAdmin godoc
//
//			@Summary		管理员修改用户信息
//			@Description	管理员在后台修改用户的信息
//			@Tags			admin
//		    @Accept		    multipart/form-data
//			@Produce		json
//	     @Param id path int64 true "用户id"
//			@Param			userName		formData		string	    true	"用户名"
//			@Param			password		formData		string	    true	"密码"
//			@Param			name		formData		string		true	"昵称"
//			@Param			sex		formData		string	    true	"性别"
//			@Param			phone		formData		string	    true	"电话号码"
//			@Param			status		formData		int	    true	"是否封禁 0未封禁 1封禁"
//			@response		200,400,500	{object}	tools.HttpCode
//			@Router			/admin/users/{id} [PUT]
func UpdateUserByAdmin(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id <= 0 {
		c.JSON(404, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "获取用户信息失败",
			Data:    struct{}{},
		})
		return
	}
	user := &model.User{}
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(400, tools.HttpCode{
			Code:    tools.UserInfoErr,
			Message: "绑定失败" + err.Error(),
			Data:    struct{}{},
		})
		return
	}
	user.Id = id
	err := model.UpdateUser(user)
	if err != nil {
		c.JSON(500, tools.HttpCode{
			Code:    tools.UserInfoErr,
			Message: "更新失败" + err.Error(),
			Data:    struct{}{},
		})
		return
	}
	c.JSON(200, tools.HttpCode{
		Code:    tools.OK,
		Message: "更新成功",
		Data:    struct{}{},
	})
	return
}

// GetUserBook godoc
//
//			@Summary		获取用户借阅记录
//			@Description	管理员获取用户的借书或还书记录
//			@Tags			admin
//			@Produce		json
//		    @Param id path int64 true "用户Id"
//	     @Param status path int true "标记是否归还字段 0未归还 1归还"
//			@response		200,400,500	{object}	tools.HttpCode{data=model.Record}
//			@Router			/admin/users/{id}/records/{status} [GET]
func GetUserBook(c *gin.Context) {
	idString := c.Param("id")
	statusString := c.Param("status")
	id, _ := strconv.ParseInt(idString, 10, 64)
	status, _ := strconv.ParseInt(statusString, 10, 64)
	ret := model.GetStatusRecord(id, status)
	c.JSON(200, tools.HttpCode{
		Code:    tools.OK,
		Message: "查找成功",
		Data:    ret,
	})
}

// GetMessage godoc
//
//			@Summary		查看消息
//			@Description	用户自己查看自己的消息
//			@Tags			users
//			@Produce		json
//	        @Param id path int64 true "用户Id"
//		    @Param Authorization header string false "Bearer 用户令牌"
//			@response		200,400,404,500	{object}	tools.HttpCode{data=[]model.Message}
//			@Router			/user/messages/{id} [GET]
func GetMessage(c *gin.Context) {
	idString := c.Param("id")
	id, _ := strconv.ParseInt(idString, 10, 64)
	if id > 0 {
		message := model.GetMessage(id)
		c.JSON(200, tools.HttpCode{
			Code:    tools.OK,
			Message: "查找成功",
			Data:    message,
		})
		return
	}
	c.JSON(404, tools.HttpCode{
		Code:    tools.DoErr,
		Message: "获取用户id失败",
		Data:    struct{}{},
	})
}

// GetCount godoc
//
//			@Summary		获取未读消息的数量
//			@Description	用户获取未读的消息的数量
//			@Tags			users
//			@Produce		json
//	        @Param id path int64 true "用户Id"
//		    @Param Authorization header string false "Bearer 用户令牌"
//			@response		200,400,404,500	{object}	tools.HttpCode{data=int}
//			@Router			/user/messages/{id}/count [GET]
func GetCount(c *gin.Context) {
	idString := c.Param("id")
	id, _ := strconv.ParseInt(idString, 10, 64)
	if id > 0 {
		count := model.GetCount(id)
		c.JSON(200, tools.HttpCode{
			Code:    tools.OK,
			Message: "查找成功",
			Data:    count,
		})
		return
	}
	c.JSON(404, tools.HttpCode{
		Code:    tools.DoErr,
		Message: "获取用户id失败",
		Data:    struct{}{},
	})
}
