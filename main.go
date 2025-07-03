package main

import (
	database "api-sample-with-echo-ddd/config"
	"api-sample-with-echo-ddd/infra"
	router "api-sample-with-echo-ddd/interface"
	"api-sample-with-echo-ddd/interface/handler"
	"api-sample-with-echo-ddd/usecase"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("failed to load .env file")
	}

	config := database.NewConfig()
	db := database.NewDB(config)
	e := echo.New()

	// user
	userRepo := infra.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUsecase)
	router.InitRouting(e, userHandler)

	e.Logger.Fatal(e.Start(":8080"))
}
