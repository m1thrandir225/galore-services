package cache

type KeyValueStore interface {
	SaveItem(name string) error
	GetItem(item string) (string, error)
	DeleteItem(name string) error
}
