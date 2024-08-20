package constants

type TokenType string

const (
	Access  TokenType = "access"
	Refresh TokenType = "refresh"
	Admin   TokenType = "admin"
	ApiKey  TokenType = "apikey"
)

const (
	USER_ROLE_CUSTOMER = 1
	USER_ROLE_ADMIN    = 2
)
