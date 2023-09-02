package exception

type Unauthenticated struct {
	*Base
}

var unauthenticatedDefaultType string = "UNAUTHENTICATED"
var unauthenticatedDefaultMsg string = "user not authenticated"

func NewUnauthenticatedFromBase(baseEx Base) Unauthenticated {
	baseExPointer := newException(&baseEx, unauthenticatedDefaultType, unauthenticatedDefaultMsg)
	return Unauthenticated{Base: baseExPointer}
}

func NewUnauthenticated() Unauthenticated {
	baseExPointer := newException(nil, unauthenticatedDefaultType, unauthenticatedDefaultMsg)
	return Unauthenticated{Base: baseExPointer}
}
