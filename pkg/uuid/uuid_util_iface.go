package uuid

type Util interface {
	GenUuidV4() (string, error)
}
