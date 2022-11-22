package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sendchamp-go-test/errors"
	_ "sendchamp-go-test/errors"
	"sendchamp-go-test/models"
	"sendchamp-go-test/server/response"
)

func (s *Server) HandleCreateTask() gin.HandlerFunc {
	return func(c *gin.Context) {
		var task models.Task
		if err := decode(c, &task); err != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, err)
			return
		}
		userResponse, err := s.TaskService.CreateTask(&task)
		if err != nil {
			err.Respond(c)
			return
		}
		response.JSON(c, "Task Creation successful", http.StatusCreated, userResponse, nil)
	}
}

func (s *Server) HandleUpdateTask() gin.HandlerFunc {
	return func(c *gin.Context) {
		//_, user, err := GetValuesFromContext(c)
		//if err != nil {
		//	err.Respond(c)
		//	return
		//}
		//id, errr := strconv.ParseUint(c.Param("companyID"), 10, 32)
		id := c.Param("taskID")
		var updateTaskRequest models.Task
		if err := decode(c, &updateTaskRequest); err != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, err)
			return
		}
		err := s.TaskService.UpdateTask(&updateTaskRequest, id)
		if err != nil {
			err.Respond(c)
			return
		}
		response.JSON(c, "task updated successfully", http.StatusOK, nil, nil)
	}
}

func (s *Server) HandleGetTaskDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("taskID")
		task, err := s.TaskService.GetAllTasks(id)
		if err != nil {
			response.JSON(c, "", http.StatusInternalServerError, nil, errors.New("internal server error", http.StatusInternalServerError))
			return
		}
		response.JSON(c, "retrieved task successfully", http.StatusOK, gin.H{"task details": task}, nil)
	}
}

func (s *Server) handleDeleteTask() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, user, err := GetValuesFromContext(c)
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