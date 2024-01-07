package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/tasuku43/go-learn-projects-hab/todo-list-app/models"
	"net/http"
)

type TaskHandler struct {
	dbConn *gorm.DB
}

func NewTaskHandler(dbConn *gorm.DB) *TaskHandler {
	return &TaskHandler{dbConn: dbConn}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.BindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.dbConn.Create(&task)
	c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) GetTasks(c *gin.Context) {
	var tasks []models.Task
	h.dbConn.Find(&tasks)
	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) GetTask(c *gin.Context) {
	var task models.Task
	if err := h.dbConn.First(&task, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	var task models.Task
	if err := h.dbConn.First(&task, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	if err := c.BindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task.UpdatedAt = task.CreatedAt
	h.dbConn.Save(&task)
	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	var task models.Task
	if err := h.dbConn.First(&task, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	h.dbConn.Delete(&task)
	c.JSON(http.StatusNoContent, gin.H{})
}
