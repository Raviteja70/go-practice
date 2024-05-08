package controller

import (
	"context"
	"gin-mongo-api/configs"
	"gin-mongo-api/models"
	"gin-mongo-api/response"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var CompanyCollection *mongo.Collection = configs.GetCollection(*configs.DB, "companies")

// var validate = validator.New()

func CreateCompanys() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		CompId := models.Company{}
		defer cancel()

		if err := c.BindJSON(&CompId); err != nil {
			c.JSON(http.StatusBadRequest, response.Response{
				Status:  http.StatusBadRequest,
				Message: "Error",
				Data:    map[string]interface{}{"data": err.Error()},
			})
			return
		}

		if ValidationErr := validate.Struct(&CompId); ValidationErr != nil {
			c.JSON(http.StatusBadRequest, response.Response{
				Status:  http.StatusBadRequest,
				Message: "Error",
				Data:    map[string]interface{}{"data": ValidationErr.Error()},
			})
			return
		}

		Comp := models.Company{
			Id:             primitive.NewObjectID(),
			CompanyName:    CompId.CompanyName,
			CompanyEmail:   CompId.CompanyEmail,
			CompanyAddress: CompId.CompanyAddress,
			CompanyNumber:  CompId.CompanyNumber,
			EmployeeId:     CompId.EmployeeId,
		}
		CreateComp, err := CompanyCollection.InsertOne(ctx, Comp)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.Response{
				Status:  http.StatusInternalServerError,
				Message: "Error",
				Data:    map[string]interface{}{"data": err.Error()},
			})
			return
		}

		log.Println("Created Employee : ", CreateComp)
		c.JSON(http.StatusCreated, response.Response{
			Status:  http.StatusCreated,
			Message: "Created Successfully",
			Data:    map[string]interface{}{"data": CreateComp},
		})
	}
}



