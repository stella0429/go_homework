package order

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func addCartApi(c *gin.Context) {
	name := c.DefaultQuery("name", "jack")
	c.String(200, fmt.Sprintf("hello %s\n, add shopping cart success!", name))
}

func orderApi(c *gin.Context) {
	name := c.DefaultQuery("name", "lily")
	c.String(200, fmt.Sprintf("hello %s\n, order goods success!", name))
}
