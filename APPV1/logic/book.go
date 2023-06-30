package logic

import (
	"books/APPV1/model"
	"books/APPV1/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"net/http"
	"strconv"
	"time"
)

// GetBook godoc
//
//				@Summary		图书详细信息
//				@Description	获取一本书的详细信息
//				@Tags			book
//				@Produce		json
//			    @Param Authorization header string false "Bearer 用户令牌"
//	            @Param id path string true "书籍id"
//				@response		200,400,404,500	{object}	tools.HttpCode{data=model.Book}
//				@Router			/user/books/{id} [GET]
func GetBook(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	fmt.Printf("id:%v", id)
	if id <= 0 {
		c.JSON(404, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "无此书籍",
			Data:    struct{}{},
		})
		return
	}
	ret := model.GetBook(id)
	if ret.Id != 0 {
		c.JSON(200, tools.HttpCode{
			Code:    tools.OK,
			Message: "查找成功",
			Data:    ret,
		})
		return
	}
	c.JSON(404, tools.HttpCode{
		Code:    tools.NotFound,
		Message: "无此书籍",
		Data:    struct{}{},
	})
}

// BorrowBook godoc
//
//				@Summary		借书
//				@Description	用户自己借书
//				@Tags			book
//				@Produce		json
//			    @Param Authorization header string false "Bearer 用户令牌"
//			    @CookieParam id  string true "用户id"
//	         @Param bookId path string true "书籍id"
//				@response		200,400,500	{object}	tools.HttpCode
//				@Router			/user/users/records/{bookId} [POST]
func BorrowBook(c *gin.Context) {
	idString, _ := c.Cookie("id")
	bookIdString := c.Param("bookId")
	id, _ := strconv.ParseInt(idString, 10, 64)
	bookId, _ := strconv.ParseInt(bookIdString, 10, 64)
	url := "user/borrow/"
	userPath := fmt.Sprintf("%s%s%s", idString, url, bookIdString)
	var redisClient *redis.Client = model.RedisConn1
	exists, err := redisClient.Exists(userPath).Result()
	if err != nil {
		// 处理Redis错误
		fmt.Printf("Redis错误：%v\n", err)
	} else {
		if exists == 1 {
			count, _ := redisClient.Incr(userPath).Result()
			if count > 3 {
				c.JSON(200, tools.HttpCode{
					Code:    tools.DoErr,
					Message: "太快了太快了!",
					Data:    struct{}{},
				})
				return
			}
		} else {
			redisClient.Set(userPath, 1, 5*time.Second)
		}
	}
	user := model.GetMyUser(id)
	book := model.GetBook(bookId)
	if user.Status == 1 {
		c.JSON(200, tools.HttpCode{
			Code:    tools.DoErr,
			Message: "您已被封禁,请及时还书",
			Data:    struct{}{},
		})
		return
	}
	if book.Count <= 0 {
		c.JSON(200, tools.HttpCode{
			Code:    tools.DoErr,
			Message: "该书已被借完",
			Data:    struct{}{},
		})
		return
	}
	record := model.Record{
		UserId:    id,
		BookId:    bookId,
		Status:    0,
		StartTime: time.Now(),
		OverTime:  time.Now().Add(tools.T),
	}
	model.BorrowBook(&record)
	c.JSON(200, tools.HttpCode{
		Code:    tools.OK,
		Message: "借书成功",
		Data:    struct{}{},
	})
}

// ReturnBook godoc
//
//				@Summary		还书
//				@Description	用户自己还书
//				@Tags			book
//				@Produce		json
//			    @Param Authorization header string false "Bearer 用户令牌"
//			    @CookieParam id  string true "用户id"
//	            @Param bookId path string true "书籍id"
//				@response		200,400,500	{object}	tools.HttpCode
//				@Router			/user/users/records/{bookId} [PUT]
func ReturnBook(c *gin.Context) {
	idString, _ := c.Cookie("id")
	bookIdString := c.Param("bookId")
	id, _ := strconv.ParseInt(idString, 10, 64)
	bookId, _ := strconv.ParseInt(bookIdString, 10, 64)
	recordCount := model.Determine(bookId, id)
	if len(recordCount) <= 0 {
		c.JSON(404, tools.HttpCode{
			Code:    tools.DoErr,
			Message: "该书已还",
			Data:    struct{}{},
		})
		return
	}
	user := model.GetMyUser(id)
	record := model.Record{
		UserId:     id,
		BookId:     bookId,
		Status:     1,
		ReturnTime: time.Now(),
	}
	model.ReturnBook(&record)
	ret := model.GetStatusRecord(id, 0)
	count := 0
	now := time.Now().UTC()
	for i := 0; i < len(ret); i++ {
		if ret[i].OverTime.UTC().Before(now) {
			count++
		}
	}
	if user.Status == 1 && count == 0 {
		model.RestoreUser(user.Id)
	}
	c.JSON(200, tools.HttpCode{
		Code:    tools.OK,
		Message: "还书成功",
		Data:    struct{}{},
	})
}

