package auth

import "context"

type Facade interface {
	Login(ctx context.Context, reqBytes []byte) ([]byte, error)
}
