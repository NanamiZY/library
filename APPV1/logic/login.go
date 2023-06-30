package logic

import (
	"books/APPV1/model"
	"books/APPV1/tools"
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// GetCode godoc
//
// @Tags		public
// @Summary		登录验证码
// @Description	用户登录页获取验证码操作
// @Produce		json
// @Success 200 {object} tools.HttpCode{data=map[string]string}
// @Failure 500 {object} tools.HttpCode
// @Router			/getCode [GET]
func GetCode(c *gin.Context) {
	fileName := func() string {

		// 设置图片大小
		width, height := 100, 50
		img := image.NewRGBA(image.Rect(0, 0, width, height))

		// 随机种子
		rand.Seed(time.Now().Unix())

		// 随机生成4位验证码
		code := fmt.Sprintf("%04d", rand.Intn(10000))
		fmt.Println("验证码:", code)
		//验证码存到redis

		var redisClient *redis.Client = model.RedisConn
		err := redisClient.Set("captcha", code, 5*time.Minute).Err()
		if err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusOK, tools.HttpCode{
				Code:    tools.OK,
				Message: err.Error(),
				Data:    nil,
			})
			return ""
		}
		// 设置字体大小
		fontSize := 30

		// 设置字体颜色
		fontColor := color.RGBA{255, 0, 0, 255}

		// 设置背景颜色
		bgColor := color.RGBA{255, 255, 255, 255}

		// 绘制背景
		draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.ZP, draw.Src)

		// 绘制验证码
		for i, c := range code {
			// 计算字体位置
			x := (width / 4) * i
			y := (height - fontSize) / 2

			// 绘制字体
			func(img *image.RGBA, s string, x, y int, c color.Color, size int) {
				f := basicfont.Face7x13
				d := &font.Drawer{
					Dst:  img,
					Src:  image.NewUniform(c),
					Face: f,
					Dot:  fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)},
				}
				d.DrawString(s)
			}(img, string(c), x, y, fontColor, fontSize)

		}
		// 将图像写入文件
		file, err := os.Create("APPV1/resource/static/img/captcha.png")
		fileUrlArr := strings.Split(file.Name(), "/")
		fileName := fileUrlArr[len(fileUrlArr)-1]
		if err != nil {
			panic(err)
			c.JSON(http.StatusInternalServerError, tools.HttpCode{
				Code:    tools.InternalServerError,
				Message: err.Error(),
				Data:    nil,
			})
		}
		defer file.Close()
		png.Encode(file, img)
		return fileName
	}()
	c.JSON(http.StatusOK, tools.HttpCode{
		Code:    tools.OK,
		Message: "获取验证码",
		Data:    map[string]string{"imgName": fileName},
	})
	return
}

// UserLogin godoc
//
//		@Summary		用户登录
//		@Description	会执行用户登录操作
//		@Tags			login
//		@Accept			multipart/form-data
//		@Produce		json
//		@Param			userName	formData	string	true	"用户名"
//		@Param			password		formData	string	true	"密码"
//	 @Param captcha formData string true "验证码"
//		@response		200,400,401,500	{object}	tools.HttpCode{data=Token}
//		@Router			/userLogin [POST]
func UserLogin(c *gin.Context) {
	var user model.User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(400, tools.HttpCode{
			Code:    tools.UserInfoErr,
			Message: "绑定失败" + err.Error(),
			Data:    struct{}{},
		})
		return
	}
	fmt.Printf("user:%v", user)
	//校验验证码
	formCode := c.PostForm("captcha")
	var redisClient *redis.Client = model.RedisConn
	redisCode, _ := redisClient.Get("captcha").Result()
	if redisCode != formCode {
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.CaptchaError,
			Message: "验证码错误",
		})
		return
	}
	dbUser := model.GetUser(user.UserName, user.Password)
	if dbUser.Id <= 0 {
		c.JSON(401, tools.HttpCode{
			Code:    tools.UserInfoErr,
			Message: "用户信息错误",
			Data:    struct{}{},
		})
		return
	}
	c.SetCookie("id", strconv.FormatInt(dbUser.Id, 10), 3600, "/", "", false, true)
	a, r, err := tools.Token.GetToken(dbUser.Id, dbUser.UserName)
	fmt.Printf("aToken:%s\n", a)
	fmt.Printf("rToken%s\n", r)
	if err != nil {
		c.JSON(500, tools.HttpCode{
			Code:    tools.UserInfoErr,
			Message: "token生成失败,错误信息:" + err.Error(),
			Data:    struct{}{},
		})
		return
	}
	c.JSON(200, tools.HttpCode{
		Code:    tools.UserInfoErr,
		Message: "生成成功,正在跳转",
		Data: Token{
			AccessToken:  a,
			RefreshToken: r,
		},
	})
}

