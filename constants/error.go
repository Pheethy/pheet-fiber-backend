package constants

const (
	ERROR_REQURED                                = "must be required"
	ERROR_USERNAME_WAS_DUPLICATE                 = "username was duplicate"
	ERROR_PRODUCTNAME_WAS_DUPLICATE              = "product name was duplicate"
	ERROR_ORGANIZE_NAME_WAS_DUPLICATE            = "organize_name was duplicate"
	ERROR_ORGANIZE_ALIAS_NAME_WAS_DUPLICATE      = "organize_alias_name was duplicate"
	ERROR_ORGANIZE_PRIVATE_TEL_NO_WAS_DUPLICATE  = "privtae_tel_no was duplicate"
	ERROR_ORGANIZE_DEPARTMENT_NAME_WAS_DUPLICATE = "department_name was duplicate"
	ERROR_KEYWORD_NAME_WAS_DUPLICATE             = "name was duplicate"
	ERROR_CAN_NOT_FIND_PRODUCT                   = "Can't find this product"
)

/* postgres */
const (
	ERROR_PQ_UNIQUE_PRODUCTNAME      = "pq: duplicate key value violates unique constraint \"unique_product_name\""
	ERROR_PQ_UNIQUE_USERNAME         = "pq: duplicate key value violates unique constraint \"unique_user_name\""
	ERROR_UNIQUE_ORGANIZE_NAME       = "pq: duplicate key value violates unique constraint \"unique_organize_name\""
	ERROR_UNIQUE_ORGANIZE_ALIAS_NAME = "pq: duplicate key value violates unique constraint \"unique_organize_alias_name\""
	ERROR_NO_ROWS                    = "sql: no rows in result set"
)

type TokenType string

const (
	Access  TokenType = "access"
	Refresh TokenType = "refresh"
	Admin   TokenType = "admin"
	APIKey  TokenType = "apikey"
)
