package main

import (
	"main/auth"
	"main/handler"
	"main/repository"
	"main/route"
	"main/service"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB
var err error

func main() {
	db, err = sqlx.Open("mysql", "root:Bizcuitware@/coffee_list")
	if err != nil {
		panic(err)
	}
	proRepo := repository.NewProductRepository(db)
	proService := service.NewProductService(proRepo)
	proHandler := handler.NewProductHandler(proService)

	app := *fiber.New()
	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	app.Get("/tokenz", auth.AccessToken("==signature=="))
	protected := app.Group("", auth.Protect([]byte("==signature==")))

	r := route.NewRoute(protected)
	r.RegisterProduct(proHandler)

	app.Listen(":8080")
}