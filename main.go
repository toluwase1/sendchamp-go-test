package main

import (
	"log"
	"net/http"
	"sendchamp-go-test/config"
	"sendchamp-go-test/db"
	"sendchamp-go-test/server"
	"time"
)

func main() {
	http.DefaultClient.Timeout = time.Second * 10
	conf, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	gormDB := db.GetDB(conf)
	authRepo := db.NewAuthRepo(gormDB)
	authService := services.NewAuthService(authRepo, conf)

	s := &server.Server{
		Config:         conf,
		AuthRepository: authRepo,
		AuthService:    authService,
	}
	s.Start()
}
