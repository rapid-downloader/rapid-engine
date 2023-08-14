package service

type Runner interface {
	Run() error
	Close() error
}
