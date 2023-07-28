package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"pheet-fiber-backend/config"
	"pheet-fiber-backend/migrations/database"
	"pheet-fiber-backend/route"

	_product_repo "pheet-fiber-backend/service/product/repository"
	_product_usecase "pheet-fiber-backend/service/product/usecase"
	_product_handler "pheet-fiber-backend/service/product/handler"

	_middle_repo "pheet-fiber-backend/middleware/repository"
	_middle_usecase "pheet-fiber-backend/middleware/usecase"
	_middle_handler "pheet-fiber-backend/middleware/handler"

	_monitor_handler "pheet-fiber-backend/service/monitor/handler"
	

	validate "pheet-fiber-backend/service/product/validator"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

func envPath() string {
	if len(os.Args) == 1 {
		return ".env"
	} else {
		return os.Args[1]
	}
}

func main() {
	var ctx = context.Background()
	var cfg = config.LoadConfig(envPath())
	var psqlDB = database.DBConnect(ctx, cfg.Db())
	defer psqlDB.Close()

	/* Init Repository */
	proRepo := _product_repo.NewProductRepository(psqlDB)
	midRepo := _middle_repo.NewMiddlewareRepository(psqlDB)

	/* Init Usecase */
	proService := _product_usecase.NewProductUsecase(proRepo)
	midUs := _middle_usecase.NewMiddlewareUsecase(midRepo)

	/* Init Handler */
	proHandler := _product_handler.NewProductHandler(proService)
	midHandler := _middle_handler.NewMiddlewareHandler(cfg, midUs)
	monHandler := _monitor_handler.NewMonitorHandler(cfg)

	/* Init Validator */
	var validate = validate.Validation{}


	/* Fiber server */
	app := fiber.New(fiber.Config{
		AppName:      cfg.App().Name(),
		BodyLimit:    cfg.App().BodyLimit(),
		ReadTimeout:  cfg.App().ReadTimeOut(),
		WriteTimeout: cfg.App().WriteTimeOut(),
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
	})

	/* middleware */
	app.Use(midHandler.Cors())

	/* HealthCheck Service */
	app.Get("/", monHandler.HealthCheck)


	v1 := app.Group("v1")
	r := route.NewRoute(v1)

	r.RegisterProduct(proHandler, validate)

	// Graceful Shutdown
	var c = make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func ()  {
		_ = <-c
		log.Println("Server is shutting down...")
		_ = app.Shutdown()
	}()

	//Listen to host:port
	log.Printf("Server is starting on %v", cfg.App().Url())
	app.Listen(cfg.App().Url())
}