// AdminLogin godoc
//
//	@Summary		管理员登录
//	@Description	会执行管理员登录操作
//	@Tags			login
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			userName	formData	string	true	"用户名"
//	@Param			password		formData	string	true	"密码"
//	@response		200,400,401,500	{object}	tools.HttpCode
//	@Router			/adminLogin [POST]
func AdminLogin(c *gin.Context) {
	var librarian model.Librarian
	if err := c.ShouldBind(&librarian); err != nil {
		c.JSON(400, tools.HttpCode{
			Code:    tools.UserInfoErr,
			Message: "绑定失败" + err.Error(),
			Data:    struct{}{},
		})
		return
	}
	fmt.Printf("admin:%v", librarian)
	dbAdmin := model.GetAdmin(librarian.UserName, librarian.Password)
	if dbAdmin.Id <= 0 {
		c.JSON(401, tools.HttpCode{
			Code:    tools.UserInfoErr,
			Message: "用户信息错误",
			Data:    struct{}{},
		})
		return
	}
	err := model.SetSession(c, dbAdmin.UserName, dbAdmin.Id)
	if err != nil {
		c.JSON(500, tools.HttpCode{
			Code:    tools.UserInfoErr,
			Message: err.Error(),
		})
		return
	}
	c.JSON(200, tools.HttpCode{
		Code:    tools.OK,
		Message: "登录成功，正在跳转~",
		Data:    struct{}{},
	})
	return
}

// Logout godoc
//
//	@Summary		管理员退出
//	@Description	会执行管理员退出操作
//	@Tags			admin
//	@Accept			json
//	@Produce		json
//	@response		500,401	{object}	tools.HttpCode
//	@Router			/admin/logout [get]
func Logout(c *gin.Context) {
	_ = model.FlushSession(c)
	c.JSON(200, tools.HttpCode{
		Code: tools.OK,
		Data: struct{}{},
	})
	return
}

