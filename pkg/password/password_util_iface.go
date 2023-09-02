package password

type Util interface {
	IsStrong(plainPassword string) (bool, error)
}
