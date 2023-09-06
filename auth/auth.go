package auth

type ServiceAuth interface {
	SignToken() string
}