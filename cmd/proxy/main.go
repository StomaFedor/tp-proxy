package main

import (
	"context"
	"flag"
	"log"
	"tp-proxy/pkg/proxy"
	"tp-proxy/pkg/repository"
	"tp-proxy/pkg/service"
)

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

	proxy := proxy.NewProxy(services)

	err = srvProxy.Serve("8080", proxy.InitHandlers())
	if err != nil {
		log.Fatal("error when starting the proxy server")
	}
}
