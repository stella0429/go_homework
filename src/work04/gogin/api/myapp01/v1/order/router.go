package order

import "github.com/gin-gonic/gin"

func Routers(e *gin.Engine) {
	e.GET("/v1/order/addcart", addCartApi)
	e.GET("/v1/order/order", orderApi)
}
