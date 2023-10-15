package exception

type AlreadyExists struct {
	*Base
}

var alreadyExistsDefaultType string = "RESOURCE.ALREADY_EXISTS"
var alreadyExistsDefaultMsg string = "resource already exists"

func NewAlreadyExistsFromBase(baseEx Base) AlreadyExists {
	baseExPointer := newException(&baseEx, alreadyExistsDefaultType, alreadyExistsDefaultMsg)
	return AlreadyExists{Base: baseExPointer}
}

func NewAlreadyExists() AlreadyExists {
	baseExPointer := newException(nil, alreadyExistsDefaultType, alreadyExistsDefaultMsg)
	return AlreadyExists{Base: baseExPointer}
}
