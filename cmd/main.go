package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/ex-rate/auth-service/internal/closer"
	"github.com/ex-rate/auth-service/internal/config"
	"github.com/ex-rate/auth-service/internal/handler"
	"github.com/ex-rate/auth-service/internal/router"
	"github.com/ex-rate/auth-service/internal/service"
	"github.com/ex-rate/auth-service/internal/service/auth"
	registration "github.com/ex-rate/auth-service/internal/service/registration"
	token "github.com/ex-rate/auth-service/internal/service/token"
	storage "github.com/ex-rate/auth-service/internal/storage/postgres"
	auth_repo "github.com/ex-rate/auth-service/internal/storage/postgres/auth"
	registration_repo "github.com/ex-rate/auth-service/internal/storage/postgres/registration"
	token_repo "github.com/ex-rate/auth-service/internal/storage/postgres/token"
)

const shutdownTimeout = 3 * time.Second

func Start(path, name string) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv := setup(path, name)

	if err := runServer(ctx, srv); err != nil {
		log.Fatal(err)
	}
}

func runServer(ctx context.Context, srv *http.Server) error {
	var (
		c = &closer.Closer{}
	)

	c.Add(srv.Shutdown)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen and serve err: %v", err)
		}
	}()

	log.Printf("listening on %s", srv.Addr)
	<-ctx.Done()

	log.Println("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := c.Close(shutdownCtx); err != nil {
		return fmt.Errorf("closer err: %v", err)
	}

	return nil
}

func setup(path, name string) *http.Server {
	// config
	conf, err := config.LoadConfig(path, name)
	if err != nil {
		log.Fatal("unable to load config: ", err)
	}

	// db connect
	dbStr := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", conf.PostgresUser, conf.PostgresDB, conf.PostgresPassword)
	conn, err := storage.Connect(dbStr)
	if err != nil {
		log.Fatal("unable to connect db: ", err)
	}

	// repositories
	registrationRepo := registration_repo.New(conn)
	tokenRepo := token_repo.New(conn)
	authRepo := auth_repo.New(conn)

	// services
	tokenSrv := token.New(conf.SecretKey, tokenRepo)
	registrationSrv := registration.New(registrationRepo, tokenSrv)
	authSrv := auth.New(authRepo, tokenSrv)

	service := service.New(registrationSrv, tokenSrv, authSrv)

	handler := handler.New(service)

	r := router.New(handler)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", conf.ServerHost, conf.ServerPort),
		Handler: r,
	}

	return srv
}
