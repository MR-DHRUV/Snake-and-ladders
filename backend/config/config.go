package config

import (
	"crypto/rsa"
	"encoding/pem"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func getRSAPrivateKey() *rsa.PrivateKey {
	privateKeyBytes, err := os.ReadFile("keys/private.pem")
	if err != nil {
		return nil
	}

	block, _ := pem.Decode(privateKeyBytes)
	if block == nil {
		return nil
	}

	privateKey, _ := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	return privateKey
}

func getRSAPublicKey() *rsa.PublicKey {
	publicKeyBytes, err := os.ReadFile("keys/public.pem")
	if err != nil {
		return nil
	}

	block, _ := pem.Decode(publicKeyBytes)
	if block == nil {
		return nil
	}

	publicKey, _ := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	return publicKey
}

var GS = os.Getenv("GS")

var GoogleClientID = os.Getenv("GoogleClientID")
var GoogleClientSecret = os.Getenv("GoogleClientSecret")
var GoogleRedirectURL = os.Getenv("GoogleRedirectURL")

var MongoURI = os.Getenv("MongoURI")
var ConnectionTimeout = 20 * time.Second

var WebSocketPort = 9999

var RSAPrivateKey = getRSAPrivateKey()
var RSAPublicKey = getRSAPublicKey()

var MinPlayers = 2
var MaxPlayers = 10
var MaxBoardRetries = 30
var MaxEntriedPerPage = 80
var BoardCells = 100

var AuthCookieExpiry = 7 * 24 * 60 * 60 // 1 week in seconds
var AuthJWTExpiry = time.Hour * 24 * 7   // 1 week
var AuthCookieSecure = false