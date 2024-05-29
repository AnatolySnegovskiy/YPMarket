package system

import (
	"fmt"
	"time"
)
import "github.com/dgrijalva/jwt-go"

func CreateToken(userID uint) (string, error) {
	var jwtKey = []byte("your_secret_key")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, err := token.SignedString(jwtKey)
	return tokenString, err
}

func ValidateToken(signedToken string) (bool, error) {
	token, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("your_secret_key"), nil
	})

	if err != nil {
		return false, err
	}
	return token.Valid, nil
}

func GetUserID(signedToken string) (int, error) {
	token, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("your_secret_key"), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var userID int
		if id, ok := claims["user_id"].(float64); ok {
			userID = int(id)
		} else {
			return 0, fmt.Errorf("user_id claim is missing or not of type float64")
		}

		return userID, nil
	}

	return 0, fmt.Errorf("invalid token")
}
