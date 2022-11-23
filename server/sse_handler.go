package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sendchamp-go-test/db"
	"sendchamp-go-test/services"
)

func (s *Server) ServerSentEventHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		services.WsUpgrader.CheckOrigin = func(r *http.Request) bool {
			return true
		}
		var err error
		conn, err := services.WsUpgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Printf("could not upgrade: %s\n", err.Error())
			return
		}
		db.NewSseRepo(conn)
	}
}
