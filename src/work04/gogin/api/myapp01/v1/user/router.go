package user

import "github.com/gin-gonic/gin"

func Routers(e *gin.Engine) {
	e.GET("/v1/user/register", registerApi)
	e.GET("/v1/user/login", loginApi)
}
