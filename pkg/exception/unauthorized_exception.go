package exception

type Unauthorized struct {
	*Base
}

var unauthorizedDefaultType string = "UNAUTHORIZED"
var unauthorizedDefaultMsg string = "user not authorized"

func NewUnauthorizedFromBase(baseEx Base) Unauthenticated {
	baseExPointer := newException(&baseEx, unauthorizedDefaultType, unauthorizedDefaultMsg)
	return Unauthenticated{Base: baseExPointer}
}

func NewUnauthorized(baseEx Base) Unauthorized {
	baseExPointer := newException(&baseEx, unauthorizedDefaultType, unauthorizedDefaultMsg)
	return Unauthorized{Base: baseExPointer}
}
