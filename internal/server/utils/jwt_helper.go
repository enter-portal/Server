package utils

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// Initialize JWT secret, bcrypt salt rounds, and JWT expiration days from environment variables
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
var bcryptSaltRounds, _ = strconv.Atoi(os.Getenv("BCRYPT_SALT_ROUNDS"))
var jwtExpDays, _ = strconv.Atoi(os.Getenv("JWT_EXPIRY_DAYS"))

// Claims represents the structure of the JWT claims
type Claims struct {
	jwt.StandardClaims
	Sign string `json:"sign"`
}

// GenerateToken creates a JWT token for the given user ID
func GenerateToken(id uint) (string, error) {
	// TODO: Remove below line if using UUIDs in future
	userId := strconv.FormatUint(uint64(id), 10)
	expires := expiresIn(jwtExpDays)
	hash, err := getHash(userId, bcryptSaltRounds)
	if err != nil {
		return "", err
	}

	claims := &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expires.Unix(),
		},
		Sign: hash,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken checks if the provided JWT token is valid for the given user ID
func ValidateToken(tokenString string, userId string) (bool, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return false, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		if claims.ExpiresAt <= time.Now().Unix() {
			return false, fmt.Errorf("token expired")
		}
		if !checkHash(claims.Sign, userId) {
			return false, fmt.Errorf("Invalid Token")
		}
		return true, nil
	}

	return false, fmt.Errorf("Invalid Token")
}

// expiresIn calculates the expiration time for the token based on the given number of days
func expiresIn(numDays int) time.Time {
	return time.Now().AddDate(0, 0, numDays)
}

// getHash generates a bcrypt hash of the user ID using the specified number of salt rounds
func getHash(userId string, saltRounds int) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(userId), saltRounds)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// checkHash verifies if the provided hash matches the user ID
func checkHash(hash string, userId string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(userId))
	return err == nil
}
