package main

import (
	"fmt"
	"log"

	"github.com/ex-rate/auth-service/internal/config"
	"github.com/ex-rate/auth-service/internal/handler"
	service "github.com/ex-rate/auth-service/internal/service/registration"
	token "github.com/ex-rate/auth-service/internal/service/token"
	storage "github.com/ex-rate/auth-service/internal/storage/postgres"
	"github.com/ex-rate/auth-service/internal/storage/postgres/registration"
	"github.com/gin-gonic/gin"
)

func main() {
	r := setup()

	if err := r.Run(); err != nil {
		log.Fatal("unable to start server: ", err)
	}
}

func setup() *gin.Engine {
	// config
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal("unable to load config: ", err)
	}

	// db connect
	dbStr := fmt.Sprintf("user=%s dbname=%s sslmode=disable", conf.PostgresUser, conf.PostgresDB)
	conn, err := storage.Connect(dbStr)
	if err != nil {
		log.Fatal("unable to connect db: ", err)
	}

	// repositories
	registrationRepo := registration.New(conn)

	// services
	tokenSrv := token.New(conf.SecretKey)
	registrationSrv := service.New(registrationRepo, tokenSrv)

	handler := handler.New(registrationSrv)

	r := gin.Default()

	r.GET("/signup", handler.Registration)

	return r
}
