package config

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

/* LoadConfig เป็นตัวดึงข้อมูลจาก env มาใส่ใน struct */
func LoadConfig(path string) Iconfig {
	envMap, err := godotenv.Read(path)
	if err != nil {
		log.Fatalf("load env failed: %v", err)
	}

	return &config{
		app: &app{
			host: envMap["APP_HOST"],
			port: func() int {
				p, err := strconv.Atoi(envMap["APP_PORT"])
				if err != nil {
					log.Fatalf("Load App Port Failed: %v", err)
				}
				return p
			}(),
			name:    envMap["APP_NAME"],
			version: envMap["APP_VERSION"],
			readTimeOut: func() time.Duration {
				t, err := strconv.Atoi(envMap["APP_READ_TIMEOUT"])
				if err != nil {
					log.Fatalf("Load Read Time Out Failed: %v", err)
				}
				return time.Duration(int64(t) * int64(math.Pow10(9)))
			}(),
			writeTimeOut: func() time.Duration {
				t, err := strconv.Atoi(envMap["APP_WRTIE_TIMEOUT"])
				if err != nil {
					log.Fatalf("Load Write Time Out Failed: %v", err)
				}
				return time.Duration(int64(t) * int64(math.Pow10(9)))
			}(),
			bodyLimit: func() int {
				limit, err := strconv.Atoi(envMap["APP_BODY_LIMIT"])
				if err != nil {
					log.Fatalf("Load Body Limit Failed: %v", err)
				}

				return limit
			}(),
			fileLimit: func() int {
				limit, err := strconv.Atoi(envMap["APP_FILE_LIMIT"])
				if err != nil {
					log.Fatalf("Load File Limit Failed: %v", err)
				}

				return limit
			}(),
			gcpBucket: envMap["APP_GCP_BUCKET"],
		},
		db: &db{
			host: envMap["DB_HOST"],
			port: func() int {
				p, err := strconv.Atoi(envMap["DB_PORT"])
				if err != nil {
					log.Fatalf("Load Port DB Failed: %v", err)
				}
				return p
			}(),
			protocol: envMap["DB_PROTOCOL"],
			username: envMap["DB_USERNAME"],
			password: envMap["DB_PASSWORD"],
			database: envMap["DB_DATABASE"],
			sslMode:  envMap["DB_SSL_MODE"],
			maxConnection: func() int {
				con, err := strconv.Atoi(envMap["DB_MAX_CONNECTIONS"])
				if err != nil {
					log.Fatalf("Load Max Connection Failed: %v", err)
				}
				return con
			}(),
		},
		jwt: &jwt{
			adminKey:  envMap["JWT_ADMIN_KEY"],
			secretKey: envMap["JWT_SECRET_KEY"],
			apiKey:    envMap["JWT_API_KEY"],
			accessExpiresAt: func() int {
				ex, err := strconv.Atoi(envMap["JWT_ACCESS_EXPIRES"])
				if err != nil {
					log.Fatalf("Load Access Expires Failed: %v", err)
				}
				return ex
			}(),
			refreshExpiresAt: func() int {
				ref, err := strconv.Atoi(envMap["JWT_REFRESH_EXPIRES"])
				if err != nil {
					log.Fatalf("Load Refresh Expires Failed: %v", err)
				}
				return ref
			}(),
		},
	}
}

// Struct
type config struct {
	app *app
	db  *db
	jwt *jwt
}

// Port Interface
type Iconfig interface {
	App() IAppConfig
	Db() IDbConfig
	Jwt() IJwtConfig
}

func (c *config) App() IAppConfig {
	return c.app
}

type IAppConfig interface {
	Url() string //host:port
	Name() string
	Version() string
	ReadTimeOut() time.Duration
	WriteTimeOut() time.Duration
	BodyLimit() int
	FileLimit() int
	GCPBucket() string
}

func (a *app) Url() string {
	return fmt.Sprintf("%s:%d", a.host, a.port)
}
func (a *app) Name() string {
	return a.name
}
func (a *app) Version() string {
	return a.version
}
func (a *app) ReadTimeOut() time.Duration {
	return a.readTimeOut
}
func (a *app) WriteTimeOut() time.Duration {
	return a.writeTimeOut
}
func (a *app) BodyLimit() int {
	return a.bodyLimit
}
func (a *app) FileLimit() int {
	return a.fileLimit
}
func (a *app) GCPBucket() string {
	return a.gcpBucket
}

type app struct {
	host         string
	port         int
	name         string
	version      string
	readTimeOut  time.Duration
	writeTimeOut time.Duration
	bodyLimit    int //bytes
	fileLimit    int //bytes
	gcpBucket    string
}

func (c *config) Db() IDbConfig {
	return c.db
}

type IDbConfig interface {
	Url() string
	MaxConns() int
}

func (d *db) Url() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", d.host, d.port, d.username, d.password, d.database, d.sslMode)
}
func (d *db) MaxConns() int {
	return d.maxConnection
}

type db struct {
	host          string
	port          int
	protocol      string
	username      string
	password      string
	database      string
	sslMode       string
	maxConnection int
}

func (c *config) Jwt() IJwtConfig {
	return c.jwt
}

type IJwtConfig interface {
	AdminKey() []byte
	SecretKey() []byte
	ApiKey() []byte
	AccessExpiresAt() int
	RefreshExpiresAt() int
}

func (j *jwt) AdminKey() []byte {
	return []byte(j.adminKey)
}
func (j *jwt) SecretKey() []byte {
	return []byte(j.secretKey)
}
func (j *jwt) ApiKey() []byte {
	return []byte(j.apiKey)
}
func (j *jwt) AccessExpiresAt() int {
	return j.accessExpiresAt
}
func (j *jwt) RefreshExpiresAt() int {
	return j.refreshExpiresAt
}

type jwt struct {
	adminKey         string
	secretKey        string
	apiKey           string
	accessExpiresAt  int //seconds
	refreshExpiresAt int //seconds
}
