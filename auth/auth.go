package auth

type ServiceAuth interface {
	SignToken() string
}

type AdminAuth interface {
	SignToken() string
}
type ApiKey interface {
	SignToken() string
}