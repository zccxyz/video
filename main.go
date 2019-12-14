package main

import (
	"github.com/gin-gonic/gin"
	"reptile/web"
)

func main() {
	app := gin.New()

	web.Routes(app)

	_ = app.Run(":3000")
}
