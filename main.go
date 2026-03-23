package main

import (
	"net/http"
	"task-management/app"
	"task-management/controller"
	"task-management/helper"
	"task-management/middleware"
	"task-management/repository"
	"task-management/service"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db := app.NewDB()
	validate := validator.New()

	taskRepository := repository.NewTaskRepository()
	taskService := service.NewTaskService(taskRepository, db, validate)
	taskController := controller.NewTaskController(taskService)

	router := app.NewRouter(taskController)

	server := http.Server{
		Addr: "localhost:8080",
		Handler: middleware.NewAuthMiddleware(router),
	}

	err := server.ListenAndServe()
	helper.PanicIfError(err)
}