package jwt

type Util interface {
	Generate(userId string, userEmail string) (jwtString string, err error)
}
