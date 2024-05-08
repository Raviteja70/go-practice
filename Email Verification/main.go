package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// User represents a user
type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token,omitempty"`
	Verified bool   `json:"verified"`
}

var client *mongo.Client
var userCollection *mongo.Collection

// SMTP configuration
var smtpServer = "smtp.gmail.com"
var smtpPort = "587"
var smtpUser = "suryak14919@gmail.com"
var smtpPassword = "cpooltnpqnpwfirg"
var fromEmail = "suryak14919@gmail.com"
var baseURL = "http://localhost:8080"

func init() {
	// Set up MongoDB connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	userCollection = client.Database("auth").Collection("users")
}

// RegisterHandler handles user registration
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate random verification token
	newUser.Token = generateToken()

	// Insert user into MongoDB
	_, err = userCollection.InsertOne(context.Background(), newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send verification email with link
	subject := "Verify your email address"
	body := fmt.Sprintf("Click the following link to verify your email address: %s/verify?email=%s&token=%s", baseURL, newUser.Email, newUser.Token)
	err = sendEmail(newUser.Email, subject, body)
	if err != nil {
		log.Printf("Failed to send verification email to %s: %v", newUser.Email, err)
		// If sending email fails, you may want to handle this gracefully
	}

	// Return success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

// VerifyHandler handles verification of user email
func VerifyHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	token := r.URL.Query().Get("token")

	// Find user by email and token
	var user User
	err := userCollection.FindOne(context.Background(), User{Email: email, Token: token}).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid verification link", http.StatusUnauthorized)
		return
	}

	// Mark user as verified
	user.Verified = true
	user.Token = "" // Clear token after verification

	// Update user in MongoDB
	_, err = userCollection.UpdateOne(context.Background(), User{Email: email}, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Email %s verified successfully", email)
}

func generateToken() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	token := make([]byte, 32)
	for i := range token {
		token[i] = charset[rand.Intn(len(charset))]
	}
	return string(token)
}

func sendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpServer)

	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body)

	err := smtp.SendMail(smtpServer+":"+smtpPort, auth, fromEmail, []string{to}, msg)
	return err
}

func main() {
	r := mux.NewRouter()

	// Register HTTP handlers
	r.HandleFunc("/register", RegisterHandler).Methods("POST")
	r.HandleFunc("/verify", VerifyHandler).Methods("GET")

	// Start server
	fmt.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
