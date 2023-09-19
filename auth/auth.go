package auth

type ServiceAuth interface {
	SignToken() string
}

type AdminAuth interface {
	SignToken() string
}
