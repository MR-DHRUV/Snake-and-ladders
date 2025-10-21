package utils

import (
	"errors"
	"time"

	"github.com/MR-DHRUV/snake_and_ladders/config"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userId string) (string, error) {

	claims := jwt.MapClaims{
		"id":  userId,
		"exp": time.Now().Add(config.AuthJWTExpiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(config.RSAPrivateKey)
}

func VerifyJWT(tokenString string) (string, string, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure that the token is signed with RS256 algorithm
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return "", errors.New("unexpected signing method")
		}
		return config.RSAPublicKey, nil
	})
	if err != nil {
		return "", "", err
	}

	// Validate the token claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check if the token is expired
		expirationTime := int64(claims["exp"].(float64))
		if expirationTime < time.Now().Unix() {
			return "", "", errors.New("unauthorized: token expired")
		}

		userId := claims["id"].(string)

		// refresh the token
		refreshedToken, err := GenerateJWT(userId)
		if err != nil {
			return refreshedToken, userId, err
		}

		return tokenString, userId, nil
	}

	return "", "", errors.New("invalid token claims")
}
