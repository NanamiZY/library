package APPV1

import (
	"books/APPV1/model"
	"books/APPV1/router"
	"books/APPV1/tools"
)

func Start() {
	model.New()
	tools.NewToken("呼呼呼")
	r := router.New()
	r.Run(":8080")
}
