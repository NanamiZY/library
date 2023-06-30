package tools

import (
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"log"
	"os"
	"sync"
)

// 定义 dysmsapi 的客户端对象作为一个全局变量
var clientInstance *dysmsapi20170525.Client
var once sync.Once

// 定义一个函数，用来创建 dysmsapi20170525 的客户端对象，并赋值给全局变量
func createClient() {
	config := &openapi.Config{
		// 必填，您的 AccessKey ID
		AccessKeyId: tea.String(os.Getenv("SMS_ACCESS_KEY_ENV")),
		// 必填，您的 AccessKey Secret
		AccessKeySecret: tea.String(os.Getenv("SMS_ACCESS_KEY_SECRET_ENV")),
	}
	// 访问的域名
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	var err error
	clientInstance, err = dysmsapi20170525.NewClient(config)
	if err != nil {
		panic(err)
	}
}

func GetCode(phoneNumbers string, code string) (_err error, resp *dysmsapi20170525.SendSmsResponse) {
	// 使用 sync.Once 的 Do 方法，传入 createClient 函数
	once.Do(func() {
		createClient()
	})
	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		SignName:      tea.String("图书管理"),
		TemplateCode:  tea.String("SMS_461425151"),
		PhoneNumbers:  tea.String(phoneNumbers),
		TemplateParam: tea.String(fmt.Sprintf("{\"code\":\"%s\"}", code)),
	}
	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		// 复制代码运行请自行打印 API 的返回值
		resp, _err = clientInstance.SendSmsWithOptions(sendSmsRequest, runtime)
		if _err != nil {
			return _err
		}
		// 打印 API 的返回值
		log.Println(resp)

		return nil
	}()

	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.Message = tea.String(tryErr.Error())
		}
		// 打印 error
		log.Println(error)

		return error, nil
	}
	return _err, resp
}
