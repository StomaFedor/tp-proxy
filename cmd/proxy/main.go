package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"tp-proxy/pkg/api"
	handler "tp-proxy/pkg/api/handlers"
	"tp-proxy/pkg/proxy"
	"tp-proxy/pkg/repository"
	"tp-proxy/pkg/service"
)

// @title Proxy API
// @version 1.0
// @description API Server for Proxy Application

// @host localhost:8000
// @BasePath /
func main() {
	ctx := context.Background()

	var protoc, cert, key string
	flag.StringVar(&protoc, "protocol", "http", "")
	flag.StringVar(&key, "key", "ca.key", "")
	flag.StringVar(&cert, "crt", "ca.crt", "")
	flag.Parse()

	srvProxy := proxy.NewServerProxy(protoc, key, cert)

	db, err := repository.NewPostgresDB(ctx, repository.PostgresConfig{
		Host:     "db_proxy",
		Port:     "5432",
		Username: "postgres",
		DBName:   "postgres",
		SSLMode:  "disable",
		Password: "1474",
	})
	if err != nil {
		log.Fatal("error when connecting to the database")
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	proxy := proxy.NewProxy(services)

	go func() {
		err = srvProxy.Serve("8080", proxy.InitHandlers())
		if err != nil {
			log.Fatal("error when starting the proxy server")
		}
	}()

	srv := new(api.Server)

	go func() {
		if err = srv.Serve("8000", handlers.InitRoutes()); err != nil {
			log.Fatal("error occurred on server shutting down")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err = srv.Shutdown(context.Background()); err != nil {
		log.Fatal("error occured on server shutting down")
	}
}
