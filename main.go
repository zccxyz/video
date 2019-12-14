package main

import (
	"github.com/gin-gonic/gin"
	"reptile/web"
	"time"
)

func main() {
	app := gin.New()

	web.Routes(app)

	err := app.Run(":3000")
	if err != nil {
		web.Echo("启动出错: " + err.Error())
		time.Sleep(time.Minute)
		return
	}
	web.Echo("请访问 http://localhost:3000")
}
