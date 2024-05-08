package main

import (
	"fmt"
	"gin-mongo-api/configs"
	"gin-mongo-api/routes"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Main Started")
	router := gin.Default()
	time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	fmt.Println(time)
	configs.ConnectDB()
	routes.EmployeeRoute(router)

	router.Run("localhost:6000")
}
