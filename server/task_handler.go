package server

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sendchamp-go-test/errors"
	_ "sendchamp-go-test/errors"
	"sendchamp-go-test/models"
	"sendchamp-go-test/server/response"
)

func (s *Server) HandleCreateTask() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, user, err := GetValuesFromContext(c)
		if err != nil {
			err.Respond(c)
			return
		}
		userId := user.ID
		var task models.Task
		if err := decode(c, &task); err != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, err)
			return
		}
		task.UserID = userId
		userResponse, errr := s.TaskService.CreateTask(&task)
		if errr != nil {
			log.Println(errr)
			errr.Respond(c)
			return
		}
		response.JSON(c, "Task Creation successful", http.StatusCreated, userResponse, nil)
	}
}

func (s *Server) HandleUpdateTask() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, user, err := GetValuesFromContext(c)
		if err != nil {
			err.Respond(c)
			return
		}
		//id, errr := strconv.ParseUint(c.Param("companyID"), 10, 32)
		id := c.Param("taskID")
		userId := user.ID
		var updateTaskRequest models.Task
		if err := decode(c, &updateTaskRequest); err != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, err)
			return
		}
		updateTaskRequest.UserID = userId
		err = s.TaskService.UpdateTask(&updateTaskRequest, id)
		if err != nil {
			err.Respond(c)
			return
		}
		response.JSON(c, "task updated successfully", http.StatusOK, nil, nil)
	}
}

func (s *Server) HandleGetTaskDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, _, err := GetValuesFromContext(c)
		if err != nil {
			err.Respond(c)
			return
		}
		id := c.Param("taskID")
		task, errr := s.TaskService.GetTaskById(id)
		if errr != nil {
			response.JSON(c, "", http.StatusInternalServerError, nil, errors.New("internal server error", http.StatusInternalServerError))
			return
		}
		response.JSON(c, "retrieved task successfully", http.StatusOK, gin.H{"task details": task}, nil)
	}
}

func (s *Server) handleDeleteTask() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, _, err := GetValuesFromContext(c)
		if err != nil {
			err.Respond(c)
			return
		}
		id := c.Param("taskID")
		if err := s.TaskService.DeleteTaskById(id); err != nil {
			err.Respond(c)
			return
		}

		response.JSON(c, "task successfully deleted", http.StatusOK, nil, nil)
	}
}

func GetValuesFromContext(c *gin.Context) (string, *models.User, *errors.Error) {
	var tokenI, userI interface{}
	var tokenExists, userExists bool
	log.Println("hello 1")
	if tokenI, tokenExists = c.Get("access_token"); !tokenExists {
		return "", nil, errors.New("forbidden", http.StatusForbidden)
	}
	log.Println("hello 2")
	if userI, userExists = c.Get("user"); !userExists {
		return "", nil, errors.New("forbidden", http.StatusForbidden)
	}
	log.Println("hello 3")
	token, ok := tokenI.(string)
	if !ok {
		return "", nil, errors.New("internal server error", http.StatusInternalServerError)
	}
	log.Println("hello 4")
	user, ok := userI.(*models.User)
	if !ok {
		return "", nil, errors.New("internal server error", http.StatusInternalServerError)
	}
	log.Println("hello 5")
	return token, user, nil
}
