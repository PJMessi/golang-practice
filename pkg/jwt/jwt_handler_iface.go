package jwt

type JwtPayload struct {
	UserId    string
	UserEmail string
}

type Handler interface {
	Generate(payload JwtPayload) (jwtString string, err error)
	Verify(jwtStr string) (valid bool, payload JwtPayload, err error)
}
