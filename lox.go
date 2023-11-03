package lox

type Error interface {
	error
	Line() int
}
