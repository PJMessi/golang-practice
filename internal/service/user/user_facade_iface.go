package user

type Facade interface {
	RegisterUser(reqBytes []byte) ([]byte, error)
}
