package adt

type constructor int

const (
	Null_ constructor = iota
	Cons_
)

type List[T any] interface {
	Dispatch() constructor
	AsNull() null[T]
	AsCons() cons[T]
}

type null[T any] struct{}

type cons[T any] struct {
	Head T
	Tail List[T]
}

func Null[T any]() List[T] {
	return &null[T]{}
}

func Cons[T any](head T, tail List[T]) List[T] {
	return &cons[T]{
		Head: head,
		Tail: tail,
	}
}

func (n *null[T]) Dispatch() constructor {
	return Null_
}

func (n *null[T]) AsNull() null[T] {
	return *n
}

func (n *null[T]) AsCons() cons[T] {
	panic("null is not cons")
}

func (c *cons[T]) Dispatch() constructor {
	return Cons_
}

func (c *cons[T]) AsNull() null[T] {
	panic("cons is not null")
}

func (c *cons[T]) AsCons() cons[T] {
	return *c
}

func length[T any](l List[T]) int {
	switch l.Dispatch() {
	case Null_:
		return 0
	case Cons_:
		c := l.AsCons()
		return 1 + length(c.Tail)
	default:
		panic("invalid type constructor")
	}
}
