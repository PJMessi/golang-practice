package exception

type FailedPrecondition struct {
	*Base
}

var failedExceptionDefaultType string = "FAILED_PRECONDITION"
var failedExceptionDefaultMsg string = "precondition not met"

func NewFailedPreconditionFromBase(baseEx Base) FailedPrecondition {
	baseExPointer := newException(&baseEx, failedExceptionDefaultType, failedExceptionDefaultMsg)
	return FailedPrecondition{Base: baseExPointer}
}

func NewFailedPrecondition() FailedPrecondition {
	baseExPointer := newException(nil, failedExceptionDefaultType, failedExceptionDefaultMsg)
	return FailedPrecondition{Base: baseExPointer}
}
