package main

type RuntimeError struct {
	token   Token
	message string
}

func (r *RuntimeError) Error() string { return r.message }

func (r RuntimeError) Is(target error) bool {
	_, ok := target.(*RuntimeError)
	return ok
}
