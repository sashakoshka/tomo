package iterator

type Iterator[ELEMENT_TYPE any] struct {
	index int
	slice []ELEMENT_TYPE
}

func New[ELEMENT_TYPE any] (
	slice []ELEMENT_TYPE,
) (
	iterator Iterator[ELEMENT_TYPE],
) {
	iterator.slice = slice
	return
}

func (iterator *Iterator[ELEMENT_TYPE]) Length () (length int) {
	return len(iterator.slice)
}

func (iterator *Iterator[ELEMENT_TYPE]) Next () (element ELEMENT_TYPE) {
	if !iterator.MoreLeft() { return }
	element = iterator.slice[iterator.index]
	iterator.index ++
	return
}

func (iterator *Iterator[ELEMENT_TYPE]) MoreLeft () (moreLeft bool) {
	return iterator.index < len(iterator.slice)
}
