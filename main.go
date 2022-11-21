package main

import (
	"log"
	"net/http"
	"sendchamp-go-test/config"
	"sendchamp-go-test/db"
	"sendchamp-go-test/server"
	"sendchamp-go-test/services"
	"time"
)

func main() {
	http.DefaultClient.Timeout = time.Second * 10
	conf, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	gormDB := db.GetDB(conf)
	userRepo := db.NewUserRepo(gormDB)
	userService := services.NewUserService(userRepo, conf)

	s := &server.Server{
		Config:         conf,
		UserRepository: userRepo,
		UserService:    userService,
	}
	s.Start()
}
