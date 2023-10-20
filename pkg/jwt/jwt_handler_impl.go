package jwt

import (
	"fmt"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/pjmessi/golang-practice/config"
	"github.com/pjmessi/golang-practice/pkg/logger"
	"github.com/pjmessi/golang-practice/pkg/strutil"
	"github.com/pjmessi/golang-practice/pkg/timeutil"
)

type HandlerImpl struct {
	secret          []byte
	jwtExpTimeInSec int64
}

func NewHandler(logService logger.Service, appConfig *config.AppConfig) (Handler, error) {
	jwtExpTimeDurationStr := appConfig.JWT_EXPIRATION_TIME
	jwtExpTimeInSec, err := timeutil.ConvertDurationStrToSec(jwtExpTimeDurationStr)
	if err != nil {
		return nil, fmt.Errorf("jwt.NewHandler(): %w", err)
	}

	logService.Debug(fmt.Sprintf("jwt expiration time set to: '%s' i.e '%d' seconds", jwtExpTimeDurationStr, jwtExpTimeInSec))

	return &HandlerImpl{
		secret:          strutil.ConvertToBytes(appConfig.JWT_SECRET),
		jwtExpTimeInSec: jwtExpTimeInSec,
	}, nil
}

func (h *HandlerImpl) Generate(userId string, userEmail string) (jwtString string, err error) {
	claims := h.createClaims(userId, userEmail)
	token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, claims)

	jwtString, err = token.SignedString(h.secret)
	if err != nil {
		return "", err
	}

	return jwtString, nil
}

func (h *HandlerImpl) createClaims(userId string, userEmail string) jwtgo.MapClaims {
	return jwtgo.MapClaims{
		"user_id":    userId,
		"user_email": userEmail,
		"exp":        timeutil.GetTimestampAfterNSec(h.jwtExpTimeInSec),
	}
}

func (h *HandlerImpl) Verify(jwtStr string) (valid bool, userId string, userEmail string, err error) {
	token, err := jwtgo.Parse(jwtStr, func(token *jwtgo.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtgo.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return []byte(h.secret), nil
	})

	if err != nil {
		return false, "", "", err
	}

	if token.Valid {
		claims, ok := token.Claims.(jwtgo.MapClaims)
		if !ok {
			return false, "", "", fmt.Errorf("jwt.HandlerImpl.Verify(): Error getting claims from token")
		}
		userId, userIdOk := claims["user_id"].(string)
		userEmail, userEmailOk := claims["user_email"].(string)

		if !userIdOk || !userEmailOk {
			return false, "", "", fmt.Errorf("jwt.HandlerImpl.Verify(): User ID or User Email not found in claims")
		}

		return true, userId, userEmail, nil
	} else {
		return false, "", "", fmt.Errorf("jwt.HandlerImpl.Verify(): Token is not valid")
	}
}
