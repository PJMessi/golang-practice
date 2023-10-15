package exception

type InvalidReq struct {
	*Base
}

var invalidReqDefaultType string = "REQUEST_DATA.INVALID"
var invalidReqDefaultMsg string = "invalid request data"

func NewInvalidReqFromBase(baseEx Base) InvalidReq {
	baseExPointer := newException(&baseEx, invalidReqDefaultType, invalidReqDefaultMsg)
	return InvalidReq{Base: baseExPointer}
}

func NewInvalidReq() InvalidReq {
	baseExPointer := newException(nil, invalidReqDefaultType, invalidReqDefaultMsg)
	return InvalidReq{Base: baseExPointer}
}
