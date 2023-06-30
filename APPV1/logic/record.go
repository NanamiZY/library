package logic

import (
	"books/APPV1/model"
	"books/APPV1/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

// GetRecords godoc
//
//			@Summary		查看借阅记录信息
//			@Description	用户自己查看自己所有的借阅记录
//			@Tags			users
//			@Produce		json
//	     @Param id path int64 true "用户Id"
//		    @Param Authorization header string false "Bearer 用户令牌"
//			@response		200,400,404,500	{object}	tools.HttpCode{data=[]model.Record}
//			@Router			/user/users/{id}/records [GET]
func GetRecords(c *gin.Context) {
	idString := c.Param("id")
	id, _ := strconv.ParseInt(idString, 10, 64)
	if id > 0 {
		ret := model.GetAllRecord(id)
		c.JSON(200, tools.HttpCode{
			Code:    tools.OK,
			Message: "查找成功",
			Data:    ret,
		})
		return
	}
	c.JSON(404, tools.HttpCode{
		Code:    tools.DoErr,
		Message: "获取用户id失败",
		Data:    struct{}{},
	})
}

// GetStatusRecords godoc
//
//			@Summary		查看借或还书信息
//			@Description	用户自己查看自己的借或还记录
//			@Tags			users
//			@Produce		json
//		    @Param id path int64 true "用户Id"
//	     @Param status path int true "标记是否归还字段"
//		    @Param Authorization header string false "Bearer 用户令牌"
//			@response		200,400,500	{object}	tools.HttpCode{data=[]model.Record}
//			@Router			/user/users/{id}/records/{status} [GET]
func GetStatusRecords(c *gin.Context) {
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

// GetStatusRecordsByAdmin godoc
//
//			@Summary		查看借阅记录信息
//			@Description	管理员查看所有的借书或还书记录
//			@Tags			admin
//			@Produce		json
//	     @Param status path int true "标记是否归还字段"
//			@response		200,400,404,500	{object}	tools.HttpCode{data=[]model.Record}
//			@Router			/admin/records/{status} [GET]
func GetStatusRecordsByAdmin(c *gin.Context) {
	statusString := c.Param("status")
	status, _ := strconv.ParseInt(statusString, 10, 64)
	ret := model.GetAllStatusRecord(status)
	fmt.Printf("ret:%v", ret)
	c.JSON(200, tools.HttpCode{
		Code:    tools.OK,
		Message: "查找成功",
		Data:    ret,
	})
}
