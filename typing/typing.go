package typing

type Optional[T any] struct {
	Ok bool
	V  T
}
