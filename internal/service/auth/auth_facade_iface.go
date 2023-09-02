package auth

type Facade interface {
	Login(reqBytes []byte) ([]byte, error)
}
