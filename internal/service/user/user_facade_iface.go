package user

import "context"

type Facade interface {
	RegisterUser(ctx context.Context, reqBytes []byte) ([]byte, error)
}
