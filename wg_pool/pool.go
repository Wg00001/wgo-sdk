package wg_pool

type Pool[T any] interface {
	Init(option Option) error
	Get() (T, error)
	CloseAll()
	Len() int
	Open() (T, error)
}

type Option interface {
}
