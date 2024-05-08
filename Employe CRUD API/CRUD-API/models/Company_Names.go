package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Company struct {
	Id             primitive.ObjectID `json:"id,omitempty"`
	CompanyName    string             `json:"C_name" bson:"C_name"`
	CompanyEmail   string             `json:"C_email" bson:"C_email"`
	CompanyAddress string             `json:"C_address" bson:"C_address"`
	CompanyNumber  string             `json:"C_Number" bson:"C_Number"`
	EmployeeId     primitive.ObjectID `json:"Emp_Id" bson:"Emp_Id"`
}
