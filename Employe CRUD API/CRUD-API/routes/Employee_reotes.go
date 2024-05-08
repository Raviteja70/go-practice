package routes

import (
	"gin-mongo-api/controller"

	"github.com/gin-gonic/gin"
)

func EmployeeRoute(router *gin.Engine) {
	router.POST("/emp", controller.CreateEmployee())
	router.GET("/employees", controller.GetAllEmployees())
	router.GET("/emp/:Id", controller.GetAEmployee())
	router.PUT("/emp/:Id", controller.UpdateAEmployee())
	router.DELETE("/delete/:Id", controller.DeleteAEmployee())
}
