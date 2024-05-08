package routes

import (
	"gin-mongo-api/controller"

	"github.com/gin-gonic/gin"
)

func CompanyRoutes(router *gin.Engine) {
	router.POST("/comp", controller.CreateCompanys())
}
