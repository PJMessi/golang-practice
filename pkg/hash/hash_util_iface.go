package hash

type Util interface {
	GenerateHash(plainString string) (hashedString string, err error)
	VerifyHash(hashString string, plainString string) bool
}
