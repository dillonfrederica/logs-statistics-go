package internal

type Filter interface {
	Filter(origin string) (string, string)
}
