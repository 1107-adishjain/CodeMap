package helper

import (
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecretKey = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	Email string `json:"email"`
	UserID string`json:"id"`
	// add other fields like UserID if needed
	jwt.RegisteredClaims
}

func IsValidEmail(email string) bool {
	if strings.Contains(email, " ") {
		return false
	}
	if email != strings.ToLower(email) {
		return false
	}
	if !strings.Contains(email, "@") || !strings.HasSuffix(email, "@gmail.com") {
		return false
	}
	return true
}

func HashPassword(password string) (string, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(pass), err
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func GenerateJWT(email, userID string) (string, string, error) {
	accessExp := time.Now().Add(15 * time.Minute)
	refreshExp := time.Now().Add(7 * 24 * time.Hour)

	accessClaims := &Claims{
		Email: email,
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	refreshClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(refreshExp),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	access_token_jwt := jwt.NewWithClaims(jwt.SigningMethodHS256,accessClaims)
	refresh_token_jwt := jwt.NewWithClaims(jwt.SigningMethodHS256,refreshClaims)

	access_token , err := access_token_jwt.SignedString(jwtSecretKey)
	if err!=nil{
		return "", "", err
	}
	refresh_token, err := refresh_token_jwt.SignedString(jwtSecretKey)
	if err!=nil{
		return "", "", err
	}
	return access_token, refresh_token, nil
}


func VerifyJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, err
	}
	return claims, nil
}