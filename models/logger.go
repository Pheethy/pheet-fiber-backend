package models

import (
	"fmt"
	"log"
	"os"
	"pheet-fiber-backend/service/logger"
	"pheet-fiber-backend/service/utils"
	"github.com/gofiber/fiber/v2"

	"strings"
	"time"

)

type Logger struct {
	Time       string `json:"time"`
	Ip         string `json:"ip"`
	Method     string `json:"method"`
	StatusCode int    `json:"status_code"`
	Path       string `json:"path"`
	Query      any    `json:"query"`
	Body       any    `json:"body"`
	Response   any    `json:"response"`
}

func (l *Logger) Print() logger.ILogger {
	utils.Debug(l)
	return l
}

func (l *Logger) Save() {
	var data = utils.OutPut(l)
	fileName := fmt.Sprintf("./assets/logs/logger_%v.json", strings.ReplaceAll(time.Now().Format("2006-10-15"), "-", ""))
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Opening file failed:%v", err)
	}
	defer file.Close()

	file.WriteString(string(data) + "\n")
}

func (l *Logger) SetQuery(c *fiber.Ctx) {
	var query any
	if err := c.QueryParser(query); err != nil {
		log.Printf("query parser error: %v", err)
	}

	l.Query = query
}

func (l *Logger) SetBody(c *fiber.Ctx) {
	var body any
	if err := c.BodyParser(body); err != nil {
		log.Printf("body parser error: %v", err)
	}

	switch l.Body {
	case "/v1/users/signup":
		l.Body = "never gonna give you up!!!"
	default:
		l.Body = body
	}
}

func (l *Logger) SetResp(resp any) {
	l.Response = resp
}