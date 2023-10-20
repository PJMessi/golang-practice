package jwt

type Handler interface {
	Generate(userId string, userEmail string) (jwtString string, err error)
	Verify(jwtStr string) (valid bool, userId string, userEmail string, err error)
}
