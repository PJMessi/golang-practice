package password

type Util interface {
	IsStrong(plainPassword string) bool
	Hash(plainPw string) (string, error)
	IsHashCorrect(hashedPw string, plainPw string) bool
}
