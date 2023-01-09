package iterator

type Iterator[ELEMENT_TYPE any] interface {
	Length() int
	Next() ELEMENT_TYPE
}
