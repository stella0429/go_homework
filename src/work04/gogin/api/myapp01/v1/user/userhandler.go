package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func registerApi(c *gin.Context) {
	name := c.DefaultQuery("name", "jack")
	c.String(200, fmt.Sprintf("hello %s\n, register success!", name))
}

func loginApi(c *gin.Context) {
	name := c.DefaultQuery("name", "lily")
	c.String(200, fmt.Sprintf("hello %s\n, login success!", name))
}
