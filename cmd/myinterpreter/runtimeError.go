package main

type RuntimeError struct {
	token   Token
	message string
}

func (r *RuntimeError) Error() string { return r.message }
