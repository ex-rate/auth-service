package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ex-rate/auth-service/internal/config"
	"github.com/ex-rate/auth-service/internal/handler"
	"github.com/ex-rate/auth-service/internal/service"
	storage "github.com/ex-rate/auth-service/internal/storage/postgres"
	"github.com/ex-rate/auth-service/internal/storage/postgres/registration"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("hello")
	r := setup()

	if err := r.Run(); err != nil {
		log.Fatal("unable to start server: ", err)
	}
}

func setup() *gin.Engine {
	if len(os.Args) == 0 {
		log.Fatal("expected command line args: path to config file")
	}

	path := os.Args[1]

	conf, err := config.LoadConfig("../..", path)
	if err != nil {
		log.Fatal("unable to load config: ", err)
	}

	dbStr := fmt.Sprintf("user=%s dbname=%s sslmode=disable", conf.PostgresUser, conf.PostgresDB)
	conn, err := storage.Connect(dbStr)
	if err != nil {
		log.Fatal("unable to connect db: ", err)
	}

	registrationRepo := registration.New(conn)

	service := service.New(registrationRepo, conf.SecretKey)
	handler := handler.New(service)

	r := gin.Default()

	r.GET("/signup", handler.Registration)

	return r
}
