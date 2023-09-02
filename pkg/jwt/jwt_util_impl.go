package jwt

import (
	"fmt"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/pjmessi/go-database-usage/config"
	"github.com/pjmessi/go-database-usage/pkg/strutil"
	"github.com/pjmessi/go-database-usage/pkg/timeutil"
)

type UtilImpl struct {
	Util
	secret              []byte
	expirationTimestamp int64
}

func NewUtil(appConfig *config.AppConfig) (Util, error) {
	expirationTime, err := timeutil.GetTimestampAfterDurationStr(appConfig.JWT_EXPIRATION_TIME)
	if err != nil {
		return nil, fmt.Errorf("jwt.NewUtil: %w", err)
	}

	return &UtilImpl{
		secret:              strutil.ConvertToBytes(appConfig.JWT_SECRET),
		expirationTimestamp: expirationTime,
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
		"exp":        u.expirationTimestamp,
	}
}
