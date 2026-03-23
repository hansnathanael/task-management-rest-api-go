package app

import (
	"task-management/controller"
	"task-management/exception"
	"github.com/julienschmidt/httprouter"
)

func NewRouter(taskController controller.TaskController) *httprouter.Router {
	router := httprouter.New()

	router.GET("/api/tasks", taskController.FindAll)
	router.GET("/api/tasks/:taskId", taskController.FindById)
	router.POST("/api/tasks", taskController.Insert)
	router.PUT("/api/tasks/:taskId", taskController.Update)
	router.DELETE("/api/tasks/:taskId", taskController.Delete)

	router.PanicHandler = exception.ErrorHandler

	return router
}