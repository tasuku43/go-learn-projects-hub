package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/tasuku43/go-learn-projects-hab/todo-list-app/models"
	"github.com/tasuku43/go-learn-projects-hab/todo-list-app/pkg/presentation/rest/handlers"
	"os"
)

func main() {
	err := godotenv.Load() // .env ファイルから環境変数を読み込む
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	r := gin.Default()
	dbConn := initDB() // データベース接続の初期化
	defer dbConn.Close()
	taskHandler := handlers.NewTaskHandler(dbConn)

	r.POST("/tasks", taskHandler.CreateTask)
	r.GET("/tasks", taskHandler.GetTasks)
	r.GET("/tasks/:id", taskHandler.GetTask)
	r.PUT("/tasks/:id", taskHandler.UpdateTask)
	r.DELETE("/tasks/:id", taskHandler.DeleteTask)

	appPort := os.Getenv("APP_PORT")
	r.Run(":" + appPort)
}

func initDB() *gorm.DB {
	connect := models.GormConnect()
	connect.LogMode(true)
	return connect
}
