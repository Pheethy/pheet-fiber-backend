package main

import (
	"main/handler"
	"main/repository"
	"main/route"
	"main/service"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"

	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB
var err error

func main() {
	db, err = sqlx.Open("mysql", "root:Bizcuitware@/coffee_list")
	if err != nil {
		panic(err)
	}
	customerRepo := repository.NewCustomerRepository(db)
	custService := service.NewCustomerService(customerRepo)
	custHandler := handler.NewCustomerHandler(custService)

	app := *fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	r := route.NewRoute(&app)
	r.RegisterProduct(custHandler)

	app.Listen(":8080")
}