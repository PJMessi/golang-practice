package jwt

import (
	"fmt"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/pjmessi/golang-practice/config"
	"github.com/pjmessi/golang-practice/pkg/logger"
	"github.com/pjmessi/golang-practice/pkg/strutil"
	"github.com/pjmessi/golang-practice/pkg/timeutil"
)

type UtilImpl struct {
	secret          []byte
	jwtExpTimeInSec int64
}

func NewUtil(loggerUtil logger.Util, appConfig *config.AppConfig) (Util, error) {
	jwtExpTimeDurationStr := appConfig.JWT_EXPIRATION_TIME
	jwtExpTimeInSec, err := timeutil.ConvertDurationStrToSec(jwtExpTimeDurationStr)
	if err != nil {
		return nil, fmt.Errorf("jwt.NewUtil: %w", err)
	}

	loggerUtil.Debug(fmt.Sprintf("jwt expiration time set to: '%s' i.e '%d' seconds", jwtExpTimeDurationStr, jwtExpTimeInSec))

	return &UtilImpl{
		secret:          strutil.ConvertToBytes(appConfig.JWT_SECRET),
		jwtExpTimeInSec: jwtExpTimeInSec,
	}, nil
}

func (u *UtilImpl) Generate(userId string, userEmail string) (jwtString string, err error) {
	claims := u.createJwtClaims(userId, userEmail)
	token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, claims)

	jwtString, err = token.SignedString(u.secret)
	if err != nil {
		return "", err
	}

	return jwtString, nil
}

func (u *UtilImpl) createJwtClaims(userId string, userEmail string) jwtgo.MapClaims {
	return jwtgo.MapClaims{
		"user_id":    userId,
		"user_email": userEmail,
		"exp":        timeutil.GetTimestampAfterNSec(u.jwtExpTimeInSec),
	}
}
