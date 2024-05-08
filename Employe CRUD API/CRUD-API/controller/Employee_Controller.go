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
	"github.com/go-playground/validator/v10"
	JSON "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var EmployeeCollection *mongo.Collection = configs.GetCollection(*configs.DB, "peoples")

var validate = validator.New()

func CreateEmployee() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		EmpId := models.Employee{}
		defer cancel()

		if err := c.BindJSON(&EmpId); err != nil {
			c.JSON(http.StatusBadRequest, response.Response{
				Status:  http.StatusBadRequest,
				Message: "Error",
				Data:    map[string]interface{}{"data": err.Error()},
			})
			return
		}

		if ValidationErr := validate.Struct(&EmpId); ValidationErr != nil {
			c.JSON(http.StatusBadRequest, response.Response{
				Status:  http.StatusBadRequest,
				Message: "Error",
				Data:    map[string]interface{}{"data": ValidationErr.Error()},
			})
			return
		}

		newEmp := models.Employee{
			Id:       primitive.NewObjectID(),
			Name:     EmpId.Name,
			Email:    EmpId.Email,
			Mobile:   EmpId.Mobile,
			Company:  EmpId.Company,
			Location: EmpId.Location,
		}
		CreateEmployee, err := EmployeeCollection.InsertOne(ctx, newEmp)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.Response{
				Status:  http.StatusInternalServerError,
				Message: "Error",
				Data:    map[string]interface{}{"data": err.Error()},
			})
			return
		}
		log.Println("createBuyer", newEmp.Id)

		log.Println("Created Company : ", CreateEmployee)

		c.JSON(http.StatusCreated, response.Response{
			Status:  http.StatusCreated,
			Message: "Created Successfully",
			Data:    map[string]interface{}{"data": CreateEmployee},
		})
	}
}

func GetAllEmployees() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		employe := []models.Employee{}
		defer cancel()
		result, err := EmployeeCollection.Find(ctx, JSON.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.Response{
				Status:  http.StatusInternalServerError,
				Message: "Error",
				Data:    map[string]interface{}{"data": err.Error()},
			})
			return
		}

		defer result.Close(ctx)
		for result.Next(ctx) {
			var singleEmplo models.Employee
			if err := result.Decode(&singleEmplo); err != nil {
				c.JSON(http.StatusInternalServerError, response.Response{
					Status:  http.StatusInternalServerError,
					Message: "Error",
					Data:    map[string]interface{}{"data": err.Error()},
				})
				return
			}

			employe = append(employe, singleEmplo)
		}
		c.JSON(http.StatusOK, response.Response{
			Status:  http.StatusOK,
			Message: "Success",
			Data:    map[string]interface{}{"data": employe},
		})
	}
}

func GetAEmployee() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		emploId := c.Param("id")
		emplo := models.Employee{}
		defer cancel()
		ObjId, _ := primitive.ObjectIDFromHex(emploId)
		err := EmployeeCollection.FindOne(ctx, JSON.M{"empId": ObjId}).Decode(&emplo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.Response{
				Status:  http.StatusInternalServerError,
				Message: "Error",
				Data:    map[string]interface{}{"data": err.Error()},
			})
			return
		}
		c.JSON(http.StatusOK, response.Response{
			Status:  http.StatusOK,
			Message: "Success",
			Data:    map[string]interface{}{"data": emplo},
		})
	}
}

func UpdateAEmployee() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		empId := c.Param("id")
		var emplo models.Employee
		defer cancel()
		ObjId, _ := primitive.ObjectIDFromHex(empId)
		if err := c.BindJSON(&emplo); err != nil {
			c.JSON(http.StatusBadRequest, response.Response{
				Status:  http.StatusBadRequest,
				Message: "Error",
				Data:    map[string]interface{}{"data": err.Error()},
			})
			return
		}

		if validationErr := validate.Struct(&emplo); validationErr != nil {
			c.JSON(http.StatusBadRequest, response.Response{
				Status:  http.StatusBadRequest,
				Message: "Error",
				Data:    map[string]interface{}{"data": validationErr.Error()},
			})
			return
		}

		Update := JSON.M{
			"Name":     emplo.Name,
			"Email":    emplo.Email,
			"Mobile":   emplo.Mobile,
			"Company":  emplo.Company,
			"Location": emplo.Location,
		}
		result, err := EmployeeCollection.UpdateOne(ctx, JSON.M{"id": ObjId}, JSON.M{"$set": Update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.Response{
				Status:  http.StatusInternalServerError,
				Message: "Error",
				Data:    map[string]interface{}{"data": err.Error},
			})
			return
		}

		var updateEmp models.Employee
		if result.MatchedCount == 1 {
			err := EmployeeCollection.FindOne(ctx, JSON.M{"id": ObjId}).Decode(&updateEmp)
			if err != nil {
				c.JSON(http.StatusInternalServerError, response.Response{
					Status:  http.StatusInternalServerError,
					Message: "error",
					Data:    map[string]interface{}{"data": err.Error},
				})

				return
			}
		}
		c.JSON(http.StatusOK, response.Response{
			Status:  http.StatusOK,
			Message: "success",
			Data:    map[string]interface{}{"data": updateEmp},
		})
	}
}

func DeleteAEmployee() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
		empId := c.Param("buyerId")
		defer cancle()
		ObjId, _ := primitive.ObjectIDFromHex(empId)
		result, err := EmployeeCollection.DeleteOne(ctx, JSON.M{"id": ObjId})
		if err != nil {
			if err != nil {
				c.JSON(http.StatusInternalServerError, response.Response{
					Status:  http.StatusInternalServerError,
					Message: "error",
					Data:    map[string]interface{}{"data": err.Error},
				})
				return
			}
		}
		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound, response.Response{
				Status:  http.StatusNotFound,
				Message: "error",
				Data:    map[string]interface{}{"data": "Buyer with Specified ID Not Found"},
			})
			return
		}
		c.JSON(http.StatusOK, response.Response{
			Status:  http.StatusOK,
			Message: "error",
			Data:    map[string]interface{}{"data": "Buyer Succefully Deleted"},
		})

	}
}
