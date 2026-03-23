package controller

import (
	"net/http"
	"strconv"
	"task-management/helper"
	"task-management/model/web"
	"task-management/service"
	"github.com/julienschmidt/httprouter"
)

type TaskControllerImpl struct {
	TaskService service.TaskService
}

func NewTaskController(taskService service.TaskService) TaskController {
	return &TaskControllerImpl{
		TaskService: taskService,
	}
}

// Insert implements [TaskController].
func (controller *TaskControllerImpl) Insert(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	taskCreateRequest := web.TaskCreateRequest{}
	helper.ReadFromRequestBody(request, &taskCreateRequest)

	taskResponse := controller.TaskService.Insert(request.Context(), taskCreateRequest)
	webResponse := web.WebResponse{
		Code: http.StatusOK,
		Status: "OK",
		Data: taskResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

// Update implements [TaskController].
func (controller *TaskControllerImpl) Update(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	taskUpdateRequest := web.TaskUpdateRequest{}
	helper.ReadFromRequestBody(request, &taskUpdateRequest)

	taskId := params.ByName("taskId")
	id, err := strconv.Atoi(taskId)
	helper.PanicIfError(err)

	taskUpdateRequest.Id = id

	taskResponse := controller.TaskService.Update(request.Context(), taskUpdateRequest)
	webResponse := web.WebResponse{
		Code: http.StatusOK,
		Status: "OK",
		Data: taskResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}


// Delete implements [TaskController].
func (controller *TaskControllerImpl) Delete(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	taskId := params.ByName("taskId")
	id, err := strconv.Atoi(taskId)
	helper.PanicIfError(err)

	controller.TaskService.Delete(request.Context(), id)
	webResponse := web.WebResponse{
		Code: http.StatusOK,
		Status: "OK",
	}

	helper.WriteToResponseBody(writer, webResponse)
}

// FindById implements [TaskController].
func (controller *TaskControllerImpl) FindById(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	taskId := params.ByName("taskId")
	id, err := strconv.Atoi(taskId)
	helper.PanicIfError(err)

	taskResponse := controller.TaskService.FindById(request.Context(), id)
	webResponse := web.WebResponse{
		Code: http.StatusOK,
		Status: "OK",
		Data: taskResponse,
	}

	helper.WriteToResponseBody(writer, webResponse)
}

// FindAll implements [TaskController].
func (controller *TaskControllerImpl) FindAll(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	taskResponses := controller.TaskService.FindAll(request.Context())
	webResponse := web.WebResponse{
		Code: http.StatusOK,
		Status: "OK",
		Data: taskResponses,
	}

	helper.WriteToResponseBody(writer, webResponse)
}