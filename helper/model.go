package helper

import (
	"task-management/model/domain"
	"task-management/model/web"
)

func ToTaskResponse(task domain.Task) web.TaskResponse {
	return web.TaskResponse{
		Id: task.Id,
		Title: task.Title,
		Description: task.Description,
		Status: task.Status,
	}
}

func ToTaskResponses(tasks []domain.Task) []web.TaskResponse {
	var taskResponses []web.TaskResponse
	for _, task := range tasks {
		taskResponses = append(taskResponses, ToTaskResponse(task))
	}
	return taskResponses
}