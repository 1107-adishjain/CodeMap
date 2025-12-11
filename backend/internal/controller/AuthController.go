package controller

import (
	helper "codemap/backend/internal/helper"
	"database/sql"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func SignUp(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Email    string `json:"email" validate:"required,email"`
			Password string `json:"password" validate:"required,min=8"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		email := req.Email
		if err := validator.New().Struct(req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		if !helper.IsValidEmail(email) {
			http.Error(w, "Email must be all lowercase, contain no spaces, and be a valid Gmail address", http.StatusBadRequest)
			return
		}
		// here we will write the logic to check if user already exists in DB
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email=$1", email).Scan(&count)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		if count > 0 {
			http.Error(w, "User already exists", http.StatusConflict)
			return
		}
		// after the user has passed the email check and does not exist in DB, we will hash the password
		hashpassword, err := helper.HashPassword(req.Password)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}

		// here we will write the logic to save the user to DB with email and hashed password
		_, err = db.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", email, hashpassword)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"message":"User signed up successfully"}`))
	})
}

func Login(db *sql.DB) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle login logic
		var req struct {
			Email    string `json:"email" validate:"required,email"`
			Password string `json:"password" validate:"required,min=8"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		if err := validator.New().Struct(req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		// Here we check whether the user exsist on DB or not? if not then retrun an error that sign up first
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email=$1", req.Email).Scan(&count)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		if count == 0 {
			http.Error(w, "User does not exist. Please sign up first.", http.StatusUnauthorized)
			return
		}

		// Here we will verify the password with the hashed password stored in DB
		var storedHashedPassword string
		err = db.QueryRow("SELECT password FROM users WHERE email=$1", req.Email).Scan(&storedHashedPassword)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		if err := helper.VerifyPassword(storedHashedPassword, req.Password); err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// after we get to know that the user is logged in successfully, we will create a JWT token and return it to the user with its user info/claims
		// get the userid from DB to be used in JWT claims
		var userID string
		err = db.QueryRow("SELECT id FROM users WHERE email=$1", req.Email).Scan(&userID)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		access_token, refresh_token, err := helper.GenerateJWT(req.Email, userID) //here we will also pass the user id
		if err != nil {
			http.Error(w, "Failed to generate tokens", http.StatusInternalServerError)
			return
		}
		// now we will pass the acess token in the response and set the refresh token in http only cookie
		// Set cookie before writing body
http.SetCookie(w, &http.Cookie{
    Name:     "refresh_token",
    Value:    refresh_token,
    HttpOnly: true,
    Secure:   true,
    Path:     "/",
    MaxAge:   7 * 24 * 60 * 60,
})

// Now send JSON
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusOK)
json.NewEncoder(w).Encode(map[string]string{
    "user_id": userID,
    "access_token": access_token,
})

	})
}
