package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Routes(app *gin.Engine) {
	app.Static("/static", "./static")
	app.LoadHTMLGlob("./view/html/*")
	app.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	app.GET("/look", func(c *gin.Context) {
		c.HTML(http.StatusOK, "look.html", nil)
	})
	app.GET("/search", func(c *gin.Context) {
		page := c.DefaultQuery("p", "1")
		word := c.DefaultQuery("word", "")
		list, total, err := Search(word, page, "")
		if err != nil {
			c.JSON(200, Msg{
				Code: 0,
				Msg:  "出错了",
				Data: "",
			})
		}
		c.JSON(200, Msg{
			Code: 0,
			Msg:  "ok",
			Data: gin.H{
				"list":  list,
				"total": total,
			},
		})
	})
}
