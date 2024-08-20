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
	ERROR_USERNAME_OR_PASSWORD_IS_INVALID        = "username or password is invalid"
	ERROR_PASSWORD_NOT_MATCH                     = "password is invalid"
	ERROR_OAUTH_NOT_FOUND                        = "oauth not found"
	ERROR_CAN_NOT_FIND_PRODUCT                   = "Can't find this product"
	ERROR_PRODUCT_NOT_FOUND                      = "product not found"
	ERROR_ROLES_NOT_FOUND                        = "roles not found"
	ERROR_USER_NOT_FOUND                         = "user not found"
)

/* postgres */
const (
	ERROR_PQ_UNIQUE_PRODUCTNAME      = "pq: duplicate key value violates unique constraint \"unique_product_name\""
	ERROR_PQ_UNIQUE_USERNAME         = "pq: duplicate key value violates unique constraint \"users_username_key\""
	ERROR_PQ_UNIQUE_EMAIL            = "pq: duplicate key value violates unique constraint \"users_email_key\""
	ERROR_UNIQUE_ORGANIZE_NAME       = "pq: duplicate key value violates unique constraint \"unique_organize_name\""
	ERROR_UNIQUE_ORGANIZE_ALIAS_NAME = "pq: duplicate key value violates unique constraint \"unique_organize_alias_name\""
	ERROR_NO_ROWS                    = "sql: no rows in result set"
	ERROR_USERNAME_HAS_BEEN_USED     = "username has been used"
	ERROR_EMAIL_HAS_BEEN_USED        = "email has been used"
	ERROR_EMAIL_PATTERN_IS_INVALID   = "email is invalid pattern"
)

/* categories */
const (
	ERROR_CATEGORIES_NOT_FOUND         = "categories not found"
	ERROR_CATEGORIES_REQUEST_WAS_EMPTY = "categories request was empty"

	ERROR_ORDERS_NOT_FOUND = "order not found"
	ERROR_BOOK_NOT_FOUND   = "book not found"
)

/* JWT */
const (
	ERROR_UNKNOW_TOKEN_TYPE         = "unknow token type"
	ERROR_SIGNING_METHOD_IS_INVALID = "signing method is invalid"
	ERROR_JWT_FORMAT_IS_INVALID     = "jwt format is invalid"
	ERROR_JWT_WAS_EXPIRED           = "jwt token was expired"
	ERROR_CLAIMS_TYPE_IS_INVALID    = "claims type is invalid"
	ERROR_PARSE_TOKEN_FAILED        = "parse token failed"
	ERROR_NO_PERMISSION_TO_ACCESS   = "no permission to access"
	ERROR_CAN_GET_YOUR_PROFILE_ONLY = "you can get your profile only"
)
