package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"net/http"

	"github.com/chud-lori/go-echo-temp/adapters/controllers"
	"github.com/chud-lori/go-echo-temp/adapters/repositories"
	"github.com/chud-lori/go-echo-temp/adapters/utils"
	"github.com/chud-lori/go-echo-temp/adapters/web"
	"github.com/chud-lori/go-echo-temp/domain/services"
	"github.com/chud-lori/go-echo-temp/infrastructure"
	"github.com/chud-lori/go-echo-temp/pkg/logger"
    "github.com/labstack/echo/v4"

	"github.com/joho/godotenv"
)


func APIKeyMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(ctx echo.Context) error {
        apiKey := ctx.Request().Header.Get("x-api-key")
		if apiKey != "secret-api-key" {
            fmt.Println("Faiuled api key")
			return ctx.JSON(http.StatusUnauthorized, map[string]string{
                "message": "unahorueh",
            })
		}

        err := next(ctx)

        return err
    }
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed load keys")
	}

	postgredb := infrastructure.NewPostgreDB()
	defer postgredb.Close()

	userRepository, _ := repositories.NewUserRepositoryPostgre(postgredb)
	userService := services.NewUserService(userRepository)
	userController := controllers.NewUserController(userService)

    e := echo.New()

	web.UserRouter(userController, e)

    e.Use(logger.LogTrafficMiddleware)
    e.Use(APIKeyMiddleware)

	// Run server in a goroutine
	go func() {
		log.Printf("Server is running on port %s", os.Getenv("APP_PORT"))
    if err := e.Start(fmt.Sprintf(":%s", os.Getenv("APP_PORT"))); err != nil && err != http.ErrServerClosed{
			e.Logger.Fatalf("HTTP server error: %v", err)
		}
	}()

	wait := utils.GracefullShutdown(context.Background(), 5*time.Second, map[string]utils.Operation{
		"database": func(ctx context.Context) error {
			return postgredb.Close()
		},
		"http-server": func(ctx context.Context) error {
			return e.Shutdown(ctx)
		},
	})

	<-wait
}

