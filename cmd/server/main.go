package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"tung-tung-money-api/internal/config"
	"tung-tung-money-api/internal/handler"
	"tung-tung-money-api/internal/model"
	"tung-tung-money-api/internal/repository"
	"tung-tung-money-api/internal/router"
	"tung-tung-money-api/internal/service"
	"tung-tung-money-api/pkg/database"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "tung-tung-money-api/docs"
)

type customValidator struct {
	validator *validator.Validate
}

func (cv *customValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	cfg := config.Load()

	db := database.Connect(cfg.DSN())

	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	e := echo.New()
	e.Validator = &customValidator{validator: validator.New()}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	router.Register(e,
		handler.NewAuthHandler(userService, cfg.JWTSecret),
		handler.NewUserHandler(userService, cfg.JWTSecret),
	)

	go func() {
		if err := e.Start(":" + cfg.Port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("shutting down the server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("server gracefully stopped")
}
