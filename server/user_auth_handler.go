package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sendchamp-go-test/errors"
	"sendchamp-go-test/models"
	"sendchamp-go-test/server/response"
)

func (s *Server) HandleSignup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := decode(c, &user); err != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, err)
			return
		}
		userResponse, err := s.UserService.SignupUser(&user)
		if err != nil {
			err.Respond(c)
			return
		}
		response.JSON(c, "Signup successful", http.StatusCreated, userResponse, nil)
	}
}

func (s *Server) handleLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginRequest models.LoginRequest
		if err := decode(c, &loginRequest); err != nil {
			response.JSON(c, "", errors.ErrBadRequest.Status, nil, err)
			return
		}
		userResponse, err := s.UserService.LoginUser(&loginRequest)
		if err != nil {
			response.JSON(c, "", err.Status, nil, err)
			return
		}
		response.JSON(c, "login successful", http.StatusOK, userResponse, nil)
	}
}