// SearchBook godoc
//
//		@Summary		显示所有书籍
//		@Description	会将数据库中的所有书籍显示到页面
//		@Tags			book
//		@Produce		json
//	    @Param pageSize  query string false "页容量"
//	    @Param bookId  query string false "书籍ID"
//		@response		200,404,500	{object}	tools.HttpCode{data=[]model.Book}
//		@Router			/books [GET]
func SearchBook(c *gin.Context) {
	//改造成前端分页
	//gzip压缩http头部
	//gzip怎么实现的压缩
	//压缩完之后体积降了多少
	//redis->经典业务场景->大key
	//缓存的一致性  localCache->goCache尽可能的缩短缓冲时间
	//多级缓存
	//缓存的几个经典场景:击穿,穿透,雪崩
	//二八原则,只存第一页
	idstr := c.DefaultQuery("bookId", "4108")
	sizeStr := c.DefaultQuery("pageSize", "100")
	sizeInt, _ := strconv.Atoi(sizeStr)
	idInt, _ := strconv.Atoi(idstr)
	ret := model.GetBooks(sizeInt, idInt)
	if ret != nil {
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.OK,
			Message: "查找成功",
			Data:    ret,
		})
		return
	}
	c.JSON(500, tools.HttpCode{
		Code:    tools.NotFound,
		Message: "请输入正确的页码",
		Data:    nil,
	})
}

// AddBook godoc
//
//	@Summary		添加书籍
//	@Description	管理员添加书籍
//	@Tags			admin
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			bn		formData		string		true	"书籍编号"
//	@Param			name		formData		string	    true	"书籍名称"
//	@Param			description		formData		string	    true	"描述"
//	@Param			count		formData		int	    true	"数量"
//	@Param			categoryId		formData		int	    true	"书籍分类"
//	@response		200,400,500	{object}	tools.HttpCode
//	@Router			/admin/books [POST]
func AddBook(c *gin.Context) {
	book := model.Book{}
	if err := c.ShouldBind(&book); err != nil {
		c.JSON(400, tools.HttpCode{
			Code:    tools.DoErr,
			Message: "绑定失败",
			Data:    struct{}{},
		})
		return
	}
	fmt.Printf("book:%v", book)
	err := model.AddBook(&book)
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

// UpdateBook godoc
//
//			@Summary		修改图书信息
//			@Description	管理员在后台修改图书的信息
//			@Tags			admin
//		    @Accept		    multipart/form-data
//			@Produce		json
//	        @Param id path int64 true "图书id"
//			@Param			bn		formData		string	    true	"编号"
//			@Param			name		formData		string	    true	"书名"
//			@Param			description		formData		string		true	"简介"
//			@Param			count		formData		int	    true	"数量"
//			@Param			categoryId		formData		int	    true	"分类"
//			@response		200,400,500	{object}	tools.HttpCode
//			@Router			/admin/books/{id} [PUT]
func UpdateBook(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id <= 0 {
		c.JSON(404, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "获取书籍信息失败",
			Data:    struct{}{},
		})
		return
	}
	book := &model.Book{}
	if err := c.ShouldBind(&book); err != nil {
		c.JSON(400, tools.HttpCode{
			Code:    tools.UserInfoErr,
			Message: "绑定失败" + err.Error(),
			Data:    struct{}{},
		})
		return
	}
	book.Id = id
	err := model.UpdateBook(book)
	if err != nil {
		c.JSON(500, tools.HttpCode{
			Code:    tools.UserInfoErr,
			Message: "修改失败" + err.Error(),
			Data:    struct{}{},
		})
		return
	}
	c.JSON(200, tools.HttpCode{
		Code:    tools.OK,
		Message: "修改成功",
		Data:    struct{}{},
	})
	return
}

// DeleteBook godoc
//
//			@Summary		删除图书信息
//			@Description	管理员删除图书信息
//			@Tags			admin
//			@Produce		json
//	     @Param id path int64 true "图书id"
//			@response		200,400,500	{object}	tools.HttpCode
//			@Router			/admin/books/{id} [DELETE]
func DeleteBook(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id <= 0 {
		c.JSON(404, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "获取图书失败",
			Data:    struct{}{},
		})
		return
	}
	if err := model.DeleteBook(id); err != nil {
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

// GetCategoryBooks godoc
//
//	@Summary		获取某一分类下的所有图书
//	@Description	获取某一分类下的所有图书
//	@Tags			book
//	@Produce		json
//	@Param			category_id	path		int	false	"int valid"	minimum(1)
//	@response		200,404,500	{object}	tools.HttpCode{data=[]model.Book}
//	@Router			/books/{category_id} [get]
func GetCategoryBooks(c *gin.Context) {
	idStr := c.Param("category_id")
	fmt.Printf("idstr:%s\n", idStr)
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id < 0 {
		c.JSON(404, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "获取分类失败",
			Data:    struct{}{},
		})
		return
	}
	ret := model.GetCategoryBook(id)
	c.JSON(200, tools.HttpCode{
		Code:    tools.OK,
		Message: "查找成功",
		Data:    ret,
	})
}

// GetBookByAdmin godoc
//
//				@Summary		图书详细信息
//				@Description	获取一本书的详细信息
//				@Tags			admin
//				@Produce		json
//	            @Param id path string true "书籍id"
//				@response		200,400,404,500	{object}	tools.HttpCode{data=model.Book}
//				@Router			/admin/books/{id} [GET]
func GetBookByAdmin(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	fmt.Printf("id:%v", id)
	if id <= 0 {
		c.JSON(404, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "无此书籍",
			Data:    struct{}{},
		})
		return
	}
	ret := model.GetBook(id)
	if ret.Id != 0 {
		c.JSON(200, tools.HttpCode{
			Code:    tools.OK,
			Message: "查找成功",
			Data:    ret,
		})
		return
	}
	c.JSON(404, tools.HttpCode{
		Code:    tools.NotFound,
		Message: "无此书籍",
		Data:    struct{}{},
	})
}
