package main

import (
	"log"
	"os"
	"pheet-fiber-backend/auth"
	"pheet-fiber-backend/route"
	"pheet-fiber-backend/service/product/handler"
	"pheet-fiber-backend/service/product/repository"
	"pheet-fiber-backend/service/product/usecase"
	validate "pheet-fiber-backend/service/product/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

var psqlDB *sqlx.DB

func main() {
	var err error
	err = godotenv.Load(".env")/*Load Env*/
	if err != nil {
		log.Printf("please consider environment variable: %s", err)
	}

	psqlDB, err = sqlx.Open("postgres", os.Getenv("PSQL_DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	proRepo := repository.NewProductRepository(psqlDB)
	proService := service.NewProductUsecase(proRepo)
	proHandler := handler.NewProductHandler(proService)
	var validate = validate.Validation{}

	err = psqlDB.Ping()
	if err != nil {
		log.Println(err)
	}

	app := *fiber.New()
	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Bizcuitware Web!!!")
	})

	app.Post("/register", proHandler.SignUp)
	app.Post("/login", proHandler.Login)

	productGroup := app.Group("", auth.Protect([]byte(os.Getenv("SIGN"))))
	r := route.NewRoute(productGroup)
	r.RegisterProduct(proHandler, validate)

	app.Listen(os.Getenv("PORT"))
}