package exception

type NotFound struct {
	*Base
}

var notFoundDefaultType string = "RESOURCE.NOT_FOUND"
var notFoundDefaultMsg string = "resource not found"

func NewNotFoundFromBase(baseEx Base) NotFound {
	baseExPointer := newException(&baseEx, notFoundDefaultType, notFoundDefaultMsg)
	return NotFound{Base: baseExPointer}
}

func NewNotFound() NotFound {
	baseExPointer := newException(nil, notFoundDefaultType, notFoundDefaultMsg)
	return NotFound{Base: baseExPointer}
}
