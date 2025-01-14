package main

func Ptr[T any](t T) *T {
	return &t
}
