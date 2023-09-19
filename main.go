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

	_middle_handler "pheet-fiber-backend/middleware/handler"
	_middle_repo "pheet-fiber-backend/middleware/repository"
	_middle_usecase "pheet-fiber-backend/middleware/usecase"

	_monitor_handler "pheet-fiber-backend/service/monitor/handler"

	_users_handler "pheet-fiber-backend/service/users/handlers"
	_users_repo "pheet-fiber-backend/service/users/repository"
	_users_usecase "pheet-fiber-backend/service/users/usecase"

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
	midRepo := _middle_repo.NewMiddlewareRepository(psqlDB)
	userRepo := _users_repo.NewUsersRepository(psqlDB)
	/* Init Usecase */
	midUs := _middle_usecase.NewMiddlewareUsecase(midRepo)
	userUs := _users_usecase.NewUsersUsecase(cfg, userRepo)

	/* Init Handler */
	middleware := _middle_handler.NewMiddlewareHandler(cfg, midUs)
	monHandler := _monitor_handler.NewMonitorHandler(cfg)
	userHandler := _users_handler.NewUsersHandler(cfg, userUs)

	/* Init Validator */

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
	app.Use(middleware.Cors())
	app.Use(middleware.Logger())

	/* HealthCheck Service */
	app.Get("/", monHandler.HealthCheck)
	
	router := app.Group("users")
	r := route.NewRoute(router)

	r.RegisterUsers(userHandler, middleware)

	// Graceful Shutdown
	var c = make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		log.Println("Server is shutting down...")
		_ = app.Shutdown()
	}()

	//Listen to host:port
	log.Printf("Server is starting on %v", cfg.App().Url())
	app.Listen(cfg.App().Url())
}
