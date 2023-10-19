package jwt

type Handler interface {
	Generate(userId string, userEmail string) (jwtString string, err error)
}
