package core

type Status[T any] struct {
	IsSuccess bool
	Body      T
	RESTCode  int
}

func NewStatus[T any](body T, success bool, restCode int) *Status[T] {
	return &Status[T]{
		IsSuccess: true,
		Body:      body,
		RESTCode:  200,
	}
}
