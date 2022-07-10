package main

import (
	"log"
	"main/auth"
	"main/handler"
	"main/repository"
	"main/route"
	"main/service"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

var db *sqlx.DB

func main() {
	var err error
	err = godotenv.Load("local.env")
	if err != nil {
		log.Printf("please consider environment variable: %s", err)
	}

	db, err = sqlx.Open("mysql", os.Getenv("DB_CONN"))
	if err != nil {
		panic(err)
	}

	proRepo := repository.NewProductRepository(db)
	proService := service.NewProductService(proRepo)
	proHandler := handler.NewProductHandler(proService)

	app := *fiber.New()
	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Bizcuitware Web!!!")
	})

	app.Post("/register", proHandler.SignUp)
	app.Post("/login", proHandler.Login)

	protected := app.Group("", auth.Protect([]byte(os.Getenv("SIGN"))))

	r := route.NewRoute(protected)
	r.RegisterProduct(proHandler)

	app.Listen(os.Getenv("PORT"))
}