// GetPhoneCode godoc
//
//	@Summary		获取手机验证码
//	@Description	向用户手机发送验证码
//	@Tags			login
//	@Produce		json
//	@Param			phone	path		string	true	"手机号"
//	@response		200,404,500	{object}	tools.HttpCode
//	@Router			/getPhoneCode/{phone} [GET]
func GetPhoneCode(c *gin.Context) {
	var specialPhones = map[string]int{}
	file, err := os.Open("specialPhones.txt")
	if err != nil {
		fmt.Println("打开文件失败：", err)
		return
	}
	defer file.Close() // 关闭文件
	// 读取文件内容
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		specialPhones[line] = 0
	}
	if err = scanner.Err(); err != nil {
		fmt.Println("读取文件失败：", err)
	}
	phone := c.Param("phone")
	ip := c.ClientIP()
	fmt.Printf("ip地址为:%s", ip)
	var redisClient *redis.Client = model.RedisConn
	// 构造计数器的键名，格式为 count:{ip}

	if specialPhones[phone] != 0 {
		countKey := fmt.Sprintf("count:%s", ip)
		// 尝试将计数器的值加 1
		countId, err := redisClient.Incr(countKey).Result()
		if err != nil {
			c.JSON(http.StatusOK, tools.HttpCode{
				Code:    tools.OK,
				Message: err.Error(),
				Data:    nil,
			})
		}
		// 如果计数器的过期时间还没有设置，则设置为一定的时间
		if ttl, err := redisClient.TTL(countKey).Result(); err != nil || ttl < 0 {
			redisClient.Expire(countKey, time.Minute) // 设置过期时间为 1 分钟
		}

		// 如果发送次数超过了限制，则返回 false
		if countId > 2 {
			c.JSON(401, tools.HttpCode{
				Code:    tools.UserInfoErr,
				Message: "发送的太快,慢一点吧",
				Data:    struct{}{},
			})
			return
		}
		if matched, err := regexp.MatchString(`^1[3-9]\d{9}$`, phone); err != nil || !matched {
			c.JSON(401, tools.HttpCode{
				Code:    tools.UserInfoErr,
				Message: "请输入正确的手机号",
				Data:    struct{}{},
			})
			return
		}
		countStr := fmt.Sprintf("count:%s", phone)
		count, err := redisClient.Get(countStr).Int()
		if err != nil && err != redis.Nil {
			c.JSON(500, tools.HttpCode{
				Code:    tools.UserInfoErr,
				Message: err.Error(),
			})
			return
		}
		if count > 10 {
			c.JSON(401, tools.HttpCode{
				Code:    tools.UserInfoErr,
				Message: "该手机号今天发送短信次数已达上限",
				Data:    struct{}{},
			})
			return
		}
		err = redisClient.Incr(countStr).Err()
		if err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusOK, tools.HttpCode{
				Code:    tools.OK,
				Message: err.Error(),
				Data:    nil,
			})
			return
		}
		redisClient.Expire(countStr, 24*time.Hour)
	}
	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	err, resp := tools.GetCode(phone, code)
	if err != nil {
		c.JSON(500, tools.HttpCode{
			Code:    tools.UserInfoErr,
			Message: err.Error(),
		})
		return
	}
	if *resp.Body.Message != "OK" {
		c.JSON(500, tools.HttpCode{
			Code:    tools.UserInfoErr,
			Message: *resp.Body.Message,
		})
		return
	}
	fmt.Printf("resp:%v", resp)
	c.JSON(200, tools.HttpCode{
		Code:    tools.OK,
		Message: "已发送验证码",
		Data:    struct{}{},
	})
	//验证码存到redis
	err = redisClient.HSet(phone, code, 1).Err()
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.OK,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	redisClient.Expire(phone, time.Minute)
}

// PhoneLogin godoc
//
//			@Summary		手机号登录
//			@Description	输入手机验证码后进行登录
//			@Tags			login
//			@Accept			multipart/form-data
//			@Produce		json
//	     @Param			phone   	formData	string	true	"手机号"
//		    @Param          captcha     formData    string  false    "验证码"
//			@response		200,400,401,500	{object}	tools.HttpCode{data=Token}
//			@Router			/phoneLogin [POST]
func PhoneLogin(c *gin.Context) {
	formCode := c.PostForm("captcha")
	phone := c.PostForm("phone")
	if model.GetUserByPhone(phone).Id <= 0 {
		c.JSON(401, tools.HttpCode{
			Code:    tools.UserInfoErr,
			Message: "该手机号未注册",
			Data:    struct{}{},
		})
		return
	}
	var redisClient *redis.Client = model.RedisConn
	redisCode, _ := redisClient.HGet(phone, formCode).Result()
	if redisCode == "" {
		c.JSON(http.StatusOK, tools.HttpCode{
			Code:    tools.CaptchaError,
			Message: "验证码错误",
		})
		return
	}
	dbUser := model.GetUserByPhone(phone)
	if dbUser.Id <= 0 {
		c.JSON(401, tools.HttpCode{
			Code:    tools.UserInfoErr,
			Message: "用户信息错误",
			Data:    struct{}{},
		})
		return
	}
	c.SetCookie("id", strconv.FormatInt(dbUser.Id, 10), 3600, "/", "", false, true)
	a, r, err := tools.Token.GetToken(dbUser.Id, dbUser.UserName)
	fmt.Printf("aToken:%s\n", a)
	fmt.Printf("rToken%s\n", r)
	if err != nil {
		c.JSON(500, tools.HttpCode{
			Code:    tools.UserInfoErr,
			Message: "token生成失败,错误信息:" + err.Error(),
			Data:    struct{}{},
		})
		return
	}
	c.SetCookie("access_token", a, 3600, "/", "", false, true)
	c.SetCookie("refresh_token", r, 3600, "/", "", false, true)
	c.JSON(200, tools.HttpCode{
		Code:    tools.UserInfoErr,
		Message: "生成成功,正在跳转",
		Data: Token{
			AccessToken:  a,
			RefreshToken: r,
		},
	})
}
