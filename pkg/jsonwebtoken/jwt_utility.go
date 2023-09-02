package jsonwebtoken

import (
	"log"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pjmessi/go-database-usage/config"
)

type JwtUtility struct {
	secret         []byte
	expirationTime int
}

func CreateJwtUtility(appConfig *config.AppConfig) *JwtUtility {
	expirationTime, err := strconv.Atoi(appConfig.JWT_EXPIRATION_TIME)
	if err != nil {
		log.Println("could not parse JWT_EXPIRATION_TIME, using default value of 3600 seconds")
		expirationTime = 3600
	}

	return &JwtUtility{
		secret:         []byte(appConfig.JWT_SECRET),
		expirationTime: expirationTime,
	}
}

func (utility *JwtUtility) CreateJwt(userId string, userEmail string) (*string, error) {
	claims := jwt.MapClaims{
		"user_id":    userId,
		"user_email": userEmail,
		"exp":        time.Now().Add(time.Second * time.Duration(utility.expirationTime)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(utility.secret)
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}
