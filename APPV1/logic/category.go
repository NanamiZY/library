package logic

import (
	"books/APPV1/model"
	"books/APPV1/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// SearchCategory godoc
//
//	@Summary		显示所有分类
//	@Description	会将数据库中的所有分类显示到页面
//	@Tags			category
//	@Produce		json
//	@response		200,500	{object}	tools.HttpCode{data=[]model.Category}
//	@Router			/categories [GET]
func SearchCategory(c *gin.Context) {
	ret := model.GetCategories()
	c.JSON(200, tools.HttpCode{
		Code:    tools.OK,
		Message: "查找成功",
		Data:    ret,
	})
}

// GetCategory godoc
//
//				@Summary		分类详细信息
//				@Description	获取分类详细信息
//				@Tags			admin
//				@Produce		json
//	            @Param id path string true "分类id"
//				@response		200,400,404,500	{object}	tools.HttpCode{data=model.Category}
//				@Router			/admin/categories/{id} [GET]
func GetCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	fmt.Printf("id:%v", id)
	if id <= 0 {
		c.JSON(404, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "无此分类",
			Data:    struct{}{},
		})
		return
	}
	ret := model.GetCategory(id)
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
		Message: "无此分类",
		Data:    struct{}{},
	})
}

// AddCategory godoc
//
//	@Summary		添加分类
//	@Description	管理员添加分类
//	@Tags			admin
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			name		formData		string		true	"分类名称"
//	@response		200,400,500	{object}	tools.HttpCode
//	@Router			/admin/categories [POST]
func AddCategory(c *gin.Context) {
	cate := &model.Category{}
	if err := c.ShouldBind(cate); err != nil {
		c.JSON(400, tools.HttpCode{
			Code:    tools.DoErr,
			Message: "绑定失败",
			Data:    struct{}{},
		})
		return
	}
	err := model.AddCategory(cate)
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

// UpdateCategory godoc
//
//			@Summary		修改分类信息
//			@Description	管理员在后台修改分类的信息
//			@Tags			admin
//		    @Accept		    multipart/form-data
//			@Produce		json
//	        @Param id path int64 true "分类id"
//			@Param			name		formData		string	    true	"分类名称"
//			@response		200,400,500	{object}	tools.HttpCode
//			@Router			/admin/categories/{id} [PUT]
func UpdateCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id <= 0 {
		c.JSON(404, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "获取分类信息失败",
			Data:    struct{}{},
		})
		return
	}
	cate := &model.Category{}
	if err := c.ShouldBind(cate); err != nil {
		c.JSON(400, tools.HttpCode{
			Code:    tools.UserInfoErr,
			Message: "绑定失败" + err.Error(),
			Data:    struct{}{},
		})
		return
	}
	cate.Id = id
	err := model.UpdateCategory(cate)
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

// DeleteCategory godoc
//
//			@Summary		删除分类信息
//			@Description	管理员删除图书信息
//			@Tags			admin
//			@Produce		json
//	     @Param id path int64 true "分类id"
//			@response		200,400,500	{object}	tools.HttpCode
//			@Router			/admin/categories/{id} [DELETE]
func DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id <= 0 {
		c.JSON(404, tools.HttpCode{
			Code:    tools.NotFound,
			Message: "获取分类失败",
			Data:    struct{}{},
		})
		return
	}
	if err := model.DeleteCategory(id); err != nil {
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
